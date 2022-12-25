package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.design/x/clipboard"
)

func v4_bytes() [16]byte {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return *(*[16]byte)(b)
}

func nil_bytes() [16]byte {
	return *(*[16]byte)(make([]byte, 16))
}

func bytes_to_uuid(bytes [16]byte) string {
	s := hex.EncodeToString(bytes[:])
	s = s[:20] + "-" + s[20:]
	s = s[:16] + "-" + s[16:]
	s = s[:12] + "-" + s[12:]
	s = s[:8] + "-" + s[8:]
	return s
}

func generate(version string, number int) []string {
	var out []string
	for i := 0; i < number; i++ {
		switch version {
		case "4":
			out = append(out, bytes_to_uuid(v4_bytes()))
			continue
		case "nil":
			out = append(out, bytes_to_uuid(nil_bytes()))
			continue
		default:
			fmt.Println("UUID v" + version + " not supported")
			break
		}
	}
	return out
}

func main() {
	app := &cli.App{
		Name:                   "uuid-gen",
		Usage:                  "Generate UUIDs",
		UseShortOptionHandling: true,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "version",
				Aliases: []string{"v"},
				Value:   "4",
				Usage:   "Specify UUID version to generate",
			},
			&cli.IntFlag{
				Name:    "number",
				Aliases: []string{"n"},
				Value:   1,
				Usage:   "Specify number of UUIDs to generate",
			},
			&cli.BoolFlag{
				Name:    "copy",
				Aliases: []string{"c"},
				Value:   false,
				Usage:   "Copy output to clipboard",
			},
			&cli.BoolFlag{
				Name:    "silent",
				Aliases: []string{"s"},
				Value:   false,
				Usage:   "Silent (no output)",
			},
		},
		Action: func(cCtx *cli.Context) error {
			output := generate(cCtx.String("version"), cCtx.Int("number"))
			textOut := strings.Join(output, "\n")

			if cCtx.Bool("copy") {
				clipboard.Write(clipboard.FmtText, []byte(textOut))
			}

			if !cCtx.Bool("silent") {
				fmt.Println(textOut)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
