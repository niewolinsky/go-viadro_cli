# Viadro - Cloud-based PDF document manager - CLI
Viadro is a cloud-native HTTP API for managing PDF documents in the cloud. It utilizes S3 storage under the hood, in a sense it is a light S3 wrapper for basic CRUD operations on files. This is the CLI part of the application. For the API check: [Viadro API](https://github.com/niewolinsky/go-viadro_api)

### CLI MODE
![CLI help](https://i.imgur.com/iytclSl.png)
### TUI MODE
![TUI mode](https://i.imgur.com/14Am4JX.png)

### Features:
- Interact with Viadro API from terminal (all features supported)
- Authenticate and save credentials
- Dynamically search through list of public or user's private documents
- Merge many PDFs and upload with single command
- Grab a PDF from the web and host it on Viadro service instead
- Autocompletion, help with examples, helpful log messages

## Running:
Download executable and run or build from source. Tested on Linux, other operating systems coming soon.

## Todo:
- prevent login in if user already logged in
- folder sync command
- show documents and user list in non-interactive table when listing from CLI
- show relevant help in TUI mode
- client-side file encryption

## Stack:
- Go 1.20 + [spf13/viper](https://github.com/spf13/viper) + [spf13/cobra](https://github.com/spf13/cobra) + [charmbracelet/log](https://github.com/charmbracelet/log) + [charmbracelet/bubbletea](https://github.com/charmbracelet/bubbletea) + [charmbracelet/lipgloss](https://github.com/charmbracelet/lipgloss)

## Additional info
Application is actively developed and it is not finished yet.
