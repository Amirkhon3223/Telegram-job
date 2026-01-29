package bot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-job/internal/domain"
)

func escapeMarkdown(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"`", "\\`",
	)
	return replacer.Replace(s)
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.cmdStart(msg)
	case "post_job":
		b.cmdPostJob(msg)
	case "cancel":
		b.cmdCancel(msg)
	case "status":
		b.cmdStatus(msg)
	case "prices":
		b.cmdPrices(msg)
	case "contact":
		b.cmdContact(msg)
	default:
		b.sendMessage(msg.Chat.ID, "Unknown command. Use /post\\_job to submit a job.")
	}
}

func (b *Bot) cmdStart(msg *tgbotapi.Message) {
	text := `🎯 *BridgeJob Bot*

Я помогу разместить вакансию в канале @BridgeJob

I help you post jobs to @BridgeJob channel

*Команды / Commands:*
/post\_job — Разместить вакансию / Post a job
/prices — Прайс-лист / Price list
/contact — Связаться / Contact us
/cancel — Отменить / Cancel

📢 Канал / Channel: @BridgeJob`

	b.sendMessage(msg.Chat.ID, text)
}

func (b *Bot) cmdPostJob(msg *tgbotapi.Message) {
	b.fsm.Reset(msg.From.ID)
	b.fsm.SetState(msg.From.ID, StateWaitLanguage)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang:en"),
		),
	)
	b.sendMessageWithKeyboard(msg.Chat.ID, "🌐 Выберите язык / Choose language:", keyboard)
}

func (b *Bot) cmdCancel(msg *tgbotapi.Message) {
	lang := b.fsm.GetLanguage(msg.From.ID)
	m := GetMessages(lang)
	b.fsm.Reset(msg.From.ID)
	b.sendMessage(msg.Chat.ID, m.Cancelled)
}

func (b *Bot) cmdStatus(msg *tgbotapi.Message) {
	b.sendMessage(msg.Chat.ID, "Status check is not implemented yet.")
}

func (b *Bot) cmdPrices(msg *tgbotapi.Message) {
	text := `💰 *Прайс-лист / Price List*

📌 *Стандартное размещение / Standard* — *$25*
1 пост в канале / 1 post in channel

⭐ *Featured* — *$65*
Пост + закреп 48ч / Post + pin 48h

📦 *Пакет 5 вакансий / 5 Jobs Pack* — *$100*
5 стандартных постов / 5 standard posts

💳 *Оплата / Payment:*
USDT / Wise / PayPal

📞 *Контакт / Contact:*
@amirichinvoker | @manizha\_ash

📢 Канал / Channel: @BridgeJob`

	b.sendMessage(msg.Chat.ID, text)
}

func (b *Bot) cmdContact(msg *tgbotapi.Message) {
	text := `📞 *Контакт / Contact*

По всем вопросам обращайтесь:
For any questions contact:

👤 @amirichinvoker
👤 @manizha\_ash

📢 Канал / Channel: @BridgeJob`

	b.sendMessage(msg.Chat.ID, text)
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	userState := b.fsm.GetState(msg.From.ID)
	lang := b.fsm.GetLanguage(msg.From.ID)
	m := GetMessages(lang)

	switch userState.State {
	case StateNone:
		b.sendMessage(msg.Chat.ID, "Use /post\\_job to submit a job.\nИспользуйте /post\\_job чтобы разместить вакансию.")
		return

	case StateWaitLanguage:
		b.sendMessage(msg.Chat.ID, "Please select language using buttons above.\nВыберите язык кнопками выше.")
		return

	case StateWaitCompany:
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.Company = msg.Text
		})
		b.fsm.SetState(msg.From.ID, StateWaitContact)
		b.sendMessage(msg.Chat.ID, m.Step2Contact)

	case StateWaitContact:
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.Contact = msg.Text
		})
		b.fsm.SetState(msg.From.ID, StateWaitTitle)
		b.sendMessage(msg.Chat.ID, m.Step3Title)

	case StateWaitTitle:
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.Title = msg.Text
		})
		b.fsm.SetState(msg.From.ID, StateWaitLevel)
		b.sendLevelKeyboard(msg.Chat.ID, lang)

	case StateWaitLevel:
		level := domain.JobLevel(strings.ToLower(msg.Text))
		if !isValidLevel(level) {
			b.sendMessage(msg.Chat.ID, m.InvalidNumber)
			return
		}
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.Level = level
		})
		b.fsm.SetState(msg.From.ID, StateWaitType)
		b.sendTypeKeyboard(msg.Chat.ID, lang)

	case StateWaitType:
		jobType := domain.JobType(strings.ToLower(msg.Text))
		if !isValidType(jobType) {
			b.sendMessage(msg.Chat.ID, m.InvalidNumber)
			return
		}
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.Type = jobType
		})
		b.fsm.SetState(msg.From.ID, StateWaitCategory)
		b.sendCategoryKeyboard(msg.Chat.ID, lang)

	case StateWaitCategory:
		category := domain.JobCategory(strings.ToLower(msg.Text))
		if !isValidCategory(category) {
			b.sendMessage(msg.Chat.ID, m.InvalidNumber)
			return
		}
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.Category = category
		})
		b.fsm.SetState(msg.From.ID, StateWaitDescription)
		b.sendMessage(msg.Chat.ID, m.Step7Description)

	case StateWaitDescription:
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.Description = msg.Text
		})
		b.fsm.SetState(msg.From.ID, StateWaitSalaryFrom)
		b.sendMessage(msg.Chat.ID, m.Step8SalaryFrom)

	case StateWaitSalaryFrom:
		if !isSkip(msg.Text) {
			salary, err := strconv.Atoi(msg.Text)
			if err != nil {
				b.sendMessage(msg.Chat.ID, m.InvalidNumber)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
				d.SalaryFrom = &salary
			})
		}
		b.fsm.SetState(msg.From.ID, StateWaitSalaryTo)
		b.sendMessage(msg.Chat.ID, m.Step9SalaryTo)

	case StateWaitSalaryTo:
		if !isSkip(msg.Text) {
			salary, err := strconv.Atoi(msg.Text)
			if err != nil {
				b.sendMessage(msg.Chat.ID, m.InvalidNumber)
				return
			}
			// Валидация: ДО не может быть меньше ОТ
			draft := b.fsm.GetDraft(msg.From.ID)
			if draft != nil && draft.SalaryFrom != nil && salary < *draft.SalaryFrom {
				b.sendMessage(msg.Chat.ID, m.SalaryToLessThanFrom)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
				d.SalaryTo = &salary
			})
		}
		b.fsm.SetState(msg.From.ID, StateWaitApplyLink)
		b.sendMessage(msg.Chat.ID, m.Step10ApplyLink)

	case StateWaitApplyLink:
		b.fsm.UpdateDraft(msg.From.ID, func(d *JobDraft) {
			d.ApplyLink = msg.Text
		})
		b.fsm.SetState(msg.From.ID, StatePreview)
		b.sendPreview(msg)
	}
}

func (b *Bot) sendLevelKeyboard(chatID int64, lang Language) {
	m := GetMessages(lang)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(m.LevelJunior, "level:junior"),
			tgbotapi.NewInlineKeyboardButtonData(m.LevelMiddle, "level:middle"),
			tgbotapi.NewInlineKeyboardButtonData(m.LevelSenior, "level:senior"),
		),
	)
	b.sendMessageWithKeyboard(chatID, m.Step4Level, keyboard)
}

func (b *Bot) sendTypeKeyboard(chatID int64, lang Language) {
	m := GetMessages(lang)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(m.TypeRemote, "type:remote"),
			tgbotapi.NewInlineKeyboardButtonData(m.TypeHybrid, "type:hybrid"),
			tgbotapi.NewInlineKeyboardButtonData(m.TypeOnsite, "type:onsite"),
		),
	)
	b.sendMessageWithKeyboard(chatID, m.Step5Type, keyboard)
}

func (b *Bot) sendCategoryKeyboard(chatID int64, lang Language) {
	m := GetMessages(lang)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(m.CategoryDev, "category:dev"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(m.CategoryWeb2, "category:web2"),
			tgbotapi.NewInlineKeyboardButtonData(m.CategoryWeb3, "category:web3"),
		),
	)
	b.sendMessageWithKeyboard(chatID, m.Step6Category, keyboard)
}

func (b *Bot) sendPreview(msg *tgbotapi.Message) {
	log.Printf("sendPreview called for user %d", msg.From.ID)

	lang := b.fsm.GetLanguage(msg.From.ID)
	m := GetMessages(lang)

	draft := b.fsm.GetDraft(msg.From.ID)
	if draft == nil {
		log.Printf("Draft is nil for user %d", msg.From.ID)
		b.sendMessage(msg.Chat.ID, "Something went wrong. Please start over with /post\\_job")
		return
	}

	log.Printf("Draft found: %+v", draft)

	salary := m.SalaryNotSpecified
	if draft.SalaryFrom != nil && draft.SalaryTo != nil {
		salary = fmt.Sprintf("$%d – $%d", *draft.SalaryFrom, *draft.SalaryTo)
	} else if draft.SalaryFrom != nil {
		salary = fmt.Sprintf(m.SalaryFrom, *draft.SalaryFrom)
	} else if draft.SalaryTo != nil {
		salary = fmt.Sprintf(m.SalaryTo, *draft.SalaryTo)
	}

	text := fmt.Sprintf(`%s

🏢 *%s:* %s
📧 *%s:* %s
💼 *%s:* %s
📊 *%s:* %s
🌍 *%s:* %s
🏷️ *%s:* %s
💰 *%s:* %s
🔗 *%s:* %s

📝 *%s:*
%s

———
%s`,
		m.PreviewTitle,
		m.Company, escapeMarkdown(draft.Company),
		m.Contact, escapeMarkdown(draft.Contact),
		m.Title, escapeMarkdown(draft.Title),
		m.Level, draft.Level,
		m.Type, draft.Type,
		m.Category, draft.Category,
		m.Salary, salary,
		m.ApplyLink, escapeMarkdown(draft.ApplyLink),
		m.Description,
		escapeMarkdown(draft.Description),
		m.PreviewConfirm,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(m.BtnSubmit, "submit"),
			tgbotapi.NewInlineKeyboardButtonData(m.BtnCancel, "cancel_submit"),
		),
	)

	log.Printf("Sending preview to chat %d", msg.Chat.ID)
	b.sendMessageWithKeyboard(msg.Chat.ID, text, keyboard)
}

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	userID := callback.From.ID
	chatID := callback.Message.Chat.ID

	// Answer callback to remove loading state
	b.api.Request(tgbotapi.NewCallback(callback.ID, ""))

	// Handle language selection
	if strings.HasPrefix(data, "lang:") {
		langStr := strings.TrimPrefix(data, "lang:")
		lang := Language(langStr)
		b.fsm.SetLanguage(userID, lang)
		// Сохраняем язык в черновик вакансии
		b.fsm.UpdateDraft(userID, func(d *JobDraft) {
			d.Language = langStr
		})
		b.fsm.SetState(userID, StateWaitCompany)
		m := GetMessages(lang)
		b.sendMessage(chatID, m.Step1Company)
		return
	}

	lang := b.fsm.GetLanguage(userID)
	m := GetMessages(lang)

	// Handle level selection
	if strings.HasPrefix(data, "level:") {
		level := domain.JobLevel(strings.TrimPrefix(data, "level:"))
		b.fsm.UpdateDraft(userID, func(d *JobDraft) {
			d.Level = level
		})
		b.fsm.SetState(userID, StateWaitType)
		b.sendTypeKeyboard(chatID, lang)
		return
	}

	// Handle type selection
	if strings.HasPrefix(data, "type:") {
		jobType := domain.JobType(strings.TrimPrefix(data, "type:"))
		b.fsm.UpdateDraft(userID, func(d *JobDraft) {
			d.Type = jobType
		})
		b.fsm.SetState(userID, StateWaitCategory)
		b.sendCategoryKeyboard(chatID, lang)
		return
	}

	// Handle category selection
	if strings.HasPrefix(data, "category:") {
		category := domain.JobCategory(strings.TrimPrefix(data, "category:"))
		b.fsm.UpdateDraft(userID, func(d *JobDraft) {
			d.Category = category
		})
		b.fsm.SetState(userID, StateWaitDescription)
		b.sendMessage(chatID, m.Step7Description)
		return
	}

	// Handle submit
	if data == "submit" {
		b.submitJob(callback)
		return
	}

	// Handle cancel
	if data == "cancel_submit" {
		b.fsm.Reset(userID)
		b.sendMessage(chatID, m.Cancelled)
		return
	}

	// Handle admin callbacks (approve/reject/delete)
	if strings.HasPrefix(data, "approve:") || strings.HasPrefix(data, "reject:") ||
		strings.HasPrefix(data, "delete:") || strings.HasPrefix(data, "confirm_delete:") ||
		strings.HasPrefix(data, "cancel_delete:") {
		b.handleAdminCallback(callback)
		return
	}
}

func (b *Bot) submitJob(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	chatID := callback.Message.Chat.ID
	lang := b.fsm.GetLanguage(userID)
	m := GetMessages(lang)

	log.Printf("submitJob called for user %d", userID)

	draft := b.fsm.GetDraft(userID)
	if draft == nil {
		log.Printf("Draft is nil in submitJob for user %d", userID)
		b.sendMessage(chatID, "Something went wrong. Please start over with /post\\_job")
		return
	}

	log.Printf("Creating job for user %d", userID)

	ctx := context.Background()
	username := callback.From.UserName

	job, err := b.jobService.CreateJob(ctx, userID, username, draft.ToCreateRequest())
	if err != nil {
		log.Printf("Error creating job: %v", err)
		b.sendMessage(chatID, m.SubmitError+err.Error())
		return
	}

	log.Printf("Job created successfully: %s", job.ID.String())

	b.fsm.Reset(userID)

	b.sendMessage(chatID, fmt.Sprintf(m.SubmitSuccess, job.ID.String()))
}

func isValidLevel(level domain.JobLevel) bool {
	return level == domain.JobLevelJunior || level == domain.JobLevelMiddle || level == domain.JobLevelSenior
}

func isValidType(t domain.JobType) bool {
	return t == domain.JobTypeRemote || t == domain.JobTypeHybrid || t == domain.JobTypeOnsite
}

func isValidCategory(c domain.JobCategory) bool {
	return c == domain.JobCategoryWeb2 || c == domain.JobCategoryWeb3 || c == domain.JobCategoryDev
}

func isSkip(text string) bool {
	lower := strings.ToLower(text)
	return lower == "skip" || lower == "скип"
}
