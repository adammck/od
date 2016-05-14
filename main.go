package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	lang := "english"
	args := os.Args
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s word\n", args[0])
		os.Exit(1)
	}

	word := args[1]
	url := fmt.Sprintf("http://www.oxforddictionaries.com/definition/%s/%s", lang, word)
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	//title := strings.TrimSpace(doc.Find(".pageTitle").Text())

	doc.Find(".senseGroup").Each(func(i int, s *goquery.Selection) {
		//partOfSpeech := s.Find(".partOfSpeechTitle .partOfSpeech").Text()
		//fmt.Printf(" - %s\n", partOfSpeech)

		s.Find(".sense").Each(func(ii int, ss *goquery.Selection) {
			//iteration := ss.Find(".iteration").Text()
			def := ss.Find(".definition").Text()

			// Definitions with examples end with a colon, but we don't show them.
			def = regexp.MustCompile(`:$`).ReplaceAllString(def, ".")
			fmt.Println(def)

			// ss.Find(".exampleGroup").Each(func(iii int, sss *goquery.Selection) {
			// 	example := sss.Find(".example").Text()
			// 	fmt.Println(example)
			// })
		})
	})
}
