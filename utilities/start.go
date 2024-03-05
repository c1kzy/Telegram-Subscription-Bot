package utilities

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"

	api "git.foxminded.ua/foxstudent106270/telegramapi.git"
)

// KeyboardButton struct for a telegram button
type KeyboardButton struct {
	Text            string `json:"text"`
	Location        bool   `json:"request_location"`
	OneTimeKeyboard bool   `json:"one_time_keyboard"`
	ResizeKeyboard  bool   `json:"resize_keyboard"`
}

// ReplyKeyboardMarkup struct for keyboard layout
type ReplyKeyboardMarkup struct {
	Keyboard [][]KeyboardButton `json:"keyboard"`
}

// StartResponse for /start command
func StartResponse(body *api.WebHookReqBody, chatID int) (url.Values, error) {
	replyMarkup := ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{KeyboardButton{Text: Subscribe}},
			{KeyboardButton{Text: Unsubscribe}},
		},
	}
	jsonData, jsonErr := json.Marshal(replyMarkup)
	if jsonErr != nil {
		return url.Values{}, fmt.Errorf("error marshaling JSON: %w", jsonErr)
	}

	return url.Values{
		"chat_id":      {strconv.Itoa(chatID)},
		"text":         {Start},
		"reply_markup": {string(jsonData)},
	}, nil
}
