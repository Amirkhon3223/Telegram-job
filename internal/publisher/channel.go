package publisher

import (
	"context"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-job/internal/domain"
)

type ChannelPublisher struct {
	bot       *tgbotapi.BotAPI
	channelID int64
}

func NewChannelPublisher(bot *tgbotapi.BotAPI, channelID int64) *ChannelPublisher {
	return &ChannelPublisher{
		bot:       bot,
		channelID: channelID,
	}
}

func (p *ChannelPublisher) Publish(ctx context.Context, job *domain.JobWithCompany) (int, error) {
	text := formatJobPost(job)

	msg := tgbotapi.NewMessage(p.channelID, text)
	msg.ParseMode = "Markdown"
	msg.DisableWebPagePreview = true

	sent, err := p.bot.Send(msg)
	if err != nil {
		return 0, err
	}
	return sent.MessageID, nil
}

func (p *ChannelPublisher) Delete(ctx context.Context, messageID int) error {
	deleteMsg := tgbotapi.NewDeleteMessage(p.channelID, messageID)
	_, err := p.bot.Request(deleteMsg)
	return err
}

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

func formatJobPost(job *domain.JobWithCompany) string {
	salary := "Не указана"
	if job.SalaryFrom != nil && job.SalaryTo != nil {
		salary = fmt.Sprintf("$%d – $%d", *job.SalaryFrom, *job.SalaryTo)
	} else if job.SalaryFrom != nil {
		salary = fmt.Sprintf("От $%d", *job.SalaryFrom)
	} else if job.SalaryTo != nil {
		salary = fmt.Sprintf("До $%d", *job.SalaryTo)
	}

	levelEmoji := map[domain.JobLevel]string{
		domain.JobLevelJunior:     "🌱",
		domain.JobLevelMiddle:     "🌿",
		domain.JobLevelSenior:     "🌳",
		domain.JobLevelInternship: "🎓",
		domain.JobLevelSkip:       "📊",
	}

	levelDisplay := string(job.Level)
	if job.Level == "" {
		levelDisplay = "Не указан"
	}

	typeEmoji := map[domain.JobType]string{
		domain.JobTypeRemote: "🌍",
		domain.JobTypeHybrid: "🏢🏠",
		domain.JobTypeOnsite: "🏢",
	}

	categoryEmoji := map[domain.JobCategory]string{
		domain.JobCategoryWeb2: "🌐",
		domain.JobCategoryWeb3: "⛓️",
		domain.JobCategoryDev:  "💻",
	}

	return fmt.Sprintf(`#vacancy #вакансия

*Вакансия: %s*

🏢 *Компания:* %s
%s *Уровень:* %s
%s *Формат:* %s
%s *Категория:* %s
💰 *Зарплата:* %s

📝 *Описание:*
%s

🔗 *Откликнуться:* %s

———
📮 _Разместить вакансию: @BridgeJobsBot_`,
		escapeMarkdown(job.Title),
		escapeMarkdown(job.CompanyName),
		levelEmoji[job.Level], levelDisplay,
		typeEmoji[job.Type], job.Type,
		categoryEmoji[job.Category], job.Category,
		salary,
		escapeMarkdown(job.Description),
		escapeMarkdown(job.ApplyLink),
	)
}
