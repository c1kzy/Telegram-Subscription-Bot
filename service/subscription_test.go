package service_test

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"subscriptionbot/db"
	"subscriptionbot/mocks"
	"subscriptionbot/service"
	"subscriptionbot/utilities"
	"testing"
	"time"

	api "git.foxminded.ua/foxstudent106270/telegramapi.git"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func requestBody(t *testing.T, text string) *api.WebHookReqBody {
	reqBody := &api.WebHookReqBody{
		Message: api.Message{
			Text: text,
			Chat: api.Chat{
				ID:       358383178,
				Username: "mopsle",
			},
		},
	}
	return reqBody
}
func ButtonMarshal(buttons utilities.ReplyKeyboardMarkup) ([]byte, error) {
	data, jsonErr := json.Marshal(buttons)
	if jsonErr != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", jsonErr)
	}
	return data, nil
}

func TestService_Subscription(t *testing.T) {
	controller := gomock.NewController(t)
	storage := mocks.NewMongoStorage(controller)
	weatherService := mocks.NewWeatherService(controller)
	telegramService := mocks.NewTelegramService(controller)

	tgService := service.NewService(storage, weatherService, telegramService)

	reqBody := requestBody(t, "user1")

	newUser := db.User{
		Username:           reqBody.Message.Chat.Username,
		SubscriptionStatus: int(db.NewUser),
		UserTime:           time.Now().UTC().Round(1 * time.Second).Format("15:04"),
		ChatID:             reqBody.Message.Chat.ID,
	}
	updateUser := bson.D{{"$set", bson.D{
		{"subscriptionStatus", db.TimeUpdated},
		{"time", time.Now().UTC().Round(1 * time.Second).Format("15:04")},
	}}}
	updateUserLocation := bson.D{{"$set", bson.D{
		{"subscriptionStatus", db.LocationProvided},
		{"city", "New York"},
	}}}

	jsonData, jsonErr := ButtonMarshal(utilities.MenuButtons)
	require.NoError(t, jsonErr)
	locationData, locationErr := ButtonMarshal(utilities.LocationButton)
	require.NoError(t, locationErr)

	tests := []struct {
		name          string
		text          string
		want          url.Values
		expectedError error
		setupMocks    func(storage *mocks.MongoStorage, weather *mocks.WeatherService, telegram *mocks.TelegramService)
	}{
		{
			name: "User created",
			text: "user1",
			want: url.Values{
				"chat_id":      {strconv.Itoa(358383178)},
				"text":         {"Welcome to weather forecast bot! Subscribe and Unsubscribe options available below"},
				"reply_markup": {string(jsonData)},
			},
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{}, db.ErrNotFound)
				storage.EXPECT().Insert(&newUser).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "User time updated",
			text: time.Now().UTC().Format("15:04"),
			want: url.Values{
				"chat_id":      {strconv.Itoa(358383178)},
				"text":         {"User time updated. Please enter city or share location to update the city for weather forecast"},
				"reply_markup": {string(locationData)},
			},
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 1,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "",
				}, nil)
				storage.EXPECT().UserSubscriptionStatus(primitive.ObjectID{1}).Return(int(db.Subscribed), nil)
				storage.EXPECT().Update(updateUser, primitive.ObjectID{1})
			},
			expectedError: nil,
		},
		{
			name: "User city updated",
			text: "New York",
			want: url.Values{
				"chat_id": {strconv.Itoa(358383178)},
				"text":    {"City updated"},
			},
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 2,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "New York",
				}, nil)
				storage.EXPECT().UserSubscriptionStatus(primitive.ObjectID{1}).Return(int(db.TimeUpdated), nil)
				weather.EXPECT().WeatherRequest(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 2,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "New York",
				}).AnyTimes()
				storage.EXPECT().Update(updateUserLocation, primitive.ObjectID{1})
			},
			expectedError: nil,
		},
		{
			name: "User unsubscribe",
			text: "Unsubscribe",
			want: url.Values{
				"chat_id": {strconv.Itoa(358383178)},
				"text":    {"You have unsubscribed from weather forecast"},
			},
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "New York",
				}, nil)
				storage.EXPECT().UserSubscriptionStatus(primitive.ObjectID{1}).Return(int(db.LocationProvided), nil)
				storage.EXPECT().Delete(primitive.ObjectID{1})
			},
			expectedError: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks(storage, weatherService, telegramService)
			got, err := tgService.AddSubscription(requestBody(t, tc.text), 358383178)
			assert.ErrorIs(t, err, tc.expectedError)
			assert.Equal(t, got, tc.want)
		})

	}
}
