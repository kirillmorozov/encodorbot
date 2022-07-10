package helpers

import (
	"regexp"

	"gopkg.in/telebot.v3"
)

// EscapeMarkdownV2 escapes telegram markup symbols.
// Use entityType argument to specify the type of text that needs to be escaped.
// For the enitity types telebot.EntityCode, telebot.EntityCodeBlock and
// telebot.EntityTextLink only certain characters need to be escaped.
// See the official API documentation for details.
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
