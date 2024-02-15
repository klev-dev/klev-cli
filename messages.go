package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/klev-dev/klev-api-go"
	"github.com/klev-dev/klev-api-go/ingress_validate"
)

func publish() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "publish <log-id>",
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
		id, err := klev.ParseLogID(args[0])
		if err != nil {
			return outputErr(err)
		}

		var t time.Time
		var key, value []byte

		if cmd.Flags().Changed("time") {
			t = time.UnixMicro(*tflag)
		}

		if cmd.Flags().Changed("key") {
			key = []byte(*keyString)
		} else if cmd.Flags().Changed("key-file") {
			if b, err := os.ReadFile(*keyFile); err != nil {
				return outputErr(err)
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
				return outputErr(err)
			} else {
				value = b
			}
		} else if valueBase64 != nil {
			value = *valueBase64
		}

		out, err := klient.Messages.Post(cmd.Context(), id, t, key, value)
		return output(klev.PostOut{NextOffset: out}, err)
	}

	return cmd
}

func consume() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "consume <log-id>",
		Short: "consumes messages",
		Args:  cobra.ExactArgs(1),
	}

	offset := cmd.Flags().Int64("offset", klev.OffsetOldest, "the starting offset")
	offsetID := cmd.Flags().String("offset-id", "", "offset to get the starting consume offset")
	size := cmd.Flags().Int32("size", 10, "max messages to consume")
	poll := cmd.Flags().Duration("poll", 0, "how long to wait for new messages")
	encoding := cmd.Flags().String("encoding", "string", "how to convert message payload")
	cont := cmd.Flags().Bool("continue", false, "continue getting messages, until interrupted")

	cmd.MarkFlagsMutuallyExclusive("offset", "offset-id")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseLogID(args[0])
		if err != nil {
			return outputErr(err)
		}

		var opts []klev.ConsumeOpt
		if cmd.Flags().Changed("offset_id") {
			if offsetID, err := klev.ParseOffsetID(*offsetID); err != nil {
				return outputErr(err)
			} else {
				opts = append(opts, klev.ConsumeOffsetID(offsetID))
			}
		} else {
			opts = append(opts, klev.ConsumeOffset(*offset))
		}
		if cmd.Flags().Changed("size") {
			opts = append(opts, klev.ConsumeLen(*size))
		}
		if cmd.Flags().Changed("poll") {
			opts = append(opts, klev.ConsumePoll(*poll))
		}
		if cmd.Flags().Changed("continue") && !cmd.Flags().Changed("poll") {
			return fmt.Errorf("continue requires polling")
		}

		var coder = klev.MessageEncodingString
		if cmd.Flags().Changed("encoding") {
			encoding, err := klev.ParseMessageEncoding(*encoding)
			if err != nil {
				return outputErr(err)
			}
			coder = encoding
			opts = append(opts, klev.ConsumeEncoding(encoding))
		}

		repeat := true
		for repeat {
			next, out, err := klient.Messages.Consume(cmd.Context(), id, opts...)
			if err != nil {
				return output("", err)
			}

			var msgs = make([]klev.ConsumeMessageOut, len(out))
			for i, m := range out {
				msgs[i] = klev.ConsumeMessageOut{
					Offset: m.Offset,
					Time:   coder.EncodeTime(m.Time),
					Key:    coder.EncodeData(m.Key),
					Value:  coder.EncodeData(m.Value),
				}
			}
			if err := output(klev.ConsumeOut{
				NextOffset: next,
				Encoding:   coder,
				Messages:   msgs,
			}, err); err != nil {
				return outputErr(err)
			}

			repeat = *cont
			opts[0] = klev.ConsumeOffset(next)
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

func cleanup() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cleanup",
		Short: "cleanup messages from a log",
		Args:  cobra.ExactArgs(1),
	}

	var in klev.CleanupIn

	cmd.Flags().Int64Var(&in.TrimSeconds, "trim-seconds", 0, "age of the log to trim")
	cmd.Flags().Int64Var(&in.TrimSize, "trim-size", 0, "size of the log to trim")
	cmd.Flags().Int64Var(&in.TrimCount, "trim-count", 0, "count of message in the log to trim")
	cmd.Flags().Int64Var(&in.CompactSeconds, "compact-seconds", 0, "age of the log to compact")
	cmd.Flags().Int64Var(&in.ExpireSeconds, "expire-seconds", 0, "age of the log to expire")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		id, err := klev.ParseLogID(args[0])
		if err != nil {
			return outputErr(err)
		}

		out, err := klient.Messages.CleanupRaw(cmd.Context(), id, in)
		return output(out, err)
	}

	return cmd
}
