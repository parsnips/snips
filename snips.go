package main

import (
	"flag"
	"fmt"
	"github.com/shurcooL/github_flavored_markdown"
	"io/ioutil"
	"os"
	"sync"
)
import s "strings"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var wg sync.WaitGroup

func main() {
	inputPtr := flag.String("i", ".", "-i=/path/to/markdown/files")
	outputPtr := flag.String("o", ".", "-o=/path/to/rendered/files")
	flag.Parse()

	//Links for the index page
	links := make([]string, 0)

	//Find all files in the input directory
	files, err := ioutil.ReadDir(*inputPtr)
	check(err)

	for _, file := range files {
		//Pretty good chance it's a markdown file
		if file.IsDir() == false && s.HasSuffix(s.ToLower(file.Name()), ".md") {
			dat, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", *inputPtr, file.Name()))
			check(err)

			//Go forth and render
			wg.Add(1)
			name := s.Replace(file.Name(), ".md", ".html", -1)
			go WriteHtml(dat, fmt.Sprintf("%s/%s", *outputPtr, name))

			//Add a link to the index page, who cares about time? I dont.
			links = append(links, fmt.Sprintf("[%s](%s) \n", s.Replace(s.Replace(name, ".html", "", -1), "_", " ", -1), name))
		}
	}

	WriteIndex(links, *outputPtr)

	wg.Wait()
}

func WriteIndex(links []string, outputDir string) {
	idx := make([]byte, 0)
	title := []byte("Parsnips.net\n-----------------\n")

	idx = append(idx, title...)

	for _, link := range links {
		idx = append(idx, []byte(link)...)
	}

	f, err := os.Create(outputDir + "/index.html")
	check(err)

	defer f.Close()
	WriteHeader(f)
	_, err = f.Write(github_flavored_markdown.Markdown(idx))
	check(err)
	WriteTrailer(f)
	f.Sync()
}

func WriteHtml(markdown []byte, filePath string) {
	output := github_flavored_markdown.Markdown(markdown)
	f, err := os.Create(filePath)
	check(err)

	defer wg.Done()
	defer f.Close()

	WriteHeader(f)
	_, err = f.Write(output)
	check(err)

	f.Sync()
}

func WriteHeader(file *os.File) {
	_, err := file.WriteString(`<html><head><meta charset="utf-8"><link href="https://dl.dropboxusercontent.com/u/8554242/temp/github-flavored-markdown.css" media="all" rel="stylesheet" type="text/css" /><link href="https://cdnjs.cloudflare.com/ajax/libs/octicons/2.1.2/octicons.css" media="all" rel="stylesheet" type="text/css" /></head><body><article class="markdown-body entry-content" style="padding: 30px;">`)
	check(err)
}

func WriteTrailer(file *os.File) {
	_, err := file.WriteString(`</article></body></html>`)
	check(err)
}
