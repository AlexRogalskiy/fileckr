package main

import (
	"bufio"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/squat/fileckr/codec"
)

func main() {
	var cmdDecode = &cobra.Command{
		Use:   "decode [PNG to decode] [output file]",
		Short: "Decodes a PNG to a file",
		Long:  "decode a fileckr encoded PNG image back into a regular file.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatal("decode takes exactly two arguments")
			}
			f, err := os.Create(args[1])
			if err != nil {
				log.Fatal(err)
			}
			w := bufio.NewWriter(f)
			err = codec.DecodeFile(args[0], w)
			if err != nil {
				log.Fatal(err)
			}
			err = w.Flush()
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
		},
	}

	var cmdEncode = &cobra.Command{
		Use:   "encode [file to encode] [output PNG]",
		Short: "Encodes a file to a PNG",
		Long:  "encode a regular file into a fileckr PNG.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				log.Fatal("encode takes exactly two arguments")
			}
			f, err := os.Create(args[1])
			if err != nil {
				log.Fatal(err)
			}
			w := bufio.NewWriter(f)
			err = codec.EncodeFile(args[0], w)
			if err != nil {
				log.Fatal(err)
			}
			err = w.Flush()
			if err != nil {
				log.Fatal(err)
			}
			f.Close()
		},
	}

	var rootCommand = &cobra.Command{
		Use: "fileckr",
	}
	rootCommand.AddCommand(cmdDecode)
	rootCommand.AddCommand(cmdEncode)
	rootCommand.Execute()
}
