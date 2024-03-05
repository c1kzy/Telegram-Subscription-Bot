package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct for DB
type User struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty"`
	Username           string             `bson:"username"`
	SubscriptionStatus int                `bson:"subscriptionStatus"`
	UserTime           string             `bson:"userTime"`
	Location           Location           `bson:"location"`
	City               string             `bson:"city"`
	ChatID             int                `bson:"chatID"`
	ForecastSentAt     time.Time          `bson:"forecastSentAt"`
}

// Config struct for DB config
type Config struct {
	DB         string `env:"DB"`
	Password   string `env:"PASSWORD"`
	DBName     string `env:"DATABASE_NAME"`
	Access     string `env:"ACCESS"`
	Collection string `env:"COLLECTION"`
}

// Location struct for lat and lon
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
