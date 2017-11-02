// template.go - Fergus In London <fergus@fergus.london> (November 2017)
//
// WIP - Currently populates a template with variables taken from the command
//  object. Ideally it would be nice to cache these variables somewhere - perhaps
//  - a yaml file? - and to take the template from a file.
package commands

import (
	"os"
	"text/template"
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
	tc.template = "The directory is {{.Directory}} and it's configured for {{.Connections}} connections."
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

	tmpl, err := template.New("example").Parse(tc.template)
	if err != nil {
		return hasError("Unable to parse template file.")
	}

	file, err := os.Create(tc.filename)
	if err != nil {
		return hasError("Unable to open target file.")
	}

	defer file.Close()
	if err = tmpl.Execute(file, req); err != nil {
		return hasError("Unable to populate and save terget file.")
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
