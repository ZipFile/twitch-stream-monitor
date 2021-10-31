# Getting Credentials

For tool to operate, we need a Twitch app. If you do not have one, follow [official documentation](https://dev.twitch.tv/docs/authentication#registration).

When asked, fill:

* **Name**: any name you want
* **OAuth Redirect URLs**: `http://localhost`
* **Category**: any would suffice, `Application Integration` more or less reflects the tool intention

After registration, you'll be presented with **Client IDs** and **Client Secret**. We need both, store them securely. **Client Secret** will be presented only once, so do not miss the opportunity.

After you received both **Client IDs** and **Client Secret**, use them to set `TWITCH_CLIENT_ID` and `TWITCH_CLIENT_SECRET` [environment variables](environment-variables.md) respectively.

## App Access Token

In order to communicate with a Twitch API, we need an App Access Token. By default app will request new token on each startup. You can override this behavior by providing one of the 2 environment variables:

* `TWITCH_APP_ACCESS_TOKEN` - If set, new token won't be requested and its value will be used instead.
* `TWITCH_APP_ACCESS_TOKEN_LOCATION` - File location to store app access token. Next time app starts, it will use token from this file.

It is advised to set `TWITCH_APP_ACCESS_TOKEN_LOCATION`, since it will prevent unnecessary token generation on each restart and will enable auto-reneval of the token upon expiration.
