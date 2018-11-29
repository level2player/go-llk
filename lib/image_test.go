package llk

import (
	"fmt"
	"github.com/Comdex/imgo"
	"github.com/disintegration/imaging"
	"image"
	"testing"
)

func Test_crop(t *testing.T) {
	img, err := imgo.DecodeImage("image/test1.png")
	if err != nil {
		t.Error(err)
	}
	nr := Cropllk(img, image.Rect(295, 244, 885, 630))

	newImage := nr.SubImage(nr.Bounds())
	err = imaging.Save(newImage, "image/out.png")
	if err != nil {
		t.Error(err)
	}
	boxWidth := nr.Bounds().Dx() / 19
	boxHeight := nr.Bounds().Dy() / 11

	emptyImg, _ := imgo.DecodeImage("image/tpl1.png")

	for i := 0; i < 19; i++ {
		for k := 0; k < 11; k++ {
			r := Cropllk(newImage, image.Rect(i*boxWidth+3, k*boxHeight+3, (i+1)*boxWidth-3, (k+1)*boxHeight-3))

			n := r.SubImage(r.Bounds())
			v, _ := CompareImage(emptyImg, n)
			if v > 0.97 {
				t.Log(fmt.Sprintf("image/mini/%dx%d.png", i, k))
			}

			err = imaging.Save(n, fmt.Sprintf("image/mini/%dx%d.png", i, k))
			if err != nil {
				t.Error(err)
			}
		}
	}

}

func Test_BuildGrid(t *testing.T) {
	img, err := imgo.DecodeImage("image/test1.png")
	if err != nil {
		t.Error(err)
	}
	//t.Log(image.Rect(295, 244, 885, 630).Min)
	grid, err := BuildLinkBox(img, image.Rect(295, 244, 885, 630))
	if err != nil {
		t.Error(err)
	}
	for _, v := range grid {
		for _, p := range v {
			t.Log(p.PointX, p.PointY, p.IsEmpty,p.ClickX,p.ClickY)
		}
	}
}
