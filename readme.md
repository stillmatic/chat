# Incredibly simple GPT CLI

[![asciicast](https://asciinema.org/a/VMw7ccnlLEnZU7uCMojeiF76r.svg)](https://asciinema.org/a/VMw7ccnlLEnZU7uCMojeiF76r)

sometimes you want to talk to ChatGPT and the web interface is down or it's just far away. This is a simple CLI for that.

Requires `OPENAI_KEY` env var to be set to your OpenAI API key.

## Features

- Operating system agnostic
- Streaming completion
- Good default prompt
- Interactive chatting

## Installation

`go install github.com/stillmatic/gpt-cli`

or, clone the repo and run `go build -o gpt .` in the root directory.

## Todo

- [ ] add a config persisted to disk (instead of setting the `OPENAI_KEY` env var)
- [ ] add flag to return a single response (instead of chat)
- [ ] add usage details / count tokens
- [ ] allow for customizing the prompt

idk feels like a lot of effort this is fine
