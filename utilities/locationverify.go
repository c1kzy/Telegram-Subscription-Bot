package utilities

import api "git.foxminded.ua/foxstudent106270/telegramapi.git"

// IsLocationEmpty checks if location values received
func IsLocationEmpty(loc api.Location) bool {
	return loc.Latitude == 0.0 && loc.Longitude == 0.0
}
