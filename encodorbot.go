package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kirillmorozov/encodor/beghilosz"
	"github.com/kirillmorozov/encodor/cmd"
	"github.com/kirillmorozov/encodor/zalgo"
	"gopkg.in/telebot.v3"
)

const (
	// tokenEnvVar is an environment variable that contains telegram bot token.
	tokenEnvVar = "TELEGRAM_BOT_TOKEN"
)

const (
	// startCommand is the first command that's send from the new user.
	startCommand = "/start"
	// beghiloszCommand is the command that's used to encode messages using
	// calculator spelling.
	beghiloszCommand = "/beghilosz"
	// zalgoCommand is the command that's used to encode messages zalgo.
	zalgoCommand = "/zalgo"
)

const (
	// botUsage is format that explains bot usage.
	botUsage = `Usage:
	[command] your message

Available Commands:
	%v – Encode your message using calculator spelling
	%v – Encode your message using zalgo`
	// beghiloszUsage is a format string that expains beghiloszCommand usage.
	beghiloszUsage = "To encode your message using calculator spelling send `%v YOUR MESSAGE`"
	// zalgoUsage is a format string that expains zalgoCommand usage.
	zalgoUsage = "To encode your message using zalgo send `%v YOUR MESSAGE`"
)

// newBot returns a configured telegram bot.
func newBot() *telebot.Bot {
	settings := telebot.Settings{
		Token:     os.Getenv(tokenEnvVar),
		ParseMode: telebot.ModeMarkdownV2,
	}
	bot, botErr := telebot.NewBot(settings)
	if botErr != nil {
		log.Fatal(botErr)
	}
	bot.Handle(startCommand, handleStart)
	bot.Handle(beghiloszCommand, handleBeghilosz)
	bot.Handle(zalgoCommand, handleZalgo)
	bot.Handle(telebot.OnText, handleText)
	return bot
}

// handleStart handles startCommand messages.
func handleStart(c telebot.Context) error {
	return c.Reply(fmt.Sprintf(botUsage, beghiloszCommand, zalgoCommand))
}

// handleBeghilosz handles beghiloszCommand messages.
func handleBeghilosz(c telebot.Context) error {
	if c.Message().Payload == "" {
		usage := fmt.Sprintf(beghiloszUsage, beghiloszCommand)
		return c.Reply(usage)
	}
	encodedText := beghilosz.Encode(c.Message().Payload)
	return c.Reply(encodedText)
}

// handleZalgo handles zalgoCommand messages.
func handleZalgo(c telebot.Context) error {
	if c.Message().Payload == "" {
		usage := fmt.Sprintf(zalgoUsage, zalgoCommand)
		return c.Reply(usage)
	}
	encodedText, encodeErr := zalgo.Encode(c.Message().Payload, 3)
	if encodeErr != nil {
		return encodeErr
	}
	return c.Reply(encodedText)
}

// handleText handles all plain text messages.
func handleText(c telebot.Context) error {
	encodorCmd := cmd.NewRoot()
	var output strings.Builder
	encodorCmd.SetArgs(strings.Fields(c.Message().Text))
	encodorCmd.SetOut(&output)
	if encodeErr := encodorCmd.Execute(); encodeErr != nil {
		return encodeErr
	}
	return c.Reply(output.String())
}

// HandleTelegramWebHook is the cloud function entry point.
//
// It decodes a telebot.Update from the http.Request, creates a pre-configured
// telebot.Bot and processes the update through this bot.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	var update telebot.Update
	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		log.Print(decodeErr)
		return
	}
	bot := newBot()
	bot.ProcessUpdate(update)
}
