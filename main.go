package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/fatih/color"
)

type GPTCmd struct {
	Client   gpt3.Client                        
	Messages []gpt3.ChatCompletionRequestMessage 
}

//go:embed prompt.txt
var systemPrompt string

const (
	roleSystem    = "system"
	roleUser      = "user"
	roleAssistant = "assistant"
)

var assistantPrinter = color.New(color.FgYellow)

// onData handles the streaming output
// here, we are simply printing the output to the console
func onData(res *gpt3.ChatCompletionStreamResponse) {
	assistantPrinter.Print(res.Choices[0].Delta.Content)
}

func (c *GPTCmd) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	req := gpt3.ChatCompletionRequest{
		Model:    gpt3.GPT3Dot5Turbo,
		Messages: c.Messages,
	}
	err := c.Client.ChatCompletionStream(ctx, req, onData)
	if err != nil {
		return err
	}
	fmt.Println("")
	for {
		fmt.Print(">  ")
		buffer := bufio.NewReader(os.Stdin)
		line, err := buffer.ReadString('\n')
		if err != nil {
			return err
		}
		input := strings.TrimSpace(line)
		if input == "" {
			fmt.Println("Please enter a response. (Type 'quit' to exit.)")
			continue
		}
		if input == "quit" || input == "exit" {
			print("Thank you, come again!")
			return nil
		}
		c.Messages = append(c.Messages, gpt3.ChatCompletionRequestMessage{
			Role:    roleUser,
			Content: input,
		})
		req := gpt3.ChatCompletionRequest{
			Model:    gpt3.GPT3Dot5Turbo,
			Messages: c.Messages,
		}
		err = c.Client.ChatCompletionStream(ctx, req, onData)
		if err != nil {
			return err
		}
		fmt.Println("")
	}
}

func main() {
	// takes _all_ args as just a string
	inputs := os.Args[1:]
	inputStr := strings.Join(inputs, " ")
	messages := []gpt3.ChatCompletionRequestMessage{
		{
			Role:    roleSystem,
			Content: systemPrompt,
		},
		{
			Role:    roleUser,
			Content: inputStr,
		},
	}
	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		panic("environment variable OPENAI_KEY is not set")
	}
	client := gpt3.NewClient(apiKey)
	gptCmd := &GPTCmd{
		Client:   client,
		Messages: messages,
	}
	ctx := context.Background()
	err := gptCmd.Run(ctx)
	if err != nil {
		panic(err)
	}
}
