package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go/ingress_validate"
	"github.com/klev-dev/klev-api-go/logs"
	"github.com/klev-dev/klev-api-go/messages"
	"github.com/klev-dev/klev-api-go/offsets"
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

		out, err := klient.Messages.Post(cmd.Context(), logs.LogID(args[0]), t, key, value)
		return output(messages.PostOut{NextOffset: out}, err)
	}

	return cmd
}

func consume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consume log_id",
		Short: "consumes messages",
		Args:  cobra.ExactArgs(1),
	}

	offset := cmd.Flags().Int64("offset", messages.OffsetOldest, "the starting offset")
	offsetID := cmd.Flags().String("offset-id", "", "offset to get the starting consume offset")
	size := cmd.Flags().Int32("size", 10, "max messages to consume")
	poll := cmd.Flags().Duration("poll", 0, "how long to wait for new messages")
	encoding := cmd.Flags().String("encoding", "string", "how to convert message payload")
	cont := cmd.Flags().Bool("continue", false, "continue getting messages, until interrupted")

	cmd.MarkFlagsMutuallyExclusive("offset", "offset-id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		var opts []messages.ConsumeOpt
		if cmd.Flags().Changed("offset_id") {
			opts = append(opts, messages.ConsumeOffsetID(offsets.OffsetID(*offsetID)))
		} else {
			opts = append(opts, messages.ConsumeOffset(*offset))
		}
		if cmd.Flags().Changed("size") {
			opts = append(opts, messages.ConsumeLen(*size))
		}
		if cmd.Flags().Changed("poll") {
			opts = append(opts, messages.ConsumePoll(*poll))
		}
		if cmd.Flags().Changed("continue") && !cmd.Flags().Changed("poll") {
			return fmt.Errorf("continue requires polling")
		}

		switch *encoding {
		case "string":
			opts = append(opts, messages.ConsumeEncoding(messages.EncodingString))
		case "base64":
			opts = append(opts, messages.ConsumeEncoding(messages.EncodingBase64))
		default:
			return fmt.Errorf("invalid encoding: %s", *encoding)
		}

		repeat := true
		for repeat {
			next, out, err := klient.Messages.Consume(cmd.Context(), logs.LogID(args[0]), opts...)
			if err != nil {
				return output("", err)
			}

			var msgs = make([]messages.ConsumeMessageOut, len(out))
			for i, m := range out {
				msgs[i] = messages.ConsumeMessageOut{
					Offset: m.Offset,
					Time:   encodeTime(m.Time),
					Key:    encoded(m.Key, *encoding),
					Value:  encoded(m.Value, *encoding),
				}
			}
			if err := output(messages.ConsumeOut{
				NextOffset: next,
				Encoding:   *encoding,
				Messages:   msgs,
			}, err); err != nil {
				return err
			}

			repeat = *cont
			opts[0] = messages.ConsumeOffset(next)
		}

		return nil
	}

	return cmd
}

func receive() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "receive",
		Short: "receives messages from a webhook",
		Args:  cobra.NoArgs,
	}

	secret := cmd.Flags().String("secret", "", "secret to validate the payload")
	cmd.MarkFlagRequired("secret")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			msg, err := ingress_validate.Message(w, r, time.Now, *secret)
			if err != nil {
				outputValue(os.Stderr, err)
			}
			fmt.Printf("Offset: %d\n Time: %v\n Key: %s\n Value: %s\n",
				msg.Offset, msg.Time, msg.Key, msg.Value)
		})
		fmt.Println("running server at :9000")
		return http.ListenAndServe(":9000", nil)
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
