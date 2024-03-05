package service

import (
	"context"
	"errors"
	"fmt"
	"subscriptionbot/db"
	"time"

	"github.com/phuslu/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Service) Notify(ctx context.Context) {
	for {
		timer := time.NewTimer(1 * time.Second)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			log.Error().Err(ctx.Err())
			return
		case <-timer.C:
			err := s.NotifySubscribers(ctx)
			if err != nil {
				log.Error().Err(err)
				continue
			}
			//If user updated his time timer restarts
			if errors.Is(err, fmt.Errorf("time not equal, perhaps user time updated")) {
				timer.Stop()
				timer = time.NewTimer(1 * time.Minute)

			}
		}
	}
}

func sendNextTime(currentTime, lastUpdatedAt time.Time) time.Time {
	nextTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), lastUpdatedAt.Hour(), lastUpdatedAt.Minute(), 0, 0, time.UTC)

	if nextTime.Before(currentTime) {
		nextTime = nextTime.Add(24 * time.Hour)
	}

	return nextTime
}

func (s *Service) NotifySubscribers(ctx context.Context) error {
	subscribers, userErr := s.DB.GetSubscribedUsers(ctx)
	if userErr != nil {
		log.Error().Err(userErr)
		return userErr
	}

	for _, sub := range subscribers {
		userTime, timeErr := time.Parse("15:04", sub.UserTime)
		if timeErr != nil {
			log.Error().Err(timeErr)
			continue
		}
		currentTime := time.Now().UTC()
		nextTrigger := sendNextTime(currentTime, userTime)

		if needtoSend(currentTime, nextTrigger, sub.ForecastSentAt) {
			subscriber, _ := s.DB.GetUser(sub.Username)
			currentUserTime, currentUserTimeError := time.Parse("15:04", subscriber.UserTime)
			if currentUserTimeError != nil {
				log.Error().Err(currentUserTimeError)
				return currentUserTimeError
			}
			//I think this line could be removed because we are using latest time anyway or leave it like this and just log that we are using latest time. Let me know
			if !currentUserTime.Equal(userTime) {
				log.Info().Msgf("User time was changed. Using the latest one")
			}
			if subscriber.SubscriptionStatus == int(db.LocationProvided) {
				forecast, _ := s.Weather.WeatherRequest(subscriber)
				s.API.SendResponse(subscriber.ChatID, forecast)
				s.DB.Update(bson.D{{"$set", bson.D{
					{"forecastSentAt", nextTrigger},
				}}}, subscriber.ID)
			}
		}
	}
	return nil
}

func needtoSend(currentTime, nextTrigger, userTime time.Time) bool {
	return currentTime.Equal(nextTrigger) || userTime.Before(nextTrigger)
}
