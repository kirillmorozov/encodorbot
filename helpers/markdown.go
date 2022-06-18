package helpers

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

// EscapeMarkdownV2 escapes telegram markup symbols.
func EscapeMarkdownV2(text string, entityType telebot.EntityType) string {
	var escapeChars string
	switch entityType {
	case telebot.EntityCode, telebot.EntityCodeBlock:
		escapeChars = "\\`"
	case telebot.EntityTextLink:
		escapeChars = "\\)"
	default:
		escapeChars = "\\_*[]()~`>#+-=|{}.!"
	}
	re := regexp.MustCompilePOSIX("([{" + regexp.QuoteMeta(escapeChars) + "}])")
	return re.ReplaceAllString(text, "\\$1")
}
