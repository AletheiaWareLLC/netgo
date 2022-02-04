netserver
=========

This is a Go implementation of a static website server.

# About

`netserver` makes it easy to deploy and maintain a static website.

# Install

```
# Install the binary (or download from https://github.com/AletheiaWareLLC/netgo/releases/latest)
go install aletheiaware.com/netgo/cmd/netserver

# Create user
adduser netserver

# Create config
cat <<EOT >> /home/netserver/config
CERTIFICATE_DIRECTORY=/etc/letsencrypt/live/example.com/
HTTPS=true
EOT

# Set permissions
chmod 600 /home/netserver/config

# Allow netserver to read config
sudo chown -R netserver:netserver /home/netserver/config

# Create netserver service
sudo cat <<EOT >> /etc/systemd/system/netserver.service
[Unit]
Description=netserver static website server
[Service]
User=netserver
WorkingDirectory=/home/netserver
EnvironmentFile=/home/netserver/config
ExecStart=$(whereis netserver) start
SuccessExitStatus=143
TimeoutStopSec=10
Restart=on-failure
RestartSec=5
[Install]
WantedBy=multi-user.target
EOT

# Reload daemon
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable netserver

# Start service
sudo systemctl start netserver
```

# Content

By default `netserver` will serve content from a subdirectory called `html\static`, this can be overriden with the environment variable `CONTENT_DIRECTORY`.

## Git Bare

If a website is stored in a git repository, a bare version on the server can be used to make deploying an update to a website as simple as `git push live`.

```
# Create a directory to house the repository
mkdir -p /var/www/example.com

# Enter directory
cd /var/www/example.com

# Initialize a bare repository
git init --bare

# Create a hook to checkout the latest version of the website
cat << EOT >> hooks/post-receive
#!/bin/sh
GIT_WORK_TREE=/var/www/example.com git checkout -f
EOT

# Make the hook executable
chmod +x hooks/post-receive
```

On the development machine setup the server as a remote;

```
git remote add live username@example.com:/var/www/example.com
git push live +master:refs/heads/master
```

Whenever you make changes to the website, commit them and run `git push live` to deploy them to the server.

# Logging

By default `netserver` will log to a subdirectory called `logs`, this can be overridden with the environment variable `LOG_DIRECTORY`.

# HTTPS

HTTPS can be enabled by setting the environment variable `HTTPS=true`.

## Certificate

SSL certificates will be loaded from a subdirectory called `certificates`, this can be overriden with the environment variable `CERTIFICATE_DIRECTORY`.

A Certificate Authority such as Let's Encrypt can be used to generate a certificate.

```
# Install certbot
sudo apt install certbot

# Generate certificate
sudo certbot certonly --standalone -d example.com

# Allow netserver to read security credentials
sudo chown -R netserver:netserver /etc/letsencrypt/

# Add cron job to renew certificate on the first day of the week
(sudo crontab -l ; echo '* * * * 0 sudo certbot renew --pre-hook "systemctl stop netserver" --post-hook "chown -R netserver:netserver /etc/letsencrypt/ && systemctl start netserver"') | sudo crontab -
```

## Firewall

A Firewall such as UFW can be used to control the open ports.

```
# Install ufw
sudo apt install ufw

# Allow http port
sudo ufw allow 80

# Allow https port
sudo ufw allow 443

# Enable firewall
sudo ufw enable

# Allow netserver to bind to port 443 (HTTPS)
# This is required each time the server binary is updated
sudo setcap CAP_NET_BIND_SERVICE=+eip $(whereis netserver)
```

## HTTP to HTTPS Redirect

`netserver` can redirect clients accessing webpages via HTTP to use HTTPS instead when the following environment variables are set;

- `HOST` - only requests matching the given host will be redirected; eg `example.com`
- `ROUTES` - only requests matching one of the given comma-separated routes will be redirected; eg `/,/index.html,/logo.png`
