package weatherAPI

// WeatherConfig struct for weather config
type WeatherConfig struct {
	GeoAPI         GeoAPI
	API            string `env:"API"`
	WeatherVersion string `env:"WEATHER_VERSION"`
	Units          string `env:"UNITS"`
}

// GeoAPI for API that returns lat, lon for city provided by user
type GeoAPI struct {
	Version string `env:"GEO_API_VERSION"`
	Limit   int    `env:"GEO_API_LIMIT"`
}

// Coord struct for lat, lon
type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

// Weather struct for weather description
type Weather struct {
	ID          int    `json:"id"`
	Forecast    string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// Main struct for weather details
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
	SeaLevel  int     `json:"sea_level"`
	GrndLevel int     `json:"grnd_level"`
}

// Wind struct for wind details
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
	Gust  float64 `json:"gust"`
}

// Rain struct for rain details
type Rain struct {
	OneHour float64 `json:"1h"`
}

// Clouds struct for cloud details
type Clouds struct {
	All int `json:"all"`
}

// Sys for sunrise and sunset
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

// WeatherData struct combines all structs in one
type WeatherData struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Rain       Rain      `json:"rain"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

// LocalNames struct for local names
type LocalNames struct {
	En string `json:"en"`
	Uk string `json:"uk"`
}

// Location struct for location details
type Location struct {
	Name       string     `json:"name"`
	LocalNames LocalNames `json:"local_names"`
	Lat        float64    `json:"lat"`
	Lon        float64    `json:"lon"`
	Country    string     `json:"country"`
	State      string     `json:"state"`
}
