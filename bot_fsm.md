# BOT FSM
## Telegram Job Submission Flow

---

## States

- START
- WAIT_COMPANY
- WAIT_TITLE
- WAIT_LEVEL
- WAIT_TYPE
- WAIT_CATEGORY
- WAIT_DESCRIPTION
- WAIT_SALARY
- WAIT_APPLY_LINK
- PREVIEW
- SUBMITTED

---

## Transitions

START → WAIT_COMPANY
WAIT_COMPANY → WAIT_TITLE
WAIT_TITLE → WAIT_LEVEL
WAIT_LEVEL → WAIT_TYPE
WAIT_TYPE → WAIT_CATEGORY
WAIT_CATEGORY → WAIT_DESCRIPTION
WAIT_DESCRIPTION → WAIT_SALARY
WAIT_SALARY → WAIT_APPLY_LINK
WAIT_APPLY_LINK → PREVIEW
PREVIEW → SUBMITTED

---

## PREVIEW FORMAT

```
Company: {{company}}
Role: {{title}}
Level: {{level}}
Type: {{type}}
Category: {{category}}
Salary: {{salary_from}}–{{salary_to}}
Apply: {{apply_link}}
```

---

## Commands

- /post_job — start FSM
- /cancel — reset state
- /status — show last submitted job status

---

## Edge Cases

- /cancel resets FSM
- invalid enum → repeat question
- empty message → repeat question

---

## ADMIN FLOW

### Уведомление админу (после submit рекрутера)

```
📋 Новая вакансия:

Company: Acme
Role: Backend Go Developer
Level: senior
Type: remote
Category: web3
Salary: $4000–$6000
Apply: https://...

[✅ Approve] [❌ Reject]
```

### Callback data

```
approve:{job_id}
reject:{job_id}
```

### После нажатия Approve

```go
// Бот вызывает API
POST /api/jobs/{id}/approve
X-Telegram-ID: {admin_telegram_id}

// Результат: вакансия публикуется в канал
```

### После нажатия Reject

Бот запрашивает причину → отправляет рекрутеру.

---

## РОЛИ

| Роль | Что делает | Что НЕ делает |
|------|------------|---------------|
| Рекрутер | отправляет вакансию | публикует, аппрувит |
| Бот | собирает данные, кнопки | принимает решения |
| Админ | approve/reject | — |

---

