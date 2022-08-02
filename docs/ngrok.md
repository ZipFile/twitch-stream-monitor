# Ngrok Support

Since Twitch requires a webhook to use a HTTPS protocol on the 443 port,
it introduces extra friction when running tool locally. In order to reduce amount PITA, [Ngrok](https://ngrok.com/) support is built-in.

Upon startup, the tool will try to access local instance of the Ngrok looking for the HTTPS tunnels to port specified in `TWITCH_MONITOR_PORT` (29177 by default). If found, it will be used as a callback URL.

Ngrok search location can be overridden by `TWITCH_MONITOR_NGROK_TUNNELS_URL` env variable. It must point to Ngrok's `/api/tunnels` API endpoint (e.g. `http://localhost:4040/api/tunnels`).

If there is multiple tunnels for a single port and protocol, the first will be used. `TWITCH_MONITOR_NGROK_TUNNEL_NAME` env variable can be used to limit search by tunnel name.

Note: If `TWITCH_EVENTSUB_CALLBACK_URL` is provided, Ngrok won't be queried.

**⚠️ WARNING**

Free version of Ngrok has a time limit for opened tunnels: after about 2 hours, tunnels are closed and URLs are invalidated. Means, no webhooks will be received after a while, rendering tool unusable. It is strongly advised to use free version only for local development and/or testing. If you intend to use the tool to constantly monitor streamer's activity, consider a proper server deployment.
