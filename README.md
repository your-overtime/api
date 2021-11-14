# Your Overtime - API

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
```
