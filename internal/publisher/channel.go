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

func (p *ChannelPublisher) Publish(ctx context.Context, post *domain.PostWithDetails) (int, error) {
	var text string
	if post.PostType == domain.PostTypeResume {
		text = formatResumePost(post)
	} else {
		text = formatJobPost(post)
	}

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

func formatJobPost(post *domain.PostWithDetails) string {
	salary := "Не указана"
	if post.SalaryFrom != nil && post.SalaryTo != nil {
		salary = fmt.Sprintf("$%d – $%d", *post.SalaryFrom, *post.SalaryTo)
	} else if post.SalaryFrom != nil {
		salary = fmt.Sprintf("От $%d", *post.SalaryFrom)
	} else if post.SalaryTo != nil {
		salary = fmt.Sprintf("До $%d", *post.SalaryTo)
	}

	levelEmoji := map[domain.JobLevel]string{
		domain.JobLevelJunior:     "🌱",
		domain.JobLevelMiddle:     "🌿",
		domain.JobLevelSenior:     "🌳",
		domain.JobLevelInternship: "🎓",
		domain.JobLevelSkip:       "📊",
	}

	levelDisplay := string(post.Level)
	if post.Level == "" {
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
		escapeMarkdown(post.Title),
		escapeMarkdown(post.CompanyName),
		levelEmoji[post.Level], levelDisplay,
		typeEmoji[post.Type], post.Type,
		categoryEmoji[post.Category], post.Category,
		salary,
		escapeMarkdown(post.Description),
		escapeMarkdown(post.ApplyLink),
	)
}

func formatResumePost(post *domain.PostWithDetails) string {
	salary := "Не указаны"
	if post.SalaryFrom != nil && post.SalaryTo != nil {
		salary = fmt.Sprintf("$%d – $%d", *post.SalaryFrom, *post.SalaryTo)
	} else if post.SalaryFrom != nil {
		salary = fmt.Sprintf("От $%d", *post.SalaryFrom)
	} else if post.SalaryTo != nil {
		salary = fmt.Sprintf("До $%d", *post.SalaryTo)
	}

	levelEmoji := map[domain.JobLevel]string{
		domain.JobLevelJunior:     "🌱",
		domain.JobLevelMiddle:     "🌿",
		domain.JobLevelSenior:     "🌳",
		domain.JobLevelInternship: "🎓",
		domain.JobLevelSkip:       "📊",
	}

	levelDisplay := string(post.Level)
	if post.Level == "" {
		levelDisplay = "Не указан"
	}

	typeEmoji := map[domain.JobType]string{
		domain.JobTypeRemote: "🌍",
		domain.JobTypeHybrid: "🏢🏠",
		domain.JobTypeOnsite: "🏢",
	}

	employmentDisplay := string(post.Employment)
	if post.Employment == "" {
		employmentDisplay = "Не указана"
	}

	experience := "Не указан"
	if post.ExperienceYears != nil {
		experience = fmt.Sprintf("%.1f лет", *post.ExperienceYears)
	}

	resumeLink := ""
	if post.ResumeLink != "" {
		resumeLink = fmt.Sprintf("\n📄 *Резюме:* %s", post.ResumeLink)
	}

	return fmt.Sprintf(`#resume #резюме

*Резюме: %s*

%s *Уровень:* %s
⏱ *Опыт:* %s
%s *Формат:* %s
🕒 *Занятость:* %s
💰 *Ожидания:* %s

🧑‍💻 *О кандидате:*
%s
%s
🔗 *Контакт:* %s

———
📮 _Разместить резюме: @BridgeJobsBot_`,
		escapeMarkdown(post.Title),
		levelEmoji[post.Level], levelDisplay,
		experience,
		typeEmoji[post.Type], post.Type,
		employmentDisplay,
		salary,
		escapeMarkdown(post.About),
		resumeLink,
		escapeMarkdown(post.Contact),
	)
}
