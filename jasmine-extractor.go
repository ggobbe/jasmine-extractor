package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	newLine = ""
)

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	var waitGroup sync.WaitGroup

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".spec.js") {
			waitGroup.Add(1)
			go func(fileName string) {
				ExtractJasmine(fileName)
				waitGroup.Done()
			}(file.Name())
		}
	}

	waitGroup.Wait()
}

func ExtractJasmine(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(file)

	re := regexp.MustCompile("^(\\s*)(\\w+)\\((.+),(\\s*)function(.*)$")

	fmt.Printf("%s::: %s :::\n", newLine, fileName)
	newLine = "\n\n"

	for line, err := Readln(reader); err == nil; line, err = Readln(reader) {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			if matches[2] == "it" {
				fmt.Printf("    ")
			} else {
				fmt.Printf("\n  ")
			}
			fmt.Printf("%s: %s\n", matches[2], matches[3])
		}
	}
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
