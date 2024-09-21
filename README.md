BruteDrop

A simple but effective - I mean <i>brute</i> -  tool written in Go as response to brute force attacks.

## 1. Install BruteDrop binary

```
sudo make install
```

## 2. Add BruteDrop configuration file

To be sure you won't lock you out, you can test and see what's going on when BruteDrop runs by setting `DryRunMode` to `true` and follow log outputs with `sudo journalctl -u brutedrop -f`.

```
IptablesBinPath: /usr/bin/iptables
JournalctlBinPath: /usr/bin/journalctl

# DryRunMode: true
DryRunMode: false

# Set Logging to file path or "stdout"
LoggingTo: stdout
# LoggingTo: /var/log/brutedrop.log
LogEntriesSince: 2

AuthorizedUsers:
 - angus
 - malcolm

AuthorizedAddresses:
 - a.b.c.d
 - w.x.y.z

```

## 3. Add the systemd BruteDrop timer

/etc/systemd/system/brutedrop.timer

```
[Unit]
Description=Launch BruteDrop every 20s
Requires=brutedrop.service

[Timer]
OnCalendar=*-*-* *:*:00,20,40
Persistent=true

[Install]
WantedBy=timers.target
```

## 4. Add the systemd BruteDrop service

`/usr/lib/systemd/system/brutedrop.service`

```
[Unit]
Description=BruteDrop
After=sshd.service

[Service]
Type=oneshot
ExecStart=/sbin/brutedrop
StandardOutput=journal
```

