package specialchars

import "log/slog"

func badSpecialChars() {
	slog.Info("server started!")         // want `log message must not contain special characters or emoji`
	slog.Error("connection failed!!!")   // want `log message must not contain special characters or emoji`
	slog.Info("server started🚀")       // want `log message must not contain special characters or emoji`
	slog.Warn("something went wrong...") // want `log message must not contain special characters or emoji`
	slog.Debug("user login (failed)")    // want `log message must not contain special characters or emoji`
	slog.Info("status: ok?")            // want `log message must not contain special characters or emoji`
}

func goodSpecialChars() {
	// одиночные разрешённые символы — не триггерят
	slog.Info("server started")
	slog.Error("connection failed")
	slog.Warn("something went wrong")
	slog.Info("path: /usr/local/bin")
	slog.Info("key-value pair received")
	slog.Info("user_id not found")
	slog.Info("status: ok")
}
