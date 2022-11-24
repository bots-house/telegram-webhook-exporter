# telegram-webhook-exporter

[![CI](https://github.com/bots-house/telegram-webhook-exporter/actions/workflows/ci.yml/badge.svg)](https://github.com/bots-house/telegram-webhook-exporter/actions/workflows/ci.yml)

## Running 

```yaml 
version: '3.9'

services:
  exporter:
    image: ghcr.io/bots-house/telegram-webhook-exporter:v1.0.0
    environment:
      TOKENS: ${BOT_API_TOKEN}
    ports:
      - 8000:8000
```

## Metrics

```bash
curl localhost:8000/metrics
```

```
# HELP telegram_bot_api_last_delivery_error time and message for the most recent error that happened when trying to deliver an update via webhook
# TYPE telegram_bot_api_last_delivery_error gauge
telegram_bot_api_last_delivery_error{bot="username",msg=""} 0
# HELP telegram_bot_api_max_conns Maximum allowed number of simultaneous HTTPS connections to the webhook for update delivery
# TYPE telegram_bot_api_max_conns gauge
telegram_bot_api_max_conns{bot="username"} 0
# HELP telegram_bot_api_pending_updates number of updates awaiting delivery
# TYPE telegram_bot_api_pending_updates gauge
telegram_bot_api_pending_updates{bot="username"} 1
# HELP telegram_bot_api_response_time gauge of bot api response time
# TYPE telegram_bot_api_response_time gauge
telegram_bot_api_response_time{bot="username"} 65
# HELP telegram_bot_api_up true, if api is up
# TYPE telegram_bot_api_up gauge
telegram_bot_api_up{bot="username"} 1
```
