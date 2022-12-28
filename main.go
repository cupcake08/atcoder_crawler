package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

const URL = "https://atcoder.jp"

type Problem struct {
	Prefix string
	Uri string
	Name string
}

func create_file(filename,content string) {
	f,err := os.Create(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	f.WriteString(content + "\n")
}

func getTestCases(uri string,c *colly.Collector) {
	d := c.Clone()
	d.OnHTML("span.lang > span.lang-en",func(h *colly.HTMLElement) {
		ic := 0
		oc := 0
		h.ForEach("div.part > section",func(i int, h *colly.HTMLElement) {
			// skip first four nodes
			if i > 3 {
			    // even case for input
			    // odd case for output
				text := h.ChildText("pre")
				if i % 2 == 0 {
					filename := fmt.Sprintf("input_%d.txt",ic)
					create_file(filename,text)
					ic++
				}else {
					filename := fmt.Sprintf("output_%d.txt",oc)
					create_file(filename,text)
					oc++
				}
			}
		})
	})
	d.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:",r.URL)
	})
	d.Visit(URL + uri)
}

func getProblems(c *colly.Collector,name string) {
	d := c.Clone()

	d.OnHTML("div.panel > table",func(h *colly.HTMLElement) {
		h.ForEach("tbody > tr",func(i int, h *colly.HTMLElement) {
			pCode := h.ChildText("td.text-center")
			pb := Problem{
				Prefix: pCode,
				Uri: h.ChildAttr("td > a","href"),
				Name: strings.TrimPrefix(h.ChildText("td > a"),pCode),
			}
			os.Mkdir(pb.Prefix,0777)
			os.Chdir(pb.Prefix)
			getTestCases(pb.Uri,d)
			os.Chdir("..")
		})
	})

	d.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:",r.URL)
	})
	d.Visit(URL + "/contests/" + name + "/tasks")
}



func main() {
	folder := flag.String("dir",".","Select Directory")
	contest := flag.String("contest","","Contest Name")
	flag.Parse()
	if *contest == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	os.Chdir(*folder)
	c := colly.NewCollector()
	os.Mkdir(*contest,0777)
	os.Chdir(*contest)
	getProblems(c,*contest)
}
