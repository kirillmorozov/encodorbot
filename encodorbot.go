package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/kirillmorozov/encodor/cmd"
	"gopkg.in/telebot.v3"
)

const (
	tokenEnvVar = "TELEGRAM_BOT_TOKEN"
)

func newBot() (*telebot.Bot, error) {
	settings := telebot.Settings{Token: os.Getenv(tokenEnvVar), Verbose: true}
	bot, botErr := telebot.NewBot(settings)
	if botErr != nil {
		return nil, botErr
	}
	bot.Handle("/start", handleStart)
	bot.Handle("", handleText)
	return bot, nil
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
	bot, botErr := newBot()
	if botErr != nil {
		panic(botErr.Error())
	}
	var update telebot.Update
	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		panic(decodeErr.Error())
	}
	bot.ProcessUpdate(update)
}
