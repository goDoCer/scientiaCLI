# Dev

To start developing run the following in the root directory

1. Create a temp config file for the dev session - `cp default-config.json config.json` (ignored by git)
2. Populate the config using `scientia-cli -C config.json login`
3. Run any CLI command with the help of the `-C` flag as shown above

## Hidden Flags

`-C`, `--config`, override default config path
`--show-config`, show the config file

