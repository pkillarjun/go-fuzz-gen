package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

var (
	apiKey   = os.Getenv("OPENAI_API_KEY")
	systemAI = `
Respond with "No" (No is a pass condition in my script) or the go-fuzz test function. Source: https://go.dev/doc/security/fuzz/.
Evaluate the function to check if it's even worth fuzzing.
Note: Good functions are those that perform some kind of parsing, decoding, or unmarshaling (this is not always true).
Don't create fuzz tests for functions (from the project I gave you) that rely on official Golang libraries (e.g., JSON, HTTP, net, etc.).
`
)

func processGPT(funcName string, functionCode string) {
	if apiKey == "" {
		fmt.Println(colorRed + "Error: OPENAI_API_KEY environment variable not set." + colorReset)
		os.Exit(1)
	}

	prompt := fmt.Sprintf("Function name:\n%s\nFunction declaration:\n%s", funcName, functionCode)

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	ctx := context.Background()

	messages := openai.F([]openai.ChatCompletionMessageParamUnion{
		// openai.SystemMessage(systemAI),
		// openai.UserMessage(prompt),
		openai.UserMessage(systemAI + "\n\n" + prompt),
	})

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       openai.F(openai.ChatModelO1Preview),
		Temperature: openai.F(1.000000),
		TopP:        openai.F(1.000000),
	})

	if err != nil {
		fmt.Println(colorRed + "Function Name: " + funcName + colorReset)
		fmt.Printf(colorCyan+"ChatCompletion error: %v\n"+colorReset, err)
		return
	}

	if completion.Choices[0].Message.Content != "No" {
		fmt.Println(colorCyan + completion.Choices[0].Message.Content + colorReset)
	}
	fmt.Println(colorMagenta + strings.Repeat("-", 79) + "\n" + colorReset)
}
