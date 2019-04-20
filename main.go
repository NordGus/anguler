package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var appPrefab = `
package app

var app = App{
	Name: "%v",
	Files: []File{
		%v
	},
}
`

var filePrefab = `
File{
	Name: "%v",
	Format: "%v",
	Contents: []byte{%v},
},
`

func main() {
	logger := log.New(os.Stdout, "anguler: ", log.LstdFlags)
	if len(os.Args) < 4 {
		logger.Println("must specify app name, orgin and target")
		return
	}
	appname, orgin, target := os.Args[1], os.Args[2], os.Args[3]
	op := strings.Replace(fmt.Sprintf("%v/dist/%v", orgin, appname), "//", "/", -1)
	tp := strings.Replace(fmt.Sprintf("%v/app", target), "//", "/", -1)
	af, err := ioutil.ReadDir(op)
	if err != nil {
		logger.Println(err)
		return
	}

	var ppf []string

	for _, f := range af {
		name := f.Name()
		format := getFileFormat(f.Name())
		contents, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", op, name))
		if err != nil {
			logger.Println(err)
			return
		}
		if format != "ico" {
			cs := contentsToString(contents)
			ppf = append(ppf, fmt.Sprintf(filePrefab, name, format, cs))
		}
	}

	fc := fmt.Sprintf(appPrefab, appname, strings.Join(ppf, "\n"))
	tf, err := os.OpenFile(fmt.Sprintf("%v/app.go", tp), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		logger.Println(err)
		return
	}
	defer tf.Close()

	_, err = tf.Write([]byte(fc))
	if err != nil {
		logger.Println(err)
		return
	}

	logger.Printf("angular app exported to %v/app.go", tp)
}

func getFileFormat(filename string) string {
	segments := strings.Split(filename, ".")
	if segments[len(segments)-1] == "js" {
		return "javascript"
	}
	if segments[len(segments)-1] == "txt" {
		return "plain"
	}
	return segments[len(segments)-1]
}

func contentsToString(contents []byte) string {
	sa := []string{}

	for _, b := range contents {
		sa = append(sa, fmt.Sprintf("%v", b))
	}

	return strings.Join(sa, ",")
}
