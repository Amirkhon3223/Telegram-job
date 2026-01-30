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
	case "post", "post_job":
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
	// Ensure user exists in DB with default 'en' language
	ctx := context.Background()
	_, _ = b.userRepo.GetOrCreate(ctx, msg.From.ID, msg.From.UserName)

	// Get user's language (defaults to 'en')
	m := b.getInterfaceMessages(msg.From.ID)
	b.sendMessage(msg.Chat.ID, m.Welcome)
}

func (b *Bot) cmdPostJob(msg *tgbotapi.Message) {
	m := b.getInterfaceMessages(msg.From.ID)
	b.fsm.Reset(msg.From.ID)
	b.fsm.SetState(msg.From.ID, StateWaitPostType)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Vacancy, "post_type:vacancy"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Resume, "post_type:resume"),
		),
	)
	b.sendMessageWithKeyboard(msg.Chat.ID, m.ChoosePostType, keyboard)
}

func (b *Bot) cmdCancel(msg *tgbotapi.Message) {
	lang := b.fsm.GetLanguage(msg.From.ID)
	if lang == "" {
		lang = b.getUserInterfaceLanguage(msg.From.ID)
		if lang == "" {
			lang = LangEN
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
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Russian, "interface_lang:ru"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.English, "interface_lang:en"),
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
	posts, err := b.jobService.GetUserJobs(ctx, msg.From.ID)
	if err != nil {
		log.Printf("Error getting user posts: %v", err)
		b.sendMessage(msg.Chat.ID, "Error / Ошибка")
		return
	}

	if len(posts) == 0 {
		b.sendMessage(msg.Chat.ID, m.NoPosts)
		return
	}

	text := m.YourPosts + "\n"
	for i, post := range posts {
		statusEmoji := getStatusEmoji(post.Status)
		statusText := getStatusText(post.Status, lang)
		postTypeEmoji := "🏢"
		if post.PostType == domain.PostTypeResume {
			postTypeEmoji = "👤"
		}
		text += fmt.Sprintf("\n%d. %s *%s*\n   %s %s\n", i+1, postTypeEmoji, escapeMarkdown(post.Title), statusEmoji, statusText)
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
	posts, err := b.jobService.GetPendingJobs(ctx)
	if err != nil {
		log.Printf("Error getting pending posts: %v", err)
		b.sendMessage(msg.Chat.ID, "Error / Ошибка")
		return
	}

	if len(posts) == 0 {
		b.sendMessage(msg.Chat.ID, m.NoPendingPosts)
		return
	}

	b.sendMessage(msg.Chat.ID, fmt.Sprintf(m.PendingPostsCount, len(posts)))

	for _, post := range posts {
		b.sendPendingPostToAdmin(msg.Chat.ID, &post)
	}
}

func (b *Bot) sendPendingPostToAdmin(chatID int64, post *domain.PostWithDetails) {
	text := formatAdminNotification(post)

	var keyboardRows [][]tgbotapi.InlineKeyboardButton

	// Contact button
	contact := post.CompanyContact
	if contact == "" {
		contact = post.Contact
	}
	if strings.HasPrefix(contact, "@") {
		keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("📞 Связаться с автором", "https://t.me/"+strings.TrimPrefix(contact, "@")),
		))
	}

	// Approve/Reject buttons
	keyboardRows = append(keyboardRows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅ Approve", "approve:"+post.ID.String()),
		tgbotapi.NewInlineKeyboardButtonData("❌ Reject", "reject:"+post.ID.String()),
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

• Total posts: %d
• Pending: %d
• Published: %d
• Rejected: %d
• Archived: %d`,
			stats.Total, stats.Pending, stats.Published, stats.Rejected, stats.Archived)
	} else {
		text = fmt.Sprintf(`📊 *Статистика сервиса*

• Всего публикаций: %d
• На модерации: %d
• Опубликовано: %d
• Отклонено: %d
• В архиве: %d`,
			stats.Total, stats.Pending, stats.Published, stats.Rejected, stats.Archived)
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

// ==================== MESSAGE HANDLERS ====================

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	userState := b.fsm.GetState(msg.From.ID)
	lang := b.fsm.GetLanguage(msg.From.ID)
	m := GetMessages(lang)
	postType := b.fsm.GetPostType(msg.From.ID)

	switch userState.State {
	case StateNone:
		b.sendMessage(msg.Chat.ID, "Use /post\\_job to submit.\nИспользуйте /post\\_job чтобы добавить публикацию.")
		return

	case StateWaitPostType:
		b.sendMessage(msg.Chat.ID, "Please select using buttons above.")
		return

	// ==================== VACANCY STATES ====================
	case StateWaitCompany:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Company = msg.Text })
		b.fsm.SetState(msg.From.ID, StateWaitContact)
		b.sendMessage(msg.Chat.ID, m.VacStep2Contact)

	case StateWaitContact:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Contact = msg.Text })
		b.fsm.SetState(msg.From.ID, StateWaitTitle)
		b.sendMessage(msg.Chat.ID, m.VacStep3Title)

	case StateWaitTitle:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Title = msg.Text })
		b.fsm.SetState(msg.From.ID, StateWaitLevel)
		b.sendLevelKeyboard(msg.Chat.ID, lang, m.VacStep4Level)

	case StateWaitLevel:
		if isSkip(msg.Text) {
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Level = domain.JobLevelSkip })
		} else {
			level := domain.JobLevel(strings.ToLower(msg.Text))
			if !isValidLevel(level) {
				b.sendMessage(msg.Chat.ID, "Select using buttons / Выберите кнопками")
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Level = level })
		}
		b.fsm.SetState(msg.From.ID, StateWaitType)
		b.sendTypeKeyboard(msg.Chat.ID, lang, m.VacStep5Type)

	case StateWaitType:
		jobType := domain.JobType(strings.ToLower(msg.Text))
		if !isValidType(jobType) {
			b.sendMessage(msg.Chat.ID, "Select using buttons / Выберите кнопками")
			return
		}
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Type = jobType })
		b.fsm.SetState(msg.From.ID, StateWaitCategory)
		b.sendCategoryKeyboard(msg.Chat.ID, lang)

	case StateWaitCategory:
		category := domain.JobCategory(strings.ToLower(msg.Text))
		if !isValidCategory(category) {
			b.sendMessage(msg.Chat.ID, "Select using buttons / Выберите кнопками")
			return
		}
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Category = category })
		b.fsm.SetState(msg.From.ID, StateWaitDescription)
		b.sendMessage(msg.Chat.ID, m.VacStep7Description)

	case StateWaitDescription:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Description = msg.Text })
		b.fsm.SetState(msg.From.ID, StateWaitSalaryFrom)
		b.sendMessage(msg.Chat.ID, m.VacStep8SalaryFrom)

	case StateWaitSalaryFrom:
		if !isSkip(msg.Text) {
			salary, err := strconv.Atoi(msg.Text)
			if err != nil {
				b.sendMessage(msg.Chat.ID, m.InvalidNumber)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.SalaryFrom = &salary })
		}
		b.fsm.SetState(msg.From.ID, StateWaitSalaryTo)
		b.sendMessage(msg.Chat.ID, m.VacStep9SalaryTo)

	case StateWaitSalaryTo:
		if !isSkip(msg.Text) {
			salary, err := strconv.Atoi(msg.Text)
			if err != nil {
				b.sendMessage(msg.Chat.ID, m.InvalidNumber)
				return
			}
			draft := b.fsm.GetDraft(msg.From.ID)
			if draft != nil && draft.SalaryFrom != nil && salary < *draft.SalaryFrom {
				b.sendMessage(msg.Chat.ID, m.SalaryToLessThanFrom)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.SalaryTo = &salary })
		}
		b.fsm.SetState(msg.From.ID, StateWaitApplyLink)
		b.sendMessage(msg.Chat.ID, m.VacStep10ApplyLink)

	case StateWaitApplyLink:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.ApplyLink = msg.Text })
		b.fsm.SetState(msg.From.ID, StatePreview)
		b.sendVacancyPreview(msg.Chat.ID, msg.From.ID)

	// ==================== RESUME STATES ====================
	case StateResumeWaitTitle:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Title = msg.Text })
		b.fsm.SetState(msg.From.ID, StateResumeWaitLevel)
		b.sendLevelKeyboard(msg.Chat.ID, lang, m.ResStep2Level)

	case StateResumeWaitLevel:
		if isSkip(msg.Text) {
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Level = domain.JobLevelSkip })
		} else {
			level := domain.JobLevel(strings.ToLower(msg.Text))
			if !isValidLevel(level) {
				b.sendMessage(msg.Chat.ID, "Select using buttons / Выберите кнопками")
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Level = level })
		}
		b.fsm.SetState(msg.From.ID, StateResumeWaitExperience)
		b.sendMessage(msg.Chat.ID, m.ResStep3Experience)

	case StateResumeWaitExperience:
		if !isSkip(msg.Text) {
			exp, err := strconv.ParseFloat(msg.Text, 64)
			if err != nil {
				b.sendMessage(msg.Chat.ID, m.InvalidExperience)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.ExperienceYears = &exp })
		}
		b.fsm.SetState(msg.From.ID, StateResumeWaitType)
		b.sendTypeKeyboard(msg.Chat.ID, lang, m.ResStep4Type)

	case StateResumeWaitType:
		jobType := domain.JobType(strings.ToLower(msg.Text))
		if !isValidType(jobType) {
			b.sendMessage(msg.Chat.ID, "Select using buttons / Выберите кнопками")
			return
		}
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Type = jobType })
		b.fsm.SetState(msg.From.ID, StateResumeWaitEmployment)
		b.sendEmploymentKeyboard(msg.Chat.ID, lang)

	case StateResumeWaitEmployment:
		emp := domain.EmploymentType(strings.ToLower(msg.Text))
		if !isValidEmployment(emp) {
			b.sendMessage(msg.Chat.ID, "Select using buttons / Выберите кнопками")
			return
		}
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.Employment = emp })
		b.fsm.SetState(msg.From.ID, StateResumeWaitSalaryFrom)
		b.sendMessage(msg.Chat.ID, m.ResStep6SalaryFrom)

	case StateResumeWaitSalaryFrom:
		if !isSkip(msg.Text) {
			salary, err := strconv.Atoi(msg.Text)
			if err != nil {
				b.sendMessage(msg.Chat.ID, m.InvalidNumber)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.SalaryFrom = &salary })
		}
		b.fsm.SetState(msg.From.ID, StateResumeWaitSalaryTo)
		b.sendMessage(msg.Chat.ID, m.ResStep7SalaryTo)

	case StateResumeWaitSalaryTo:
		if !isSkip(msg.Text) {
			salary, err := strconv.Atoi(msg.Text)
			if err != nil {
				b.sendMessage(msg.Chat.ID, m.InvalidNumber)
				return
			}
			draft := b.fsm.GetDraft(msg.From.ID)
			if draft != nil && draft.SalaryFrom != nil && salary < *draft.SalaryFrom {
				b.sendMessage(msg.Chat.ID, m.SalaryToLessThanFrom)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.SalaryTo = &salary })
		}
		b.fsm.SetState(msg.From.ID, StateResumeWaitAbout)
		b.sendMessage(msg.Chat.ID, m.ResStep8About)

	case StateResumeWaitAbout:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.About = msg.Text })
		b.fsm.SetState(msg.From.ID, StateResumeWaitContact)
		b.sendMessage(msg.Chat.ID, m.ResStep9Contact)

	case StateResumeWaitContact:
		b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.ResumeContact = msg.Text })
		b.fsm.SetState(msg.From.ID, StateResumeWaitLink)
		b.sendResumeLinkStep(msg.Chat.ID, lang)

	case StateResumeWaitLink:
		// Check if it's a file (reject files)
		if msg.Document != nil || msg.Photo != nil {
			b.sendMessage(msg.Chat.ID, m.OnlyLinksAllowed)
			return
		}
		if !isSkip(msg.Text) {
			// Validate it looks like a URL
			text := strings.TrimSpace(msg.Text)
			if !strings.HasPrefix(text, "http://") && !strings.HasPrefix(text, "https://") && !strings.HasPrefix(text, "www.") {
				b.sendMessage(msg.Chat.ID, m.OnlyLinksAllowed)
				return
			}
			b.fsm.UpdateDraft(msg.From.ID, func(d *PostDraft) { d.ResumeLink = text })
		}
		b.fsm.SetState(msg.From.ID, StateResumePreview)
		b.sendResumePreview(msg.Chat.ID, msg.From.ID)

	default:
		// Handle preview states - they wait for button clicks
		if postType == domain.PostTypeResume && userState.State == StateResumePreview {
			b.sendMessage(msg.Chat.ID, "Use buttons below / Используйте кнопки ниже")
		} else if userState.State == StatePreview {
			b.sendMessage(msg.Chat.ID, "Use buttons below / Используйте кнопки ниже")
		}
	}
}

// ==================== KEYBOARDS ====================

func (b *Bot) sendLevelKeyboard(chatID int64, lang Language, prompt string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Junior, "level:junior"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Middle, "level:middle"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Senior, "level:senior"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Internship, "level:internship"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.SkipLevel, "level:skip"),
		),
	)
	b.sendMessageWithKeyboard(chatID, prompt, keyboard)
}

func (b *Bot) sendTypeKeyboard(chatID int64, lang Language, prompt string) {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Remote, "type:remote"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Hybrid, "type:hybrid"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Onsite, "type:onsite"),
		),
	)
	b.sendMessageWithKeyboard(chatID, prompt, keyboard)
}

func (b *Bot) sendCategoryKeyboard(chatID int64, lang Language) {
	m := GetMessages(lang)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Other, "category:dev"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Web2, "category:web2"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Web3, "category:web3"),
		),
	)
	b.sendMessageWithKeyboard(chatID, m.VacStep6Category, keyboard)
}

func (b *Bot) sendEmploymentKeyboard(chatID int64, lang Language) {
	m := GetMessages(lang)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.FullTime, "employment:full-time"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.PartTime, "employment:part-time"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Contract, "employment:contract"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Freelance, "employment:freelance"),
		),
	)
	b.sendMessageWithKeyboard(chatID, m.ResStep5Employment, keyboard)
}

func (b *Bot) sendResumeLinkStep(chatID int64, lang Language) {
	m := GetMessages(lang)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Skip, "resume_link:skip"),
		),
	)
	b.sendMessageWithKeyboard(chatID, m.ResStep10Link, keyboard)
}

// ==================== PREVIEWS ====================

func (b *Bot) sendVacancyPreview(chatID int64, userID int64) {
	lang := b.fsm.GetLanguage(userID)
	m := GetMessages(lang)
	draft := b.fsm.GetDraft(userID)

	if draft == nil {
		b.sendMessage(chatID, "Error. Please start over with /post\\_job")
		return
	}

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
		m.VacPreviewTitle,
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
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Submit, "submit"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Cancel, "cancel_submit"),
		),
	)

	b.sendMessageWithKeyboard(chatID, text, keyboard)
}

func (b *Bot) sendResumePreview(chatID int64, userID int64) {
	lang := b.fsm.GetLanguage(userID)
	m := GetMessages(lang)
	draft := b.fsm.GetDraft(userID)

	if draft == nil {
		b.sendMessage(chatID, "Error. Please start over with /post\\_job")
		return
	}

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

	experience := m.NotSpecifiedLabel
	if draft.ExperienceYears != nil {
		if lang == LangEN {
			experience = fmt.Sprintf("%.1f years", *draft.ExperienceYears)
		} else {
			experience = fmt.Sprintf("%.1f лет", *draft.ExperienceYears)
		}
	}

	resumeLink := m.NotSpecifiedLabel
	if draft.ResumeLink != "" {
		resumeLink = draft.ResumeLink
	}

	text := fmt.Sprintf(`%s

💼 *%s:* %s
📊 *%s:* %s
⏱ *%s:* %s
🌍 *%s:* %s
🕒 *%s:* %s
💰 *%s:* %s

🧑‍💻 *%s:*
%s

📄 *%s:* %s
🔗 *%s:* %s

———
%s`,
		m.ResPreviewTitle,
		m.TitleLabel, escapeMarkdown(draft.Title),
		m.LevelLabel, levelDisplay,
		m.ExperienceLabel, experience,
		m.TypeLabel, draft.Type,
		m.EmploymentLabel, draft.Employment,
		m.ExpectationsLabel, salary,
		m.AboutLabel,
		escapeMarkdown(draft.About),
		m.ResumeLinkLabel, resumeLink,
		m.ContactLabel, escapeMarkdown(draft.ResumeContact),
		m.PreviewConfirm,
	)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Submit, "submit_resume"),
			tgbotapi.NewInlineKeyboardButtonData(ButtonLabels.Cancel, "cancel_submit"),
		),
	)

	b.sendMessageWithKeyboard(chatID, text, keyboard)
}

// Keep old sendPreview for backward compatibility
func (b *Bot) sendPreview(msg *tgbotapi.Message) {
	b.sendVacancyPreview(msg.Chat.ID, msg.From.ID)
}

// ==================== CALLBACKS ====================

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	userID := callback.From.ID
	chatID := callback.Message.Chat.ID

	b.api.Request(tgbotapi.NewCallback(callback.ID, ""))

	// Interface language selection
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

	// Post type selection - go directly to first step (use interface_language for post)
	if strings.HasPrefix(data, "post_type:") {
		postTypeStr := strings.TrimPrefix(data, "post_type:")
		postType := domain.PostType(postTypeStr)
		b.fsm.SetPostType(userID, postType)

		// Use interface language for both FSM and post
		interfaceLang := b.getUserInterfaceLanguage(userID)
		if interfaceLang == "" {
			interfaceLang = LangEN
		}
		b.fsm.SetLanguage(userID, interfaceLang)
		b.fsm.UpdateDraft(userID, func(d *PostDraft) { d.Language = string(interfaceLang) })

		m := GetMessages(interfaceLang)

		if postType == domain.PostTypeResume {
			b.fsm.SetState(userID, StateResumeWaitTitle)
			b.sendMessage(chatID, m.ResStep1Title)
		} else {
			b.fsm.SetState(userID, StateWaitCompany)
			b.sendMessage(chatID, m.VacStep1Company)
		}
		return
	}

	lang := b.fsm.GetLanguage(userID)
	m := GetMessages(lang)
	postType := b.fsm.GetPostType(userID)

	// Level selection
	if strings.HasPrefix(data, "level:") {
		levelStr := strings.TrimPrefix(data, "level:")
		var level domain.JobLevel
		if levelStr == "skip" {
			level = domain.JobLevelSkip
		} else {
			level = domain.JobLevel(levelStr)
		}
		b.fsm.UpdateDraft(userID, func(d *PostDraft) { d.Level = level })

		if postType == domain.PostTypeResume {
			b.fsm.SetState(userID, StateResumeWaitExperience)
			b.sendMessage(chatID, m.ResStep3Experience)
		} else {
			b.fsm.SetState(userID, StateWaitType)
			b.sendTypeKeyboard(chatID, lang, m.VacStep5Type)
		}
		return
	}

	// Type selection
	if strings.HasPrefix(data, "type:") {
		jobType := domain.JobType(strings.TrimPrefix(data, "type:"))
		b.fsm.UpdateDraft(userID, func(d *PostDraft) { d.Type = jobType })

		if postType == domain.PostTypeResume {
			b.fsm.SetState(userID, StateResumeWaitEmployment)
			b.sendEmploymentKeyboard(chatID, lang)
		} else {
			b.fsm.SetState(userID, StateWaitCategory)
			b.sendCategoryKeyboard(chatID, lang)
		}
		return
	}

	// Category selection (vacancy only)
	if strings.HasPrefix(data, "category:") {
		category := domain.JobCategory(strings.TrimPrefix(data, "category:"))
		b.fsm.UpdateDraft(userID, func(d *PostDraft) { d.Category = category })
		b.fsm.SetState(userID, StateWaitDescription)
		b.sendMessage(chatID, m.VacStep7Description)
		return
	}

	// Employment selection (resume only)
	if strings.HasPrefix(data, "employment:") {
		emp := domain.EmploymentType(strings.TrimPrefix(data, "employment:"))
		b.fsm.UpdateDraft(userID, func(d *PostDraft) { d.Employment = emp })
		b.fsm.SetState(userID, StateResumeWaitSalaryFrom)
		b.sendMessage(chatID, m.ResStep6SalaryFrom)
		return
	}

	// Resume link skip
	if data == "resume_link:skip" {
		b.fsm.SetState(userID, StateResumePreview)
		b.sendResumePreview(chatID, userID)
		return
	}

	// Submit vacancy
	if data == "submit" {
		b.submitVacancy(callback)
		return
	}

	// Submit resume
	if data == "submit_resume" {
		b.submitResume(callback)
		return
	}

	// Cancel
	if data == "cancel_submit" {
		b.fsm.Reset(userID)
		b.sendMessage(chatID, m.Cancelled)
		return
	}

	// Admin callbacks
	if strings.HasPrefix(data, "approve:") || strings.HasPrefix(data, "reject:") ||
		strings.HasPrefix(data, "delete:") || strings.HasPrefix(data, "confirm_delete:") ||
		strings.HasPrefix(data, "cancel_delete:") {
		b.handleAdminCallback(callback)
		return
	}
}

// ==================== SUBMIT ====================

func (b *Bot) submitVacancy(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	chatID := callback.Message.Chat.ID
	lang := b.fsm.GetLanguage(userID)
	m := GetMessages(lang)

	draft := b.fsm.GetDraft(userID)
	if draft == nil {
		b.sendMessage(chatID, "Error. Please start over with /post\\_job")
		return
	}

	ctx := context.Background()
	username := callback.From.UserName

	job, err := b.jobService.CreateJob(ctx, userID, username, draft.ToCreateJobRequest())
	if err != nil {
		log.Printf("Error creating vacancy: %v", err)
		b.sendMessage(chatID, m.SubmitError+err.Error())
		return
	}

	b.fsm.Reset(userID)
	b.sendMessage(chatID, fmt.Sprintf(m.SubmitVacancySuccess, job.ID.String()))
}

func (b *Bot) submitResume(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID
	chatID := callback.Message.Chat.ID
	lang := b.fsm.GetLanguage(userID)
	m := GetMessages(lang)

	draft := b.fsm.GetDraft(userID)
	if draft == nil {
		b.sendMessage(chatID, "Error. Please start over with /post\\_job")
		return
	}

	ctx := context.Background()
	username := callback.From.UserName

	resume, err := b.jobService.CreateResume(ctx, userID, username, draft.ToCreateResumeRequest())
	if err != nil {
		log.Printf("Error creating resume: %v", err)
		b.sendMessage(chatID, m.SubmitError+err.Error())
		return
	}

	b.fsm.Reset(userID)
	b.sendMessage(chatID, fmt.Sprintf(m.SubmitResumeSuccess, resume.ID.String()))
}

// Keep old submitJob for backward compatibility
func (b *Bot) submitJob(callback *tgbotapi.CallbackQuery) {
	b.submitVacancy(callback)
}

// ==================== VALIDATORS ====================

func isValidLevel(level domain.JobLevel) bool {
	return level == domain.JobLevelJunior || level == domain.JobLevelMiddle ||
		level == domain.JobLevelSenior || level == domain.JobLevelInternship
}

func isValidType(t domain.JobType) bool {
	return t == domain.JobTypeRemote || t == domain.JobTypeHybrid || t == domain.JobTypeOnsite
}

func isValidCategory(c domain.JobCategory) bool {
	return c == domain.JobCategoryWeb2 || c == domain.JobCategoryWeb3 || c == domain.JobCategoryDev
}

func isValidEmployment(e domain.EmploymentType) bool {
	return e == domain.EmploymentFullTime || e == domain.EmploymentPartTime ||
		e == domain.EmploymentContract || e == domain.EmploymentFreelance
}

func isSkip(text string) bool {
	lower := strings.ToLower(strings.TrimSpace(text))
	return lower == "skip" || lower == "скип" || lower == "пропустить"
}
