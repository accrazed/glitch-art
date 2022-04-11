package main

import (
	"image/png"
	"log"
	"os"

	ps "github.com/accrazed/glitch-art/src/pixelsort"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:        "Glitch Art",
		Description: "glitch your images!",
		Commands: []*cli.Command{
			{
				Name:    "pixelsort",
				Aliases: []string{"ps"},
				Usage:   "Runs a pixelsort algorithm on your image",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Usage:    "image to pixel sort",
						Required: true,
					},
					&cli.Int64Flag{
						Name:     "seed",
						Aliases:  []string{"s"},
						Usage:    "seed to base random shit on",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "threshold",
						Aliases:  []string{"t"},
						Usage:    "threshold to use for saturation sorting",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "direction",
						Aliases:  []string{"d"},
						Usage:    "direction to sort in",
						Value:    "horizontal",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "filename to output as",
						Value:    "output",
						Required: false,
					},
				},
				Action: func(ctx *cli.Context) error {
					var sortDir ps.SortDir = ps.Horizontal
					if ctx.String("direction") == "vertical" {
						sortDir = ps.Vertical
					}

					pixSort := ps.Must(ps.New(ctx.String("path"),
						ps.WithSortDir(sortDir),
						ps.WithSeed(ctx.Int64("seed")),
						ps.WithThreshold(ctx.Int("threshold")),
					))

					img := pixSort.Sort()

					f, err := os.Create(ctx.String("output") + ".png")
					if err != nil {
						return err
					}
					err = png.Encode(f, img)
					if err != nil {
						return err
					}
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
