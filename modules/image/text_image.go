package image

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/draw"
)

var (
    width 		  = 100
    height 		  = 100
	textTopMargin = 78
    opt = truetype.Options{
        Size:              80,
        DPI:               0,
        Hinting:           0,
        GlyphCacheEntries: 0,
        SubPixelsX:        0,
        SubPixelsY:        0,
    }
)

func CreateImage(text string, filename string) error {
    // フォントファイルを読み込み
    ftBinary, err := ioutil.ReadFile("./modules/image/Koruri-Bold.ttf")
    if err != nil {
        return err
    }

    ft, err := truetype.Parse(ftBinary)
    if err != nil {
        return err
    }

    img := colorImage()

    face := truetype.NewFace(ft, &opt)

    dr := &font.Drawer{
        Dst:  img,
        Src:  image.Black,
        Face: face,
        Dot:  fixed.Point26_6{},
    }
	initial := string(getRuneAt(text, 0))
    dr.Dot.X = (fixed.I(width) - dr.MeasureString(initial)) / 2
    dr.Dot.Y = fixed.I(textTopMargin)

    dr.DrawString(initial)

    buf := &bytes.Buffer{}
    err = png.Encode(buf, img)

    if err != nil {
        return err
    }

    file, err := os.Create("./upload/user/" + filename)
    if err != nil {
        return err
    }
    defer file.Close()

    file.Write(buf.Bytes())
	return nil
}

func ResizeImage() {
	input, _ := os.Open("abc.png")
	defer input.Close()

	output, _ := os.Create("a.png")
	defer output.Close()

	// Decode the image (from PNG to image.Image):
	src, _ := png.Decode(input)

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:
	png.Encode(output, dst)
}

func colorImage() *image.RGBA {
    rand.Seed(time.Now().UnixNano())
	result := [] uint8{}
    for i := 0; i < 3; i++ {
		var randInt uint8
		for {
			randInt = uint8(rand.Intn(256))
			if 75 < randInt {
				break
			}
		}
		result = append(result, randInt) //sliceに"Ruby"を追加
    }
    img := image.NewRGBA(image.Rect(0, 0, width, height))
	randColor := color.RGBA{result[0], result[1], result[2], 255}
    for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, randColor)
        }
    }
	return img
    // f, _ := os.Create("./image.png")
    // defer f.Close()
    // png.Encode(f, img)
}

func getRuneAt(s string, i int) rune {
    rs := []rune(s)
    return rs[i]
}

