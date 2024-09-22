# BruteDrop

A simple but effective - I mean <i>brute</i> -  tool written in Go as response to brute force attacks.

The idea, very common instead, is to block all ports to each IP address from which someone try to gain access to an SSH session by brute force attack.

## 1. Basic pre-required sshd_config configuration against SSH attack

Configure your SSH daemon with those advices in mind:

- For sure don't use common user names as admin, mysql, etc.
- No password authentication
- No root login
- Use key access only
- Use the best key type of the moment. Currently `ssh-ed25519`
- Limit SSH access to a list of users

Which gives in `sshd_config` file directives:

```
PasswordAuthentication no

PermitRootLogin no

PubkeyAcceptedKeyTypes ssh-ed25519

AllowUsers angus@* malcom@e.f.g.h
```

## 2. Install BruteDrop binary

```
sudo make install
```

## 3. Add BruteDrop configuration file

```
IptablesBinPath: /usr/bin/iptables
JournalctlBinPath: /usr/bin/journalctl

#DryRunMode: true
DryRunMode: false

#Set Logging to file path or "stdout"
LoggingTo: stdout
#LoggingTo: /var/log/brutedrop.log
LogEntriesSince: 2

AuthorizedUsers:
 - angus
 - malcolm

AuthorizedAddresses:
 - a.b.c.d
 - w.x.y.z
```

## 4. Add the systemd BruteDrop timer

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

## 5. Add the systemd BruteDrop service

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

## 6. Enable and start the systemd BruteDrop service

To be sure you won't lock you out, you can test your configuration and see what's going on when BruteDrop runs by setting `DryRunMode` to `true` and follow log outputs with `sudo journalctl -u brutedrop -f`.

```
[angus@box ~]$ sudo systemdctl enable brutedrop
[angus@box ~]$ sudo systemdctl start brutedrop
```

