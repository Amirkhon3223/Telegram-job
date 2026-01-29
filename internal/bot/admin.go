package bot

import (
	"context"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"telegram-job/internal/domain"
)

// AdminNotifier sends notifications to admins about new jobs
type AdminNotifier struct {
	bot      *tgbotapi.BotAPI
	adminIDs map[int64]bool
}

func NewAdminNotifier(bot *tgbotapi.BotAPI, adminIDs map[int64]bool) *AdminNotifier {
	return &AdminNotifier{
		bot:      bot,
		adminIDs: adminIDs,
	}
}

func (n *AdminNotifier) NotifyNewJob(ctx context.Context, post *domain.PostWithDetails) error {
	log.Printf("NotifyNewJob called for post %s (type: %s)", post.ID.String(), post.PostType)
	log.Printf("Admin IDs to notify: %v", n.adminIDs)

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

	for adminID := range n.adminIDs {
		log.Printf("Sending notification to admin %d", adminID)
		msg := tgbotapi.NewMessage(adminID, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = keyboard
		resp, err := n.bot.Send(msg)
		if err != nil {
			log.Printf("Error sending to admin %d: %v", adminID, err)
		} else {
			log.Printf("Sent to admin %d, message ID: %d", adminID, resp.MessageID)
		}
	}

	return nil
}

func (n *AdminNotifier) NotifyAuthor(authorTelegramID int64, approved bool, postTitle string, postLanguage string, postType domain.PostType) {
	var text string
	isResume := postType == domain.PostTypeResume

	if postLanguage == "en" {
		if approved {
			if isResume {
				text = fmt.Sprintf("✅ *Your resume has been approved!*\n\n"+
					"Resume *%s* is now published in @BridgeJob channel\n\n"+
					"📢 View: https://t.me/BridgeJob\n\n"+
					"Thank you for using @BridgeJobsBot!", escapeMarkdownAdmin(postTitle))
			} else {
				text = fmt.Sprintf("✅ *Your job has been approved!*\n\n"+
					"Job *%s* is now published in @BridgeJob channel\n\n"+
					"📢 View: https://t.me/BridgeJob\n\n"+
					"Thank you for using @BridgeJobsBot!", escapeMarkdownAdmin(postTitle))
			}
		} else {
			if isResume {
				text = fmt.Sprintf("❌ *Your resume has been rejected*\n\n"+
					"Resume *%s* did not pass moderation.\n\n"+
					"Please try again with correct data: /post\\_job\n\n"+
					"📢 Channel: @BridgeJob", escapeMarkdownAdmin(postTitle))
			} else {
				text = fmt.Sprintf("❌ *Your job has been rejected*\n\n"+
					"Job *%s* did not pass moderation.\n\n"+
					"Please try again with correct data: /post\\_job\n\n"+
					"📢 Channel: @BridgeJob", escapeMarkdownAdmin(postTitle))
			}
		}
	} else {
		if approved {
			if isResume {
				text = fmt.Sprintf("✅ *Ваше резюме одобрено!*\n\n"+
					"Резюме *%s* опубликовано в канале @BridgeJob\n\n"+
					"📢 Смотреть: https://t.me/BridgeJob\n\n"+
					"Спасибо за использование @BridgeJobsBot!", escapeMarkdownAdmin(postTitle))
			} else {
				text = fmt.Sprintf("✅ *Ваша вакансия одобрена!*\n\n"+
					"Вакансия *%s* опубликована в канале @BridgeJob\n\n"+
					"📢 Смотреть: https://t.me/BridgeJob\n\n"+
					"Спасибо за использование @BridgeJobsBot!", escapeMarkdownAdmin(postTitle))
			}
		} else {
			if isResume {
				text = fmt.Sprintf("❌ *Ваше резюме отклонено*\n\n"+
					"Резюме *%s* не прошло модерацию.\n\n"+
					"Попробуйте отправить заново с корректными данными: /post\\_job\n\n"+
					"📢 Канал: @BridgeJob", escapeMarkdownAdmin(postTitle))
			} else {
				text = fmt.Sprintf("❌ *Ваша вакансия отклонена*\n\n"+
					"Вакансия *%s* не прошла модерацию.\n\n"+
					"Попробуйте отправить заново с корректными данными: /post\\_job\n\n"+
					"📢 Канал: @BridgeJob", escapeMarkdownAdmin(postTitle))
			}
		}
	}

	msg := tgbotapi.NewMessage(authorTelegramID, text)
	msg.ParseMode = "Markdown"
	n.bot.Send(msg)
}

func (n *AdminNotifier) NotifyAuthorDeleted(authorTelegramID int64, postTitle string, postLanguage string, postType domain.PostType) {
	var text string
	isResume := postType == domain.PostTypeResume

	if postLanguage == "en" {
		if isResume {
			text = fmt.Sprintf("🗑 *Your resume has been removed from channel*\n\n"+
				"Resume *%s* was removed from @BridgeJob\n\n"+
				"To post again: /post\\_job", escapeMarkdownAdmin(postTitle))
		} else {
			text = fmt.Sprintf("🗑 *Your job has been removed from channel*\n\n"+
				"Job *%s* was removed from @BridgeJob\n\n"+
				"To post a new job: /post\\_job", escapeMarkdownAdmin(postTitle))
		}
	} else {
		if isResume {
			text = fmt.Sprintf("🗑 *Ваше резюме удалено из канала*\n\n"+
				"Резюме *%s* было удалено из @BridgeJob\n\n"+
				"Если хотите разместить новое: /post\\_job", escapeMarkdownAdmin(postTitle))
		} else {
			text = fmt.Sprintf("🗑 *Ваша вакансия удалена из канала*\n\n"+
				"Вакансия *%s* была удалена из @BridgeJob\n\n"+
				"Если хотите разместить новую вакансию: /post\\_job", escapeMarkdownAdmin(postTitle))
		}
	}

	msg := tgbotapi.NewMessage(authorTelegramID, text)
	msg.ParseMode = "Markdown"
	n.bot.Send(msg)
}

func escapeMarkdownAdmin(s string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"`", "\\`",
	)
	return replacer.Replace(s)
}

func formatAdminNotification(post *domain.PostWithDetails) string {
	salary := "Не указана"
	if post.SalaryFrom != nil && post.SalaryTo != nil {
		salary = fmt.Sprintf("$%d – $%d", *post.SalaryFrom, *post.SalaryTo)
	} else if post.SalaryFrom != nil {
		salary = fmt.Sprintf("От $%d", *post.SalaryFrom)
	} else if post.SalaryTo != nil {
		salary = fmt.Sprintf("До $%d", *post.SalaryTo)
	}

	levelDisplay := string(post.Level)
	if post.Level == "" {
		levelDisplay = "Не указан"
	}

	langDisplay := "🇷🇺 RU"
	if post.Language == "en" {
		langDisplay = "🇬🇧 EN"
	}

	// Resume format
	if post.PostType == domain.PostTypeResume {
		experience := "Не указан"
		if post.ExperienceYears != nil {
			experience = fmt.Sprintf("%.1f лет", *post.ExperienceYears)
		}

		contact := post.Contact
		resumeLink := "Не указана"
		if post.ResumeLink != "" {
			resumeLink = post.ResumeLink
		}

		return fmt.Sprintf(`👤 *Новое резюме на модерацию*

🌐 *Язык:* %s
💼 *Позиция:* %s
📊 *Уровень:* %s
⏱ *Опыт:* %s
🌍 *Формат:* %s
🕒 *Занятость:* %s
💰 *Ожидания:* %s
📄 *Ссылка на резюме:* %s
📞 *Контакт:* %s

🧑‍💻 *О кандидате:*
%s

———
Resume ID: `+"`%s`",
			langDisplay,
			escapeMarkdownAdmin(post.Title),
			levelDisplay,
			experience,
			post.Type,
			post.Employment,
			salary,
			resumeLink,
			escapeMarkdownAdmin(contact),
			escapeMarkdownAdmin(post.About),
			post.ID.String(),
		)
	}

	// Vacancy format
	return fmt.Sprintf(`🏢 *Новая вакансия на модерацию*

🌐 *Язык:* %s
🏢 *Компания:* %s
💼 *Должность:* %s
📊 *Уровень:* %s
🌍 *Формат:* %s
🏷️ *Категория:* %s
💰 *Зарплата:* %s
🔗 *Ссылка для кандидатов:* %s
📞 *Контакт автора:* %s

📝 *Описание:*
%s

———
Job ID: `+"`%s`",
		langDisplay,
		escapeMarkdownAdmin(post.CompanyName),
		escapeMarkdownAdmin(post.Title),
		levelDisplay,
		post.Type,
		post.Category,
		salary,
		escapeMarkdownAdmin(post.ApplyLink),
		escapeMarkdownAdmin(post.CompanyContact),
		escapeMarkdownAdmin(post.Description),
		post.ID.String(),
	)
}

func (b *Bot) handleAdminCallback(callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	adminID := callback.From.ID
	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	log.Printf("handleAdminCallback: data=%s, adminID=%d, chatID=%d, messageID=%d", data, adminID, chatID, messageID)

	// Check admin permission
	if !b.cfg.IsAdmin(adminID) {
		log.Printf("Admin check failed for ID %d", adminID)
		b.api.Request(tgbotapi.NewCallback(callback.ID, "You are not authorized"))
		return
	}

	ctx := context.Background()

	// Handle approve
	if strings.HasPrefix(data, "approve:") {
		jobIDStr := strings.TrimPrefix(data, "approve:")
		jobID, err := uuid.Parse(jobIDStr)
		if err != nil {
			b.sendMessage(chatID, "Invalid job ID")
			return
		}

		// Получаем данные о вакансии до approve
		jobInfo, _ := b.jobService.GetJobWithCompany(ctx, jobID)

		err = b.jobService.ApproveJob(ctx, jobID, adminID)
		if err != nil {
			// Если вакансия уже обработана - не показываем ошибку
			if err.Error() == "invalid status transition" {
				log.Printf("Job %s already processed", jobIDStr)
				return
			}
			b.sendMessage(chatID, "Failed to approve: "+err.Error())
			return
		}

		// Обновляем сообщение с кнопкой удаления
		newText := callback.Message.Text + "\n\n✅ ОДОБРЕНО И ОПУБЛИКОВАНО"
		deleteKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить из канала", "delete:"+jobIDStr),
			),
		)
		edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, newText, deleteKeyboard)
		_, err = b.api.Send(edit)
		if err != nil {
			log.Printf("Error editing message after approve: %v", err)
		}

		// Уведомляем автора
		if jobInfo != nil {
			notifier := NewAdminNotifier(b.api, b.cfg.AdminTelegramIDs)
			notifier.NotifyAuthor(jobInfo.AuthorTelegramID, true, jobInfo.Title, jobInfo.Language, jobInfo.PostType)
		}

		return
	}

	// Handle reject
	if strings.HasPrefix(data, "reject:") {
		jobIDStr := strings.TrimPrefix(data, "reject:")
		jobID, err := uuid.Parse(jobIDStr)
		if err != nil {
			b.sendMessage(chatID, "Invalid job ID")
			return
		}

		// Получаем данные о публикации до reject
		jobInfo, _ := b.jobService.GetJobWithCompany(ctx, jobID)

		err = b.jobService.RejectJob(ctx, jobID, adminID, "Rejected by admin")
		if err != nil {
			b.sendMessage(chatID, "Failed to reject: "+err.Error())
			return
		}

		// Обновляем сообщение без кнопок
		newText := callback.Message.Text + "\n\n❌ ОТКЛОНЕНО"
		emptyKeyboard := tgbotapi.NewInlineKeyboardMarkup()
		edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, newText, emptyKeyboard)
		_, err = b.api.Send(edit)
		if err != nil {
			log.Printf("Error editing message after reject: %v", err)
		}

		// Уведомляем автора
		if jobInfo != nil {
			notifier := NewAdminNotifier(b.api, b.cfg.AdminTelegramIDs)
			notifier.NotifyAuthor(jobInfo.AuthorTelegramID, false, jobInfo.Title, jobInfo.Language, jobInfo.PostType)
		}

		return
	}

	// Handle delete (show confirmation)
	if strings.HasPrefix(data, "delete:") {
		jobIDStr := strings.TrimPrefix(data, "delete:")

		confirmKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("✅ Да, удалить", "confirm_delete:"+jobIDStr),
				tgbotapi.NewInlineKeyboardButtonData("❌ Отмена", "cancel_delete:"+jobIDStr),
			),
		)

		// Добавляем предупреждение в текст сообщения
		newText := callback.Message.Text + "\n\n⚠️ Вы уверены, что хотите удалить?"
		edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, newText, confirmKeyboard)
		_, err := b.api.Send(edit)
		if err != nil {
			log.Printf("Error showing delete confirmation: %v", err)
		}
		return
	}

	// Handle confirm delete
	if strings.HasPrefix(data, "confirm_delete:") {
		jobIDStr := strings.TrimPrefix(data, "confirm_delete:")
		jobID, err := uuid.Parse(jobIDStr)
		if err != nil {
			log.Printf("Invalid job ID: %s", jobIDStr)
			b.sendMessage(chatID, "Invalid job ID")
			return
		}

		// Получаем данные о вакансии ДО удаления
		jobInfo, _ := b.jobService.GetJobWithCompany(ctx, jobID)

		log.Printf("Archiving job %s", jobID.String())
		err = b.jobService.ArchiveJob(ctx, jobID, adminID)
		if err != nil {
			log.Printf("Failed to archive job: %v", err)
			b.sendMessage(chatID, "Failed to delete: "+err.Error())
			return
		}

		log.Printf("Job %s archived successfully", jobID.String())

		// Убираем предупреждение из текста и добавляем статус удаления
		originalText := callback.Message.Text
		originalText = strings.Replace(originalText, "\n\n⚠️ Вы уверены, что хотите удалить?", "", 1)
		newText := originalText + "\n\n🗑 УДАЛЕНО ИЗ КАНАЛА"

		// Редактируем сообщение без кнопок
		edit := tgbotapi.NewEditMessageText(chatID, messageID, newText)
		_, err = b.api.Send(edit)
		if err != nil {
			log.Printf("Error updating message after delete: %v", err)
		}

		// Уведомляем автора об удалении
		if jobInfo != nil {
			notifier := NewAdminNotifier(b.api, b.cfg.AdminTelegramIDs)
			notifier.NotifyAuthorDeleted(jobInfo.AuthorTelegramID, jobInfo.Title, jobInfo.Language, jobInfo.PostType)
		}

		return
	}

	// Handle cancel delete
	if strings.HasPrefix(data, "cancel_delete:") {
		jobIDStr := strings.TrimPrefix(data, "cancel_delete:")

		// Убираем предупреждение из текста
		originalText := callback.Message.Text
		originalText = strings.Replace(originalText, "\n\n⚠️ Вы уверены, что хотите удалить?", "", 1)

		// Возвращаем кнопку удаления
		deleteKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🗑 Удалить из канала", "delete:"+jobIDStr),
			),
		)
		edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, originalText, deleteKeyboard)
		_, err := b.api.Send(edit)
		if err != nil {
			log.Printf("Error cancelling delete: %v", err)
		}

		return
	}
}
