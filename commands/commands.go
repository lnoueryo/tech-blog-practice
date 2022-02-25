package commands

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var verb string
var object string
var name string
var arg1 string

func CMD() {
	verb = flag.Arg(0)
	object = flag.Arg(1)
	name = flag.Arg(2)
	arg1 = flag.Arg(3)

	if verb == "" {
		RunServer()
		return
	}

	if verb == "deploy" {
		Deploy()
	}

	if verb == "watch" && object == "" {
		Watch()
		return
	}
	if verb == "make" {
		MakeCMD()
		return
	}
	errMessage := fmt.Sprintf(`no command "go run main.go %v"`, CreateArgsText())
	fmt.Println(errMessage)
}

func MakeCMD() {
	if object == "" {
		fmt.Print("make what? ")
		object = AskRequiredThing()
	}

	if object == "controller" {
		err := MakeController(name, arg1)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	if object == "db" {
		err := MakeDBData(name, arg1)
		if err != nil {
			fmt.Println(err)
		}
		return
	}

	errMessage := fmt.Sprintf(`no command "go run main.go %v"`, CreateArgsText())
	fmt.Println(errMessage)
}

func AskRequiredThing() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func CreateArgsText() string {
	var argsText string
	for i, v := range flag.Args() {
		if i == 0 {
			argsText += v
			continue
		}
		argsText += " " + v
	}
	return argsText
}