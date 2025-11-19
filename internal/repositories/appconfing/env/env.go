package env

import (
	"os"
	"todo/internal/repositories/contracts"

	"github.com/joho/godotenv"
)

const (
	appURL                  = "APP_URL"
	appPort                 = "PORT"
	environmentVariableName = "STORAGE_TYPE"
	frontendTypeEnvName     = "FRONTEND"

	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbDatabase = "DB_DATABASE"
	dbUsername = "DB_USERNAME"
	dbPassword = "DB_PASSWORD"

	smtpHost     = "SMTP_HOST"
	smtpPort     = "SMTP_PORT"
	smtpUsername = "SMTP_USERNAME"
	smtpPassword = "SMTP_PASSWORD"
	emailFrom    = "EMAIL_FROM"
)

func New() contracts.AppConfig {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return &appEnv{}
}

type appEnv struct {
}

func (a *appEnv) AppPort() string {
	return ":" + os.Getenv(appPort)
}

// AppURL implements contracts.AppConfig.
func (a *appEnv) AppURL() string {
	return os.Getenv(appURL)
}

// EmailFrom implements contracts.AppConfig.
func (a *appEnv) EmailFrom() string {
	return os.Getenv(emailFrom)
}

// SmtpHost implements contracts.AppConfig.
func (a *appEnv) SmtpHost() string {
	return os.Getenv(smtpHost)
}

// SmtpPassword implements contracts.AppConfig.
func (a *appEnv) SmtpPassword() string {
	return os.Getenv(smtpPassword)
}

// SmtpPort implements contracts.AppConfig.
func (a *appEnv) SmtpPort() string {
	return os.Getenv(smtpPort)
}

// SmtpUsername implements contracts.AppConfig.
func (a *appEnv) SmtpUsername() string {
	return os.Getenv(smtpUsername)
}

// DBDatabase implements contracts.AppConfig.
func (a *appEnv) DBDatabase() string {
	return os.Getenv(dbDatabase)
}

// DBHost implements contracts.AppConfig.
func (a *appEnv) DBHost() string {
	return os.Getenv(dbHost)
}

// DBPassword implements contracts.AppConfig.
func (a *appEnv) DBPassword() string {
	return os.Getenv(dbPassword)
}

// DBPort implements contracts.AppConfig.
func (a *appEnv) DBPort() string {
	return os.Getenv(dbPort)
}

// DBUsername implements contracts.AppConfig.
func (a *appEnv) DBUsername() string {
	return os.Getenv(dbUsername)
}

// StorageType implements contracts.AppConfig.
func (a *appEnv) StorageType() string {
	return os.Getenv(environmentVariableName)
}

// FrontendType implements contracts.AppConfig.
func (a *appEnv) FrontendType() string {
	return os.Getenv(frontendTypeEnvName)
}
