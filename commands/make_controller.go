package commands

import (
	"fmt"
	"os"
	"strings"
)

func MakeController(filename string, arg1 string) error {
	if arg1 != "" {
        err := fmt.Errorf(fmt.Sprintf(`no command "go run main.go %v"`, CreateArgsText()))
		return err
	}
	if filename == "" {
		fmt.Print("controller name? ")
		filename = AskRequiredThing()
	}
	structName := strings.Replace(filename, string(filename[0]), strings.ToUpper(string(filename[0])), 1)
	filepath := fmt.Sprintf("./controller/%v_controller.go", filename)
    _, err := os.Stat(filepath)
    if err == nil {
        err = fmt.Errorf(fmt.Sprintf("the file, %v_controller.go exists", filename))
		return err
    }
	fp, err := os.Create(filepath)
	if err != nil {
        err = fmt.Errorf(fmt.Sprintf("couldn't make %v_controller.go", filename))
		return err
	}
	defer fp.Close()
	content := fmt.Sprintf(fileContent, structName, string(filename[0]), structName, filename)
	fp.WriteString(content)
	return nil
}


const fileContent =
`package controller

import (
	"net/http"
)

type %v struct{}

func (%v *%v) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}

	RenderTemplate(w, r, "%v.html", &TemplateData{})
}`