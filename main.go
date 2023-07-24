package main

import (
	"encoding/json"
	"log"
	"os"

	logger "github.com/deepset/tweether/logger"
	utils "github.com/deepset/tweether/utils"
)

type InputMap struct {
	inputDataMap map[string]map[string]interface{}
}

type OutputMap struct {
	outputDataMap map[string]interface{}
}

func main() {

	// Open our jsonFile
	jsonFile, err := os.Open("json/input.json")
	if err != nil {
		logger.ErrorLogger.Fatalln("Error opening the file...")
	}
	defer jsonFile.Close()

	// Reading the input data to InputMap{}
	input := InputMap{}
	err = json.NewDecoder(jsonFile).Decode(&input.inputDataMap)
	if err != nil {
		log.Fatalf("File reading error %v", err)
	}

	// Reading the output data in OutputMap{}
	output := OutputMap{}
	output.outputDataMap = make(map[string]interface{})

	for key, inner := range input.inputDataMap {
		mapKey, err := utils.ParseKey(key)
		if err != nil {
			logger.ErrorLogger.Printf("Invalid map key %q ", key)
			continue
		}

		for typeKey, value := range inner {

			// parsing and validating the map key
			typeKey, err = utils.ParseKey(typeKey)
			if err != nil {
				logger.ErrorLogger.Printf("Invalid Type key %q ", typeKey)
				continue
			}

			// Parsing and validating different data types
			switch typeKey {

			case "S":
				if utils.CheckType(value) != "STRING" {
					continue
				}
				newValue, err := utils.ParseString(value.(string))
				if err == nil {
					output.outputDataMap[mapKey] = newValue
				}
			case "N":
				if utils.CheckType(value) != "STRING" {
					continue
				}
				newValue, err := utils.ParseNumber(value.(string))
				if err == nil {
					output.outputDataMap[mapKey] = newValue
				}

			case "BOOL":
				if utils.CheckType(value) != "STRING" {
					continue
				}
				newValue, err := utils.ParseBoolean(value.(string))
				if err == nil {
					output.outputDataMap[mapKey] = newValue
				}
			case "NULL":
				if utils.CheckType(value) != "STRING" {
					continue
				}
				newValue, err := utils.ParseBoolean(value.(string))
				if err == nil && newValue.(bool) {
					output.outputDataMap[mapKey] = newValue
				}

			case "L":
				if utils.CheckType(value) != "LIST" {
					continue
				}
				newList, err := utils.ParseList(value.([]interface{}))
				if err == nil && len(newList) != 0 {
					output.outputDataMap[mapKey] = newList
				}

			case "M":
				if utils.CheckType(value) != "MAP" {
					continue
				}
				newMap, err := utils.ParseMap(value.(map[string]interface{}))
				if err == nil && len(newMap) != 0 {
					output.outputDataMap[mapKey] = newMap
				}

			}
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(output.outputDataMap); err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

}
