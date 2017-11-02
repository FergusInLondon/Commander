// echo.go - Fergus In London <fergus@fergus.london> (November 2017)
//
// A very simple demonstration of a Command object; all this command does is
//  echo the user input back.
package commands

import "encoding/json"

type EchoCommand struct{}

type EchoResponse struct {
	Message string `json:"message"`
}

func (ec *EchoCommand) Init() {

}

func (ec *EchoCommand) Identifier() string {
	return "echo"
}

func (ec *EchoCommand) Object() interface{} {
	return &EchoResponse{}
}

func (ec *EchoCommand) Handle(command interface{}) []byte {
	jsonObject, err := json.Marshal(command)
	if err != nil {
		jsonObject = []byte("{ \"success\" : false ")
	}

	return jsonObject
}

func (ec *EchoCommand) Description() CommandDescription {
	desc := new(CommandDescription)
	desc.Name = "Echo"
	desc.Description = "Echos back a message from the daemon"
	desc.Command = "echo"

	return *desc
}
