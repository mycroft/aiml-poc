package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func makeRegexp(content string) *regexp.Regexp {
	content = strings.TrimSpace(content)

	re := regexp.MustCompile("(\\?|\\(|\\)|\\[|\\]|\\||\\+)")
	content = re.ReplaceAllString(content, "\\$1")

	content = strings.Replace(content, "*", "(.*)", -1)
	content = strings.Replace(content, "_", ".*", -1)

	content = fmt.Sprintf("(?iU)^%s$", content)

	return regexp.MustCompile(content)
}

func (context *Context) run(input string) (string, error) {
	categories := context.aimlRoot.Categories

	// Select current topic
	if context.Topic != "" {
		for _, topic := range context.aimlRoot.Topics {
			if topic.Name == context.Topic {
				categories = topic.Categories
				break
			}
		}
	}

	for _, category := range categories {
		if category.Pattern.regexp == nil {
			category.Pattern.regexp = makeRegexp(category.Pattern.Content)
		}

		if category.That.Content != "" && category.That.Content != context.LastSent {
			continue
		}

		is_matching := category.Pattern.regexp.Match([]byte(input))
		if false == is_matching {
			continue
		}

		matches := category.Pattern.regexp.FindAllStringSubmatch(input, -1)
		if len(matches) >= 0 {
			template, is_srai := context.Process(category.Template.Content, matches[0][1:])
			if is_srai {
				return context.run(template)
			}

			context.LastSent = template

			return context.LastSent, nil
		}

	}

	return "", errors.New(fmt.Sprintf("No match for '%s'", input))
}
