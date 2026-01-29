# DEPLOYMENT
## Telegram Job Platform

---

## Environment Variables

```
BOT_TOKEN=
CHANNEL_ID=
DATABASE_URL=
API_PORT=8080
BOT_WEBHOOK_URL=
ADMIN_TELEGRAM_IDS=123456,987654
```

> ⚠️ `ADMIN_TELEGRAM_IDS` — whitelist админов для approve вакансий

---

## docker-compose.yml (MVP)

```yaml
version: '3.9'
services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: jobs
      POSTGRES_USER: jobs
      POSTGRES_PASSWORD: jobs
    ports:
      - "5432:5432"

  api:
    build: .
    command: cmd/api/main.go
    depends_on:
      - postgres

  bot:
    build: .
    command: cmd/bot/main.go
    depends_on:
      - api
```

---

## Local Run

```bash
docker-compose up --build
```

---

## Production Plan

- VPS (Hetzner / DO)
- Docker
- Nginx
- HTTPS (Let's Encrypt)
- Telegram Webhook

---

