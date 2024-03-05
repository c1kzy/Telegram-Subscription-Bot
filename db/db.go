package db

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/caarlos0/env/v10"
	"github.com/phuslu/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrNotFound = errors.New("user not found")

type Storage interface {
	Insert(user *User) error
	Update(user bson.D, id primitive.ObjectID) error
	Delete(id primitive.ObjectID) error
	GetUser(userName string) (User, error)
	GetSubscribedUsers(ctx context.Context) ([]User, error)
	UserSubscriptionStatus(id primitive.ObjectID) (int, error)
}

// DB struct for database name and Client
type DB struct {
	Database   string
	Collection string
	Client     *mongo.Client
}

var (
	lock     = sync.Mutex{}
	singleDB *DB
)

// SubscriptionStatus for user subscription level
type SubscriptionStatus int

// iota constants for subscription levels
const (
	NewUser SubscriptionStatus = iota + 1
	Subscribed
	TimeUpdated
	LocationProvided
)

// GetDB initializing DB config
func GetDB() *DB {
	cfg := &Config{}
	if singleDB == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleDB == nil {
			if err := env.Parse(cfg); err != nil {
				log.Error().Err(err)
			}
			dbConnect, connectErr := connectToDB(cfg)
			if connectErr != nil {
				log.Error().Err(connectErr)
			}
			singleDB = &DB{
				Client:     dbConnect,
				Database:   cfg.DBName,
				Collection: cfg.Collection,
			}
			log.Info().Msg("DB API created")
		}
	}
	return singleDB
}

func connectToDB(cfg *Config) (*mongo.Client, error) {
	url := fmt.Sprintf("mongodb+srv://c1kzy:%v@%v.gehij8o.mongodb.net/?retryWrites=true&w=majority", cfg.Password, cfg.DBName)

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(url).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	if err := client.Database(cfg.Access).RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		panic(err)
	}
	log.Info().Msg("Connected to MongoDB!")

	return client, nil
}

// Insert creates a new user in DB
func (db *DB) Insert(user *User) error {
	collection := db.Client.Database(db.Database).Collection(db.Collection)
	log.Debug().Msgf("Connected to database to insert: %v", db.Collection)
	_, insertErr := collection.InsertOne(context.TODO(), user)
	if insertErr != nil {
		return insertErr
	}

	return nil
}

// Update updates user's information in DB
func (db *DB) Update(user bson.D, id primitive.ObjectID) error {
	collection := db.Client.Database(db.Database).Collection(db.Collection)
	log.Debug().Msgf("Connected to database to update: %v", db.Collection)
	filter := bson.D{{"_id", id}}
	_, updateErr := collection.UpdateOne(context.TODO(), filter, user)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

// Delete deletes user from DB
func (db *DB) Delete(id primitive.ObjectID) error {
	collection := db.Client.Database(db.Database).Collection(db.Collection)
	log.Debug().Msgf("Connected to database to delete: %v", db.Collection)
	filter := bson.D{{"_id", id}}

	_, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	return nil
}

// GetUser returns single user from DB
func (db *DB) GetUser(userName string) (User, error) {
	var result User
	collection := db.Client.Database(db.Database).Collection(db.Collection)
	filter := bson.D{{"username", userName}}
	findErr := collection.FindOne(context.TODO(), filter).Decode(&result)
	if findErr != nil {
		return User{}, db.convertErr(findErr)
	}

	log.Debug().Msgf("User found in collection: %v", db.Collection)
	return result, nil
}

func (db *DB) convertErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNotFound
	}

	return err
}

// GetSubscribedUsers returns subscribed users from DB
func (db *DB) GetSubscribedUsers(ctx context.Context) ([]User, error) {
	var subscribers []User
	filter := bson.D{{"subscriptionStatus", 4}}
	collection := db.Client.Database(db.Database).Collection(db.Collection)
	cursor, findErr := collection.Find(ctx, filter)
	if findErr != nil {
		if errors.Is(findErr, mongo.ErrNoDocuments) {
			return []User{}, findErr
		}
		return []User{}, findErr
	}
	for cursor.Next(ctx) {
		var user User
		err := cursor.Decode(&user)
		if err != nil {
			log.Error().Msgf("unable to decode users at :%v. Error:%v", cursor.Current, err)
			return []User{}, err
		}
		subscribers = append(subscribers, user)
	}
	log.Debug().Msgf("Users found in collection: %v", db.Collection)

	return subscribers, nil
}

// UserSubscriptionStatus returns user's status of subscription
func (db *DB) UserSubscriptionStatus(id primitive.ObjectID) (int, error) {
	var result User
	collection := db.Client.Database(db.Database).Collection(db.Collection)
	filter := bson.D{{"_id", id}}
	findErr := collection.FindOne(context.TODO(), filter).Decode(&result)
	if findErr != nil {
		return 0, findErr
	}

	log.Debug().Msgf("User subscription status is : %v", result.SubscriptionStatus)

	return result.SubscriptionStatus, nil
}
