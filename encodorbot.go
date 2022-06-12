package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kirillmorozov/encodor/cmd"
	"gopkg.in/telebot.v3"
)

const (
	tokenEnvVar = "TELEGRAM_BOT_TOKEN"
)

func newBot() *telebot.Bot {
	settings := telebot.Settings{Token: os.Getenv(tokenEnvVar)}
	bot, botErr := telebot.NewBot(settings)
	if botErr != nil {
		log.Fatal(botErr)
	}
	bot.Handle("/start", handleStart)
	bot.Handle(telebot.OnText, handleText)
	return bot
}

func handleStart(c telebot.Context) error {
	encodorCmd := cmd.NewRoot()
	var output strings.Builder
	encodorCmd.SetArgs([]string{"--help"})
	encodorCmd.SetOut(&output)
	if encodeErr := encodorCmd.Execute(); encodeErr != nil {
		return encodeErr
	}
	return c.Reply(output.String())
}

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

// HandleTelegramWebHook sends a message back to the chat in encoded form
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	bot := newBot()
	var update telebot.Update
	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		log.Print(decodeErr)
		return
	}
	bot.ProcessUpdate(update)
}
