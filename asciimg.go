package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"

	// Side-effect import.
	// Сайд-эффект — добавление декодера PNG в пакет image.
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/buger/goterm"
	colorTerm "github.com/fatih/color"
	"golang.org/x/image/draw"
)

var (
	fileRead = flag.String("o", "", "Output path")
	noScale  = flag.Bool("noscale", false, "Scales picture with <w> and <h> as separate flags")
	w        = flag.Int("w", goterm.Width(), "width of ascii picture")
	h        = flag.Int("h", goterm.Height(), "height of ascii picture")
	isColor  = flag.Bool("c", false, "make ascii picture colorful")
)

func scale(img image.Image, w int, h int) image.Image {
	dstImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dstImg, dstImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dstImg
}

func decodeImageFile(imgName string) (image.Image, error) {
	imgFile, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(imgFile)
	return img, err
}

func processPixel(c color.Color) rune {
	symbols := []rune("@80GCLft1i;:,. ")
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8
	if r == 255 {
		//Костыль, но не могу придумать, как пофиксить.
		r = 250
	}
	return symbols[int(r/(255/uint32(len(symbols))))%(len(symbols))]
}

func convertToASCII(img image.Image) [][]rune {
	textImg := make([][]rune, img.Bounds().Dy())
	// i := 0; i < i.Bound().Dy; i++
	for i := range textImg {
		textImg[i] = make([]rune, img.Bounds().Dx())
	}

	for i := range textImg {
		for j := range textImg[i] {
			textImg[i][j] = processPixel(img.At(j, i))
		}
	}
	return textImg
}

func main() {
	flag.Parse()

	var (
		Cyan   = colorTerm.New(colorTerm.FgCyan)
		Red    = colorTerm.New(colorTerm.FgRed)
		Yellow = colorTerm.New(colorTerm.FgYellow)
	)
	if flag.NArg() == 0 {
		fmt.Println("Usage: asciimg <image.jpg>")
		os.Exit(0)
	}

	img := flag.Arg(0)
	image, err := decodeImageFile(img)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	if *fileRead == "" && !*noScale {
		image = scale(image, *w, *h)
	}
	textImg := convertToASCII(image)
	if *fileRead != "" {
		f, err := os.Create(*fileRead)
		defer f.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		for i := range textImg {
			for j := range textImg[i] {

				fmt.Fprintf(f, "%c", textImg[i][j])
			}
		}
		fmt.Fprintln(f)

	} else {
		for i := range textImg {
			for j := range textImg[i] {
				if *isColor {
					if textImg[i][j] == '@' || textImg[i][j] == '8' || textImg[i][j] == '0' || textImg[i][j] == 'G' {
						Cyan.Printf("%c", textImg[i][j])
					} else if textImg[i][j] == 'C' || textImg[i][j] == 'L' || textImg[i][j] == 'f' || textImg[i][j] == 't' {
						Red.Printf("%c", textImg[i][j])
					} else {
						Yellow.Printf("%c", textImg[i][j])
					}
				} else {
					fmt.Printf("%c", textImg[i][j])
				}
			}
			fmt.Println()
		}
	}
}
