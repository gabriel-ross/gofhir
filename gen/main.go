package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

type data struct {
	Package  string
	Resource string
}

func main() {
	// var d data
	// flag.StringVar(&d.Type, "type", "", "The subtype used for the queue being generated")
	// flag.StringVar(&d.Name, "name", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
	// flag.Parse()

	// t := template.Must(template.New("queue").Parse(queueTemplate))
	// //t.Execute(os.Stdout, d)
	// rc, wc, _ := pipe.Commands(
	// 	exec.Command("gofmt"),
	// 	exec.Command("goimports"),
	// )
	// t.Execute(wc, d)
	// wc.Close()
	// io.Copy(os.Stdout, rc)

	data := data{
		Package:  "packageName",
		Resource: "resourceName",
	}

	outPath := "./gen/out"
	err := os.MkdirAll(outPath, os.ModePerm)
	if err != nil {
		log.Fatal("error creating directory path for output: ", err)
	}

	tmpl := template.Must(template.New("http.tmpl").ParseFiles("gen/templates/http.tmpl"))
	outFileName := "demo.go"
	var processed bytes.Buffer
	// err := tmpl.ExecuteTemplate(&processed, outFileName, data)
	err = tmpl.Execute(&processed, data)
	if err != nil {
		log.Fatalf("Unable to parse data into template: %v\n", err)
	}
	formatted, err := format.Source(processed.Bytes())
	if err != nil {
		log.Fatalf("Could not format processed template: %v\n", err)
	}

	outFullPath := strings.Join([]string{outPath, outFileName}, "/")
	fmt.Println("Writing file: ", outFullPath)
	f, _ := os.Create(outFullPath)
	w := bufio.NewWriter(f)
	w.WriteString(string(formatted))
	w.Flush()
}

var queueTemplate = `
package queue

import (
  "container/list"
)

func New{{.Name}}() *{{.Name}} {
  return &{{.Name}}{list.New()}
}

type {{.Name}} struct {
  list *list.List
}

func (q *{{.Name}}) Len() int {
  return q.list.Len()
}

func (q *{{.Name}}) Enqueue(i {{.Type}}) {
  q.list.PushBack(i)
}

func (q *{{.Name}}) Dequeue() {{.Type}} {
  if q.list.Len() == 0 {
    panic(ErrEmptyQueue)
  }
  raw := q.list.Remove(q.list.Front())
  if typed, ok := raw.({{.Type}}); ok {
    return typed
  }
  panic(ErrInvalidType)
}
`
