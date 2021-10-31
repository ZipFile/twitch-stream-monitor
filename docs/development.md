# Development

Before you can start working on the code, you need the following tooling:

* Git.
* Go 1.17 or newer.
* Python 3. Used for running tests, should be installed in all modern distros.
* GNU Make. `build-essential` package in Ubuntu, for other distros refer to respective distro documentation.

## Clone the Code

```sh
git clone git@github.com:ZipFile/twitch-stream-monitor.git
cd twitch-stream-monitor
```

## Preparing the Environment

If your intention is to build the code, skip this step.

Otherwise, in order to run tests and experiment with actual stream recording you need a working Python 3 installation and a Streamlink.

If you have both installed - great. If not, you can setup a [Python's virtualenv](https://docs.python.org/3/library/venv.html):

```sh
python -m venv env
. env/bin/activate
```

And then install the Streamlink:

```sh
pip install streamlink
```

Next time whenever you're back to development, use:

```sh
. env/bin/activate
```

In case your streamlink located somewhere else or you do not want to use venv, set a `TWITCH_MONITOR_STREAMLINK_PATH` environment variable to the path pointing to the streamlink executable.

## Build

```sh
make
```

### Windows

```sh
GOOS=windows make
```

## Run Tests

```sh
make tests
```

## Format Code

```sh
make imports
make fmt
```

## .env Support

See [environment variables](environment-variables.md#env).

## Debug Logging

Default logging configuration is to output messages of the level INFO and above in a [CBOR](https://cbor.io/) format so it can be processed by structured-log-aggregating software like [fluentd](https://docs.fluentd.org/parser/json).

For debugging purposes you likely want more verbosity and human-readability. This can be achieved by exporting following environment variables:

```sh
TWITCH_MONITOR_LOG_LEVEL=trace
TWITCH_MONITOR_LOG_PRETTY=true
```

Either export them directly, or put into the `.env` file.

## Testing Webhooks

See [Ngrok Support](ngrok.md).
