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
	tokenEnvVar = "TELEGRAM_BOT_TOKEN"
)

const (
	startCommand     = "/start"
	beghiloszCommand = "/beghilosz"
	zalgoCommand     = "/zalgo"
)

const (
	botUsage       = `Henlo`
	beghiloszUsage = `To encode your message using calculator spelling send %v YOUR MESSAGE`
	zalgoUsage     = `To encode your message using zalgo send %v YOUR MESSAGE`
)

func newBot() *telebot.Bot {
	settings := telebot.Settings{Token: os.Getenv(tokenEnvVar)}
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

func handleStart(c telebot.Context) error {
	return c.Reply(botUsage)
}

func handleBeghilosz(c telebot.Context) error {
	if c.Message().Text == "" {
		return c.Reply(fmt.Sprintf(beghiloszUsage, beghiloszCommand))
	}
	encodedText := beghilosz.Encode(c.Message().Text)
	return c.Reply(encodedText)
}

func handleZalgo(c telebot.Context) error {
	if c.Message().Text == "" {
		return c.Reply(fmt.Sprintf(zalgoUsage, zalgoCommand))
	}
	encodedText, encodeErr := zalgo.Encode(c.Message().Text, 3)
	if encodeErr != nil {
		return encodeErr
	}
	return c.Reply(encodedText)
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

func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	bot := newBot()
	var update telebot.Update
	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		log.Print(decodeErr)
		return
	}
	bot.ProcessUpdate(update)
}
