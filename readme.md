# `chat`: a CLI tool for ChatGPT

[![asciicast](https://asciinema.org/a/hhhibsaGweztuB14IPPzG2SXL.svg)](https://asciinema.org/a/hhhibsaGweztuB14IPPzG2SXL)

`chat` is a lightweight command-line interface (CLI) tool designed to streamline interaction with OpenAI's ChatGPT. Built with pure Go, this tool is operating system agnostic and requires no dependencies, making it easy to use across various platforms. With its small build size and streaming completion feature, `chat` offers a smooth chatting experience, even when ChatGPT's frontend is down.

(Description written by GPT-4, with some minor editing)

## Usage

Start an interactive chat session with ChatGPT:
```bash
chat
```
This will initiate a conversation with the default prompt. You can now interact with ChatGPT by typing your input and pressing Enter.

You can also directly type in what you want to say.

```
chat What's the abilty of sharks to sense their prey's electric fields called?
``` 

## Features

- Operating system agnostic: Compatible with any OS, thanks to being built with pure Go.
- No dependencies: No need to install additional libraries, not even glibc.
- Small build size: The typical build size is less than 4 MB.
- Streaming completion: The only CLI with streaming completion as of 17 March 2023.
- Good default prompt: Start conversations easily with a built-in default prompt.
- Interactive chatting: Engage in interactive chat sessions with ChatGPT.
- Resilient: Continues to function even if ChatGPT's CLI is down.

## Installation

### Go

If you already have [go](https://go.dev/) installed:
```bash
go install github.com/stillmatic/chat@latest
```
Make sure your go binary packages are in your `PATH` - eg add the following to `~/.zshrc` or `~/.bashrc`, etc:
```
export PATH=${PATH}:`go env GOPATH`/bin
```

### Binary
You can download binary artifacts from the [releases page](https://github.com/stillmatic/chat/releases) directly. Make sure to move the files to your path.

### From scratch

```bash
git clone https://github.com/stillmatic/chat.git
cd gpt-cli
go build .
mv chat /usr/local/bin/
```

## License

`chat` is distributed under the MIT License. See the LICENSE file for more information.

## Todo

- [ ] add a config persisted to disk (instead of setting the `OPENAI_KEY` env var)
- [ ] add flag to return a single response (instead of chat)
- [ ] add usage details / count tokens
- [ ] allow for customizing the prompt

