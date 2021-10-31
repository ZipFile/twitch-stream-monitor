# Webhook Configuration

In order to catch streams, the tool relies on the Twitch [EventSub](https://dev.twitch.tv/docs/eventsub) feature.

The way how it works is essentially Twitch sending us a HTTP request when streamer goes online.

## Webhook Secret

Due to specifics of the subscription process, it is required to provide a secret value that later will be used by Twitch to verify validity of the webhook. This value can be generated using following command:

```sh
python -c 'from secrets import token_urlsafe; print(token_urlsafe(30))'
```

Use this value as a value for `TWITCH_WEBHOOK_SECRET` environment variable.

## Setting Up a Callback URL

For Twitch to send us an event, it needs to know where to send it to. Upon startup, app tries to create a subscription for each channel, during this stage it provides Twitch a URL where to deliver notifications to. The value of the URL is configured by `TWITCH_EVENTSUB_CALLBACK_URL` environment variable. If you're using [Ngrok](ngrok.md), the value will be pulled in automatically. Otherwise, you need to provide a HTTPS URL manually.

When app starts, it listens to HTTP requests on the port `29177`. Since this does not meet Twitch requirements, you need to put it behind reverse proxy. There are number of ways to do so. The easiest way is to use combo of [Nginx](https://nginx.org/en/docs/http/configuring_https_servers.html) + [Certbot](https://certbot.eff.org/instructions).

## Testing

To check if you configured everything properly, you can use [Twitch CLI](https://dev.twitch.tv/docs/cli) tool.

```sh
# Check that verification works:
twitch event verify-subscription streamup --secret TWITCH_WEBHOOK_SECRET -F https://HOSTNAME
# Trigger stream.online event:
twitch event trigger streamup --secret TWITCH_WEBHOOK_SECRET -F https://HOSTNAME -t CHANNEL_ID
```

Replace `TWITCH_WEBHOOK_SECRET` with the value from [Webhook Secret](#Webhook Secret) section, `HOSTNAME` with the hostname you have exposed the app (e.g. `43bb-123-45-67-89.ngrok.io` or `tsm.example.com`) and `CHANNEL_ID` with ID of the channel (e.g. `141981764`).

## ⚠️ WARNING

Twitch sends online event with approximately 1 minute delay, so if your streamer does not include "starting soon" section, there is a risk of losing content.
