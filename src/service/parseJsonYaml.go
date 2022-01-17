package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

func ParseJson(file *string, pvalue *string) {
	value := *pvalue

	dat, err := ioutil.ReadFile(*file)
	check(err)

	if strings.HasPrefix(*pvalue, "[]") {
		var result map[interface{}]interface{}

		json.Unmarshal(dat, &result)

		fmt.Print(result[value[2:]])
	} else if strings.HasPrefix(*pvalue, ".") {
		var result map[string]interface{}

		json.Unmarshal(dat, &result)

		fmt.Print(result[value[1:]])
	} else {
		var result map[string]interface{}

		json.Unmarshal(dat, &result)

		fmt.Print(result[value])

	}
}

func ParseYaml(file *string, pvalue *string) {
	value := *pvalue

	dat, err := ioutil.ReadFile(*file)
	check(err)

	if strings.HasPrefix(*pvalue, "[]") {
		var result map[interface{}]interface{}

		yaml.Unmarshal(dat, &result)

		fmt.Print(result[value[2:]])
	} else if strings.HasPrefix(*pvalue, ".") {
		var result map[string]interface{}

		yaml.Unmarshal(dat, &result)

		fmt.Print(result[value[1:]])
	} else {
		var result map[string]interface{}

		yaml.Unmarshal(dat, &result)

		fmt.Print(result[value])

	}
}
