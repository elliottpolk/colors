package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
	cli "github.com/urfave/cli/v2"
)

type rgb struct {
	r, g, b int

	format string
}

const (
	allFmt = "all"
	rgbFmt = "rgb"
	hexFmt = "hex"
)

var (
	version  string
	compiled string = fmt.Sprint(time.Now().Unix())
	githash  string
)

func (c *rgb) parseHEX(hex string) error {
	if len(hex) < 3 || len(hex) > 6 {
		return errors.New("invalid hex value specified")
	}

	// expand from 3 to 6
	if len(hex) == 3 {
		hex = fmt.Sprintf("%s%s%s%s%s%s", hex[0:1], hex[0:1], hex[1:2], hex[1:2], hex[2:], hex[2:])
	}

	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return err
	}

	g, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return err
	}

	b, err := strconv.ParseInt(hex[4:], 16, 64)
	if err != nil {
		return err
	}

	c.r = int(r)
	c.g = int(g)
	c.b = int(b)
	c.format = hexFmt

	return nil
}

func (c *rgb) invert() *rgb {
	return &rgb{
		r: 255 - c.r,
		g: 255 - c.g,
		b: 255 - c.b,
	}
}

func (c *rgb) _rgb() string {
	return fmt.Sprintf("rgb(%d,%d,%d)", c.r, c.g, c.b)
}

func (c *rgb) _hex() string {
	return fmt.Sprintf("%02x%02x%02x", c.r, c.g, c.b)
}

func (c *rgb) valid() bool {
	if c.r < 0 || c.r > 255 {
		return false
	}

	if c.g < 0 || c.g > 255 {
		return false
	}

	if c.b < 0 || c.b > 255 {
		return false
	}

	return true
}

func (c *rgb) String() string {
	var (
		i = c.invert()

		// need the foreground to be the opposite color in hopes that it's legible
		fg = color.RGB(uint8(i.r), uint8(i.g), uint8(i.b))
		bg = color.RGB(uint8(c.r), uint8(c.g), uint8(c.b))

		out string
	)

	switch c.format {
	case rgbFmt:
		out = c._rgb()
	case hexFmt:
		out = c._hex()
	default: // default to all
		out = fmt.Sprintf("%s | %s", c._rgb(), c._hex())
	}

	return color.NewRGBStyle(fg, bg).Sprintf(" %s ", out)
}

func main() {
	ct, err := strconv.ParseInt(compiled, 0, 0)
	if err != nil {
		panic(err)
	}

	var (
		f       string
		r, g, b int

		fmts = []string{
			allFmt,
			rgbFmt,
			hexFmt,
		}
	)

	app := cli.App{
		Name:      "hc",
		Copyright: "Copyright © 2021",
		Version:   fmt.Sprintf("%s | compiled %s | commit %s", version, time.Unix(ct, -1).Format(time.RFC3339), githash),
		Compiled:  time.Unix(ct, -1),
		Usage:     "hc [arguments...] [hex_value]",
		//UsageText: "",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "format",
				Aliases:     []string{"fmt", "f"},
				Usage:       fmt.Sprintf("display with one or all formats (e.g. %+v)", fmts),
				Destination: &f,
			},
			&cli.IntFlag{
				Name:        "red",
				Aliases:     []string{"r"},
				Value:       -1,
				Destination: &r,
			},
			&cli.IntFlag{
				Name:        "green",
				Aliases:     []string{"g"},
				Value:       -1,
				Destination: &g,
			},
			&cli.IntFlag{
				Name:        "blue",
				Aliases:     []string{"b"},
				Value:       -1,
				Destination: &b,
			},
		},
		Action: func(ctx *cli.Context) error {
			clr := &rgb{
				r:      r,
				g:      g,
				b:      b,
				format: rgbFmt,
			}

			if !clr.valid() {
				h := ctx.Args().First()
				if len(h) < 3 {
					if err := cli.ShowAppHelp(ctx); err != nil {
						return cli.Exit(err, 1)
					}

					return nil
				}

				if err := clr.parseHEX(h); err != nil {
					return cli.Exit(err, 1)
				}
			}

			f = strings.ToLower(f)
			if len(f) > 0 {
				for _, val := range fmts {
					if val == f {
						clr.format = val
						break
					}
				}
			}

			fmt.Print(clr)
			return nil
		},
	}

	app.Run(os.Args)
}
