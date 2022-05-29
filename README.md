# Simple XQUARE Resource Management: xctl

## Table of Contents

- Installation

  - macOS
  - Linux
  - Windows

- Usage
  - Create application
  - Get application list
  - Delete application

# Installation

## macOS

1. Get latest version tar.gz archive using curl

```bash
# macbook pro
sudo curl -L /tmp https://github.com/team-xquare/xctl/releases/download/v0.2.0/xctl-v0.2.0-darwin-amd64.tar.gz > /tmp/xctl.tar.gz

# macbook air
sudo curl -L /tmp https://github.com/team-xquare/xctl/releases/download/v0.2.0/xctl-v0.2.0-darwin-arm64.tar.gz > /tmp/xctl.tar.gz
```

2. Extract tar.gz file

```bash
sudo tar -zxvf /tmp/xctl.tar.gz -C /tmp
sudo rm /tmp/xctl.tar.gz
```

3. Install and setting credential

```bash
sudo mv /tmp/xctl-v0.2.0-darwin-amd64 /usr/local/bin/xctl

or

sudo mv /tmp/xctl-v0.2.0-darwin-arm64 /usr/local/bin/xctl

mkdir $HOME/.xctl
sudo chmod 0777 -R $HOME/.xctl
sudo mv /tmp/template $HOME/.xctl/
# You can get github token from the XQUARE notion page.
xctl set credential -t <xquare-admin_github_token>
```

## Linux

1. Get latest version tar.gz archive using curl

```bash
sudo curl -L /tmp https://github.com/team-xquare/xctl/releases/download/v0.2.0/xctl-v0.2.0-linux-amd64.tar.gz > /tmp/xctl.tar.gz
```

2. Extract tar.gz file

```bash
sudo tar -zxvf /tmp/xctl.tar.gz -C /tmp
sudo rm /tmp/xctl.tar.gz
```

3. Install and setting credential

```bash
sudo mv /tmp/xctl-v0.2.0-linux-amd64 /usr/local/bin/xctl
mkdir $HOME/.xctl
sudo chmod 0777 -R $HOME/.xctl
sudo mv /tmp/template $HOME/.xctl/
# You can get github token from the XQUARE notion page.
xctl set credential -t <xquare-admin_github_token>
```

## Windows

1. Get latest version tar.gz archive from release asset <br>
   and unzip the archive

2. Setting environment and credential

```bash
mkdir %USERPROFILE%\.xctl
# move template folder to %USERPROFILE%\.xctl
move .\template  %USERPROFILE%\.xctl
# You can get github token from the XQUARE notion page.
xctl set credential -t <xquare-admin_github_token>
```

# Usage

## create app

- Default Options

```
	type:          "backend",
	host:          "api.xquare.app",
	image registry: "registry.hub.docker.com",
	image tag:      "latest",
	container port: 8080,
	environment:   "staging",
	prefix:        "/",
```

- Example

```bash
# Create frontend application named "eungyeol" to staging environment
# host name is "eungyeol.xquare.app" and port number 3000
xctl create app eungyeol -t frontend --host eungyeol.xquare.app --port 3000

# Create backend application named "notification" to production environment
# host name is "api.xquare.app", prefix is "/notification",
xctl create app notification -p /notification -e production -t backend

or

xctl create app notification -p /notification -e prod
```

## Get application list

- Example

```bash
# Get application list from staging environment
xctl get app -e stag or xctl get app

#result
Environemnt: staging
frontend applications
 | Name                  | Base Url             | Image Version |
 | app-eungyeol-frontend | eungyeol.xquare.app/ | latest        |
backend applications
 | Name | Base Url | Image Version |

# Get application list from production environment
xctl get app -e production
Environemnt: production
frontend applications
 | Name | Base Url | Image Version |
backend applications
 | Name                     | Base Url                   | Image Version |
 | app-notification-backend | api.xquare.app/noticiation | latest        |
```

## Delete application

```bash
# Delete frontend application named "eungyeol" on staging environment
xctl delete app notification -e staging -t frontend

# Delete backend application named "notification" on production environment
xctl delete app notification -e production

or

xctl delete app app-notification-backend -e product
```
