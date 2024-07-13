package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/relby/diva.back/internal/app"
	"github.com/relby/diva.back/internal/closer"

	"github.com/go-co-op/gocron/v2"
)

var (
	BACKUPS_DIR_PATH             = filepath.Join(".", "backups")
	BACKUP_FILE_NAME_TIME_FORMAT = "20060102150405"
)

func main() {
	diContainer, err := app.NewDIContainer()
	if err != nil {
		panic(fmt.Errorf("failed to create DI container: %w", err))
	}

	postgresConfig, err := diContainer.PostgresConfig()
	if err != nil {
		panic(fmt.Errorf("failed to get postgres config: %w", err))
	}

	os.MkdirAll(BACKUPS_DIR_PATH, os.ModePerm)

	createBackup := func() {
		backupFileName := fmt.Sprintf("%s.sql", time.Now().Format(BACKUP_FILE_NAME_TIME_FORMAT))
		backupFilePath := filepath.Join(BACKUPS_DIR_PATH, backupFileName)
		cmd := exec.Command(
			"pg_dump",
			"-U", postgresConfig.User(),
			"-h", postgresConfig.Host(),
			"-p", postgresConfig.Port(),
			"-F", "c",
			postgresConfig.DB(),
			">", backupFilePath,
		)

		cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", postgresConfig.DB()))

		cmdOutput, err := cmd.CombinedOutput()
		fmt.Println("pg_dump output: ", string(cmdOutput))
		if err != nil {
			fmt.Printf("failed to execute pg_dump: %v\n", err)
			return
		}

		fmt.Printf("successfully created backup `%s`\n", backupFileName)
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
