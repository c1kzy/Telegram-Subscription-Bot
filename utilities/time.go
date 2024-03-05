package utilities

import (
	"fmt"
	"time"
)

// ConvertTime returns a string with time based on layout
func ConvertTime(inputTime string) (string, error) {
	layout := "15:04"
	parsedTime, err := time.Parse(layout, inputTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return "", err
	}

	return parsedTime.Format(layout), nil

}
