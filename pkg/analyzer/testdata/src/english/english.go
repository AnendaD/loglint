package english

import "log/slog"

func badEnglish() {
	slog.Info("запуск сервера")                // want `log message must be in English \(found: .\)`
	slog.Error("ошибка подключения к базе")    // want `log message must be in English \(found: .\)`
	slog.Warn("внимание: мало места на диске") // want `log message must be in English \(found: .\)`
	slog.Debug("получен запрос")               // want `log message must be in English \(found: .\)`
}

func goodEnglish() {
	slog.Info("starting server")
	slog.Error("failed to connect to database")
	slog.Warn("disk space low")
	slog.Debug("received request")
}
