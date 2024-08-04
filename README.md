# Introduction

This is a simple app that send a discord message via webhook when an event is occuring in google calendar

# Requirements

- Google Cloud Platform (GCP) Service Account with rights to the Calendar API
  - Token json file of the Service Account
- Google Calendar
  - Give read access to the Service Account
- Discord webhook url

# Configuration

## Docker

`docker-compose.yaml`

```yaml
version: "3.8"
services:
  redis:
    image: redis:latest
    container_name: redis_discordevent
    #volumes:
    #  - ./data:/data
  discordevents:
    image: ghcr.io/markaplay-game-hosting/goeventbot/app:latest
    restart: unless-stopped
    env_file: .env
    volumes:
      - ./config/service_account.json:/config/service_account.json
    depends_on:
      - redis
networks: {}
```

`.env`

```text
CALENDAR_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx@group.calendar.google.com
WEBHOOK_URL=https://discord.com/api/webhooks/11111111111111111111111111/tokenidblabla
REDIS_ADDR=redis_discordevent:6379
TIMESPAN=1h
POLLING=1m
```
