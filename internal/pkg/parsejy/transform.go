package parsejy

import (
	"github.com/ghodss/yaml"
)

func TransformJsonToYaml(jso []byte) (y []byte, err error) {
	y, err = yaml.JSONToYAML(jso)
	if err != nil {
		return
	}
	return
}

func TransformYamlToJson(yam []byte) (jso []byte, err error) {
	jso, err = yaml.YAMLToJSON(yam)
	if err != nil {
		return
	}
	return
}
