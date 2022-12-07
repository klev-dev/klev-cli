package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	api "github.com/klev-dev/klev-api-go"
)

func publish() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish",
		Short: "publish a message",
	}

	tflag := cmd.Flags().Int64("time", 0, "unix time")

	keyString := cmd.Flags().String("key", "", "key as a string value")
	keyFile := cmd.Flags().String("key-file", "", "a file to read the key from")
	keyBase64 := cmd.Flags().BytesBase64("key-bytes", nil, "key as a base64 encoded bytes")
	valueString := cmd.Flags().String("value", "", "value as a string value")
	valueFile := cmd.Flags().String("value-file", "", "a file to read the value from")
	valueBase64 := cmd.Flags().BytesBase64("value-bytes", nil, "value as a base64 encoded bytes")

	cmd.MarkFlagsMutuallyExclusive("key", "key-file", "key-bytes")
	cmd.MarkFlagsMutuallyExclusive("value", "value-file", "value-bytes")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var t time.Time
		var key, value []byte

		if cmd.Flags().Changed("time") {
			t = time.UnixMicro(*tflag)
		}

		if cmd.Flags().Changed("key") {
			key = []byte(*keyString)
		} else if cmd.Flags().Changed("key-file") {
			if b, err := os.ReadFile(*keyFile); err != nil {
				return err
			} else {
				key = b
			}
		} else {
			key = *keyBase64
		}

		if cmd.Flags().Changed("value") {
			value = []byte(*valueString)
		} else if cmd.Flags().Changed("value-file") {
			if b, err := os.ReadFile(*valueFile); err != nil {
				return err
			} else {
				value = b
			}
		} else if valueBase64 != nil {
			value = *valueBase64
		}

		out, err := klient.Post(cmd.Context(), api.LogID(args[0]), t, key, value)
		return output(api.PostOut{NextOffset: out}, err)
	}

	return cmd
}

func consume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consume",
		Short: "consumes messages",
	}

	encoding := cmd.Flags().String("encoding", "base64", "how to convert message payload")
	offset := cmd.Flags().Int64("offset", -1, "the starting offset")
	size := cmd.Flags().Int32("size", 10, "max messages to consume")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !(*encoding == "string" || *encoding == "base64") {
			return fmt.Errorf("invalid encoding: %s", *encoding)
		}

		next, out, err := klient.Consume(cmd.Context(), api.LogID(args[0]), *offset, *size)
		if err != nil {
			return output("", err)
		}

		var msgs = make([]api.ConsumeMessageOut, len(out))
		for i, m := range out {
			msgs[i] = api.ConsumeMessageOut{
				Offset: m.Offset,
				Time:   encodeTime(m.Time),
				Key:    encoded(m.Key, *encoding),
				Value:  encoded(m.Value, *encoding),
			}
		}
		return output(api.ConsumeOut{
			NextOffset: next,
			Encoding:   *encoding,
			Messages:   msgs,
		}, err)
	}

	return cmd
}

func encodeTime(t time.Time) int64 {
	return t.UnixMicro()
}

func encoded(b []byte, encoding string) *string {
	if b == nil {
		return nil
	}
	var s string
	switch encoding {
	case "base64":
		s = base64.StdEncoding.EncodeToString(b)
	case "string":
		s = string(b)
	}
	return &s
}
