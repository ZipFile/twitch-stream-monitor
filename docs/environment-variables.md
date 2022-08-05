# Environment Variables

Currently, the only way to configure the app is through [environment variables](https://en.wikipedia.org/wiki/Environment_variable).

## Setting Variables

Generally, it is expected for you to be familiar on how to set environment variables. Otherwise, refer to these guides:

* Linux: [Tutorial](https://www.digitalocean.com/community/tutorials/how-to-read-and-set-environmental-and-shell-variables-on-linux)
* Windows: [set](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/set_1), [setx](https://docs.microsoft.com/en-us/windows-server/administration/windows-commands/setx)

## .env

The tool will detect the presence of `.env` file in the current working directory and will load environment variables from there.

File format is simple: list of name-value pairs in a format `ENV_VARIABLE_NAME=value`, one per line.

Example:

```
TWITCH_BROADCASTERS=123,456
TWITCH_CLIENT_ID=abcdefghijklmnopqrstuvwxyz1234
TWITCH_CLIENT_SECRET=zyxwvutsrqponmlkjihgfedcba4321
TWITCH_WEBHOOK_SECRET=mywebhooksecret
TWITCH_APP_ACCESS_TOKEN_LOCATION=/mnt/data/twitch/app_access_token
TWITCH_MONITOR_STREAMLINK_PATH=/mnt/data/twitch/env/bin/streamlink
TWITCH_MONITOR_STREAMLINK_FILE_DIR=/mnt/data/twitch/recordings
TWITCH_MONITOR_STREAMLINK_LOG_DIR=/mnt/data/twitch/logs
TWITCH_MONITOR_LOG_PRETTY=true
```

Env variables that are already present in the environment will take the precedence. I.e. if you already have `TWITCH_BROADCASTERS` in your environ, its value will be used instead of the one from `.env`.

## Windows

Since Windows is picky to file naming, `.env` usage is not always possible. Alternatively, you can create `twitch-stream-monitor.bat` and set environment variables explicitly:

```bat
set TWITCH_BROADCASTERS=123,456
set TWITCH_CLIENT_ID=abcdefghijklmnopqrstuvwxyz1234
set TWITCH_CLIENT_SECRET=zyxwvutsrqponmlkjihgfedcba4321
set TWITCH_WEBHOOK_SECRET=mywebhooksecret
set TWITCH_APP_ACCESS_TOKEN_LOCATION=D:\twitch\app_access_token
set TWITCH_MONITOR_STREAMLINK_PATH=C:\Program Files (x86)\Streamlink\bin\streamlink.exe
set TWITCH_MONITOR_STREAMLINK_FILE_DIR=D:\twitch\recordings
set TWITCH_MONITOR_STREAMLINK_LOG_DIR=D:\twitch\logs
set TWITCH_MONITOR_LOG_PRETTY=true
start twitch-stream-monitor.exe monitor
```

## List

### TWITCH_BROADCASTERS

* Required: No
* Type: string list

List of channel IDs (not usernames) separated by comma.

IDs can be found using direct API query to Twitch:

```sh
curl 'https://api.twitch.tv/helix/users?login=USERNAME' -H 'Client-ID: MY_CLIENT_ID' -H "Authorization: Bearer APP_ACCESS_TOKEN"
```

### TWITCH_MONITOR_KEEP_EXISTING_SUBS

* Required: No
* Type: bool
* Default: `false`

Define behavior in case there is already a `stream.online` event subscription present for given broadcaster id.

Delete (`false`) or preserve (`true`) such subscriptions at exit.

### TWITCH_MONITOR_KEEP_NEW_SUBS

* Required: No
* Type: bool
* Default: `false`

Define behavior for newly-created `stream.online` event subscriptions.

Delete (`false`) or preserve (`true`) such subscriptions at exit.

### TWITCH_MONITOR_IGNORE_START_ERRORS

* Required: No
* Type: bool
* Default: `false`

Start listening for `stream.online` events even if it is known that they won't be handled properly.

Use at your own risk, mainly useful for debugging.

### TWITCH_MONITOR_IGNORE_SUB_ERRORS

* Required: No
* Type: bool
* Default: `false`

Keep going if creation of `stream.online` event subscription fails.

Use at your own risk, mainly useful for debugging.

### TWITCH_MONITOR_HANDLER

* Required: No
* Type: string enum
* Possible values: `streamlink`, `http`, `noop`
* Default: `streamlink`

See [Handlers](handlers.md).

### TWITCH_MONITOR_CHECK_TIMEOUT

* Required: No
* Type: duration
* Default: `5s`

Timeout for pre-startup checks. If check takes longer than specified value, consider it failed and abort the startup.

### TWITCH_CLIENT_ID

* Required: **Yes**
* Type: string

See [Credentials](credentials.md).

### TWITCH_CLIENT_SECRET

* Required: **Yes**
* Type: string

See [Credentials](credentials.md).

### TWITCH_APP_ACCESS_TOKEN

* Required: No
* Type: string

See [Credentials](credentials.md).

### TWITCH_APP_ACCESS_TOKEN_LOCATION

* Required: No
* Type: string

See [Credentials](credentials.md).

### TWITCH_EVENTSUB_CALLBACK_URL

* Required: **Yes** (if no Ngrok detected), No otherwise
* Type: string

See [Webhook](webhook.md).

### TWITCH_WEBHOOK_SECRET

* Required: Yes
* Type: string

See [Webhook](webhook.md).

### TWITCH_MONITOR_HOST

* Required: No
* Type: string

Hostname to bind to, if not set - listens to all the interfaces.

### TWITCH_MONITOR_PORT

* Required: No
* Type: unsigned int
* Default: `29177`

Port to bind to.

### TWITCH_MONITOR_STREAMLINK_PATH

* Required: No
* Type: string
* Default: `streamlink`

Path to the `streamlink` executable.

### TWITCH_MONITOR_STREAMLINK_FILE_DIR

* Required: No
* Type: string
* Default: `.` (current working directory)

Directory to store stream recordings.

### TWITCH_MONITOR_STREAMLINK_LOG_DIR

* Required: No
* Type: string
* Default: `.` (current working directory)

Directory to store stream recordings logs.

### TWITCH_MONITOR_STREAMLINK_CONFIG

* Required: No
* Type: string

Path to the streamlink configuration file.

### TWITCH_MONITOR_STREAMLINK_KILL_TIMEOUT

* Required: No
* Type: duration
* Default: `60s`

When shutting down the tool, the time for running streamlink recordings to gracefully shutdown. I.e. time after sending `SIGTERM` before issuing `SIGKILL`.

### TWITCH_MONITOR_HTTP_NOTIFICATOR_URL

* Required: No
* Type: string

URL to send callback to by `http` handler.

### TWITCH_MONITOR_HTTP_NOTIFICATOR_USERNAME

* Required: No
* Type: string

Username to use with HTTP Basic Auth for `http` handler.

### TWITCH_MONITOR_HTTP_NOTIFICATOR_PASSWORD

* Required: No
* Type: string

Password to use with HTTP Basic Auth for `http` handler.

### TWITCH_MONITOR_NGROK_TUNNELS_URL

* Required: No
* Type: string

See [Ngrok](ngrok.md).

### TWITCH_MONITOR_NGROK_TUNNEL_NAME

* Required: No
* Type: string

See [Ngrok](ngrok.md).

### TWITCH_MONITOR_LOG_LEVEL

* Required: No
* Type: string enum
* Default: `info`
* Possible values: `trace`, `debug`, `info`, `warn`, `error`, `fatal`, `panic`

Log verbosity.

### TWITCH_MONITOR_LOG_PRETTY

* Required: No
* Type: bool
* Default: `false`

Enable human-friendly logging.

### TWITCH_MONITOR_LOG_STDOUT

* Required: No
* Type: bool
* Default: `false`

Output logs to stdout instead of the stderr. Useful when running in docker.
