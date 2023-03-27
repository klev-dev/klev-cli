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
		Use:   "publish log_id",
		Short: "publish a message",
		Args:  cobra.ExactArgs(1),
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
		Use:   "consume log_id",
		Short: "consumes messages",
		Args:  cobra.ExactArgs(1),
	}

	offset := cmd.Flags().Int64("offset", api.OffsetOldest, "the starting offset")
	offsetID := cmd.Flags().String("offset-id", "", "offset to get the starting consume offset")
	size := cmd.Flags().Int32("size", 10, "max messages to consume")
	poll := cmd.Flags().Duration("poll", 0, "how long to wait for new messages")
	encoding := cmd.Flags().String("encoding", "string", "how to convert message payload")

	cmd.MarkFlagsMutuallyExclusive("offset", "offset-id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if !(*encoding == "string" || *encoding == "base64") {
			return fmt.Errorf("invalid encoding: %s", *encoding)
		}

		var opts []api.ConsumeOpt
		opts = append(opts, api.ConsumeOffset(*offset))
		if cmd.Flags().Changed("offset_id") {
			opts = append(opts, api.ConsumeOffsetID(api.OffsetID(*offsetID)))
		}
		if cmd.Flags().Changed("size") {
			opts = append(opts, api.ConsumeLen(*size))
		}
		if cmd.Flags().Changed("poll") {
			opts = append(opts, api.ConsumePoll(*poll))
		}
		opts = append(opts, api.ConsumeEncoding(*encoding))

		next, out, err := klient.Consume(cmd.Context(), api.LogID(args[0]), opts...)
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
