package main

//partially adapted from https://blog.serverbooter.com/post/parsing-nested-json-in-go/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
)

func iterate(data interface{}) interface{} {

	fmt.Printf("data kind is %v", reflect.ValueOf(data).Kind())
	if reflect.ValueOf(data).Kind() == reflect.Slice {
		d := reflect.ValueOf(data)
		tmpData := make([]interface{}, d.Len())
		returnSlice := make([]interface{}, d.Len())
		for i := 0; i < d.Len(); i++ {
			tmpData[i] = d.Index(i).Interface()
		}
		for i, v := range tmpData {
			returnSlice[i] = iterate(v)
		}
		return returnSlice
	} else if string(reflect.ValueOf(data).Kind()) == "map" {
		d := reflect.ValueOf(data)
		tmpData := make(map[string]interface{})
		for _, k := range d.MapKeys() {
			typeOfValue := reflect.TypeOf(d.MapIndex(k).Interface()).Kind()
			if typeOfValue == reflect.Map || typeOfValue == reflect.Slice {
				tmpData[k.String()] = iterate(d.MapIndex(k).Interface())
			} else {
				tmpData[k.String()] = d.MapIndex(k).Interface()
			}
		}
		return tmpData
	}
	return data
}

func main() {
	// Open our jsonFile
	jsonFile, err := os.Open("familytree.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	data, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		fmt.Println(err)
	}
	//myJSON := string(file)

	var myJSON interface{}
	err = json.Unmarshal(data, &myJSON)
	fmt.Printf("%v: ", myJSON)
	if err != nil {
		fmt.Println(err)
	}

	m, _ := myJSON.(map[string]interface{})
	newM := iterate(m)
	jsonBytes, err := json.Marshal(newM)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonBytes))

}
