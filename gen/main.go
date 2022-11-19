package main

import (
	"bufio"
	"bytes"
	"go/format"
	"log"
	"os"
	"strings"
	"text/template"
)

type Resource struct {
	Upper string
	Lower string
}

func newResource(name string) Resource {
	return Resource{
		Upper: strings.Title(name),
		Lower: strings.ToLower(name),
	}
}

type Service struct {
	TemplateName   string
	OutputFileName string
}

func newService(name string) Service {
	return Service{
		TemplateName:   name + ".tmpl",
		OutputFileName: name + ".go",
	}
}

var SupportedServices = []string{"http", "routes", "storage", "service"}
var SupportedResources = []string{"foo", "bar"}

func main() {

	Services := []Service{}
	for _, service := range SupportedServices {
		Services = append(Services, newService(service))
	}

	Resources := []Resource{}
	for _, res := range SupportedResources {
		Resources = append(Resources, newResource(res))
	}

	tmpls := []string{}
	tmplDir := "gen/templates"
	for _, b := range Services {
		tmpls = append(tmpls, tmplDir+"/"+b.TemplateName)
	}

	baseOutputDir := "./gen/out"

	tmpl := template.Must(template.New("").ParseFiles(tmpls...))
	var processed bytes.Buffer
	var formatted []byte
	var outputPath string
	var f *os.File
	var w *bufio.Writer
	var outputDir string
	for _, resource := range Resources {
		for _, val := range Services {
			outputDir = strings.Join([]string{baseOutputDir, resource.Lower}, "/")
			err := os.MkdirAll(outputDir, os.ModePerm)
			if err != nil {
				log.Fatal("error creating directory path for output: ", err)
			}
			processed = bytes.Buffer{}
			formatted = []byte{}
			err = tmpl.ExecuteTemplate(&processed, val.TemplateName, resource)
			if err != nil {
				log.Fatalf("Unable to parse data into template: %v\n", err)
			}
			formatted, err = format.Source(processed.Bytes())
			if err != nil {
				log.Fatalf("Could not format processed template: %v\n", err)
			}

			outputPath = strings.Join([]string{outputDir, val.OutputFileName}, "/")
			f, _ = os.Create(outputPath)
			w = bufio.NewWriter(f)
			w.WriteString(string(formatted))
			w.Flush()
		}
	}

}

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
