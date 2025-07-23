package main

import (
    "log"
    "net/http"
    "os"
    "strconv"
    "time"

    "github.com/joho/godotenv"
    _ "github.com/lib/pq"

    "github.com/nurzhanova2/ci-platform/configs"
    "github.com/nurzhanova2/ci-platform/internal/handler"
    "github.com/nurzhanova2/ci-platform/internal/logger"
    "github.com/nurzhanova2/ci-platform/internal/repository"
    "github.com/nurzhanova2/ci-platform/internal/service"
)

func main() {
    // Load .env (if present)
    if err := godotenv.Load(); err != nil {
        log.Println(".env file not found. Using system environment variables.")
    }

    // Initialize logger
    logger.Init()
    logger.Info.Println("Starting CI/CD platform...")

    // Load configuration
    cfg, err := config.Load("configs/config.yaml")
    if err != nil {
        logger.Error.Fatalf("Failed to load configuration: %v", err)
    }

    // Override config values with env vars if not set
    if cfg.Database.DSN == "" {
        cfg.Database.DSN = os.Getenv("DB_DSN")
    }
    if cfg.Database.Driver == "" {
        cfg.Database.Driver = os.Getenv("DB_DRIVER")
    }

    // Connect to the database with retries
    var jobRepo *repository.DBRepository
    for i := 1; i <= 10; i++ {
        jobRepo, err = repository.NewJobDB(cfg.Database.DSN, cfg.Database.Driver)
        if err == nil {
            logger.Info.Println("✅ Successfully connected to the database.")
            break
        }
        logger.Error.Printf("❌ Connection attempt #%d failed: %v", i, err)
        time.Sleep(3 * time.Second)
    }
    if jobRepo == nil || err != nil {
        logger.Error.Fatalf("Failed to connect to the database after 10 attempts: %v", err)
    }

    // Initialize repositories
    gitRepo := repository.NewGitRepository()
    docker := repository.NewDockerRunner()

    // Notification service
    notifierSvc := service.NewNotifierService(
        os.Getenv("TELEGRAM_TOKEN"),
        os.Getenv("TELEGRAM_CHAT_ID"),
        os.Getenv("SLACK_WEBHOOK_URL"),
        true,
        false,
    )

    // Pipeline log file
    logWriter, err := logger.NewLogWriter("pipeline.log")
    if err != nil {
        logger.Error.Fatalf("Failed to create log file: %v", err)
    }
    defer logWriter.Close()

    // Initialize Pipeline service
    pipelineSvc := service.NewPipelineService(gitRepo, docker, notifierSvc, logWriter)

    // Webhook handler
    webhookHandler := handler.NewWebhookHandler(pipelineSvc)

    addr := ":" + strconv.Itoa(cfg.Server.Port)
    http.HandleFunc("/webhook", webhookHandler.HandleWebhook)

    logger.Info.Printf("🚀 HTTP server listening on %s", addr)
    if err := http.ListenAndServe(addr, nil); err != nil {
        logger.Error.Fatalf("Server startup error: %v", err)
    }
}