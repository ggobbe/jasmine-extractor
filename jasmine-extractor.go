package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

func main() {
	var folder string = "." // current folder if no args
	if len(os.Args) > 2 {
		fmt.Println("Usage: jasmine-extractor [folder]")
		return
	}
	if len(os.Args) > 1 {
		folder = os.Args[1]
	}

	files, err := ioutil.ReadDir(folder)
	check(err)

	var waitGroup sync.WaitGroup

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".spec.js") {
			waitGroup.Add(1)
			go func(filename string) {
				ExtractJasmine(path.Join(folder, filename))
				waitGroup.Done()
			}(file.Name())
		}
	}

	waitGroup.Wait()
}

func ExtractJasmine(filename string) {
	jsFile, err := os.Open(filename)
	check(err)
	defer jsFile.Close()

	specFile, err := os.Create(filename[0:len(filename)-len(path.Ext(filename))] + ".txt")
	check(err)
	defer specFile.Close()

	reader := bufio.NewReader(jsFile)
	writer := bufio.NewWriter(specFile)

	re := regexp.MustCompile("^(\\s*)(\\w+)\\((.+),(\\s*)function(.*)$")

	fmt.Printf("Extracting %s...\n", path.Base(filename))

	for line, err := Readln(reader); err == nil; line, err = Readln(reader) {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			if matches[2] == "it" {
				writer.WriteString("    ")
			} else {
				writer.WriteString("\n  ")
			}
			writer.WriteString(fmt.Sprintf("%s: %s\n", matches[2], matches[3]))
		}
	}

	writer.Flush()
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned if there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
