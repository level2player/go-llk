package llk

import (
	"github.com/Comdex/imgo"
	"github.com/disintegration/imaging"
	"image"
	"math"
	"log"
	
)

//切图
func Cropllk(img image.Image, rect image.Rectangle) *image.NRGBA {
	return imaging.Crop(img, rect)
}

//创建连连看表格
func BuildLinkBox(img image.Image, rectangle image.Rectangle) ([19][11]*Node, error) {
		var nodeArry [19][11]*Node
	emptyImg, err := imgo.DecodeImage(MainConfig.EmptyImgPath)
	if err != nil {
		return nodeArry, err
	}	

	nr := Cropllk(img, rectangle)
	boxWidth := nr.Bounds().Dx() / 19
	boxHeight := nr.Bounds().Dy() / 11
	newImage := nr.SubImage(nr.Bounds())

	for i := 0; i < 19; i++ {
		for k := 0; k < 11; k++ {
			r := Cropllk(newImage, image.Rect(i*boxWidth+MainConfig.BoxInterval, k*boxHeight+MainConfig.BoxInterval, (i+1)*boxWidth-MainConfig.BoxInterval, (k+1)*boxHeight-MainConfig.BoxInterval))
			u := r.SubImage(r.Bounds())
			v, err := CompareImage(emptyImg, u)
			if err!=nil{
				log.Println(err)
				continue
			}
			isEmpty := false
			if v > MainConfig.ImageAcquaintance {
				isEmpty = true
			}
			nodeArry[i][k] = &Node{
				PointX: i,
				PointY:  k,
				Img:     u,
				IsEmpty: isEmpty,
				ClickX:rectangle.Min.X+(i*boxWidth+MainConfig.BoxInterval)+boxWidth/MainConfig.ClickPointInterval,
				ClickY:rectangle.Min.Y+(k*boxHeight+MainConfig.BoxInterval)+boxHeight/MainConfig.ClickPointInterval}
		}
	}
	return nodeArry, nil

}

//图片比较
func CompareImage(img1, img2 image.Image) (cossimi float64, err error) {
	matrix1, err2 := imgo.Read(img1)
	matrix2, err2 := imgo.Read(img2)
	if err2!=nil{
		err=err2
		return
	}

	myx := imgo.Matrix2Vector(matrix1)
	myy := imgo.Matrix2Vector(matrix2)

	cos1 := imgo.Dot(myx, myy)
	cos21 := math.Sqrt(imgo.Dot(myx, myx))
	cos22 := math.Sqrt(imgo.Dot(myy, myy))
	cossimi = cos1 / (cos21 * cos22)
	return
}
