package helpers

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

// escapeMarkdownV2 escapes telegram markup symbols.
func escapeMarkdownV2(text string, entityType telebot.EntityType) string {
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
