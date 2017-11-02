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
		return []byte("{ \"success\" : false }")
	}

	file, err := os.Create(tc.filename)
	if err != nil {
		return []byte("{ \"success\" : false }")
	}
	defer file.Close()

	err = tmpl.Execute(file, req)

	if err != nil {
		return []byte("{ \"success\" : false }")
	}

	return []byte("{ \"success\" : true }")
}

//
func (tc *TemplateCommand) Description() CommandDescription {
	desc := new(CommandDescription)
	desc.Name = "Template"
	desc.Description = "Displays a system notification via the Desktop Environment."
	desc.Command = "template"

	return *desc
}
