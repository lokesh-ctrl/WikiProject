package main

import (
	"awesomeProject/config"
	"encoding/xml"
	"fmt"
	"github.com/pkg/profile"
	"io/ioutil"
	"os"
	"regexp"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type MediaWiki struct {
	XMLName xml.Name `xml:"mediawiki"`
	Pages   []Page   `xml:"page"`
}

type Page struct {
	Title    string   `xml:"title"`
	Redirect string   `xml:"redirect"`
	Revision Revision `xml:"revision"`
}

type Revision struct {
	XMLName xml.Name `xml:"revision"`
	Text    string   `xml:"text"`
}

type FormattedPage struct {
	Title    string
	Redirect string
	Links    []string
}

func main() {
	//Profiling to check cpu usage
	start := time.Now()
	defer profile.Start().Stop()
	config := config.GetConfig()
	// Open our xmlFile
	xmlFile, err := os.Open(config.InputFile)
	check(err)

	fmt.Println("Successfully Opened test.xml")

	defer xmlFile.Close()
	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var pages MediaWiki
	var formattedPage FormattedPage
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'pages' which we defined above
	xml.Unmarshal(byteValue, &pages)

	// we iterate through every page within our pages array
	for i := 0; i < len(pages.Pages); i++ {
		formattedPage.Title = pages.Pages[i].Title
		formattedPage.Redirect = pages.Pages[i].Redirect
		formattedPage.Links = findLinksInPage(pages.Pages[i].Revision.Text)
		savePage(formattedPage)
	}
	fmt.Println("Time taken")
	fmt.Println(time.Since(start))
}

func savePage(page FormattedPage) {
	insertPage(page)
}

func findLinksInPage(text string) []string {
	//REGEX to find the link in page
	const REGEX = `\[\[([[:alnum:]]|[[:blank:]]|\|)*\]\]`
	var Links []string
	r := regexp.MustCompile(REGEX)
	matches := r.FindAllStringSubmatch(text, -1)
	for _, v := range matches {
		Links = append(Links, v[0])
	}
	return Links
}
