# ChatGPT CLI

[![asciicast](https://asciinema.org/a/VMw7ccnlLEnZU7uCMojeiF76r.svg)](https://asciinema.org/a/VMw7ccnlLEnZU7uCMojeiF76r)

Sometimes you want to talk to ChatGPT and the web interface is down or it's just far away. This is a simple CLI for that, with streaming completion.

Requires `OPENAI_KEY` env var to be set to your OpenAI API key.

## Features

- Operating system agnostic with minimal dependencies
  - currently Go, should be _none_
- Streaming completion (only CLI as of 17 March 2023 with this) 
- Good default prompt
- Interactive chatting

## Installation

`go install github.com/stillmatic/gpt-cli@latest`

or, clone the repo and run `go build -o gpt .` in the root directory.

## Todo

- [ ] add builds in CI for different OS
- [ ] add a config persisted to disk (instead of setting the `OPENAI_KEY` env var)
- [ ] add flag to return a single response (instead of chat)
- [ ] add usage details / count tokens
- [ ] allow for customizing the prompt

idk feels like a lot of effort this is fine
