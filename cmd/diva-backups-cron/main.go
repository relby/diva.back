package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	S3_REGION_ENV_NAME           = "S3_REGION"
	S3_ENDPOINT_ENV_NAME         = "S3_ENDPOINT"
	S3_TENANT_ID_ENV_NAME        = "S3_TENANT_ID"
	S3_KEY_ID_ENV_NAME           = "S3_KEY_ID"
	S3_KEY_SECRET_ENV_NAME       = "S3_KEY_SECRET"
	S3_BACKUPS_BUCKET_NAME       = "diva-backups"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("failed to read .env file: %v", err)
	}

	s3Region := os.Getenv(S3_REGION_ENV_NAME)
	if s3Region == "" {
		panic(fmt.Errorf("provide `%s`", S3_REGION_ENV_NAME))
	}

	s3Endpoint := os.Getenv(S3_ENDPOINT_ENV_NAME)
	if s3Endpoint == "" {
		panic(fmt.Errorf("provide `%s`", S3_ENDPOINT_ENV_NAME))
	}

	s3TenantID := os.Getenv(S3_TENANT_ID_ENV_NAME)
	if s3TenantID == "" {
		panic(fmt.Errorf("provide `%s`", S3_TENANT_ID_ENV_NAME))
	}

	s3KeyID := os.Getenv(S3_KEY_ID_ENV_NAME)
	if s3KeyID == "" {
		panic(fmt.Errorf("provide `%s`", S3_KEY_ID_ENV_NAME))
	}

	s3KeySecret := os.Getenv(S3_KEY_SECRET_ENV_NAME)
	if s3KeySecret == "" {
		panic(fmt.Errorf("provide `%s`", S3_KEY_SECRET_ENV_NAME))
	}

	diContainer, err := app.NewDIContainer()
	if err != nil {
		panic(fmt.Errorf("failed to create DI container: %w", err))
	}

	postgresConfig, err := diContainer.PostgresConfig()
	if err != nil {
		panic(fmt.Errorf("failed to get postgres config: %w", err))
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
	s3Uploader := s3manager.NewUploader(s3Session)

	os.MkdirAll(BACKUPS_DIR_PATH, os.ModePerm)

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
			gocron.NewWeekdays(time.Friday),
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
