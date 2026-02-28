package sensitive

import "log/slog"

// Правило срабатывает на имена переменных и функций.

func badSensitive() {
	password := "hunter2"
	slog.Info("user login attempt", password) // want `log message contains sensitive keyword: "password"`

	tokenValue := "eyJhbG..."
	slog.Debug("auth completed", tokenValue) // want `log message contains sensitive keyword: "token"`

	userSecret := "shhh"
	slog.Warn("config loaded", userSecret) // want `log message contains sensitive keyword: "secret"`

	// keyword в имени ключа структурированного лога
	apiKey := "key-abc-123"
	slog.Info("request sent", "api_key", apiKey) // want `log message contains sensitive keyword: "apiKey"`
}

func goodSensitive() {
	// Переменные с нейтральными именами — не триггерят
	value := "some-value"
	slog.Info("user authenticated successfully", value)

	code := 200
	slog.Debug("response received", code)

	slog.Info("token validated")
	slog.Info("api request completed")
}
