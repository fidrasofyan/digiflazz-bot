package util

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var print = message.NewPrinter(language.Indonesian)

func Sprintf(key message.Reference, args ...any) string {
	return print.Sprintf(key, args...)
}
