package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"subscriptionbot/db"
	"subscriptionbot/utilities"
	weatherAPI "subscriptionbot/weather"
	"time"
	"unicode"

	api "github.com/c1kzy/Telegram-API"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SubscriptionService interface {
	AddSubscription(body *api.WebHookReqBody, chatID int) (url.Values, error)
	Notify(ctx context.Context) error
	TickUser(ctx context.Context) error
}

// Service struct for DB, weather, api
type Service struct {
	DB      db.Storage
	Weather weatherAPI.WeatherService
	API     api.TelegramService
}

func NewService(DB db.Storage, weather weatherAPI.WeatherService, API api.TelegramService) *Service {
	return &Service{DB: DB, Weather: weather, API: API}
}

// AddSubscription function handles user subscriptions
func (s *Service) AddSubscription(body *api.WebHookReqBody, chatID int) (url.Values, error) {
	currentTime := time.Now().UTC()
	var (
		user    db.User
		userErr error
	)
	jsonData, jsonErr := utilities.ButtonMarshal(utilities.MenuButtons)
	if jsonErr != nil {
		return url.Values{}, fmt.Errorf("error marshaling JSON: %w", jsonErr)
	}

	user, userErr = s.DB.GetUser(body.Message.Chat.Username)
	if errors.Is(userErr, db.ErrNotFound) {
		newUser := db.User{
			Username:           body.Message.Chat.Username,
			SubscriptionStatus: int(db.NewUser),
			UserTime:           currentTime.Round(1 * time.Second).Format("15:04"),
			ChatID:             body.Message.Chat.ID,
		}
		insertErr := s.DB.Insert(&newUser)
		if insertErr != nil {
			return nil, insertErr
		}
		return url.Values{
			"chat_id":      {strconv.Itoa(chatID)},
			"text":         {"Welcome to weather forecast bot! Subscribe and Unsubscribe options available below"},
			"reply_markup": {string(jsonData)},
		}, nil
	}

	userSubscriptionStatus, statusErr := s.DB.UserSubscriptionStatus(user.ID)
	if statusErr != nil {
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"unable to retrieve user subscription status"},
		}, nil
	}

	switch userSubscriptionStatus {
	case int(db.NewUser):
		return s.userSubscribe(body, chatID, user.ID)
	case int(db.Subscribed):
		return s.userTimeRequest(body, chatID, user.ID)
	case int(db.TimeUpdated):
		return s.userLocationRequest(body, chatID, user)
	default:
		return s.answerHandle(body, chatID, user)

	}
}

func (s *Service) userTimeRequest(body *api.WebHookReqBody, chatID int, id primitive.ObjectID) (url.Values, error) {
	jsonData, jsonErr := utilities.ButtonMarshal(utilities.LocationButton)
	if jsonErr != nil {
		return url.Values{}, fmt.Errorf("error marshaling JSON: %w", jsonErr)
	}
	//Unsubscribe option in case user decided to unsubscribe
	if body.Message.Text == utilities.Unsubscribe {
		if _, err := s.userUnsubscribe(id, chatID); err != nil {
			return url.Values{}, err
		}

	}
	userTime, timeErr := utilities.ConvertTime(body.Message.Text)
	if timeErr != nil {
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"invalid time, try again.Example: 12:00"},
		}, timeErr
	}
	update := bson.D{{"$set", bson.D{
		{"subscriptionStatus", db.TimeUpdated},
		{"userTime", userTime},
	}}}

	updateErr := s.DB.Update(update, id)
	if updateErr != nil {
		return nil, updateErr
	}
	return url.Values{
		"chat_id":      {strconv.Itoa(chatID)},
		"text":         {"User time updated. Please enter city or share location to update the city for weather forecast"},
		"reply_markup": {string(jsonData)},
	}, nil

}

func (s *Service) userSubscribe(body *api.WebHookReqBody, chatID int, id primitive.ObjectID) (url.Values, error) {
	if body.Message.Text == utilities.Subscribe {
		currentTime := fmt.Sprintf("%02d:%02d", time.Now().Hour(), time.Now().Minute())
		update := bson.D{{"$set", bson.D{
			{"subscriptionStatus", db.Subscribed},
			{"userTime", currentTime},
		}}}
		err := s.DB.Update(update, id)
		if err != nil {
			return nil, err
		}
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"You have subscribed to weather forecast! Please enter time to provide time in 24H format for weather forecast every day.Example: /time 15:00.\nTime when subscribed is used by default"},
		}, nil
	}

	return url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"text":    {"Please subscribe to continue"},
	}, nil
}

func (s *Service) userUnsubscribe(id primitive.ObjectID, chatID int) (url.Values, error) {
	deleteErr := s.DB.Delete(id)
	if deleteErr != nil {
		return nil, deleteErr
	}
	return url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"text":    {"You have unsubscribed from weather forecast"},
	}, nil
}

func (s *Service) userLocationRequest(body *api.WebHookReqBody, chatID int, user db.User) (url.Values, error) {
	jsonData, jsonErr := utilities.ButtonMarshal(utilities.LocationButton)
	if jsonErr != nil {
		return url.Values{}, fmt.Errorf("error marshaling JSON: %w", jsonErr)
	}
	//Left unsubscribe option in case user decided to unsubscribe at this point
	if body.Message.Text == utilities.Unsubscribe {
		if _, err := s.userUnsubscribe(user.ID, chatID); err != nil {
			return url.Values{}, err
		}

	}

	if !utilities.IsLocationEmpty(body.Message.Location) {
		update := bson.D{{"$set", bson.D{
			{"location", body.Message.Location},
			{"subscriptionStatus", db.LocationProvided},
		}}}

		err := s.DB.Update(update, user.ID)
		if err != nil {
			return nil, err
		}
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"Location updated"},
		}, nil

	}

	if body.Message.Text != "" {
		update := bson.D{{"$set", bson.D{
			{"subscriptionStatus", db.LocationProvided},
			{"city", body.Message.Text},
		}}}

		updateErr := s.DB.Update(update, user.ID)
		if updateErr != nil {
			return nil, updateErr
		}
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"City updated"},
		}, nil
	}

	//Checking if user location can be used in weather request
	_, weatherError := s.Weather.WeatherRequest(user)
	if weatherError != nil {
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"invalid location provided"},
		}, weatherError
	}

	return url.Values{
		"chat_id":      {strconv.Itoa(chatID)},
		"text":         {"please enter city or share location to continue\nExample: New York"},
		"reply_markup": {string(jsonData)},
	}, nil
}

func (s *Service) answerHandle(body *api.WebHookReqBody, chatID int, user db.User) (url.Values, error) {
	subscribedButtons, _ := utilities.ButtonMarshal(utilities.SubscribedMenu)

	if body.Message.Text == utilities.Unsubscribe {
		return s.userUnsubscribe(user.ID, chatID)
	}

	if unicode.IsDigit(rune(body.Message.Text[0])) {
		return s.timeUpdate(body.Message.Text, user.ID, chatID)
	}

	if unicode.IsLetter(rune(body.Message.Text[0])) || !utilities.IsLocationEmpty(body.Message.Location) {
		return s.locationUpdate(body.Message.Text, body, user, chatID)
	}

	return url.Values{
		"chat_id":      {strconv.Itoa(chatID)},
		"text":         {utilities.SubscribedOptions},
		"reply_markup": {string(subscribedButtons)},
	}, nil
}

func (s *Service) timeUpdate(time string, id primitive.ObjectID, chatID int) (url.Values, error) {
	userTime, timeErr := utilities.ConvertTime(time)
	if timeErr != nil {
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"invalid time, try again.Example: 12:00"},
		}, timeErr
	}
	update := bson.D{{"$set", bson.D{
		{"time", userTime},
	}}}

	updateErr := s.DB.Update(update, id)
	if updateErr != nil {
		return nil, updateErr
	}
	return url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"text":    {"User time updated"},
	}, nil
}

func (s *Service) locationUpdate(city string, body *api.WebHookReqBody, user db.User, chatID int) (url.Values, error) {
	//Checking if user location can be used in weather request
	_, weatherError := s.Weather.WeatherRequest(user)
	if weatherError != nil {
		return url.Values{
			"chat_id": {strconv.Itoa(chatID)},
			"text":    {"invalid weather input provided"},
		}, weatherError
	}

	if !utilities.IsLocationEmpty(body.Message.Location) {
		update := bson.D{{"$set", bson.D{
			{"location", body.Message.Location},
		}}}

		err := s.DB.Update(update, user.ID)
		if err != nil {
			return nil, err
		}

	}
	update := bson.D{{"$set", bson.D{
		{"city", city},
	}}}

	updateErr := s.DB.Update(update, user.ID)
	if updateErr != nil {
		return nil, updateErr
	}

	return url.Values{
		"chat_id": {strconv.Itoa(chatID)},
		"text":    {"Location status updated!"},
	}, nil
}
