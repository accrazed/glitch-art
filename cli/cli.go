package cli

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
	"time"

	cs "github.com/accrazed/glitch-art/src/channelshift"
	"github.com/accrazed/glitch-art/src/corrupt/jpg"
	"github.com/accrazed/glitch-art/src/lib"
	ps "github.com/accrazed/glitch-art/src/pixelsort"

	"github.com/urfave/cli/v2"
)

func RunCLI() {
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
						Required: false,
					},
					&cli.IntFlag{
						Name:     "threshold",
						Aliases:  []string{"t"},
						Usage:    "threshold to use for saturation sorting",
						Required: true,
					},
					&cli.IntFlag{
						Name:     "chunklim",
						Aliases:  []string{"cl"},
						Usage:    "max chunk length a chunk can be be before sorting",
						Required: false,
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
						Name:    "sortfunc",
						Aliases: []string{"sf"},
						Usage:   "Sorting func to use when pixel sorting. Refer to sorts.go",
						Value:   "MeanComp",
					},
					&cli.StringFlag{
						Name:    "thresholdfunc",
						Aliases: []string{"tf"},
						Usage:   "Threshold func to use when pixel sorting. Refer to thresholds.go",
						Value:   "OutThresholdColorMean",
					},
					&cli.StringFlag{
						Name:     "output",
						Aliases:  []string{"o"},
						Usage:    "filename to output as",
						Value:    "output",
						Required: false,
					},
				},
				Action: DoPixelSort,
			}, {
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
					&cli.StringFlag{
						Name:     "chunkdir",
						Aliases:  []string{"d"},
						Usage:    "Which direction to chunk the image in",
						Required: false,
					},
					&cli.IntFlag{
						Name:    "chunkvol",
						Aliases: []string{"cv"},
						Usage:   "How volatile to adjust the chunk width in every iteration",
						Value:   0,
					},
					&cli.IntFlag{
						Name:     "offsetvol",
						Aliases:  []string{"ov"},
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
				},
				Action: DoChannelShift,
			},
			{
				Name:        "corrupt",
				Description: "corrupt your images :)",
				Aliases:     []string{"c"},
				Subcommands: []*cli.Command{
					{
						Name:  "jpg",
						Usage: "corrupt jpg images. May corrupt to the point of image error",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "path",
								Aliases:  []string{"p"},
								Usage:    "path to your jpg image",
								Required: true,
							},
							&cli.Int64Flag{
								Name:    "seed",
								Aliases: []string{"s"},
								Value:   time.Now().Unix(),
							},
							&cli.IntFlag{
								Name:    "corruptStrength",
								Aliases: []string{"cs"},
								Value:   1e4,
							},
							&cli.StringFlag{
								Name:     "output",
								Aliases:  []string{"o"},
								Usage:    "filename to output as",
								Value:    "output",
								Required: false,
							},
						},
						Action: DoCorruptJPEG,
					},
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

func DoPixelSort(ctx *cli.Context) error {
	sortDir := lib.Vertical
	if strings.ToLower(ctx.String("direction")) == "horizontal" {
		sortDir = lib.Horizontal
	}

	pixSort := ps.Must(ps.New(
		ctx.String("path"),
		ps.WithDirection(sortDir),
		ps.WithSeed(ctx.Int64("seed")),
		ps.WithThreshold(ctx.Int("threshold")),
		ps.WithInvert(ctx.Bool("invert")),
		ps.WithSortFuncString(ctx.String("sortfunc")),
		ps.WithThresholdFuncString(ctx.String("thresholdfunc")),
		ps.WithChunkLimit(ctx.Int("chunklim")),
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
}

func DoChannelShift(ctx *cli.Context) error {
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

	sortDir := lib.Vertical
	if strings.ToLower(ctx.String("chunkdir")) == "horizontal" {
		sortDir = lib.Horizontal
	}

	chanShift := cs.Must(cs.New(ctx.String("path"),
		cs.WithRedShift(rX, rY),
		cs.WithGreenShift(gX, gY),
		cs.WithBlueShift(bX, bY),
		cs.WithAlphaShift(aX, aY),
		cs.WithChunks(ctx.Int("chunk")),
		cs.WithOffsetVolatility(ctx.Int("offsetvol")),
		cs.WithSeed(ctx.Int64("seed")),
		cs.WithAnimate(ctx.Int("animate")),
		cs.WithChunkVolatility(ctx.Int("chunkvol")),
		cs.WithDirection(sortDir),
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

		var palette color.Palette = append(palette.WebSafe, color.Transparent)

		// Prepare gif options
		delay := make([]int, len(imgs))
		disposal := make([]byte, len(imgs))
		oImgs := make([]*image.Paletted, len(imgs))
		for i, img := range imgs {
			bounds := img.Bounds()
			dst := image.NewPaletted(bounds, palette)
			draw.Draw(dst, bounds, img, bounds.Min, draw.Src)
			oImgs[i] = dst

			delay[i] = 10
			disposal[i] = gif.DisposalBackground
		}

		gif.EncodeAll(f, &gif.GIF{
			Image:    oImgs,
			Delay:    delay,
			Disposal: disposal,
		})
	}

	return nil
}

func DoCorruptJPEG(ctx *cli.Context) error {
	jc := jpg.Must(jpg.New(
		ctx.String("path"),
		jpg.WithSeed(ctx.Int64("seed")),
		jpg.WithCorruptStrength(ctx.Int("corruptStrength")),
	))

	jpeg, err := jc.Corrupt()
	if err != nil {
		return err
	}

	outf, err := os.Create(ctx.String("output"))
	if err != nil {
		panic(err)
	}
	_, err = outf.Write(jpeg)
	if err != nil {
		panic(err)
	}

	return nil
}
