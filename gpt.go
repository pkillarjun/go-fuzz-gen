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
Respond with "No (pass condition in my script)" or the go-fuzz test function (I need only the function).
You have to write a go-fuzz test for the given function.
Make sure to evaluate the function, to check even if it's worth fuzzing.
Usually, I find internal APIs of a program (target under test) are sometimes good and sometimes useless.
Note: I think the good functions are when they do some kind of parsing, decoding, unmarshaling.
Request: For now, you can ignore the functions that heavily rely on Go-lang official libraries (for example, JSON, HTTP, and others).
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

	fmt.Println(colorYellow + strings.Repeat("+", 79) + colorReset)
	fmt.Println(colorGreen + "Function Name: " + funcName + colorReset)
	fmt.Println(colorCyan + completion.Choices[0].Message.Content + colorReset)
	fmt.Println(colorMagenta + strings.Repeat("-", 79) + "\n" + colorReset)
}
