package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	_ "embed"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/fatih/color"
)

type GPTCmd struct {
	Client         gpt3.Client
	Messages       []gpt3.ChatCompletionRequestMessage
	currentMessage *strings.Builder
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
// and appending it to the current message
func (c *GPTCmd) onData(res *gpt3.ChatCompletionStreamResponse) {
	assistantPrinter.Print(res.Choices[0].Delta.Content)
	c.currentMessage.WriteString(res.Choices[0].Delta.Content)
}

func (c *GPTCmd) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	req := gpt3.ChatCompletionRequest{
		Model:    gpt3.GPT3Dot5Turbo,
		Messages: c.Messages,
	}
	c.currentMessage = &strings.Builder{}
	err := c.Client.ChatCompletionStream(ctx, req, c.onData)
	c.Messages = append(c.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    roleAssistant,
		Content: c.currentMessage.String(),
	})
	if err != nil {
		return err
	}
	fmt.Println("")
	for {
		shouldContinue, err := c.runSingle(ctx)
		if err != nil {
			return err
		}
		if !shouldContinue {
			break
		}
	}
	return nil
}

func (c *GPTCmd) runSingle(ctx context.Context) (shouldContinue bool, err error) {
	fmt.Print(">  ")
	buffer := bufio.NewReader(os.Stdin)
	line, err := buffer.ReadString('\n')
	if err != nil {
		return false, err
	}
	input := strings.TrimSpace(line)
	if input == "" {
		fmt.Println("Please enter a response. (Type 'quit' to exit.)")
		return true, nil
	}
	if input == "quit" || input == "exit" {
		print("Thank you, come again!")
		return false, nil
	}
	c.Messages = append(c.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    roleUser,
		Content: input,
	})
	c.currentMessage.Reset()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	req := gpt3.ChatCompletionRequest{
		Model:    gpt3.GPT3Dot5Turbo,
		Messages: c.Messages,
	}
	err = c.Client.ChatCompletionStream(ctx, req, c.onData)
	if err != nil {
		return false, err
	}
	c.Messages = append(c.Messages, gpt3.ChatCompletionRequestMessage{
		Role:    roleAssistant,
		Content: c.currentMessage.String(),
	})
	fmt.Println("")
	return true, nil
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
