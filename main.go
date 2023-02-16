package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

func ParseEnv(inputStruct interface{}) {
	metaData := getStructMetaData(inputStruct)

	addDefaultData(metaData, inputStruct)

	updateEnvData(metaData, inputStruct)
}

type StructMetaData struct {
	Key     string
	Type    string
	Env     string
	Default string
}

func getStructMetaData(inputStruct interface{}) []StructMetaData {
	result := []StructMetaData{}

	elements := reflect.ValueOf(inputStruct).Elem()
	for i := 0; i < elements.NumField(); i++ {
		field := elements.Type().Field(i)
		result = append(result, StructMetaData{
			Key:     field.Name,
			Type:    field.Type.Name(),
			Env:     field.Tag.Get("env"),
			Default: field.Tag.Get("default"),
		})
	}

	return result
}

func updateEnvData(structFieldInfo []StructMetaData, inputStruct interface{}) {
	for _, data := range structFieldInfo {
		envValue := os.Getenv(data.Env)

		if envValue != "" {
			field := getField(data.Key, inputStruct)

			switch data.Type {
			case "string":
				field.SetString(envValue)
			case "int":
				if value, err := strconv.Atoi(envValue); err == nil {
					field.SetInt(int64(value))
				} else {
					fmt.Printf("Parsing %s (%s) as %s failed\n", data.Key, envValue, data.Type)
				}
			case "bool":
				field.SetBool(envValue == "true")
			}
		}
	}
}

func addDefaultData(structFieldInfo []StructMetaData, inputStruct interface{}) {
	for _, data := range structFieldInfo {

		if data.Default != "" {
			field := getField(data.Key, inputStruct)

			switch data.Type {
			case "string":
				field.SetString(data.Default)
			case "int":
				value, _ := strconv.Atoi(data.Default)
				field.SetInt(int64(value))
			case "bool":
				field.SetBool(data.Default == "true")
			}
		}
	}
}

func getField(fieldKey string, inputStruct interface{}) reflect.Value {
	r := reflect.ValueOf(inputStruct)
	field := reflect.Indirect(r).FieldByName(fieldKey)

	return field
}
