package ngrok

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"

	tsm "github.com/ZipFile/twitch-stream-monitor/internal"
)

var defaultTunnelsURL = url.URL{
	Scheme: "http",
	Host:   "localhost:4040",
	Path:   "/api/tunnels",
}

type tunnelConfig struct {
	Addr string `json:"addr"`
}

type tunnel struct {
	PublicUrl string       `json:"public_url"`
	Name      string       `json:"name"`
	Proto     string       `json:"proto"`
	Config    tunnelConfig `json:"config"`
}

func (t *tunnel) MarshalZerologObject(e *zerolog.Event) {
	e.Str("public_url", t.PublicUrl).Str("name", t.Name).Str("proto", t.Proto).Str("addr", t.Config.Addr)
}

type tunnelsResponse struct {
	Tunnels []tunnel `json:"tunnels"`
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type URLGetter struct {
	URL    url.URL
	Name   string
	Port   uint64
	Client HTTPClient
	Log    zerolog.Logger
}

func New(apiUrl, name string, port uint64, log zerolog.Logger) (tsm.CallbackURLGetter, error) {
	parsedUrl, err := url.Parse(apiUrl)

	if err != nil {
		return nil, err
	}

	parsedUrl = defaultTunnelsURL.ResolveReference(parsedUrl)

	return &URLGetter{
		URL:    *parsedUrl,
		Name:   name,
		Port:   port,
		Client: &http.Client{},
		Log:    log.With().Str("component", "ngrok_callback_url_getter").Logger(),
	}, nil
}

func (n *URLGetter) match(t tunnel) bool {
	if t.Proto != "https" {
		return false
	}

	name := n.Name == "" || t.Name == n.Name
	port := n.Port == 0 || t.Config.Addr == fmt.Sprintf("http://localhost:%d", n.Port)

	return name && port
}

func (n *URLGetter) GetCallbackURL(ctx context.Context) (string, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, n.URL.String(), nil)

	if err != nil {
		return "", err
	}

	response, err := n.Client.Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	tunnels := tunnelsResponse{}
	err = json.NewDecoder(response.Body).Decode(&tunnels)

	if err != nil {
		return "", err
	}

	for _, tunnel := range tunnels.Tunnels {
		n.Log.Trace().Object("tunnel", &tunnel).Msg("Tunnel found")

		if n.match(tunnel) {
			return tunnel.PublicUrl, nil
		}
	}

	return "", nil
}
