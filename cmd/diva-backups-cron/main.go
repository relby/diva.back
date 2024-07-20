package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"github.com/relby/diva.back/internal/app"
	"github.com/relby/diva.back/internal/closer"

	"github.com/go-co-op/gocron/v2"
)

var (
	BACKUPS_DIR_PATH = filepath.Join(".", "backups")
)

const (
	BACKUP_FILE_NAME_TIME_FORMAT = "20060102150405"

	S3_REGION_ENV_NAME     = "S3_REGION"
	S3_ENDPOINT_ENV_NAME   = "S3_ENDPOINT"
	S3_TENANT_ID_ENV_NAME  = "S3_TENANT_ID"
	S3_KEY_ID_ENV_NAME     = "S3_KEY_ID"
	S3_KEY_SECRET_ENV_NAME = "S3_KEY_SECRET"
	S3_BACKUPS_BUCKET_NAME = "diva-backups"

	TELEGRAM_TOKEN_ENV_NAME   = "TELEGRAM_TOKEN"
	TELEGRAM_CHAT_ID_ENV_NAME = "TELEGRAM_CHAT_ID"
)

func getS3Session() *session.Session {
	s3Region := os.Getenv(S3_REGION_ENV_NAME)
	if s3Region == "" {
		panic(fmt.Errorf("provide `%s` in env", S3_REGION_ENV_NAME))
	}

	s3Endpoint := os.Getenv(S3_ENDPOINT_ENV_NAME)
	if s3Endpoint == "" {
		panic(fmt.Errorf("provide `%s` in env", S3_ENDPOINT_ENV_NAME))
	}

	s3TenantID := os.Getenv(S3_TENANT_ID_ENV_NAME)
	if s3TenantID == "" {
		panic(fmt.Errorf("provide `%s` in env", S3_TENANT_ID_ENV_NAME))
	}

	s3KeyID := os.Getenv(S3_KEY_ID_ENV_NAME)
	if s3KeyID == "" {
		panic(fmt.Errorf("provide `%s` in env", S3_KEY_ID_ENV_NAME))
	}

	s3KeySecret := os.Getenv(S3_KEY_SECRET_ENV_NAME)
	if s3KeySecret == "" {
		panic(fmt.Errorf("provide `%s` in env", S3_KEY_SECRET_ENV_NAME))
	}

	s3SessionConfig := aws.NewConfig()
	s3SessionConfig.Region = aws.String(s3Region)
	s3SessionConfig.Endpoint = aws.String(s3Endpoint)
	s3SessionConfig.Credentials = credentials.NewStaticCredentials(
		fmt.Sprintf("%s:%s", s3TenantID, s3KeyID),
		s3KeySecret,
		"",
	)

	s3Session := session.Must(session.NewSession(s3SessionConfig))

	return s3Session
}

func getS3Uploader(s3Session *session.Session) *s3manager.Uploader {
	return s3manager.NewUploader(s3Session)
}

func uploadBackupToS3(backupFileName string, backupFile *os.File, s3Uploader *s3manager.Uploader) {
	result, err := s3Uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(S3_BACKUPS_BUCKET_NAME),
		Key:    aws.String(backupFileName),
		Body:   backupFile,
	})
	if err != nil {
		fmt.Printf("failed to upload backup to s3: %v\n", err)
		return
	}

	fmt.Printf("successfully uploaded backup to s3: %s\n", result.Location)
}

type telegramConfig struct {
	token  string
	chatId int64
}

func getTelegramConfig() *telegramConfig {
	chatIDString := os.Getenv(TELEGRAM_CHAT_ID_ENV_NAME)
	if chatIDString == "" {
		panic(fmt.Errorf("provide `%s` in env", TELEGRAM_CHAT_ID_ENV_NAME))
	}

	chatID, err := strconv.ParseInt(chatIDString, 10, 64)
	if err != nil {
		panic(fmt.Errorf("couldn't parse `%s`", TELEGRAM_CHAT_ID_ENV_NAME))
	}

	token := os.Getenv(TELEGRAM_TOKEN_ENV_NAME)
	if token == "" {
		panic(fmt.Errorf("provide `%s` in env", TELEGRAM_TOKEN_ENV_NAME))
	}

	return &telegramConfig{
		token:  token,
		chatId: chatID,
	}
}

func getTelegramBot(config *telegramConfig) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(config.token)
	if err != nil {
		panic(fmt.Errorf("failed to connect to telegram: %v", err))
	}

	return bot
}

func uploadBackupToTelegram(telegramBot *tgbotapi.BotAPI, telegramChatID int64, backupFileName string, backupFile *os.File) {
	backupFileBytes, err := io.ReadAll(backupFile)
	if err != nil {
		fmt.Printf("failed to read backup file to bytes: %v\n", err)
		return
	}

	fileBytes := tgbotapi.FileBytes{
		Name:  backupFileName,
		Bytes: backupFileBytes,
	}

	_, err = telegramBot.Send(tgbotapi.NewDocument(telegramChatID, fileBytes))
	if err != nil {
		fmt.Printf("failed to upload backup file to telegram: %v\n", err)
		return
	}

	fmt.Println("successfully uploaded backup to telegram")
}

type BackupUploadDestinations []string

func (i *BackupUploadDestinations) String() string {
	return strings.Join(*i, ", ")
}

func (i *BackupUploadDestinations) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var backupUploadDestinations BackupUploadDestinations
	backupUploadDestinationFlagName := "upload-destination"
	flag.Var(&backupUploadDestinations, backupUploadDestinationFlagName, "")
	flag.Parse()

	var isS3UploadEnabled, isTelegramUploadEnabled bool
	if len(backupUploadDestinations) == 0 {
		panic(fmt.Errorf("provide `%s` flags", backupUploadDestinationFlagName))
	}

	for _, backupUploadDestination := range backupUploadDestinations {
		switch backupUploadDestination {
		case "s3":
			isS3UploadEnabled = true
		case "telegram":
			isTelegramUploadEnabled = true
		default:
			fmt.Printf("Invalid `%s` flag. Must be one of `s3`, `telegram`.\n", backupUploadDestinationFlagName)
			return
		}
	}

	if isS3UploadEnabled {
		fmt.Println("Uploading to s3")
	}

	if isTelegramUploadEnabled {
		fmt.Println("Uploading to telegram")
	}

	if err := godotenv.Load(); err != nil {
		fmt.Printf("failed to read .env file: %v", err)
	}

	diContainer, err := app.NewDIContainer()
	if err != nil {
		panic(fmt.Errorf("failed to create DI container: %w", err))
	}

	postgresConfig, err := diContainer.PostgresConfig()
	if err != nil {
		panic(fmt.Errorf("failed to get postgres config: %w", err))
	}

	os.MkdirAll(BACKUPS_DIR_PATH, os.ModePerm)

	var s3Session *session.Session
	var s3Uploader *s3manager.Uploader
	if isS3UploadEnabled {
		s3Session = getS3Session()
		s3Uploader = getS3Uploader(s3Session)
	}

	var telegramConfig *telegramConfig
	var telegramBot *tgbotapi.BotAPI
	if isTelegramUploadEnabled {
		telegramConfig = getTelegramConfig()
		telegramBot = getTelegramBot(telegramConfig)
	}

	createBackup := func() {
		backupFileName := fmt.Sprintf("diva-backup_%s.sql", time.Now().Format(BACKUP_FILE_NAME_TIME_FORMAT))
		backupFilePath := filepath.Join(BACKUPS_DIR_PATH, backupFileName)
		cmd := exec.Command(
			"pg_dump",
			"-U", postgresConfig.User(),
			"-h", postgresConfig.Host(),
			"-p", postgresConfig.Port(),
			"-f", backupFilePath,
			"-F", "c",
			postgresConfig.DB(),
		)

		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", postgresConfig.DB()))

		cmdOutput, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("pg_dump output:\n%s\n", cmdOutput)
			fmt.Printf("failed to execute pg_dump: %v\n", err)
			return
		}

		fmt.Printf("successfully created backup `%s`\n", backupFileName)

		backupFile, err := os.Open(backupFilePath)
		if err != nil {
			fmt.Printf("failed to open backup file: %v\n", err)
			return
		}
		defer backupFile.Close()

		if isS3UploadEnabled {
			uploadBackupToS3(backupFileName, backupFile, s3Uploader)
		}
		if isTelegramUploadEnabled {
			uploadBackupToTelegram(telegramBot, telegramConfig.chatId, backupFileName, backupFile)
		}
	}

	cronScheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(fmt.Errorf("failed to create cron scheduler: %w", err))
	}
	closer.Add(func() error {
		err := cronScheduler.Shutdown()
		if err != nil {
			return fmt.Errorf("failed to shutdown cron scheduler: %w", err)
		}
		return nil
	})

	_, err = cronScheduler.NewJob(
		gocron.WeeklyJob(
			1,
			gocron.NewWeekdays(time.Monday, time.Wednesday, time.Friday),
			gocron.NewAtTimes(gocron.NewAtTime(23, 0, 0)),
		),
		gocron.NewTask(createBackup),
	)
	if err != nil {
		panic(fmt.Errorf("failed to create cron: %w", err))
	}

	cronScheduler.Start()

	closer.Wait()
}
