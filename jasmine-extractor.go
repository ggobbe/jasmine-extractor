package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".spec.js") {
			ExtractJasmine(file.Name())
		}
	}
}

func ExtractJasmine(fileName string) {
	file, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(file)

	re := regexp.MustCompile("^(\\s*)(\\w+)\\((.+),(.*)function(.*)$")

	for line, err := Readln(reader); err == nil; line, err = Readln(reader) {
		matches := re.FindStringSubmatch(line)
		if len(matches) > 0 {
			if matches[2] == "it" {
				fmt.Printf("  ")
			} else {
				fmt.Printf("\n")
			}
			fmt.Printf("%s: %s\n", matches[2], matches[3])
		}
	}
}

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
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
