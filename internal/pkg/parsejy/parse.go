package parsejy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itchyny/gojq"
)

type Format string

const JsonSyntax Format = "JSON"
const YamlSyntax Format = "YAML"

func ParseFile(queryString string, file string, format Format) (err error) {
	d, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = Parse(queryString, d, format)

	return
}

func Parse(queryString string, input []byte, format Format) (err error) {
	query, err := gojq.Parse(queryString)
	if err != nil {
		return err
	}

	if format == YamlSyntax {
		input, err = TransformYamlToJson(input)
		if err != nil {
			return err
		}
	}

	var v map[string]interface{}
	err = json.Unmarshal(input, &v)
	if err != nil {
		return
	}

	iter := query.Run(v)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return err
		}
		fmt.Printf("%#v\n", v)
	}
	return nil
}
