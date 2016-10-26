package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func (context *Context) Process(template string, matches []string) (string, bool) {
	template = context.ReplaceStars(template, matches)
	template = context.ProcessGet(template)
	template = context.ProcessBot(template)
	template = context.ProcessThat(template)
	template = context.ProcessThatstar(template)
	template = context.ProcessInput(template)
	template = context.ProcessSet(template)

	template = context.ProcessRandom(template)
	template = context.ProcessThink(template)
	template = context.ProcessCondition(template)

	template, has_srai := context.ProcessSRAI(template)

	return strings.TrimSpace(template), has_srai
}

func (context *Context) ReplaceStars(template string, matches []string) string {
	starStruct := struct {
		XMLName xml.Name `xml:"star"`
		Index   int      `xml:"index,attr"`
	}{}

	current_index := 0

	for {
		// Grab <star ... /> element.
		re := regexp.MustCompile("(?msU)<star.*/>")

		match := re.FindString(template)
		if "" == match {
			return template
		}

		err := xml.Unmarshal([]byte(match), &starStruct)
		if err != nil {
			return template
		}

		var replacement = ""

		// Replace
		if starStruct.Index == 0 {
			// No index: using current_index
			if len(matches) <= current_index {
				return template
			}

			replacement = string(matches[current_index])

			current_index = current_index + 1
		} else {
			if len(matches) < starStruct.Index {
				return template
			}

			replacement = matches[starStruct.Index-1]
		}

		replacement = strings.TrimSpace(replacement)
		template = strings.Replace(template, match, fmt.Sprintf("%s", replacement), 1)
	}
}

func (context *Context) ProcessThat(template string) string {
	// thatStruct := struct {
	// 	XMLName xml.Name `xml:"that"`
	// 	Index   string   `xml:"index,attr"`
	// }{}

	for {
		re := regexp.MustCompile("(?msU)<that[^a-z].*/>")
		set_string := re.FindString(template)
		if "" == set_string {
			return template
		}

		template = strings.Replace(template, set_string, context.LastSent, 1)
	}
}

func (context *Context) ProcessInput(template string) string {
	// inputStruct := struct {
	// 	XMLName xml.Name `xml:"input"`
	// 	Index   string   `xml:"index,attr"`
	// }{}

	for {
		re := regexp.MustCompile("(?msU)<input.*/>")
		set_string := re.FindString(template)
		if "" == set_string {
			return template
		}

		template = strings.Replace(template, set_string, context.LastRecv, 1)
	}
}

func (context *Context) ProcessSet(template string) string {
	setStruct := struct {
		XMLName xml.Name `xml:"set"`
		Name    string   `xml:"name,attr"`
		Content string   `xml:",innerxml"`
	}{}

	for {
		re := regexp.MustCompile("(?msU)<set.*</set>")
		set_string := re.FindString(template)
		if "" == set_string {
			return template
		}

		err := xml.Unmarshal([]byte(set_string), &setStruct)
		if err != nil {
			return template
		}

		template = strings.Replace(template, set_string, setStruct.Content, 1)

		if setStruct.Name == "topic" {
			context.Topic = setStruct.Content
		} else {
			context.Memory[setStruct.Name] = setStruct.Content
		}
	}
}

func (context *Context) ProcessGet(template string) string {
	getStruct := struct {
		XMLName xml.Name `xml:"get"`
		Name    string   `xml:"name,attr"`
	}{}

	for {
		err := xml.Unmarshal([]byte(template), &getStruct)
		if err != nil {
			return template
		}

		value := context.Memory[getStruct.Name]

		orig_string := fmt.Sprintf("<get name=\"%s\" />", getStruct.Name)
		template = strings.Replace(template, orig_string, value, 1)
	}
}

func (context *Context) ProcessRandom(template string) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for {
		randomStruct := struct {
			XMLName  xml.Name `xml:"random"`
			Contents []string `xml:"li"`
			/*
			   List    []struct {
			       XMLName xml.Name `xml:"li"`
			       Content string   `xml:",innerxml"`
			   } `xml:"li"`
			*/
		}{}

		err := xml.Unmarshal([]byte(template), &randomStruct)
		if err != nil {
			return template
		}

		re := regexp.MustCompile("(?ms)[[:space:]]*<random>.*?</random>[[:space:]]*")
		orig_string := re.FindString(template)

		random_string := fmt.Sprintf(" %s ", randomStruct.Contents[r.Int()%len(randomStruct.Contents)])

		template = strings.Replace(template, orig_string, random_string, 1)
	}
}

func (context *Context) ProcessThink(template string) string {
	thinkStruct := struct {
		XMLName xml.Name `xml:"think"`
	}{}

	for {
		re := regexp.MustCompile("(?msU)<think.*</think>")
		orig_string := re.FindString(template)
		if "" == orig_string {
			return template
		}

		err := xml.Unmarshal([]byte(template), &thinkStruct)
		if err != nil {
			return template
		}

		// Remove contents (set are already handled)
		template = strings.Replace(template, orig_string, "", 1)
	}
}

func (context *Context) ProcessSRAI(template string) (string, bool) {
	sraiStruct := struct {
		XMLName xml.Name `xml:"srai"`
		Content string   `xml:",innerxml"`
	}{}

	re := regexp.MustCompile("(?msU)<srai.*</srai>")
	orig_string := re.FindString(template)
	if "" == orig_string {
		return template, false
	}

	err := xml.Unmarshal([]byte(template), &sraiStruct)
	if err != nil {
		return template, false
	}

	return sraiStruct.Content, true
}

func (context *Context) ProcessCondition(template string) string {
	conditionStruct := struct {
		XMLName xml.Name `xml:"condition"`
		Name    string   `xml:"name,attr"`
		Value   string   `xml:"value,attr"`
		Content string   `xml:",innerxml"`
	}{}

	for {
		re := regexp.MustCompile("(?msU)<condition .*</condition>")
		orig_string := re.FindString(template)
		if "" == orig_string {
			return template
		}

		err := xml.Unmarshal([]byte(template), &conditionStruct)
		if err != nil {
			return template
		}

		// If matching, remplace with content, else, remove.
		memory_value, in_set := context.Memory[conditionStruct.Name]

		if in_set && memory_value == conditionStruct.Value {
			// Replace by content
			template = strings.Replace(template, orig_string, conditionStruct.Content, 1)
		} else {
			// Remove
			template = strings.Replace(template, orig_string, "", 1)
		}
	}
}

func (context *Context) ProcessThatstar(template string) string {
	for {
		re := regexp.MustCompile("(?msU)<thatstar.*/>")
		orig_string := re.FindString(template)
		if "" == orig_string {
			return template
		}

		template = strings.Replace(template, orig_string, context.ThatMatches[0], 1)
	}
}

func (context *Context) ProcessBot(template string) string {
	botStruct := struct {
		XMLName xml.Name `xml:"bot"`
		Name    string   `xml:"name,attr"`
	}{}

	for {
		orig_str, err := GrabPart(template, &botStruct, "(?msU)<bot .*/>")
		if err != nil {
			return template
		}

		var ok bool
		replace_value := ""

		if replace_value, ok = context.Bot[botStruct.Name]; !ok {
			replace_value = "#!%#"
		}

		template = strings.Replace(template, orig_str, replace_value, 1)

		return template
	}
}

func GrabPart(template string, node interface{}, regexpstr string) (string, error) {
	re := regexp.MustCompile(regexpstr)

	match := re.FindString(template)
	if "" == match {
		return "", errors.New("No match")
	}

	err := xml.Unmarshal([]byte(match), node)
	if err != nil {
		return "", err
	}

	return match, nil
}
