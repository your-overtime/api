# Your Overtime - API

[![Go Report Card](https://goreportcard.com/badge/github.com/your-overtime/api)](https://goreportcard.com/report/github.com/your-overtime/api)[![Go](https://github.com/your-overtime/api/actions/workflows/go.yml/badge.svg)](https://github.com/your-overtime/api/actions/workflows/go.yml)[![Publish Docker image](https://github.com/your-overtime/api/actions/workflows/api.yml/badge.svg)](https://github.com/your-overtime/api/actions/workflows/api.yml)

Swagger documentation can be found on https://your-overtime.de/api/v1/swagger/index.html.

```yml
version: "3.3"
services:
  api:
    image: ghcr.io/your-overtime/api:1
    environment:
      - HOST=0.0.0.0:8080
      - DB_USER=overtime
      - DB_PASSWORD=secret
      - DB_HOST=db:3306
      - DB_NAME=overtime
      - DEBUG=false
      - ADMIN_TOKEN=secret
      - TZ=Europe/Berlin
    restart: unless-stopped
    port: 8080:8080
  db:
    image: mariadb
    restart: unless-stopped
    volumes:
      - ./db:/var/lib/mysql
    environment:
      - MARIADB_INITDB_SKIP_TZINFO=true
      - TZ=Europe/Berlin
      - MYSQL_RANDOM_ROOT_PASSWORD=true
      - MYSQL_DATABASE=overtime
      - MYSQL_PASSWORD=secret
```
