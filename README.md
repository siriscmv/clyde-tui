# Clyde TUI

> **Warning**
> This tool uses a Discord user account, automating user accounts is against Discord's TOS. Use at your own discretion.

- A TUI app that uses Discord's Clyde AI functionality to effectively use ChatGPT for free
- This is useful because:
  - ChatGPT website has heavy bot checks, to the point where it is annoying to use even as an actual user
  - The official ChatGPT API is free only for the first 3 months

## Demo

[![asciicast](https://asciinema.org/a/HeEq8E9ku7CwS7SKdRZ6v3HJF.svg)](https://asciinema.org/a/HeEq8E9ku7CwS7SKdRZ6v3HJF)

## Usage

- Create a new Discord user account and get its token
- Create a private Discord server, enable Clyde, create a thread for Clyde and copy its ID
- Create an .env file:

```env
TOKEN="user_token_here"
CLYDE_CHANNEL_ID="channel_id_here"
```

- Run using `go run .`

## Features

- Paste content from clipboard using `@cb` in your prompt
