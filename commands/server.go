package commands

import (
	_ "embed"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

// path, err := os.Getwd()
// if err != nil {
// 	fmt.Println(err)
// }


type Stdout struct {
	IsError   bool
	Messages []string
}

var (
	ignoredFolderNames = []string{"commands", "public", "data", "session", "templates", "upload"}
	ignoredFileNames = []string{"main.go"}
	cmd *exec.Cmd
	fileCache map[string][]byte
	wg sync.WaitGroup
	compileStdErr Stdout
)

func RunServer() {
	for _, text := range startOutputMessage {
		fmt.Println("\x1b[36m", text, "\x1b[0m")
	}
	cmd := exec.Command("go", "run", "./server/main.go")
	stdout, _ := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	err := cmd.Start(); if err != nil {
		tmp := make([]byte, 1024)
		stdout.Read(tmp)
		fmt.Print(string(tmp))
		return
	}
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
}

func Watch() {
	for _, text := range startOutputMessage {
		fmt.Println("\x1b[36m", text, "\x1b[0m")
	}
	wg.Add(1)
    cancel := make(chan struct{})
	go loading(cancel, &wg)
	go preRunServer(&wg)
	wg.Wait()
	close(cancel)
	CheckGoFiles()
}

func CheckGoFiles() {
	fileCache = map[string][]byte{}
	paths := dirwalk("./")
	for _, path := range paths {
		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}
		fileCache[path] = file
	}

	reconfirmation := false
	var cancel chan struct{}
	// counter := 0
	for {
		isChanged := false
		for _, path := range paths {
			file, err := os.ReadFile(path)
			if err != nil {
				fmt.Println(err)
			}
			if string(fileCache[path]) != string(file) {
				fileCache[path] = file
				isChanged = true
				if reconfirmation {
					continue
				}
				cancel = make(chan struct{})
				wg.Add(1)
				go loading(cancel, &wg)
				continue
			}
		}
		if isChanged {
			reconfirmation = true
			time.Sleep(time.Second*2)
			continue
		}
		if reconfirmation {
			reconfirmation = false
			go preRunServer(&wg)
			wg.Wait()
			close(cancel)
		}
		time.Sleep(time.Second)
	}
}

func preRunServer(wg *sync.WaitGroup){
	// 現在のファイルを元にコンパイル
	compileStdErr = compileMain();if compileStdErr.IsError {
		if _, err := os.Stat("main"); err != nil {
			fmt.Println("error: main file doesn't exist😩")
			os.Exit(2)
			return
		}
	}

	// 初回起動時
	if cmd != nil {
		cmd.Process.Kill()
		cmd.Wait()
		runBinaryFile(wg)
		return
	}
	runBinaryFile(wg)
}

func loading(cancel chan struct{}, wg *sync.WaitGroup) bool {
	fmt.Println("")
	fmt.Printf("\x1b[36mLoading Go\x1b[0m  ")
    for {
        select {
        case <-cancel:
			fmt.Printf("\x1b[36mStart Watching!!\x1b[0m \n\n")
			if compileStdErr.IsError {
				display(compileStdErr)
			}
            return true
        default:
            fmt.Printf("\x1b[36m→\x1b[0m  ")
            time.Sleep(time.Second / 3)
        }
    }
}

func runBinaryFile(wg *sync.WaitGroup) {
	wg.Done()
	main := exec.Command("./main")
	output, _ := main.StdoutPipe()
	main.Start()
	cmd = main
	time.Sleep(time.Second)
	for {
		tmp := make([]byte, 1024)
		_, err := output.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
}

func compileMain() Stdout {
	build := exec.Command("go", "build", "-v", "-o", "main", "./server/main.go")
	stderr, err := build.StderrPipe()
	if err != nil {
		fmt.Println(err)
	}
	build.Start()
	output := outputArray(stderr)
	build.Wait()
	return output

}

func outputArray(std io.ReadCloser) Stdout {
	var messages []string
	stdout := Stdout{false, messages}
	for {
		tmp := make([]byte, 1024)
		_, err := std.Read(tmp)
		message := string(tmp)
		if err != nil {
			break
		}
		stdout.Messages = append(messages, message)
		if !stdout.IsError {
			// r := regexp.MustCompile("error")
			// stdout.IsError = r.MatchString(message)
			r := regexp.MustCompile(".go:")
			stdout.IsError = r.MatchString(message)
		}
	}
	return stdout
}

func display(stdout Stdout) bool {
	if len(stdout.Messages) != 0 {
		fmt.Println("")
		fmt.Println("---something happened😶---\n")
		if stdout.IsError {
			for _, message := range stdout.Messages {
				fmt.Println("\x1b[31m", message, "\x1b[0m")
			}
		} else {
			for _, message := range stdout.Messages {
				fmt.Println(message)
			}
		}
		fmt.Println("----------------------------\n")
		time.Sleep(5 * time.Second)
		return true
	}
	return false
}

func dirwalk(dir string) []string {
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        panic(err)
    }

    var paths []string
    for _, file := range files {
        if file.IsDir() {
			if !contains(file.Name(), ignoredFolderNames) {
				paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
				continue
			}
        }
		if contains(file.Name(), ignoredFileNames) {
			continue
		}
		if filepath.Ext(file.Name()) == ".go" {
			paths = append(paths, filepath.Join(dir, file.Name()))
		}
    }

    return paths
}

func contains(dirName string, names []string) bool {
    for _, name := range names {
        if name == dirName {
            return true
        }
    }
    return false
}


var startOutputMessage = []string{
	"",
	" ######  ########    ###    ########  ########     ######    #######",
	"##    ##    ##      ## ##   ##     ##    ##       ##    ##  ##     ##",
	"##          ##     ##   ##  ##     ##    ##       ##        ##     ## ",
	" ######     ##    ##     ## ########     ##       ##   #### ##     ## ",
	"      ##    ##    ######### ##   ##      ##       ##    ##  ##     ## ",
	"##    ##    ##    ##     ## ##    ##     ##       ##    ##  ##     ## ",
	" ######     ##    ##     ## ##     ##    ##        ######    ####### ",
}

