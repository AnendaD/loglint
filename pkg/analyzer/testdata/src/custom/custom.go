package custom

import "log/slog"

func badCustomPatterns() {
	// AWS ключ прямо в строковом литерале
	slog.Info("connecting with key AKIAIOSFODNN7EXAMPLE") // want `AWS access key detected`

	// AWS ключ в конкатенации
	prefix := "key: "
	_ = prefix
	slog.Debug("access AKIAIOSFODNN7EXAMPLE granted") // want `AWS access key detected`
}

func goodCustomPatterns() {
	// Строки без паттернов — не триггерят
	slog.Info("connecting to database")
	slog.Debug("access granted")
	slog.Info("key not found")
}
