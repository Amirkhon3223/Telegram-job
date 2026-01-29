package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"telegram-job/internal/bot"
	"telegram-job/internal/config"
	"telegram-job/internal/publisher"
	"telegram-job/internal/repository"
	"telegram-job/internal/service"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := repository.NewDB(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	jobRepo := repository.NewJobRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Initialize bot first (to get bot API)
	telegramBot, err := bot.New(cfg, nil)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Initialize publisher and notifier with the same bot API
	channelPublisher := publisher.NewChannelPublisher(telegramBot.GetAPI(), cfg.ChannelID)
	adminNotifier := bot.NewAdminNotifier(telegramBot.GetAPI(), cfg.AdminTelegramIDs)

	// Initialize service with publisher and notifier
	jobService := service.NewJobService(cfg, jobRepo, companyRepo, userRepo, channelPublisher, adminNotifier)

	// Set service to bot (use same bot instance!)
	telegramBot.SetJobService(jobService)

	// Start cleanup service (auto-archive old jobs)
	cleanupService := bot.NewCleanupService(jobRepo, channelPublisher, cfg.JobMaxDays)
	go cleanupService.Start(ctx)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	log.Printf("Bot starting... (auto-cleanup after %d days)", cfg.JobMaxDays)
	telegramBot.Start()
}
