# ⚡ Wake up babe, the server is on! ⚡

So your SMP is basically a **WhatsApp graveyard** now? ☠️

Server&#39;s got a better sleep schedule than you, cuz **nobody actually plays together**? 😴

AND you&#39;re still rocking that **pay-as-you-go Exaroton struggle**?

<h3> Bro, just start the server from the group chat. </h3>

Literally, just @ the bot and type:

`/start [server-id]`

✨ BOOM. Server wakes up. Your friends see it. Gaming resumes.
No more "Yo is the server up?" messages at 2 AM.

Current features: 
- Start
- Stop
- List servers
- List players on a server
- Getting a server info

## 🚀 Installation guide

### 🧱 Prerequisites
make sure you have:

- Docker installed
- A valid `config.yml`. (e.g: [config.yml.example](config.yml.example))

### 🐳 Run using Docker

```sh
docker run -it \
  -p 8080:8080 \
  --name exaroton-wa-bot \
  -v {YOUR config.yml PATH HERE}:/app/config.yml \
  dnabil/exaroton-wa-bot:latest
```

### Getting Started
visit localhost:8080 (or port of your choice)
- Login with default username (admin) and password (admin)
- Login whatsapp via QRCode
- Click the burger menu on the left top corner of your screen, go to exaroton settings page
- Fill your exaroton api token (can get it [here](https://exaroton.com/account/settings/))
- Go to whatsapp settings, and whitelist the group of your choice
- You're done :D

To start using it, @ the bot (the logged in whatsapp account in this app) then follow it with /help

e.g
```text
@UserExample /help
```

use /help [command] to explore its usage :) 
