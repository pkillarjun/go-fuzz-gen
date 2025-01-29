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
Respond with "No" or a go-fuzz(https://go.dev/doc/security/fuzz/) test function without any explanation or how-to guide.
Evaluate the given function to check "Is it even worth fuzzing?"
Note: Usually, a good target function does parsing, decoding, deserialization, unmarshaling, etc. of a given input/inputs.
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
		openai.UserMessage(systemAI + "\n\n" + prompt),
	})

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages:    messages,
		Model:       openai.F(openai.ChatModelO1Mini),
		Temperature: openai.F(1.000000),
		TopP:        openai.F(1.000000),
	})

	if err != nil {
		fmt.Println(colorMagenta + strings.Repeat("-", 79) + "\n" + colorReset)
		fmt.Println(colorRed + "Function Name: " + funcName + colorReset)
		fmt.Printf(colorCyan+"ChatCompletion error: %v\n"+colorReset, err)
		fmt.Println(colorMagenta + strings.Repeat("-", 79) + "\n" + colorReset)
		return
	}

	if completion.Choices[0].Message.Content != "No" {
		fmt.Println(colorMagenta + strings.Repeat("-", 79) + "\n" + colorReset)
		fmt.Println(colorRed + "Function Name: " + funcName + colorReset)
		fmt.Println(colorCyan + completion.Choices[0].Message.Content + colorReset)
		fmt.Println(colorMagenta + strings.Repeat("-", 79) + "\n" + colorReset)
	}
}
