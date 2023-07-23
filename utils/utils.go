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
	// the number of seconds elapsed since January 1, 1970 UTC.
	// fmt.Println(t.Unix())
	return t.Unix(), nil

}

func ParseKey(s string) (string, error) {

	// Must sanitize the keys of trailing and leading whitespace before processing.
	// Must represent keys with String data type.
	// Must omit fields with empty keys.
	// Must omit all invalid fields.

	newString := strings.TrimSpace(s)
	if len(newString) == 0 {
		return "", ErrEmptyString
	}
	return newString, nil
}

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

	// It stores any Numeric value (positive, negative, int, float, etc.)
	// Transformation criteria.
	// Must be transformed to the relevant Numeric data type.
	// Must sanitize the value of trailing and leading whitespace before processing.
	// Must strip the leading zeros.
	// Must omit fields with invalid Numeric values.
	return -1, nil

}
func ParseBoolean(s string) (interface{}, error) {

	//Remove empty space from front and back of string
	newString := strings.TrimSpace(s)

	b, err := strconv.ParseBool(newString)
	if err != nil {
		return -1, ErrInvalidBoolean
	}

	return b, nil

}

func ParseList(l []interface{}) []interface{} {

	innerList := make([]interface{}, 0)

	for _, inner := range l {
		// fmt.Printf("%v %T %v \n", mapKey, inner, inner)

		//fmt.Println(inner.(map[string]interface{}))
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

		/*if value, found := inner.(map[string]interface{})["S"]; found {
			//fmt.Println(value)
			newValue, err := ParseString(value.(string))
			if err == nil {
				innerList = append(innerList, newValue)
			}
			continue
		}
		if value, found := inner.(map[string]interface{})["N"]; found {
			//fmt.Println(value)
			newValue, err := ParseNumber(value.(string))
			if err == nil {
				innerList = append(innerList, newValue)
			}
			continue
		}
		if value, found := inner.(map[string]interface{})["BOOL"]; found {
			//fmt.Println(value)
			newValue, err := ParseBoolean(value.(string))
			if err == nil {
				innerList = append(innerList, newValue)
			}
			continue
		}
		if value, found := inner.(map[string]interface{})["NULL"]; found {
			//fmt.Println(value)

			newValue, err := ParseBoolean(value.(string))
			if err == nil && newValue.(bool) {
				innerList = append(innerList, newValue)
			}
			continue
		}*/
	}

	// fmt.Println(innerList)
	return innerList
}

func ParseMap(m map[string]interface{}) map[string]interface{} {

	innerMap := make(map[string]interface{}, 0)

	// fmt.Println()
	// fmt.Println()
	// fmt.Println("Here::::", m)

	for mapKey, inner := range m {

		// fmt.Printf("%v %T %v \n", mapKey, inner, inner)
		// fmt.Println(inner.(map[string]interface{}))

		mapKey, err := ParseKey(mapKey)
		if err != nil {
			logger.ErrorLogger.Printf("Invalid map key %q ", mapKey)
			continue
		}

		// How to parse type key from the map
		// Removing white space from type key

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
				newMap := ParseMap(value.(map[string]interface{}))
				if len(newMap) != 0 {
					innerMap[mapKey] = newMap
				}
			}
		}

		/*if value, found := inner.(map[string]interface{})["S"]; found {
			//fmt.Println(value)
			newValue, err := ParseString(value.(string))
			if err == nil {
				innerMap[key] = newValue
			}
			continue
		}
		if value, found := inner.(map[string]interface{})["N"]; found {
			//fmt.Println(value)
			newValue, err := ParseNumber(value.(string))
			if err == nil {
				innerMap[key] = newValue
			}
			continue
		}
		if value, found := inner.(map[string]interface{})["BOOL"]; found {
			//fmt.Println(value)
			newValue, err := ParseBoolean(value.(string))
			if err == nil {
				innerMap[key] = newValue
			}
			continue
		}
		if value, found := inner.(map[string]interface{})["NULL"]; found {
			//fmt.Println(value)

			newValue, err := ParseBoolean(value.(string))

			if err == nil && newValue.(bool) {
				innerMap[key] = newValue
			}
			continue
		}*/
	}

	// fmt.Println(innerMap)
	return innerMap
}
