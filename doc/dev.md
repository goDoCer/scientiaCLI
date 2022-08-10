## Development

To start development install Go from [here](https://go.dev/doc/install).

1. In the source directory run: `go mod tidy`.
2. (Recommended) Install [`cobra-cli`](https://github.com/spf13/cobra-cli) - `go install github.com/spf13/cobra-cli@latest`.

### Hidden Flags

`go run main.go -C` - define custom path to the config file

`go run main.go -s` - show the config file variables
