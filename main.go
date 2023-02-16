package main

import (
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/spf13/viper"
)

type Dummy struct {
	PropertyA string `mapstructure:"IDENTIFIER_A" default:"Dummy"`
	PropertyB int    `mapstructure:"identifier_B" default:"15"`
	PropertyC bool   `mapstructure:"identifierC"  default:"false"`
}

func Initialize(useConfigFile bool) *Dummy {
	dummy := Dummy{}

	bindData := getStructKeys(&dummy)
	for _, data := range bindData {
		viper.BindEnv(data.Mapstructure)
	}

	addDefaultData(bindData, &dummy)

	if useConfigFile {
		viper.SetConfigType("env")
		viper.SetConfigFile("./test.conf")
		fmt.Printf("Using config: %s\n", viper.ConfigFileUsed())
		viper.ReadInConfig()
	
		if err := viper.Unmarshal(&dummy); err != nil {
			fmt.Printf("Unmarshal failed: %s\n", err.Error())
			return nil
		}
	
		fmt.Printf("After config file - Some dummy: %+v\n", dummy)
	}

	updateEnvData(bindData, &dummy)

	fmt.Printf("After env vars - Some dummy: %+v\n", dummy)

	return &dummy
}

type StructFieldInfo struct {
	Key          string
	Type         string
	Mapstructure string
	Default      string
}

func getStructKeys(dummy *Dummy) []StructFieldInfo {
	result := []StructFieldInfo{}

	elements := reflect.ValueOf(dummy).Elem()
	for i := 0; i < elements.NumField(); i++ {
		field := elements.Type().Field(i)
		result = append(result, StructFieldInfo{
			Key:          field.Name,
			Type:         field.Type.Name(),
			Mapstructure: field.Tag.Get("mapstructure"),
			Default:      field.Tag.Get("default"),
		})
	}

	return result
}

func updateEnvData(structFieldInfo []StructFieldInfo, dummy *Dummy) {
	for _, data := range structFieldInfo {
		envValue := os.Getenv(data.Mapstructure)

		if envValue != "" {
			field := getField(data.Key, dummy)

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

func addDefaultData(structFieldInfo []StructFieldInfo, dummy *Dummy) {
	for _, data := range structFieldInfo {

		if data.Default != "" {
			field := getField(data.Key, dummy)

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

func isValueFilled(fieldKey, fieldType string, dummy *Dummy) bool {
	field := getField(fieldKey, dummy)

	return field.IsValid()
}

func getField(fieldKey string, dummy *Dummy) reflect.Value {
	r := reflect.ValueOf(dummy)
	field := reflect.Indirect(r).FieldByName(fieldKey)

	return field
}
