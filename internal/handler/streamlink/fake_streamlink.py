#!/usr/bin/env python

import signal
import sys
from argparse import ArgumentParser


def noop(_, __):
    pass


def main():
    parser = ArgumentParser(description="Simulate Streamlink CLI.")

    parser.add_argument("--config")
    parser.add_argument("--logfile")
    parser.add_argument("--twitch-disable-ads", action="store_true")
    parser.add_argument("-o", "--output", required=True)
    parser.add_argument("url")
    parser.add_argument("format")

    args = parser.parse_args()

    if "error" in args.url:
        print("failed")
        return 1

    if not args.url.startswith("twitch.tv/"):
        print("twitch only")
        return 2

    if args.format != "best":
        print("best format only")
        return 3

    if not args.twitch_disable_ads:
        print("duh")
        return 4

    if args.config:
        with open(args.config, "rt", errors="strict") as f:
            print(f.read())

    if args.logfile:
        with open(args.logfile, "wt") as f:
            f.write("test log")

    with open(args.output, "wt") as f:
        f.write("test stream")

    print("recording")

    if "wait" not in args.url:
        return 0

    signal.signal(signal.SIGTERM, noop)
    signal.signal(signal.SIGINT, noop)

    if "forever" in args.url:
        while True:
            signal.pause()
    else:
        signal.pause()

    return 0


if __name__ == "__main__":
    sys.exit(main())
