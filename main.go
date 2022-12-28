package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	URL = "https://atcoder.jp"
	// your c++ template file path
	TEMPLATE = "/home/ankit/CP/template.cpp"
)

type Problem struct {
	Prefix string
	Uri    string
	Name   string
}

type Contest struct {
	Colly *colly.Collector
	Name  string
	Base  string
}

func create_file(filename, content string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()
	f.WriteString(content + "\n")
}

func (p *Problem) InitTemplateFile(contestName string) {
	f, err := os.Open(TEMPLATE)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err.Error())
	}

	currentTime := time.Now()
	// initial header
	header := fmt.Sprintf(`/**
		*
		* Author      : Ankit Bhankharia
		* Created At  : %d-%d-%d %d:%d:%d
		* Contest     : %s
		* Problem     : %s
		*
		**/
		`,
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day(),
		currentTime.Hour(),
		currentTime.Minute(),
		currentTime.Second(),
		contestName,
		p.Name,
	)

	idx := len(header)
	idx++

	file, err := os.Create(p.Prefix + ".cpp")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	file.WriteString(header + "\n")

	_, err = file.WriteAt(bytes, int64(idx))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = file.Sync()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (pb *Problem) GetTestCases(uri string, c *colly.Collector) {
	d := c.Clone()

	d.OnHTML("span.lang > span.lang-en", func(h *colly.HTMLElement) {
		ic := 0
		oc := 0
		h.ForEach("div.part > section", func(i int, h *colly.HTMLElement) {
			// skip first four nodes
			if i > 3 {
				// even case for input
				// odd case for output
				text := h.ChildText("pre")
				if i%2 == 0 {
					filename := fmt.Sprintf("input_%d.txt", ic)
					create_file(filename, text)
					ic++
				} else {
					filename := fmt.Sprintf("output_%d.txt", oc)
					create_file(filename, text)
					oc++
				}
			}
		})
	})

	d.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	d.Visit(URL + uri)
}

func (contest *Contest) ScrapeAtcoder() {
	d := contest.Colly.Clone()

	d.OnHTML("nav > div.container-fluid > div#navbar-collapse", func(h *colly.HTMLElement) {
		contest.Name = h.ChildText("ul > li > a.contest-title")
	})

	d.OnHTML("div.panel > table", func(h *colly.HTMLElement) {
		h.ForEach("tbody > tr", func(i int, h *colly.HTMLElement) {
			pCode := h.ChildText("td.text-center")
			pb := Problem{
				Prefix: pCode,
				Uri:    h.ChildAttr("td > a", "href"),
				Name:   strings.TrimPrefix(h.ChildText("td > a"), pCode),
			}
			os.Mkdir(pb.Prefix, 0777)
			os.Chdir(pb.Prefix)
			pb.GetTestCases(pb.Uri, d)
			pb.InitTemplateFile(contest.Name)
			os.Chdir("..")
		})
	})

	d.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting:", r.URL)
	})

	d.Visit(URL + "/contests/" + contest.Base + "/tasks")
}

func Init(contestName string) *Contest {
	err := os.Mkdir(contestName, 0777)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	err = os.Chdir(contestName)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}
	return &Contest{
		Name:  "",
		Base:  contestName,
		Colly: colly.NewCollector(),
	}
}

func main() {
	folder := flag.String("dir", ".", "Select Directory")
	contest := flag.String("contest", "", "Contest Name")
	flag.Parse()
	if *contest == "" {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(1)
	}
	os.Chdir(*folder)
	atcoder := Init(*contest)
	if atcoder != nil {
		atcoder.ScrapeAtcoder()
	}
}
