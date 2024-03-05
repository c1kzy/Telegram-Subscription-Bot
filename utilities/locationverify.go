package utilities

import api "github.com/c1kzy/Telegram-API"

// IsLocationEmpty checks if location values received
func IsLocationEmpty(loc api.Location) bool {
	return loc.Latitude == 0.0 && loc.Longitude == 0.0
}
