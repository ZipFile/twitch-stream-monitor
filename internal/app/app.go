package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
	app_access_token "github.com/ZipFile/twitch-stream-monitor/internal/app_access_token"
	helix_app_access_token "github.com/ZipFile/twitch-stream-monitor/internal/app_access_token/helix"
	ngrok_callback_url_getter "github.com/ZipFile/twitch-stream-monitor/internal/callback_url_getter/ngrok"
	config "github.com/ZipFile/twitch-stream-monitor/internal/config"
	http_handler "github.com/ZipFile/twitch-stream-monitor/internal/handler/http"
	noop_handler "github.com/ZipFile/twitch-stream-monitor/internal/handler/noop"
	streamlink_handler "github.com/ZipFile/twitch-stream-monitor/internal/handler/streamlink"
	monitor "github.com/ZipFile/twitch-stream-monitor/internal/monitor"
	helix_subscription_service "github.com/ZipFile/twitch-stream-monitor/internal/subscription_service/helix"
	username_resolver "github.com/ZipFile/twitch-stream-monitor/internal/username_resolver"
	helix_username_resolver "github.com/ZipFile/twitch-stream-monitor/internal/username_resolver/helix"
	utils "github.com/ZipFile/twitch-stream-monitor/internal/utils"
)

type app struct {
	Now               func() time.Time
	Config            *config.Config
	ConfigLoader      config.Loader
	TokenStore        app_access_token.Store
	TokenManager      app_access_token.Manager
	EventHandler      tsm.TwitchStreamOnlineEventHandler
	EventListener     tsm.TwitchOnlineSubscriptionService
	CallbackURLGetter tsm.CallbackURLGetter
	UsernameResolver  username_resolver.UsernameResolver
	Log               *zerolog.Logger
	HelixClient       *helix.Client
	KeepTokenUpToDate bool
}

func (app *app) loadConfig() error {
	if app.Config != nil {
		return nil
	}

	config, err := app.ConfigLoader.Load()

	if err != nil {
		return err
	}

	config.SetDefaults()

	app.Config = config

	return nil
}

func (app *app) initLogger() error {
	if app.Log != nil {
		return nil
	}

	log, err := NewLogger(
		app.Config.Log.Level,
		app.Config.Log.Pretty,
		app.Config.Log.Stdout,
	)

	if err != nil {
		return err
	}

	app.Log = log

	return nil
}

func (app *app) initHelixClient() error {
	if app.HelixClient != nil {
		return nil
	}

	helixClient, err := helix.NewClient(&helix.Options{
		ClientID:     app.Config.Twitch.ClientID,
		ClientSecret: app.Config.Twitch.ClientSecret,
	})

	if err != nil {
		return err
	}

	app.HelixClient = helixClient

	return nil
}

func (app *app) initTokenStore() error {
	if app.TokenStore != nil {
		return nil
	}

	tokenPath := app.Config.Twitch.AppAccessTokenLocation

	if tokenPath == "" {
		app.Log.Warn().Msg("Not requesting new app access token")
		app.TokenStore = &app_access_token.NoopStore{
			Token:   app.Config.Twitch.AppAccessToken,
			Expires: app.Now().AddDate(1, 0, 0),
		}
	} else {
		app.KeepTokenUpToDate = true
		app.TokenStore = app_access_token.NewFileStore(tokenPath)
	}

	return nil
}

func (app *app) initTokenManager() error {
	if app.TokenManager != nil {
		return nil
	}

	app.TokenManager = helix_app_access_token.NewManager(app.HelixClient, *app.Log)

	return nil
}

func (app *app) loadAppAccessToken() error {
	token, err := app_access_token.Get(
		app.TokenStore,
		app.TokenManager,
		app.Now(),
	)

	if err != nil {
		return err
	}

	app.HelixClient.SetAppAccessToken(token)

	return nil
}

func (app *app) refreshToken(ctx context.Context) {
	// TODO: Make token refresh more reactive
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := app.loadAppAccessToken()

			if err != nil {
				app.Log.Error().Err(err).Msg("Failed to refresh app access token")
			}
		case <-ctx.Done():
			break
		}
	}
}

func (app *app) initCallbackURLGetter() error {
	if app.CallbackURLGetter != nil {
		return nil
	}

	callbackURLGetter, err := ngrok_callback_url_getter.New(
		app.Config.Ngrok.TunnelsURL,
		app.Config.Ngrok.TunnelName,
		app.Config.Subscription.Port,
		*app.Log,
	)

	if err != nil {
		return err
	}

	app.CallbackURLGetter = callbackURLGetter

	return nil
}

func (app *app) initUsernameResolver() error {
	if app.UsernameResolver != nil {
		return nil
	}

	app.UsernameResolver = helix_username_resolver.New(app.HelixClient, *app.Log)

	return nil
}

func (app *app) initEventListener() error {
	if app.EventListener != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	callbackUrl, err := app.GetCallbackURL(ctx)

	if err != nil {
		return err
	}

	if callbackUrl == "" {
		return errors.New("Callback URL was not found")
	} else {
		app.Log.Info().Msgf("Using %s as a callback URL", utils.ObfuscateUrl(callbackUrl))
	}

	svc, err := helix_subscription_service.New(
		app.HelixClient,
		callbackUrl,
		app.Config.Subscription.WebhookSecret,
		utils.MakeAddr(app.Config.Subscription.Host, app.Config.Subscription.Port),
		*app.Log,
	)

	if err != nil {
		return err
	}

	app.EventListener = svc

	return nil
}

func (app *app) initEventHandler() error {
	switch app.Config.EventHandlerType {
	case "streamlink":
		return app.initStreamlinkEventHandler()
	case "http":
		return app.initHTTPNotificatorEventHandler()
	case "noop":
		return app.initNoopEventHandler()
	default:
		return fmt.Errorf("Unknown event handler: %s", app.Config.EventHandlerType)
	}
}

func (app *app) initStreamlinkEventHandler() error {
	if app.EventHandler != nil {
		return nil
	}

	handler, err := streamlink_handler.New(
		app.Config.Streamlink.Path,
		app.Config.Streamlink.FileDir,
		app.Config.Streamlink.LogDir,
		app.Config.Streamlink.ConfigPath,
		app.Config.Streamlink.KillTimeout,
		*app.Log,
	)

	if err != nil {
		return err
	}

	app.EventHandler = handler

	return nil
}

func (app *app) initHTTPNotificatorEventHandler() error {
	if app.EventHandler != nil {
		return nil
	}

	app.EventHandler = http_handler.New(
		app.Config.HTTPNotificator.URL,
		app.Config.HTTPNotificator.UserName,
		app.Config.HTTPNotificator.Password,
	)

	return nil
}

func (app *app) initNoopEventHandler() error {
	if app.EventHandler != nil {
		return nil
	}

	app.EventHandler = noop_handler.New()

	return nil
}

func (app *app) GetCallbackURL(ctx context.Context) (string, error) {
	url, err := app.Config.GetCallbackURL(ctx)

	if err != nil {
		return "", err
	}

	if url != "" {
		return url, nil
	}

	return app.CallbackURLGetter.GetCallbackURL(ctx)
}

type InitFunc func() error

func (app *app) Init() error {
	var err error

	initFuncs := []InitFunc{
		app.loadConfig,
		app.initLogger,
		app.initHelixClient,
		app.initTokenStore,
		app.initTokenManager,
		app.loadAppAccessToken,
		app.initCallbackURLGetter,
		app.initUsernameResolver,
		app.initEventListener,
		app.initEventHandler,
	}

	for _, initFunc := range initFuncs {
		if err = initFunc(); err != nil {
			return err
		}
	}

	return nil
}

func (app *app) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)

	signal.Notify(
		sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	go func() {
		<-sig
		cancel()
	}()

	if app.KeepTokenUpToDate {
		go app.refreshToken(ctx)
	}

	return monitor.Monitor(
		ctx,
		app.EventListener,
		app.EventHandler,
		monitor.Settings{
			UserIDs:                  app.Config.BroadcasterUserIDs,
			KeepExistingSubs:         app.Config.KeepExistingSubs,
			KeepNewSubs:              app.Config.KeepNewSubs,
			IgnoreStartupErrors:      app.Config.IgnoreStartupErrors,
			IgnoreSubscriptionErrors: app.Config.IgnoreSubscriptionErrors,
		},
		*app.Log,
	)
}

func (app *app) GetTwitchOnlineSubscriptionService() tsm.TwitchOnlineSubscriptionService {
	return app.EventListener
}

func (app *app) GetLogger() *zerolog.Logger {
	return app.Log
}

func (app *app) GetTokenManager() app_access_token.Manager {
	return app.TokenManager
}

func (app *app) GetTokenStore() app_access_token.Store {
	return app.TokenStore
}

func (app *app) GetUsernameResolver() username_resolver.UsernameResolver {
	return app.UsernameResolver
}

func New() App {
	return &app{
		Now:          time.Now,
		ConfigLoader: config.NewEnvironLoader(nil),
	}
}
