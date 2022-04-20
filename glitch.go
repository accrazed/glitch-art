package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	cs "github.com/accrazed/glitch-art/src/channelshift"
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
					&cli.BoolFlag{
						Name:     "invert",
						Aliases:  []string{"i"},
						Usage:    "invert sorting algorithm direction",
						Value:    false,
						Required: false,
					},
					&cli.StringFlag{
						Name:     "direction",
						Aliases:  []string{"d"},
						Usage:    "direction to sort in",
						Value:    "vertical",
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
					var sortDir ps.SortDir = ps.Vertical
					if strings.ToLower(ctx.String("direction")) == "horizontal" {
						sortDir = ps.Horizontal
					}

					pixSort := ps.Must(ps.New(ctx.String("path"),
						ps.WithSortDir(sortDir),
						ps.WithSeed(ctx.Int64("seed")),
						ps.WithThreshold(ctx.Int("threshold")),
						ps.WithInvert(ctx.Bool("invert")),
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
			{
				Name:    "channelshift",
				Aliases: []string{"cs"},
				Usage:   "Translates color channels on your image",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "path",
						Aliases:  []string{"p"},
						Usage:    "image to pixel sort",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "filename to output as",
						Value:    "output",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "red",
						Aliases:  []string{"r"},
						Usage:    "x,y pair of red translate. e.g. \"15,0\"",
						Value:    "0,0",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "green",
						Aliases:  []string{"g"},
						Usage:    "x,y pair of green translate. e.g. \"15,0\"",
						Value:    "0,0",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "blue",
						Aliases:  []string{"b"},
						Usage:    "x,y pair of blue translate. e.g. \"15,0\"",
						Value:    "0,0",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "alpha",
						Aliases:  []string{"a"},
						Usage:    "x,y pair of alpha translate. e.g. \"15,0\"",
						Value:    "0,0",
						Required: false,
					},
					&cli.IntFlag{
						Name:     "chunk",
						Aliases:  []string{"c"},
						Usage:    "How many pixels to chunk volatility by",
						Required: false,
					},
					&cli.IntFlag{
						Name:     "volatility",
						Aliases:  []string{"v"},
						Usage:    "How strongly to shift pixels per chunk",
						Required: false,
					},
					&cli.Int64Flag{
						Name:    "seed",
						Aliases: []string{"s"},
						Usage:   "Seed to base chunking/volatility off of",
					},
					&cli.IntFlag{
						Name:    "animate",
						Aliases: []string{"gif"},
						Usage:   "Animate to several frames :)",
						Value:   1,
					},
					&cli.IntFlag{
						Name:    "anivolatility",
						Aliases: []string{"aniv"},
						Usage:   "How volatile to adjust the chunk value in every iteration",
						Value:   0,
					},
				},
				Action: func(ctx *cli.Context) error {
					rX, rY, err := parseCoord(ctx.String("red"))
					if err != nil {
						return err
					}
					gX, gY, err := parseCoord(ctx.String("green"))
					if err != nil {
						return err
					}
					bX, bY, err := parseCoord(ctx.String("blue"))
					if err != nil {
						return err
					}
					aX, aY, err := parseCoord(ctx.String("alpha"))
					if err != nil {
						return err
					}

					chanShift := cs.Must(cs.New(ctx.String("path"),
						cs.RedShift(rX, rY),
						cs.GreenShift(gX, gY),
						cs.BlueShift(bX, bY),
						cs.AlphaShift(aX, aY),
						cs.WithChunks(ctx.Int("chunk")),
						cs.WithVolatility(ctx.Int("volatility")),
						cs.WithSeed(ctx.Int64("seed")),
						cs.WithAnimate(ctx.Int("animate")),
						cs.WithAnimateVolatility(ctx.Int("anivolatility")),
					))

					imgs := chanShift.ShiftIterate()

					// Single image channelshift
					if len(imgs) == 1 {
						img := imgs[0]
						f, err := os.Create(ctx.String("output") + ".png")
						if err != nil {
							return err
						}

						err = png.Encode(f, img)
						if err != nil {
							return err
						}
					} else {
						// animated gif channelshift
						f, err := os.Create(ctx.String("output") + ".gif")
						if err != nil {
							return err
						}

						for i, img := range imgs {
							var palette color.Palette = append(palette.WebSafe, color.Transparent)
							bounds := img.Bounds()
							dst := image.NewPaletted(bounds, palette)
							draw.Draw(dst, bounds, img, bounds.Min, draw.Src)

							imgs[i] = dst
						}

						delay := make([]int, len(imgs))
						for i := range delay {
							delay[i] = 10
						}

						disposal := make([]byte, len(imgs))
						for i := range disposal {
							disposal[i] = gif.DisposalBackground
						}

						oImgs := make([]*image.Paletted, 0)
						for _, img := range imgs {
							oImgs = append(oImgs, img.(*image.Paletted))
						}

						gif.EncodeAll(f, &gif.GIF{
							Image:    oImgs,
							Delay:    delay,
							Disposal: disposal,
						})
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

func parseCoord(coord string) (int, int, error) {
	xStr, yStr := strings.Split(coord, ",")[0], strings.Split(coord, ",")[1]
	x, err := strconv.Atoi(xStr)
	if err != nil {
		return 0, 0, err
	}
	y, err := strconv.Atoi(yStr)
	if err != nil {
		return 0, 0, err
	}
	return x, y, nil
}
