package bot

import (
	"context"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-job/internal/config"
	"telegram-job/internal/repository"
	"telegram-job/internal/service"
)

type Bot struct {
	api        *tgbotapi.BotAPI
	cfg        *config.Config
	jobService *service.JobService
	userRepo   *repository.UserRepository
	fsm        *FSM
}

func New(cfg *config.Config, jobService *service.JobService, userRepo *repository.UserRepository) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	return &Bot{
		api:        api,
		cfg:        cfg,
		jobService: jobService,
		userRepo:   userRepo,
		fsm:        NewFSM(),
	}, nil
}

func NewWithService(cfg *config.Config, jobService *service.JobService, userRepo *repository.UserRepository) (*Bot, error) {
	return New(cfg, jobService, userRepo)
}

func (b *Bot) SetJobService(jobService *service.JobService) {
	b.jobService = jobService
}

func (b *Bot) GetAPI() *tgbotapi.BotAPI {
	return b.api
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			b.handleCallback(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.handleCommand(update.Message)
			continue
		}

		b.handleMessage(update.Message)
	}
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

func (b *Bot) sendMessageWithKeyboard(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

// getUserInterfaceLanguage returns the user's interface language or empty string if not set
func (b *Bot) getUserInterfaceLanguage(telegramID int64) Language {
	ctx := context.Background()
	user, err := b.userRepo.GetByTelegramID(ctx, telegramID)
	if err != nil || user.InterfaceLanguage == nil {
		return "" // Not set
	}
	return Language(*user.InterfaceLanguage)
}

// setUserInterfaceLanguage sets the user's interface language
func (b *Bot) setUserInterfaceLanguage(telegramID int64, username string, lang Language) error {
	ctx := context.Background()
	// Ensure user exists
	_, err := b.userRepo.GetOrCreate(ctx, telegramID, username)
	if err != nil {
		return err
	}
	return b.userRepo.SetInterfaceLanguage(ctx, telegramID, string(lang))
}

// getInterfaceMessages returns messages in user's interface language
func (b *Bot) getInterfaceMessages(telegramID int64) Messages {
	lang := b.getUserInterfaceLanguage(telegramID)
	if lang == "" {
		lang = LangEN // fallback to English
	}
	return GetMessages(lang)
}
