package management

import (
	"errors"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrParseFailed       = errors.New("error during float parsing")
	ErrDecimalLength     = errors.New("more than 8 digits after decimal point")
	ErrLatitudeBoundary  = errors.New("latitude is not withing -90 90 degrees range")
	ErrLongitudeBoundary = errors.New("longitude is not withing -180 180 degrees range")
	ErrRadiusNegative    = errors.New("radius is not positive")
)

func ValidateLatitudeAndGet(latitude string) (float64, error) {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		log.Print(err)
		return 0, ErrParseFailed
	}
	err = ValidateDigitsOfCoordAfterPoint(latitude)
	if lat < -90 || lat > 90 {
		return 0, ErrLatitudeBoundary
	}
	return lat, err
}

func ValidateLongitudeAndGet(longitude string) (float64, error) {
	lon, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		log.Print(err)
		return 0, ErrParseFailed
	}
	err = ValidateDigitsOfCoordAfterPoint(longitude)
	if lon < -180 || lon > 180 {
		return 0, ErrLongitudeBoundary
	}
	return lon, err
}

func ValidateUserName(userName string) bool {
	userNameRegex := regexp.MustCompile(`^[a-zA-Z0-9]{4,16}$`)
	return userNameRegex.MatchString(userName)
}

func ValidateDigitsOfCoordAfterPoint(coord string) error {
	split := strings.Split(coord, ".")
	if len(split) == 2 && len(split[1]) > 8 {
		return ErrDecimalLength
	}
	return nil
}

func ValidateRadiusAndGet(radius string) (float64, error) {
	rad, err := strconv.ParseFloat(radius, 64)
	if err != nil {
		log.Print(err)
		return 0, ErrParseFailed
	}
	if rad < 0 {
		return 0, ErrRadiusNegative
	} else {
		return rad, nil
	}
}
