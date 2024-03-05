package utilities

import (
	"encoding/json"
	"fmt"
)

// LocationButton for telegram location
var LocationButton = ReplyKeyboardMarkup{
	Keyboard: [][]KeyboardButton{
		{KeyboardButton{Text: "Share my location", Location: true, OneTimeKeyboard: true, ResizeKeyboard: true}},
	},
}

// MenuButtons sends a menu for unsubscribed user
var MenuButtons = ReplyKeyboardMarkup{
	Keyboard: [][]KeyboardButton{
		{KeyboardButton{Text: Subscribe}},
		{KeyboardButton{Text: Unsubscribe}},
	},
}

// SubscribedMenu sends a menu for subscribed users
var SubscribedMenu = ReplyKeyboardMarkup{
	Keyboard: [][]KeyboardButton{
		{KeyboardButton{Text: "Share my location", Location: true, OneTimeKeyboard: true, ResizeKeyboard: true}},
		{KeyboardButton{Text: Unsubscribe}},
	},
}

// ButtonMarshal wraps a button into JSON
func ButtonMarshal(buttons ReplyKeyboardMarkup) ([]byte, error) {
	data, jsonErr := json.Marshal(buttons)
	if jsonErr != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", jsonErr)
	}
	return data, nil
}
