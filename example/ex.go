package main

import (
	"log/slog"
	"os"

	"go.uber.org/zap"
)

func main() {
	// --- slog через переменную ---
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// ❌ Правило 1: заглавная буква
	logger.Info("starting server on port 8080")

	// ❌ Правило 2: не английский язык
	logger.Error("ошибка подключения к базе данных")

	// ❌ Правило 3: спецсимволы и эмодзи
	logger.Warn("server started ")

	// ❌ Правило 4: чувствительные данные в имени переменной
	password := "supersecret"
	logger.Info("user login", "pass", password)

	// --- zap через переменную ---
	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()

	// ❌ Правило 1: заглавная буква
	zapLogger.Info("failed to process request")

	// ❌ Правило 2: не английский язык
	zapLogger.Error("запуск воркера завершён")

	// ❌ Правило 3: эмодзи
	zapLogger.Warn("connection lost ")

	// ❌ Правило 4: токен в имени переменной
	userToken := "eyJhbGciOiJIUzI1NiJ9.payload.signature"
	zapLogger.Info("request received", zap.String("auth", userToken))
}
