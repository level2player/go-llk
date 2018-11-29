package llk

import (
	"github.com/disintegration/imaging"
	"testing"
	"time"
)
var(
	 n1=&Node{PointX:0,PointY:0,IsEmpty:false}
	 n2=&Node{PointX:2,PointY:0,IsEmpty:true}
	 n3=&Node{PointX:2,PointY:2,IsEmpty:false}

	 grid=[3][3]*Node{  
		{n1, &Node{PointX:0,PointY:1,IsEmpty:true}, &Node{PointX:0,PointY:2,IsEmpty:false}} ,   
		{&Node{PointX:1,PointY:0,IsEmpty:false}, &Node{PointX:1,PointY:1,IsEmpty:true}, &Node{PointX:1,PointY:2,IsEmpty:true}} ,   
		{n2, &Node{PointX:2,PointY:1,IsEmpty:true}, n3}}

	n4=&Node{PointX:0,PointY:3,IsEmpty:false}
	n5=&Node{PointX:3,PointY:0,IsEmpty:false}

	grid2=[4][4]*Node{  
		{&Node{PointX:0,PointY:0,IsEmpty:false}, &Node{PointX:0,PointY:1,IsEmpty:false}, &Node{PointX:0,PointY:2,IsEmpty:true},n4} ,   
		{&Node{PointX:1,PointY:0,IsEmpty:false}, &Node{PointX:1,PointY:1,IsEmpty:false}, &Node{PointX:1,PointY:2,IsEmpty:true},&Node{PointX:1,PointY:3,IsEmpty:false}} ,   
		{&Node{PointX:2,PointY:0,IsEmpty:false}, &Node{PointX:2,PointY:1,IsEmpty:false}, &Node{PointX:2,PointY:2,IsEmpty:true},&Node{PointX:2,PointY:3,IsEmpty:false}},
		{n5, &Node{PointX:3,PointY:1,IsEmpty:true}, &Node{PointX:3,PointY:2,IsEmpty:true},&Node{PointX:3,PointY:3,IsEmpty:false}}}
)
  
func Test_ScreenshotWindows(t *testing.T){
	img,err:= ScreenshotWindows()
	if err!=nil{
		t.Error(err)
	}
	imaging.Save(img,"image/test_screenshotWindows.png")
}

func Test_GetWindowRect(t *testing.T){
	rect:= getTopWindowRect()
	t.Log(rect)
}

func Test_MatchBox(t *testing.T){
	//t.Log(matchBox(n1,n2,grid))	
	//t.Log(matchBoxOne(n1,n3,grid))	
	//t.Log(matchBoxTwo(n4,n5,grid2))
	t.Log(int(time.Millisecond))
}


func Test_LoadConfig(t *testing.T){

	t.Log(MainConfig)
}



func Test_Run(t *testing.T){

	v:=Start()
	t.Log(v)
}







