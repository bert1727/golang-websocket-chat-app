package logger

import (
	"os"
	"path/filepath"
	"strconv"

	middleware "github.com/gofiber/contrib/v3/zerolog"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLoggerWithConsoleWriter(app *fiber.App) {
	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: "24",
	}

	logger := zerolog.New(consoleWriter).
		With().
		Timestamp().
		Stack().
		Logger()
	log.Logger = zerolog.New(consoleWriter).With().
		Stack().
		Caller().
		Timestamp().
		Logger()

	// for beautiful caller
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	app.Use(middleware.New(middleware.Config{
		Logger: &logger,
	}))
}

func SetupLoggerWithFileConsoleWriter(app *fiber.App, fileName string) {
	logFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	consoleWriter := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    false,
		TimeFormat: "24",
	}

	multiLogger := zerolog.MultiLevelWriter(
		consoleWriter,
		logFile)

	logger := zerolog.New(multiLogger).
		With().
		Timestamp().
		Logger()

	log.Logger = logger

	app.Use(middleware.New(middleware.Config{
		Logger: &logger,
	}))
}
