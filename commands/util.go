// util.go - Fergus In London <fergus@fergus.london> (November 2017)
package commands

import (
	"encoding/json"
	"os"
	"text/template"
)

// A Command is an object that is capable of handling an instruction from the
//  JSON API.
type Command interface {
	Init()
	Identifier() string
	Object() interface{}
	Handle(interface{}) []byte
	Description() CommandDescription
}

// CommandDescription holds basic information about a command which is supplied
//  via the HTTP JSON API when a client hits the "/listing" endpoint.
type CommandDescription struct {
	Name        string `json:"name"`
	Command     string `json:"command"`
	Description string `json:"description"`
}

func hasError(msg string) (jsonObject []byte) {
	payload := make(map[string]interface{})
	payload["success"] = false
	payload["error"] = msg

	jsonObject, err := json.Marshal(payload)
	if err != nil {
		jsonObject = []byte("{ \"success\" : false }")
	}

	return
}

func ranSuccessfully(result map[string]interface{}) (jsonObject []byte) {
	payload := make(map[string]interface{})
	payload["success"] = true
	payload["result"] = result

	jsonObject, err := json.Marshal(payload)
	if err != nil {
		jsonObject = []byte("{ \"success\" : true }")
	}

	return
}

func handleTemplate(data interface{}, tplFile, destination string) (err error) {
	tmpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return
	}

	file, err := os.Create(destination)
	if err != nil {
		return
	}
	defer file.Close()

	return tmpl.Execute(file, data)
}
