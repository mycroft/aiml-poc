package main

import (
	"fmt"
	"testing"
)

type MatcherSample struct {
	Input  string
	Output string
}

func LoadContext(t *testing.T, filename string) (*AIMLRoot, error) {
	aimlRoot, err := Parse(filename)
	if err != nil {
		t.Error("Could not parse", filename, ":", err)
	}
	return aimlRoot, err
}

func RunTest(t *testing.T, context Context, MatcherSamples []MatcherSample) {
	for _, sample := range MatcherSamples {
		output, err := context.run(sample.Input)
		if err != nil {
			t.Error(fmt.Sprintf("Error while testing '%s': %s", sample.Input, err))
		}

		if sample.Output != output {
			t.Error(fmt.Sprintf(
				"Error while testing '%s': '%s' does not match expected '%s'.",
				sample.Input,
				output,
				sample.Output,
			))
		}
	}
}

func TestMatcher(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test.aiml")
	context.Memory = make(map[string]string)

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "Hi",
			Output: "Hi!",
		},
		MatcherSample{
			Input:  "My name is abc",
			Output: "abc",
		},
		MatcherSample{
			Input:  "My name is abc and I'm 123 years old",
			Output: "abc / 123",
		},
		MatcherSample{
			Input:  "I'm 111 years old and my name is edc",
			Output: "edc / 111",
		},
	}

	RunTest(t, context, MatcherSamples)
}

func TestMemory(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test.aiml")
	context.Memory = make(map[string]string)

	input1 := "Record 1 user patrick"
	expected_output1 := "User patrick recorded."
	output1, err := context.run(input1)
	if err != nil {
		t.Error(fmt.Sprintf("Error while testing '%s': %s", input1, err))
	}

	if expected_output1 != output1 {
		t.Error(fmt.Sprintf("'%s' is not matching expected output '%s'", output1, expected_output1))
	}

	val, ok := context.Memory["user"]
	if !ok {
		t.Error("Key 'user' should exist")
	}

	if val != "patrick" {
		t.Error("Invalid value for 'user' key in memory:", val)
	}

	input2 := "Record 2 user bob with val 123"
	expected_output2 := "Val 123 recorded for user bob"

	output2, err := context.run(input2)
	if err != nil {
		t.Error(fmt.Sprintf("Error while testing '%s': %s", input2, err))
	}

	if expected_output2 != output2 {
		t.Error(fmt.Sprintf("'%s' is not matching expected output '%s'", output2, expected_output2))
	}

	val, ok = context.Memory["user"]
	if !ok {
		t.Error("Key 'user' should exist")
	}

	if val != "bob" {
		t.Error("Invalid value for 'user' key in memory:", val)
	}

	input3 := "What is my name ?"
	expected_output3 := "Your name is 'bob'."

	output3, err := context.run(input3)
	if err != nil {
		t.Error(fmt.Sprintf("Error while testing '%s': %s", input3, err))
	}

	if expected_output3 != output3 {
		t.Error(fmt.Sprintf("'%s' is not matching expected output '%s'", output3, expected_output3))
	}

}

func TestThat(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test.aiml")
	context.Memory = make(map[string]string)

	input1 := "Test THAT"
	expected_output1 := "Repeat yourself"
	output1, err := context.run(input1)
	if err != nil {
		t.Error(fmt.Sprintf("Error while testing '%s': %s", input1, err))
	}

	if expected_output1 != output1 {
		t.Error(fmt.Sprintf("'%s' is not matching expected output '%s'", output1, expected_output1))
	}

	input2 := "Test THAT"
	expected_output2 := "Nailed it"
	output2, err := context.run(input2)
	if err != nil {
		t.Error(fmt.Sprintf("Error while testing '%s': %s", input2, err))
	}

	if expected_output2 != output2 {
		t.Error(fmt.Sprintf("'%s' is not matching expected output '%s'", output2, expected_output2))
	}
}

func TestThink(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test.aiml")
	context.Memory = make(map[string]string)

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "Test THINK 'My name is GameOver'",
			Output: "Hello",
		},
		MatcherSample{
			Input:  "Test THINK 'What is my name ?'",
			Output: "Your name is 'GameOver'.",
		},
	}

	RunTest(t, context, MatcherSamples)
}

func TestSRAI(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test_srai.aiml")
	context.Memory = make(map[string]string)

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "DO YOU KNOW WHO ALBERT EINSTEIN IS",
			Output: "Albert Einstein was a German physicist.",
		},
		MatcherSample{
			Input:  "WHO IS ALBERT EINSTEIN",
			Output: "Albert Einstein was a German physicist.",
		},
		MatcherSample{
			Input:  "DO YOU KNOW WHO Isaac NEWTON IS",
			Output: "Isaac Newton was a English physicist and mathematician.",
		},
		MatcherSample{
			Input:  "Bye",
			Output: "Good Bye!",
		},
		MatcherSample{
			Input:  "Bye Alice!",
			Output: "Good Bye!",
		},
		MatcherSample{
			Input:  "Factory",
			Output: "Development Center!",
		},
		MatcherSample{
			Input:  "Industry",
			Output: "Development Center!",
		},

		MatcherSample{
			Input:  "I love going to school daily.",
			Output: "School is an important institution in a child's life.",
		},
		MatcherSample{
			Input:  "I like my school.",
			Output: "School is an important institution in a child's life.",
		},
	}

	RunTest(t, context, MatcherSamples)
}

func TestCondition(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test.aiml")
	context.Memory = make(map[string]string)

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "HOW ARE YOU FEELING TODAY",
			Output: "I am happy!",
		},
	}

	RunTest(t, context, MatcherSamples)
}

func TestTemplateThat(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test_template_that.aiml")
	context.Memory = make(map[string]string)

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "Hello",
			Output: "Hello. What's up ?",
		},
		MatcherSample{
			Input:  "Repeat ?",
			Output: "Hello. What's up ?",
		},
		MatcherSample{
			Input:  "Bye",
			Output: "Bye...",
		},
		MatcherSample{
			Input:  "Repeat ?",
			Output: "Bye...",
		},
	}

	RunTest(t, context, MatcherSamples)
}

func TestInput(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test_input.aiml")
	context.Memory = make(map[string]string)

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "Hello",
			Output: "You just said 'Hello'",
		},
		MatcherSample{
			Input:  "Repeat ?",
			Output: "You just said 'Repeat ?'",
		},
	}

	RunTest(t, context, MatcherSamples)
}

func TestThatStar(t *testing.T) {
	var context Context
	context.aimlRoot, _ = LoadContext(t, "tests/test_thatstar.aiml")
	context.Memory = make(map[string]string)

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "Hello",
			Output: "What is your name ?",
		},
		MatcherSample{
			Input:  "My name is IAMROBOT",
			Output: "Hello IAMROBOT!",
		},
		MatcherSample{
			Input:  "Well.... (1)",
			Output: "IAMROBOT, how are you ?",
		},
		MatcherSample{
			Input:  "Well.... (2)",
			Output: "I forgot your name...",
		},
	}

	RunTest(t, context, MatcherSamples)
}

func TestBot(t *testing.T) {
	context := InitContext()
	context.aimlRoot, _ = LoadContext(t, "tests/test_bot.aiml")

	MatcherSamples := []MatcherSample{
		MatcherSample{
			Input:  "Your name ?",
			Output: "My name is StupidBot",
		},
		MatcherSample{
			Input:  "Your age ?",
			Output: "My age is #!%#",
		},
	}

	RunTest(t, context, MatcherSamples)
}
