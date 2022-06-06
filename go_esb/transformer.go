package main

import (
	"encoding/json"
	"encoding/xml"

	"gopkg.in/yaml.v2"
)

func JSONTransformer(message Message) ([]byte, error) {
	return json.Marshal(message)
}

func XMLTransformer(message Message) ([]byte, error) {
	return xml.Marshal(message)
}

func YMLTransformer(message Message) ([]byte, error) {
	return yaml.Marshal(message)
}

func TSVTransformer() {

}
