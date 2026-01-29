# API CONTRACT
## Telegram Job Platform

---

## Общие принципы

- REST
- JSON
- HTTP status codes
- Authorization через Telegram ID (MVP)

---

## Архитектурное решение

> **Вакансия публикуется в канал ТОЛЬКО после ручного approve админом.**

### Кто что делает

| Роль | Права |
|------|-------|
| Рекрутер | отправляет вакансию, НЕ может публиковать/аппрувить |
| Бот | собирает данные, показывает кнопки, НЕ принимает решений |
| Админ | видит pending, нажимает Approve/Reject |
| Backend | проверяет роль, меняет статус, публикует в канал |

### Допустимые переходы статусов

```
draft → pending       (рекрутер отправил)
pending → approved    (админ одобрил)
approved → published  (система опубликовала)
pending → rejected    (админ отклонил)
```

⚠️ `pending → published` напрямую **ЗАПРЕЩЁН**

---

## POST /api/jobs

### Request
```json
{
  "company": "Acme",
  "title": "Backend Go Developer",
  "level": "senior",
  "type": "remote",
  "category": "web3",
  "salary_from": 4000,
  "salary_to": 6000,
  "description": "Job description",
  "apply_link": "https://..."
}
```

### Response
```json
{
  "id": "uuid",
  "status": "pending"
}
```

---

## GET /api/jobs?status=pending

### Response
```json
[
  {
    "id": "uuid",
    "company": "Acme",
    "title": "Backend Go Developer",
    "status": "pending",
    "created_at": "2026-01-01T10:00:00Z"
  }
]
```

---

## POST /api/jobs/{id}/approve

> ⚠️ Только для админов. На MVP: **Approve = Publish immediately**

### Headers
```
X-Telegram-ID: 123456
```

### Logic
```go
func (s *JobService) ApproveJob(id uuid.UUID, adminID int64) error {
    if !s.auth.IsAdmin(adminID) {
        return errors.New("forbidden")
    }
    s.repo.UpdateStatus(id, "approved")
    return s.publisher.Publish(id)  // сразу публикуем
}
```

### Response
```json
{
  "status": "published",
  "published_at": "2026-01-01T12:00:00Z"
}
```

---

## POST /api/jobs/{id}/reject

### Request
```json
{
  "reason": "Invalid description"
}
```

---

## POST /api/jobs/{id}/publish

> ⚠️ На MVP этот endpoint НЕ используется отдельно.
> Публикация происходит автоматически при approve.

### Response
```json
{
  "status": "published",
  "published_at": "2026-01-01T12:00:00Z"
}
```

---

## Будущее расширение (v2)

```
[✅ Approve & Publish]     — сразу публикует
[🕒 Approve (publish later)] — только аппрув
[❌ Reject]                 — отклонить
```

---

## ERROR FORMAT

```json
{
  "error": "validation_error",
  "message": "salary_from must be <= salary_to"
}
```

---

