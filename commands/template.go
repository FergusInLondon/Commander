// template.go - Fergus In London <fergus@fergus.london> (November 2017)
//
// WIP - Currently populates a template with variables taken from the command
//  object. Ideally it would be nice to cache these variables somewhere - perhaps
//  - a yaml file? - and to take the template from a file.
package commands

import (
	"os"
	"path/filepath"
)

type TemplateCommand struct {
	filename string
	template string
}

type TemplateRequest struct {
	Connections int16  `json:"connections"`
	Directory   string `json:"webroot"`
}

func (tc *TemplateCommand) Init() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	tc.template = filepath.Dir(ex) + "/example/commander_templating.tmpl"
	tc.filename = "/tmp/commander_example.conf"
}

//
func (tc *TemplateCommand) Identifier() string {
	return "template"
}

//
func (tc *TemplateCommand) Object() interface{} {
	return &TemplateRequest{}
}

//
func (tc *TemplateCommand) Handle(command interface{}) []byte {
	req := command.(*TemplateRequest)

	if err := handleTemplate(req, tc.template, tc.filename); err != nil {
		return hasError("Unable to update configuration file")
	}

	returnData := make(map[string]interface{})
	returnData["file_updated"] = tc.filename
	returnData["new_values"] = req

	return ranSuccessfully(returnData)
}

//
func (tc *TemplateCommand) Description() CommandDescription {
	desc := new(CommandDescription)
	desc.Name = "Template"
	desc.Description = "Updates a template (i.e a configuration file)."
	desc.Command = "template"

	return *desc
}
