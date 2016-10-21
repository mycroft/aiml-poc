package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"regexp"
)

type AIMLRoot struct {
	XMLName    xml.Name       `xml:"aiml"`
	Categories []AIMLCategory `xml:"category"`
	Topics     []AIMLTopic    `xml:"topic"`
}

type AIMLTopic struct {
	XMLName    xml.Name       `xml:"topic"`
	Name       string         `xml:"name,attr"`
	Categories []AIMLCategory `xml:"category"`
}

type AIMLCategory struct {
	XMLName  xml.Name     `xml:"category"`
	Pattern  AIMLPattern  `xml:"pattern"`
	Template AIMLTemplate `xml:"template"`
	That     AIMLThat     `xml:"that"`
}

type AIMLPattern struct {
	XMLName xml.Name `xml:"pattern"`
	Content string   `xml:",innerxml"`
	regexp  *regexp.Regexp
}

type AIMLTemplate struct {
	XMLName xml.Name `xml:"template"`
	Content string   `xml:",innerxml"`
}

type AIMLThat struct {
	XMLName xml.Name `xml:"that"`
	Content string   `xml:",innerxml"`
}

func Parse(filename string) (*AIMLRoot, error) {
	var root AIMLRoot

	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	err = xml.Unmarshal(bytes, &root)
	if err != nil {
		return nil, err
	}

	return &root, nil
}
