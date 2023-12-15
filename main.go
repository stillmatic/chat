package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
)

type GPTCmd struct {
	Client         *openai.Client
	Messages       []openai.ChatCompletionMessage
	currentMessage *strings.Builder
	Model          string
}

type ModelCfg struct {
	APIKey  string `json:"api_key"`
	Model   string `json:"model"`
	BaseURL string `json:"base_url"`
}

type GPTCmdCfg struct {
	Configs map[string]ModelCfg `json:"configs"`
	Active  string              `json:"active"`
}

// saveConfig saves the GPTCmdCfg struct to a JSON file
func saveConfig(config GPTCmdCfg, filepath string) error {
	file, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, file, 0644)
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
func (c *GPTCmd) onData(res *openai.ChatCompletionStream) {
	for {
		msg, err := res.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			panic(err)
		}
		assistantPrinter.Print(msg.Choices[0].Delta.Content)
		c.currentMessage.WriteString(msg.Choices[0].Delta.Content)
	}
}

func (c *GPTCmd) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	req := openai.ChatCompletionRequest{
		Model:    c.Model,
		Messages: c.Messages,
	}
	c.currentMessage = &strings.Builder{}
	stream, err := c.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return err
	}
	c.onData(stream)
	c.Messages = append(c.Messages, openai.ChatCompletionMessage{
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
	c.Messages = append(c.Messages, openai.ChatCompletionMessage{
		Role:    roleUser,
		Content: input,
	})
	c.currentMessage.Reset()
	req := openai.ChatCompletionRequest{
		Model:    c.Model,
		Messages: c.Messages,
	}
	stream, err := c.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return false, err
	}
	c.onData(stream)
	c.Messages = append(c.Messages, openai.ChatCompletionMessage{
		Role:    roleAssistant,
		Content: c.currentMessage.String(),
	})
	fmt.Println("")
	return true, nil
}

func main() {
	// config path
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configPath := fmt.Sprintf("%s/.gptcmd.json", homeDir)
	gptConfig := GPTCmdCfg{
		Configs: make(map[string]ModelCfg),
		Active:  "",
	}
	//check if config exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		//	make config
		gptConfig.Configs["default"] = ModelCfg{
			APIKey:  os.Getenv("OPENAI_API_KEY"),
			Model:   openai.GPT3Dot5Turbo,
			BaseURL: "https://api.openai.com/v1",
		}
		gptConfig.Active = "default"
		err := saveConfig(gptConfig, configPath)
		if err != nil {
			panic(err)
		}
	}
	// load config
	configFile, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&gptConfig)
	if err != nil {
		panic(err)
	}
	// check if active config exists
	if _, ok := gptConfig.Configs[gptConfig.Active]; !ok {
		panic(fmt.Errorf("active config %s does not exist", gptConfig.Active))
	}
	// set active config
	activeConfig := gptConfig.Configs[gptConfig.Active]
	clientCfg := openai.DefaultConfig(activeConfig.APIKey)
	clientCfg.BaseURL = activeConfig.BaseURL
	client := openai.NewClientWithConfig(clientCfg)

	// takes _all_ args as just a string
	inputs := os.Args[1:]
	inputStr := strings.Join(inputs, " ")
	messages := []openai.ChatCompletionMessage{
		{
			Role:    roleSystem,
			Content: systemPrompt,
		},
		{
			Role:    roleUser,
			Content: inputStr,
		},
	}
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("environment variable OPENAI_API_KEY is not set")
	}

	gptCmd := &GPTCmd{
		Client:   client,
		Messages: messages,
	}
	ctx := context.Background()
	err = gptCmd.Run(ctx)
	if err != nil {
		panic(err)
	}
}
