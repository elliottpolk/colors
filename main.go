package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	cli "github.com/urfave/cli/v2"
)

var (
	version  string
	compiled string = fmt.Sprint(time.Now().Unix())
	githash  string
)

func main() {
	ct, err := strconv.ParseInt(compiled, 0, 0)
	if err != nil {
		panic(err)
	}

	app := cli.App{
		Name:      "colors",
		Copyright: "Copyright © 2021",
		Version:   fmt.Sprintf("%s | compiled %s | commit %s", version, time.Unix(ct, -1).Format(time.RFC3339), githash),
		Compiled:  time.Unix(ct, -1),
		Usage:     "TODO...",
		UsageText: "colors [options] [arguments...]",
		Flags:     []cli.Flag{},
		Action: func(ctx *cli.Context) error {
			bg := ctx.Args().First()
			if strings.HasPrefix(bg, "#") {
				bg = bg[1:]
			}

			if l := len(bg); l < 3 || l > 6 {
				return cli.Exit("invalid hex color value", 1)
			}

			// expand
			if len(bg) == 3 {
				bg = fmt.Sprintf("%s%s%s%s%s%s", bg[0:1], bg[0:1], bg[1:2], bg[1:2], bg[2:], bg[2:])
			}

			r, err := strconv.ParseInt(bg[0:2], 16, 64)
			if err != nil {
				return cli.Exit(err, 1)
			}

			g, err := strconv.ParseInt(bg[2:4], 16, 64)
			if err != nil {
				return cli.Exit(err, 1)
			}

			b, err := strconv.ParseInt(bg[4:], 16, 64)
			if err != nil {
				return cli.Exit(err, 1)
			}

			// invert
			r = 255 - r
			g = 255 - g
			b = 255 - b

			color.HEXStyle(fmt.Sprintf("%x%x%x", r, g, b), bg).Println(" ", bg, " ")

			return nil
		},
	}

	app.Run(os.Args)
}
