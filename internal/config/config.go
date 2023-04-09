package config

import (
	"os"
	"time"
)

const TimeOutImageDownloading = 2 * time.Minute

type configName string

const (
	EmailAccount = configName("email_account")
	EmailPassword = configName("email_password")
	DbDsn = configName("db_dsn")
	RetryDuration = configName("http_retry_duration")
)

func GetValue(cnfg configName) string {
	return os.Getenv(string(cnfg))
}