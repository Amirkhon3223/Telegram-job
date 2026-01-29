# DATABASE SCHEMA
## Telegram Job Platform

---

## Общие принципы

- PostgreSQL 15+
- UUID в качестве PK
- snake_case
- soft-delete не используется на MVP
- все временные поля в UTC

---

## Архитектурное решение

> **Вакансия публикуется в канал ТОЛЬКО после ручного approve админом.**
> Approve на MVP = немедленная публикация.

⚠️ Переход `pending → published` напрямую **ЗАПРЕЩЁН**

---

## ENUM TYPES

```sql
CREATE TYPE job_level AS ENUM ('junior', 'middle', 'senior');
CREATE TYPE job_type AS ENUM ('remote', 'hybrid', 'onsite');
CREATE TYPE job_category AS ENUM ('web2', 'web3');
-- Статусы вакансии:
-- draft     — ввод данных (FSM)
-- pending   — отправлена, ждёт админа
-- approved  — одобрена админом
-- published — опубликована в канале
-- rejected  — отклонена
CREATE TYPE job_status AS ENUM ('draft', 'pending', 'approved', 'published', 'rejected');
CREATE TYPE user_role AS ENUM ('admin', 'recruiter');
CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed');
CREATE TYPE payment_type AS ENUM ('single', 'package', 'subscription');
```

---

## TABLE: users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id BIGINT UNIQUE NOT NULL,
    username TEXT,
    role user_role NOT NULL DEFAULT 'recruiter',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## TABLE: companies

```sql
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    contact TEXT NOT NULL,
    telegram TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## TABLE: jobs

```sql
CREATE TABLE jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    level job_level NOT NULL,
    type job_type NOT NULL,
    category job_category NOT NULL,
    salary_from INT,
    salary_to INT,
    description TEXT NOT NULL,
    apply_link TEXT NOT NULL,
    status job_status NOT NULL DEFAULT 'draft',
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_created_at ON jobs(created_at);
```

---

## TABLE: payments (future)

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id),
    amount INT NOT NULL,
    currency TEXT NOT NULL,
    type payment_type NOT NULL,
    status payment_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

---

## RELATIONSHIPS

- users 1—1 companies (logical)
- companies 1—N jobs
- companies 1—N payments

---

## ADMIN AUTH (MVP)

Админы определяются через whitelist в `.env`:

```
ADMIN_TELEGRAM_IDS=123456,987654
```

Проверка роли:

```go
func IsAdmin(id int64) bool {
    return adminSet[id]
}
```

Позже — через таблицу `users.role`.

---

