package lowercase

import "log/slog"

func badLowercase() {
	slog.Info("Starting server on port 8080")   // want `Log must start with a lowercase letter`
	slog.Error("Failed to connect to database") // want `Log must start with a lowercase letter`
	slog.Warn("Warning: disk space low")        // want `Log must start with a lowercase letter`
	slog.Debug("Received request from client")  // want `Log must start with a lowercase letter`
}

func goodLowercase() {
	slog.Info("starting server on port 8080")
	slog.Error("failed to connect to database")
	slog.Warn("disk space low")
	slog.Debug("received request from client")
}
