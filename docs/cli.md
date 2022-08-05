# CLI

## get-app-access-token

```sh
twitch-stream-monitor get-app-access-token
```

Retrieve latest App Access Token (a.k.a. Bearer) from the storage.

> :warning: This command does not refresh token, if service is down it is possible that returned token will be expired.

## monitor

```sh
twitch-stream-monitor monitor
```

Launch Twitch Stream Monitor. Core functionality of the app, listens for incoming stream.online events and launches stream recorder.

See also:

* [Getting Credentials](credentials.md)
* [Webhook Configuration](webhook.md)
* [Environment Variables](environment-variables.md)

## resolve-username

```sh
twitch-stream-monitor resolve-username USERNAME
```

Find twitch user id by its username. Where `USERNAME` is Twitch channel name (case insensitive).

Example:

```
twitch-stream-monitor resolve-username twitch

12826
```

## list

```sh
twitch-stream-monitor list
```

Return list of all stream.online event subscriptions.

Example:

```
$ twitch-stream-monitor list

2022-08-05T19:35:00Z INF callback_url=https://example.com/ id=e41d5649-2a53-40e5-9747-bf97a438275e status=enabled user_id=123
2022-08-05T19:35:00Z INF callback_url=https://example.com/ id=b41750a5-edd2-4c16-ac19-df14a4a18620 status=enabled user_id=456
2022-08-05T19:35:00Z INF callback_url=https://example.com/ id=1baaf3c1-a0d1-405c-ac6f-356a2c67246e status=enabled user_id=789
2022-08-05T19:35:00Z INF callback_url=https://example.com/ id=95c2441d-e6a6-4d04-808d-846ef0ecb623 status=enabled user_id=159
```

## subscribe

```sh
twitch-stream-monitor subscribe USERNAME ...
twitch-stream-monitor subscribe BROADCASTER_ID ...
```

Subscribe to stream.online events for given twitch users.

Where `USERNAME` is Twitch channel name (case insensitive) and `BROADCASTER_ID` is an ID of the Twitch user (can be fetched via [resolve-username subcommand](#resolve-username`)).

Example:

```
$ twitch-stream-monitor subscribe twitch

2022-08-05T19:35:00Z INF Using https://example.com/ as a callback URL
2022-08-05T19:35:00Z INF Succesfully subscribed broadcasterID=12826 subID=369ee40b-7f52-4911-9509-c3bb68ef0229
```

## unsubscribe

```sh
twitch-stream-monitor unsubscribe SUBSCRIPTION_ID
```

Unsubscribe from stream.online by subscription id.

Where `SUBSCRIPTION_ID` is an id of the subscription returned by [subscribe](#subscribe) or [list](#list) subcommands.

Example:

```
$ twitch-stream-monitor unsubscribe 369ee40b-7f52-4911-9509-c3bb68ef0229

2022-08-05T19:35:00Z INF Succesfully unsubscribed subID=369ee40b-7f52-4911-9509-c3bb68ef0229
```
