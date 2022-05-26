package encodorbot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/kirillmorozov/encodor/cmd"
	"go.uber.org/zap"
)

var logger *zap.Logger

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Update is a Telegram object that the handler receives every time an user
// interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

const startCommand = "/start"

func init() {
	logger, _ = zap.NewProduction()
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if decodeErr := json.NewDecoder(r.Body).Decode(&update); decodeErr != nil {
		logger.Error(
			"Could not decode incoming update",
			zap.String("severity", "ERROR"),
			zap.Error(decodeErr),
		)
		return nil, decodeErr
	}
	return &update, nil
}

// HandleTelegramWebHook sends a message back to the chat in encoded form
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	defer logger.Sync() //nolint:errcheck
	// Parse incoming request
	update, parseErr := parseTelegramRequest(r)
	if parseErr != nil {
		logger.Error(
			"Error parsing update",
			zap.String("severity", "ERROR"),
			zap.Error(parseErr),
		)
		return
	}
	logger.Info(
		"New message received",
		zap.Int("update_id", update.UpdateId),
		zap.String("text", update.Message.Text),
		zap.Int("chat_id", update.Message.Chat.Id),
		zap.String("severity", "DEBUG"),
	)
	if update.Message.Text == startCommand {
		update.Message.Text = "--help"
	}
	var output strings.Builder
	if ecnodeErr := cmd.ExecuteCustomIO(strings.Fields(update.Message.Text), &output); ecnodeErr != nil {
		logger.Error(
			"Error encodor command execution",
			zap.String("args", update.Message.Text),
			zap.String("severity", "WARNING"),
			zap.Error(ecnodeErr),
		)
	}
	telegramResponseBody, errTelegram := sendTextToTelegramChat(update.Message.Chat.Id, output.String())
	if errTelegram != nil {
		logger.Error(
			"Error from Telegram",
			zap.Int("chat_id", update.Message.Chat.Id),
			zap.String("text", output.String()),
			zap.String("response", telegramResponseBody),
			zap.String("severity", "ERROR"),
			zap.Error(errTelegram),
		)
		return
	}
	logger.Info(
		"Successfully sent encoded message",
		zap.Int("chat_id", update.Message.Chat.Id),
		zap.String("args", update.Message.Text),
		zap.String("text", output.String()),
		zap.String("severity", "INFO"),
	)
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified
// by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {
	var telegramApi string = "https://api.telegram.org/bot" + os.Getenv("TELEGRAM_BOT_TOKEN") + "/sendMessage"
	response, postErr := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		},
	)
	if postErr != nil {
		logger.Error(
			"Error when posting text to the chat",
			zap.String("severity", "ERROR"),
			zap.Error(postErr),
		)
		return "", postErr
	}
	defer response.Body.Close()
	bodyBytes, errRead := ioutil.ReadAll(response.Body)
	if errRead != nil {
		logger.Error(
			"Error when parsing telegram answer",
			zap.String("severity", "ERROR"),
			zap.Error(errRead),
		)
		return "", errRead
	}
	bodyString := string(bodyBytes)
	logger.Info(
		"Telegram response body",
		zap.String("response", bodyString),
		zap.String("severity", "DEBUG"),
	)
	return bodyString, nil
}
