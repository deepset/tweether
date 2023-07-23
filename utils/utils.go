package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	logger "github.com/deepset/tweether/logger"
)

var (
	ErrEmptyString    = errors.New("empty string")
	ErrInvalidTime    = errors.New("invalid time")
	ErrInvalidNumber  = errors.New("invalid number")
	ErrInvalidBoolean = errors.New("invalid Boolean")
)
var digitCheck = regexp.MustCompile(`^[+-]?([0-9]*[.])?[0-9]+$`)

func ParseTime(s string) (int64, error) {

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		logger.ErrorLogger.Printf("string %q is not in RFC3339 format ", s)
		return 0, ErrInvalidTime
	}
	// Unix returns t as a Unix time,
	return t.Unix(), nil

}

// Parsing the map key
func ParseKey(s string) (string, error) {

	newString := strings.TrimSpace(s)
	if len(newString) == 0 {
		return "", ErrEmptyString
	}
	return newString, nil
}

// Parsing the string type
func ParseString(s string) (interface{}, error) {

	//Remove empty space from front and back of string
	newString := strings.TrimSpace(s)

	if len(newString) == 0 {
		return "", ErrEmptyString
	}

	//check if string is in RFC3339 format
	if t, err := ParseTime(newString); err == nil {
		return t, nil
	}
	return newString, nil
}

// Parsing the number type
func ParseNumber(s string) (interface{}, error) {

	//Remove empty space from front and back of string
	newString := strings.TrimSpace(s)
	// validate if string is a number
	if !digitCheck.MatchString(newString) {
		return 0, ErrInvalidNumber
	}

	//Number conversion to integer
	if num, err := strconv.Atoi(newString); err == nil {
		return int64(num), nil
	}

	//Number conversion to float
	if num, err := strconv.ParseFloat(newString, 64); err == nil {
		return num, nil
	}
	return -1, nil

}

// Parsing the Boolean type
func ParseBoolean(s string) (interface{}, error) {

	//Remove empty space from front and back of string
	newString := strings.TrimSpace(s)

	b, err := strconv.ParseBool(newString)
	if err != nil {
		return -1, ErrInvalidBoolean
	}

	return b, nil

}

// Parsing the List type
func ParseList(l []interface{}) []interface{} {

	innerList := make([]interface{}, 0)
	for _, inner := range l {
		for typeKey, value := range inner.(map[string]interface{}) {

			typeKey, err := ParseKey(typeKey)
			if err != nil {
				logger.ErrorLogger.Printf("Invalid Type key %q ", typeKey)
				continue
			}
			switch typeKey {

			case "S":
				newValue, err := ParseString(value.(string))
				if err == nil {
					innerList = append(innerList, newValue)
				}
			case "N":
				newValue, err := ParseNumber(value.(string))
				if err == nil {
					innerList = append(innerList, newValue)
				}

			case "BOOL":
				newValue, err := ParseBoolean(value.(string))
				if err == nil {
					innerList = append(innerList, newValue)
				}
			case "NULL":
				newValue, _ := ParseBoolean(value.(string))
				if err == nil && newValue.(bool) {
					innerList = append(innerList, newValue)
				}
			}
		}
	}

	return innerList
}

// Parsing the Map type
func ParseMap(m map[string]interface{}) map[string]interface{} {

	innerMap := make(map[string]interface{}, 0)
	for mapKey, inner := range m {
		mapKey, err := ParseKey(mapKey)
		if err != nil {
			logger.ErrorLogger.Printf("Invalid map key %q ", mapKey)
			continue
		}

		for typeKey, value := range inner.(map[string]interface{}) {

			typeKey, err = ParseKey(typeKey)
			if err != nil {
				logger.ErrorLogger.Printf("Invalid Type key %q ", typeKey)
				continue
			}
			switch typeKey {

			case "S":
				newValue, err := ParseString(value.(string))
				if err == nil {
					innerMap[mapKey] = newValue
				}
			case "N":
				newValue, err := ParseNumber(value.(string))
				if err == nil {
					innerMap[mapKey] = newValue
				}

			case "BOOL":
				newValue, err := ParseBoolean(value.(string))
				if err == nil {
					innerMap[mapKey] = newValue
				}
			case "NULL":
				newValue, _ := ParseBoolean(value.(string))
				if err == nil && newValue.(bool) {
					innerMap[mapKey] = newValue
				}

			case "L":
				newList := ParseList(value.([]interface{}))
				if len(newList) != 0 {
					innerMap[mapKey] = newList
				}

			case "M":
				// Recursively calling the ParseMap
				newMap := ParseMap(value.(map[string]interface{}))
				if len(newMap) != 0 {
					innerMap[mapKey] = newMap
				}
			}
		}
	}

	return innerMap
}
