# klev-cli

A cli to interact with [klev](https://klev.dev)

## Authentication

To interact with [klev](https://klev.dev) you need an authtoken. You can get one after registering at the [dashboard](https://dash.klev.dev). Pass it to the cli either via `--authtoken` flag or through `KLEV_TOKEN` environment variable. For example:

```bash
$ klev --authtoken "XXX_YYY" paths
{
  "/log": "get/delete log",
  "/logs": "list/create logs",
  "/message": "post/get message",
  "/messages": "publish/consume messages",
  "/offset": "get/set/delete offset",
  "/offsets": "list offsets",
  "/token": "get/delete token",
  "/tokens": "list/create tokens",
  "/webhook": "get/delete webhook",
  "/webhooks": "list/create webhooks"
}
```

## Basic usage

`klev` gives access to most of the functionality available through the [api](https://klev.dev/api).

```bash
$ klev 
cli to interact with klev

Usage:
  klev [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  consume     consumes messages
  help        Help about any command
  logs        interact with logs
  offsets     interact with offsets
  paths       get paths in klev; validate token
  publish     publish a message
  tokens      interact with tokens
  webhooks    interact with webhooks

Flags:
      --authtoken string   token to use for authorization
  -h, --help               help for klev

Use "klev [command] --help" for more information about a command.
```

### Publishing messages

To publish a message with values as a string use:

```bash
$ klev publish log_2IKrqtEBeYobBAM2gkuFNB6pBFL --value "hello world"
{
  "next_offset": 1
}
```

### Consuming messages

To consume messages and render them as strings use:

```bash
$ klev consume log_2IKrqtEBeYobBAM2gkuFNB6pBFL --encoding string
{
  "next_offset": 1,
  "encoding": "string",
  "messages": [
    {
      "offset": 0,
      "time": 0,
      "value": "hello world"
    }
  ]
}
```
