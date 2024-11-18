package log

import (
	"log"
)

type TelegramLogger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}
