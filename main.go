package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"strings"
)

type Entry struct {
	Title          string       `json:"title,omitempty"`
	Pronounciation string       `json:"pronounciation,omitempty"`
	SenseGroups    []SenseGroup `json:"sense_groups,omitempty"`
}

type SenseGroup struct {
	PartOfSpeech string  `json:"part_of_speech,omitempty"`
	Senses       []Sense `json:"senses,omitempty"`
	SubSenses    []Sense `json:"subsenses,omitempty"`
}

type Sense struct {
	Iteration           string `json:"iteration,omitempty"`
	TransivityStatement string `json:"transivity_statement,omitempty"`
	WordForm            string `json:"word_form,omitempty"`
	Definition          string `json:"definition,omitempty"`
	ExampleGroups       []ExampleGroup
}

type ExampleGroup struct {
	TransivityStatement string `json:"transivity_statement,omitempty"`
	Example             string `json:"example,omitempty"`
}

type Example struct {
	Text string `json:"text,omitempty"`
}

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

	entries := []Entry{}
	doc.Find(".entryPageContent").Each(func(i int, domEntryPageContent *goquery.Selection) {

		senseGroups := []SenseGroup{}
		domEntryPageContent.Find(".senseGroup").Each(func(i int, domSenseGroup *goquery.Selection) {

			senses := []Sense{}
			domSenseGroup.Find(".sense").Each(func(ii int, domSense *goquery.Selection) {
				senses = append(senses, getSense(domSense))
			})

			subSenses := []Sense{}
			domSenseGroup.Find(".subsense").Each(func(ii int, domSubSense *goquery.Selection) {
				subSenses = append(subSenses, getSense(domSubSense))
			})

			senseGroups = append(senseGroups, SenseGroup{
				PartOfSpeech: domSenseGroup.Find(".partOfSpeechTitle .partOfSpeech").Text(),
				Senses:       senses,
				SubSenses:    subSenses,
			})
		})

		entries = append(entries, Entry{
			Title:          domEntryPageContent.Find(".entryHeader > .definitionOf > em").Text(),
			Pronounciation: domEntryPageContent.Find(".entryHeader > .headpron").Text(),
			SenseGroups:    senseGroups,
		})
	})

	b, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		fmt.Println("Error:", err)
	}

	os.Stdout.Write(b)
}

func getSense(node *goquery.Selection) Sense {

	exampleGroups := []ExampleGroup{}
	node.Find(".senseInnerWrapper .exampleGroup").Each(func(ii int, domExampleGroup *goquery.Selection) {
		exampleGroups = append(exampleGroups, ExampleGroup{
			TransivityStatement: strings.TrimSpace(domExampleGroup.Find(".transivityStatement").Text()),
			Example:             strings.TrimSpace(domExampleGroup.Find(".example").Text()),
		})
	})

	return Sense{
		Iteration:           strings.TrimSpace(node.Find(".senseInnerWrapper > .iteration").Text()),
		TransivityStatement: strings.TrimSpace(node.Find(".senseInnerWrapper > .transivityStatement").Text()),
		WordForm:            strings.TrimSpace(node.Find(".senseInnerWrapper > .wordForm").Text()),
		Definition:          strings.TrimSpace(node.Find(".senseInnerWrapper > .definition").Text()),
		ExampleGroups:       exampleGroups,
	}
}
