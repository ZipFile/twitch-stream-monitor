# Add a systemd user service to automatically run twitch-stream-monitor

After installing twitch-stream-monitor you might want to have the service start automatically. You can use a systemd unit, that runs in the user context.

## Allow user services to run without user logged in
According to https://wiki.archlinux.org/title/Systemd/User#Automatic_start-up_of_systemd_user_instances you need to enable the possibility to have user services running, even when the user is not logged in:

```sh
sudo loginctl enable-linger _username_
```

## Create necessary directory
In the user context of the user running twitch-stream-monitor create the following directory:
```sh
mkdir -p ~/.config/systemd/user/
```

## Add service unit
edit with your preferred editor the following file:
$EDITOR ~/.config/systemd/user/twitch-stream-monitor.service

Add the following content (adapt the paths to your setup (WorkingDirectory, ExecStart)):
```
[Unit]
Description=twitch-stream-monitor
Documentation=https://github.com/ZipFile/twitch-stream-monitor
After=network-online.target
Wants=network-online.target

[Service]
WorkingDirectory=/home/streamsaver/twitch
ExecStart=/home/streamsaver/go/bin/twitch-stream-monitor monitor
Restart=on-failure
RestartSec=3
RestartPreventExitStatus=3

[Install]
WantedBy=default.target
```

## Enable and start service
After saving the file, you need to reload the system daemon, enable and start the service:
```sh
systemctl --user daemon-reload 
systemctl --user enable twitch-stream-monitor.service 
systemctl --user start twitch-stream-monitor.service 
```

## query service status
To check if the service is running, you can query systemd:
```sh
systemctl --user status twitch-stream-monitor.service 
```

To show all logfiles from your service, you can use:
```sh
journalctl --user-unit twitch-stream-monitor.service
```

