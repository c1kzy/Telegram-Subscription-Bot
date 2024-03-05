//go:generate go run github.com/golang/mock/mockgen -destination=storage.go -package=mocks -mock_names=Storage=MongoStorage subscriptionbot/db Storage
//go:generate go run github.com/golang/mock/mockgen -destination=weather.go -package=mocks -mock_names=WeatherService=WeatherService subscriptionbot/weather WeatherService
//go:generate go run github.com/golang/mock/mockgen -destination=api.go -package=mocks -mock_names=TelegramService=TelegramService git.foxminded.ua/foxstudent106270/telegramapi.git TelegramService

package mocks
