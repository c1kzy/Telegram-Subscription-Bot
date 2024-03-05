package service

import (
	"context"
	"subscriptionbot/db"
	"subscriptionbot/mocks"
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
				Username: "elon",
			},
		},
	}
	return reqBody
}

func TestTickUser(t *testing.T) {
	controller := gomock.NewController(t)
	storage := mocks.NewMongoStorage(controller)
	telegram := mocks.NewTelegramService(controller)
	weather := mocks.NewWeatherService(controller)

	tgService := NewService(storage, weather, telegram)

	currentTime := time.Now().UTC()
	missedTime := currentTime.Add(-3 * time.Hour)

	reqBody := requestBody(t, "user1")
	cityUpdated := bson.D{{"$set", bson.D{
		{"city", "Toronto"},
	}}}

	newSubscriber := db.User{
		Username:           reqBody.Message.Chat.Username,
		SubscriptionStatus: int(db.LocationProvided),
		UserTime:           time.Now().UTC().Round(1 * time.Second).Format("15:04"),
		ChatID:             reqBody.Message.Chat.ID,
		City:               "New York",
		ForecastSentAt:     currentTime,
	}

	tests := []struct {
		name       string
		setupMocks func(storage *mocks.MongoStorage, weather *mocks.WeatherService, telegram *mocks.TelegramService)
	}{
		{
			name: "subscribed users",
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetSubscribedUsers(gomock.Any()).Return([]db.User{
					{
						ID:                 primitive.ObjectID{1},
						Username:           "mopsle",
						SubscriptionStatus: 4,
						UserTime:           time.Now().UTC().Format("15:04"),
						Location:           db.Location{},
						City:               "",
						ForecastSentAt:     missedTime,
					},
				}, nil)
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "",
					ForecastSentAt:     missedTime,
				}, nil).AnyTimes()
				forecast := weather.EXPECT().WeatherRequest(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "New York",
					ForecastSentAt:     missedTime,
				}).AnyTimes()
				telegram.EXPECT().SendResponse(reqBody.Message.Chat.ID, forecast).AnyTimes()
			},
		},
		{
			name: "new subscriber",
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetSubscribedUsers(gomock.Any()).Return([]db.User{
					{
						ID:                 primitive.ObjectID{1},
						Username:           "mopsle",
						SubscriptionStatus: 4,
						UserTime:           time.Now().UTC().Format("15:04"),
						Location:           db.Location{},
						City:               "Berlin",
					},
				}, nil)
				storage.EXPECT().Insert(&newSubscriber).Return(nil).AnyTimes()
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{3},
					Username:           "elon",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "Berlin",
				}, nil).AnyTimes()
				forecast := weather.EXPECT().WeatherRequest(db.User{
					ID:                 primitive.ObjectID{3},
					Username:           "elon",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "Berlin",
				}).AnyTimes()
				telegram.EXPECT().SendResponse(reqBody.Message.Chat.ID, forecast).AnyTimes()
			},
		},
		{
			name: "city added",
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetSubscribedUsers(gomock.Any()).Return([]db.User{
					{
						ID:                 primitive.ObjectID{1},
						Username:           "mopsle",
						SubscriptionStatus: 4,
						UserTime:           time.Now().UTC().Format("15:04"),
						Location:           db.Location{},
						City:               "New York",
					},
				}, nil)
				storage.EXPECT().Update(cityUpdated, primitive.ObjectID{1}).AnyTimes()
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "Toronto",
				}, nil).AnyTimes()
				forecast := weather.EXPECT().WeatherRequest(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "Toronto",
				}).AnyTimes()
				telegram.EXPECT().SendResponse(reqBody.Message.Chat.ID, forecast).AnyTimes()
			},
		},
		{
			name: "missed and not missed forecast",
			setupMocks: func(
				storage *mocks.MongoStorage,
				weather *mocks.WeatherService,
				telegram *mocks.TelegramService,
			) {
				storage.EXPECT().GetSubscribedUsers(gomock.Any()).Return([]db.User{
					{
						ID:                 primitive.ObjectID{1},
						Username:           "mopsle",
						SubscriptionStatus: 4,
						UserTime:           time.Now().UTC().Format("15:04"),
						Location:           db.Location{},
						City:               "New York",
						ForecastSentAt:     missedTime,
					},
					{
						ID:                 primitive.ObjectID{2},
						Username:           "Maria",
						SubscriptionStatus: 4,
						UserTime:           time.Now().UTC().Format("15:04"),
						Location:           db.Location{},
						City:               "London",
					},
				}, nil)
				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "New York",
					ForecastSentAt:     missedTime,
				}, nil).AnyTimes()
				forecastUser := weather.EXPECT().WeatherRequest(db.User{
					ID:                 primitive.ObjectID{1},
					Username:           "mopsle",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "New York",
					ForecastSentAt:     missedTime,
				}).AnyTimes()
				telegram.EXPECT().SendResponse(reqBody.Message.Chat.ID, forecastUser).AnyTimes()

				storage.EXPECT().GetUser(reqBody.Message.Chat.Username).Return(db.User{
					ID:                 primitive.ObjectID{2},
					Username:           "Maria",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "London",
				}, nil).AnyTimes()
				forecast := weather.EXPECT().WeatherRequest(db.User{
					ID:                 primitive.ObjectID{2},
					Username:           "Maria",
					SubscriptionStatus: 4,
					UserTime:           time.Now().UTC().Format("15:04"),
					Location:           db.Location{},
					City:               "London",
				}).AnyTimes()
				telegram.EXPECT().SendResponse(reqBody.Message.Chat.ID, forecast).AnyTimes()
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.setupMocks(storage, weather, telegram)
			err := tgService.NotifySubscribers(context.Background())
			require.NoError(t, err)
		})
	}
}

func Test_needtoSend(t *testing.T) {
	timeNow := time.Now().UTC()

	type args struct {
		currentTime time.Time
		nextTrigger time.Time
		sentAt      time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "current time equal",
			args: args{
				currentTime: timeNow,
				nextTrigger: timeNow,
				sentAt:      timeNow,
			},
			want: true,
		},
		{
			name: "current time before",
			args: args{
				currentTime: timeNow,
				nextTrigger: timeNow.Add(1 * time.Hour),
				sentAt:      timeNow.Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "current time after",
			args: args{
				currentTime: timeNow.Add(1 * time.Hour),
				nextTrigger: timeNow,
				sentAt:      timeNow,
			},
			want: false,
		},
		{
			name: "ForecastSentAt before next trigger",
			args: args{
				currentTime: timeNow.Add(1 * time.Hour),
				nextTrigger: timeNow,
				sentAt:      timeNow.Add(-1 * time.Hour),
			},
			want: true,
		},
		{
			name: "ForecastSentAt after next trigger",
			args: args{
				currentTime: timeNow.Add(1 * time.Hour),
				nextTrigger: timeNow,
				sentAt:      timeNow.Add(1 * time.Hour),
			},
			want: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equalf(t, tc.want, needtoSend(tc.args.currentTime, tc.args.nextTrigger, tc.args.sentAt), "needtoSend(%v, %v, %v)", tc.args.currentTime, tc.args.nextTrigger, tc.args.sentAt)
		})
	}
}
