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
	// if we os.Open returns an error then handle it
	if err != nil {
		logger.ErrorLogger.Fatalln("Error opening the file...")
	}
	//fmt.Println("Successfully Opened users.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	//fmt.Println(jsonFile)

	// byteValue, err := ioutil.ReadAll(jsonFile)
	// if err != nil {
	// 	fmt.Println("Error of ReadAll :", err)
	// }

	//var result map[string]interface{}
	//err = json.Unmarshal([]byte(byteValue), &result)

	input := InputMap{}
	//input := map[string]map[string]interface{}{}
	err = json.NewDecoder(jsonFile).Decode(&input.inputDataMap)
	if err != nil {
		log.Fatalf("File reading error %v", err)
	}

	output := OutputMap{}
	output.outputDataMap = make(map[string]interface{})

	for key, inner := range input.inputDataMap {
		//fmt.Printf("%s : %v \n", key, inner)

		mapKey, err := utils.ParseKey(key)
		if err != nil {
			logger.ErrorLogger.Printf("Invalid map key %q ", key)
			continue
		}

		for typeKey, value := range inner {
			//fmt.Printf("%s \n", keyType)

			//fmt.Printf("key , type , value : (%s) : (%T) : %+v \n", key, value, value)

			//Removing white space from type key
			typeKey, err = utils.ParseKey(typeKey)
			if err != nil {
				logger.ErrorLogger.Printf("Invalid Type key %q ", typeKey)
				continue
			}
			switch typeKey {

			case "S":
				newValue, err := utils.ParseString(value.(string))
				if err == nil {
					output.outputDataMap[mapKey] = newValue
				}
			case "N":
				newValue, err := utils.ParseNumber(value.(string))
				if err == nil {
					output.outputDataMap[mapKey] = newValue
				}

			case "BOOL":
				newValue, err := utils.ParseBoolean(value.(string))
				if err == nil {
					output.outputDataMap[mapKey] = newValue
				}
			case "NULL":
				newValue, _ := utils.ParseBoolean(value.(string))
				if err == nil && newValue.(bool) {
					output.outputDataMap[mapKey] = newValue
				}

			case "L":
				newList := utils.ParseList(value.([]interface{}))
				if len(newList) != 0 {
					output.outputDataMap[mapKey] = newList
				}

			case "M":
				newMap := utils.ParseMap(value.(map[string]interface{}))
				if len(newMap) != 0 {
					output.outputDataMap[mapKey] = newMap
				}
			}
		}
	}

	//json.NewEncoder(os.Stdout).Encode(input)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(output.outputDataMap); err != nil {
		logger.ErrorLogger.Fatalln(err)
	}

}
