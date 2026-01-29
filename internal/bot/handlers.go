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
	case "help":
		b.cmdHelp(msg)
	case "post_job":
		b.cmdPostJob(msg)
	case "myjobs":
		b.cmdMyJobs(msg)
	case "pricing", "prices":
		b.cmdPrices(msg)
	case "faq":
		b.cmdFAQ(msg)
	case "about":
		b.cmdAbout(msg)
	case "contact":
		b.cmdContact(msg)
	case "cancel":
		b.cmdCancel(msg)
	case "language":
		b.cmdLanguage(msg)
	// Admin commands
	case "pending":
		b.cmdPending(msg)
	case "stats":
		b.cmdStats(msg)
	case "admins":
		b.cmdAdmins(msg)
	default:
		m := b.getInterfaceMessages(msg.From.ID)
		b.sendMessage(msg.Chat.ID, m.UnknownCommand)
	}
}

func (b *Bot) cmdStart(msg *tgbotapi.Message) {
	// Check if user has interface language set
	lang := b.getUserInterfaceLanguage(msg.From.ID)

	if lang == "" {
		// First time - show bilingual welcome and ask to choose language
		text := `👋 *Добро пожаловать! / Welcome!*

Это сервис для публикации вакансий с ручной модерацией.
This is a job posting service with manual moderation.

Пожалуйста, выберите язык / Please choose your language:`

		keyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "interface_lang:ru"),
				tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "interface_lang:en"),
			),
		)
		b.sendMessageWithKeyboard(msg.Chat.ID, text, keyboard)
		return
	}

	// User has language set - show normal welcome
	m := GetMessages(lang)
	b.sendMessage(msg.Chat.ID, m.Welcome)
}

func (b *Bot) cmdPostJob(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	b.fsm.Reset(msg.From.ID)
	b.fsm.SetState(msg.From.ID, StateWaitLanguage)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "lang:ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "lang:en"),
		),
	)
	b.sendMessageWithKeyboard(msg.Chat.ID, m.ChooseJobLanguage, keyboard)
}

func (b *Bot) cmdCancel(msg *tgbotapi.Message) {
	// Use FSM language if in job creation flow, otherwise interface language
	lang := b.fsm.GetLanguage(msg.From.ID)
	if lang == "" {
		lang = b.getUserInterfaceLanguage(msg.From.ID)
		if lang == "" {
			lang = LangRU
		}
	}
	m := GetMessages(lang)
	b.fsm.Reset(msg.From.ID)
	b.sendMessage(msg.Chat.ID, m.Cancelled)
}

func (b *Bot) cmdLanguage(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🇷🇺 Русский", "interface_lang:ru"),
			tgbotapi.NewInlineKeyboardButtonData("🇬🇧 English", "interface_lang:en"),
		),
	)
	b.sendMessageWithKeyboard(msg.Chat.ID, m.ChooseLanguage, keyboard)
}

func (b *Bot) cmdPrices(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	b.sendMessage(msg.Chat.ID, m.Pricing)
}

func (b *Bot) cmdContact(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	b.sendMessage(msg.Chat.ID, m.Contact)
}

func (b *Bot) cmdHelp(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	text := m.Help

	// Добавляем админские команды для админов
	if b.cfg.IsAdmin(msg.From.ID) {
		text += m.HelpAdmin
	}

	b.sendMessage(msg.Chat.ID, text)
}

func (b *Bot) cmdFAQ(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	b.sendMessage(msg.Chat.ID, m.FAQ)
}

func (b *Bot) cmdAbout(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	b.sendMessage(msg.Chat.ID, m.About)
}

func (b *Bot) cmdMyJobs(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	lang := b.getUserInterfaceLanguage(msg.From.ID)

	ctx := context.Background()
	jobs, err := b.jobService.GetUserJobs(ctx, msg.From.ID)
	if err != nil {
		log.Printf("Error getting user jobs: %v", err)
		b.sendMessage(msg.Chat.ID, "Error / Ошибка")
		return
	}

	if len(jobs) == 0 {
		b.sendMessage(msg.Chat.ID, m.NoJobs)
		return
	}

	text := m.YourJobs + "\n"
	for i, job := range jobs {
		statusEmoji := getStatusEmoji(job.Status)
		statusText := getStatusText(job.Status, lang)
		text += fmt.Sprintf("\n%d. *%s*\n   %s %s\n", i+1, escapeMarkdown(job.Title), statusEmoji, statusText)
	}

	b.sendMessage(msg.Chat.ID, text)
}

// Admin commands

func (b *Bot) cmdPending(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)

	if !b.cfg.IsAdmin(msg.From.ID) {
		b.sendMessage(msg.Chat.ID, m.NoPermission)
		return
	}

	ctx := context.Background()
	jobs, err := b.jobService.GetPendingJobs(ctx)
	if err != nil {
		log.Printf("Error getting pending jobs: %v", err)
		b.sendMessage(msg.Chat.ID, "Error / Ошибка")
		return
	}

	if len(jobs) == 0 {
		b.sendMessage(msg.Chat.ID, m.NoPendingJobs)
		return
	}

	b.sendMessage(msg.Chat.ID, fmt.Sprintf(m.PendingJobsCount, len(jobs)))

	// Отправляем каждую вакансию только этому админу
	for _, job := range jobs {
		b.sendPendingJobToAdmin(msg.Chat.ID, &job)
	}
}

func (b *Bot) sendPendingJobToAdmin(chatID int64, job *domain.JobWithCompany) {
	text := formatAdminNotification(job)

	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Кнопка связи с автором
	contact := job.CompanyContact
	if strings.HasPrefix(contact, "@") {
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📞 Связаться с автором", "https://t.me/"+strings.TrimPrefix(contact, "@")),
		))
	}

	// Кнопки Approve/Reject
	keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Approve", "approve:"+job.ID.String()),
		tgbotapi.NewInlineKeyboardButtonData("❌ Reject", "reject:"+job.ID.String()),
	))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(keyboardRows...)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

func (b *Bot) cmdStats(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	lang := b.getUserInterfaceLanguage(msg.From.ID)

	if !b.cfg.IsAdmin(msg.From.ID) {
		b.sendMessage(msg.Chat.ID, m.NoPermission)
		return
	}

	ctx := context.Background()
	stats, err := b.jobService.GetStats(ctx)
	if err != nil {
		log.Printf("Error getting stats: %v", err)
		b.sendMessage(msg.Chat.ID, "Error / Ошибка")
		return
	}

	var text string
	if lang == LangEN {
		text = fmt.Sprintf(`📊 *Service Statistics*

• Total jobs: %d
• Pending: %d
• Published: %d
• Rejected: %d
• Archived: %d`,
			stats.Total,
			stats.Pending,
			stats.Published,
			stats.Rejected,
			stats.Archived,
		)
	} else {
		text = fmt.Sprintf(`📊 *Статистика сервиса*

• Всего вакансий: %d
• На модерации: %d
• Опубликовано: %d
• Отклонено: %d
• В архиве: %d`,
			stats.Total,
			stats.Pending,
			stats.Published,
			stats.Rejected,
			stats.Archived,
		)
	}

	b.sendMessage(msg.Chat.ID, text)
}

func (b *Bot) cmdAdmins(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	lang := b.getUserInterfaceLanguage(msg.From.ID)

	if !b.cfg.IsAdmin(msg.From.ID) {
		b.sendMessage(msg.Chat.ID, m.NoPermission)
		return
	}

	var text string
	if lang == LangEN {
		text = "👮 *Service Administrators:*\n"
	} else {
		text = "👮 *Администраторы сервиса:*\n"
	}
	for adminID := range b.cfg.AdminTelegramIDs {
		text += fmt.Sprintf("\n• ID: `%d`", adminID)
	}

	b.sendMessage(msg.Chat.ID, text)
}

func getStatusEmoji(status domain.JobStatus) string {
	switch status {
	case domain.JobStatusPending:
		return "🕒"
	case domain.JobStatusApproved:
		return "✅"
	case domain.JobStatusPublished:
		return "📢"
	case domain.JobStatusRejected:
		return "❌"
	case domain.JobStatusArchived:
		return "🗑"
	default:
		return "📝"
	}
}

func getStatusText(status domain.JobStatus, lang Language) string {
	if lang == LangEN {
		switch status {
		case domain.JobStatusPending:
			return "Pending"
		case domain.JobStatusApproved:
			return "Approved"
		case domain.JobStatusPublished:
			return "Published"
		case domain.JobStatusRejected:
			return "Rejected"
		case domain.JobStatusArchived:
			return "Archived"
		default:
			return "Draft"
		}
	}
	// Russian (default)
	switch status {
	case domain.JobStatusPending:
		return "На модерации"
	case domain.JobStatusApproved:
		return "Одобрена"
	case domain.JobStatusPublished:
		return "Опубликована"
	case domain.JobStatusRejected:
		return "Отклонена"
	case domain.JobStatusArchived:
		return "В архиве"
	default:
		return "Черновик"
	}
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
		var level domain.JobLevel
		if isSkip(msg.Text) {
			level = domain.JobLevelSkip
		} else {
			level = domain.JobLevel(strings.ToLower(msg.Text))
			if !isValidLevel(level) {
				b.sendMessage(msg.Chat.ID, "Выберите уровень кнопками или введите 'skip' / 'скип'")
				return
			}
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
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(m.LevelInternship, "level:internship"),
			tgbotapi.NewInlineKeyboardButtonData(m.LevelSkip, "level:skip"),
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
		salary = fmt.Sprintf(m.SalaryFromLabel, *draft.SalaryFrom)
	} else if draft.SalaryTo != nil {
		salary = fmt.Sprintf(m.SalaryToLabel, *draft.SalaryTo)
	}

	levelDisplay := string(draft.Level)
	if draft.Level == "" {
		levelDisplay = m.LevelNotSpecified
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
		m.CompanyLabel, escapeMarkdown(draft.Company),
		m.ContactLabel, escapeMarkdown(draft.Contact),
		m.TitleLabel, escapeMarkdown(draft.Title),
		m.LevelLabel, levelDisplay,
		m.TypeLabel, draft.Type,
		m.CategoryLabel, draft.Category,
		m.SalaryLabel, salary,
		m.ApplyLinkLabel, escapeMarkdown(draft.ApplyLink),
		m.DescriptionLabel,
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

	// Handle interface language selection (for user preferences)
	if strings.HasPrefix(data, "interface_lang:") {
		langStr := strings.TrimPrefix(data, "interface_lang:")
		lang := Language(langStr)
		username := callback.From.UserName
		if err := b.setUserInterfaceLanguage(userID, username, lang); err != nil {
			log.Printf("Error setting interface language: %v", err)
		}
		m := GetMessages(lang)
		b.sendMessage(chatID, m.LanguageSet+"\n\n"+m.Welcome)
		return
	}

	// Handle job language selection (for job creation flow)
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
		levelStr := strings.TrimPrefix(data, "level:")
		var level domain.JobLevel
		if levelStr == "skip" {
			level = domain.JobLevelSkip // empty string
		} else {
			level = domain.JobLevel(levelStr)
		}
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
	return level == domain.JobLevelJunior || level == domain.JobLevelMiddle || level == domain.JobLevelSenior || level == domain.JobLevelInternship
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
