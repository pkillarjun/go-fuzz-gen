package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
)

var APIKEY = os.Getenv("OPENAI_API_KEY")
var SYSTEM_AI = `
Respond with Yes and No.
Check if the given function/functions can be fuzz tested using go test fuzz.
When evaluating functions, check if they implement parsing, decoding, 
decrypting, unmarshaling, etc. because they are good targets for fuzz testing.
Usually, internal APIs are a waste of time, so keep that in mind.
`

func run_gpt(func_name string, functions string) {

	gptinput := "Function name:\n" + func_name + "\nAll Function declaration:\n" + functions

	client := openai.NewClient(APIKEY)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "gpt-4o",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: SYSTEM_AI,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: gptinput,
				},
			},
			Temperature: 0.5,
		},
	)

	if err != nil {
		fmt.Println(red_color + "Function Name: " + func_name + reset_color)
		fmt.Printf(cyan_color+"ChatCompletion error: %v\n"+reset_color, err)
		return
	}

	if resp.Choices[0].Message.Content == "No" {
		return
	}

	fmt.Println(yellow_color + strings.Repeat("+", 79) + "" + reset_color)
	fmt.Println(green_color + "Function Name: " + func_name + reset_color)
	fmt.Println(cyan_color + resp.Choices[0].Message.Content + reset_color)
	fmt.Println(magenta_color + strings.Repeat("-", 79) + "\n" + reset_color)
}
