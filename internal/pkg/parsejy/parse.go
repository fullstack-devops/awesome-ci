package parsejy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/itchyny/gojq"
)

type Syntax string

const JsonSyntax Syntax = "JSON"
const YamlSyntax Syntax = "YAML"

func ParseFile(queryString string, file string, syntax Syntax) (err error) {
	d, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = Parse(queryString, d, syntax)

	return
}

func Parse(queryString string, input []byte, syntax Syntax) (err error) {
	query, err := gojq.Parse(queryString)
	if err != nil {
		return err
	}

	if syntax == YamlSyntax {
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
