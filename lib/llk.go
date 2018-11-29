package llk

import (
	"image"
	"github.com/lxn/win"
	"github.com/kbinani/screenshot"
	"time"
	"syscall"
	"errors"
	"log"
	"io/ioutil"
	"encoding/json"
)



type LLKConfig struct {
	ImageAcquaintance  float64 `json:"imageAcquaintance"`
	Interval           int     `json:"interval"`
	BoxInterval        int     `json:"boxInterval"`
	ClickPointInterval int     `json:"clickPointInterval"`
	MarginX            int     `json:"marginX"`
	MarginY            int     `json:"marginY"`
	PaddingX           int     `json:"paddingX"`
	PaddingY           int     `json:"paddingY"`
	EmptyImgPath       string  `json:"emptyImgPath"`
}


type Node struct {
	// LeftNode *Node
	// RightNode *Node
	// TopNode *Node
	// BottomNode *Node

	Img     image.Image
	IsEmpty bool
	PointX  int
	PointY  int
	ClickX  int
	ClickY  int
}

type point struct{
	x,y int
}

var (
	libuser32 uintptr
	mouse_event uintptr
	getWindowThreadProcessId uintptr
	MainConfig *LLKConfig
)

func init() {
	libuser32 = win.MustLoadLibrary("user32.dll")
	mouse_event = win.MustGetProcAddress(libuser32, "mouse_event")
	getWindowThreadProcessId=win.MustGetProcAddress(libuser32,"GetWindowThreadProcessId")
	MainConfig=&LLKConfig{}
	LoadConfig()
}

func LoadConfig()*LLKConfig{
	buffer,err:=ioutil.ReadFile("../config.json")
	if err!=nil{
		panic(err)
	}
	err=json.Unmarshal(buffer,MainConfig)
	if err!=nil{
		panic(err)
	}
	return MainConfig
}

func (node *Node) SetEmpty() {
	node.IsEmpty = true
}


func (node *Node) Click() {
	r := win.SetCursorPos(int32(node.ClickX), int32(node.ClickY))
	if r {
		mouseEvent(win.MOUSEEVENTF_LEFTDOWN|win.MOUSEEVENTF_LEFTUP, 0, 0)
	}
}


func ScreenshotWindows()(image.Image, error) {
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			return nil,err
		}
		return img,nil
	}
	return nil,errors.New("screenshot fail")
}


func getTopWindowRect() image.Rectangle {
	var rect win.RECT
	win.GetWindowRect(win.GetForegroundWindow(), &rect)
	px0 := int(rect.Left) +MainConfig.PaddingX
	py0 := int(rect.Top) +MainConfig.PaddingY
	return image.Rect(px0,py0,px0+MainConfig.MarginX,py0+MainConfig.MarginY)
}

func BytesToUint16(array []byte) uint16 {
	var data uint16 =0
	for i:=0;i< len(array);i++  {
	   data = data+uint16(uint(array[i])<<uint(8*i))
	}
	
	return data
 }


func mouseEvent(msg, px, py int) {
	if mouse_event == 0 {
		return
	}

	syscall.Syscall(mouse_event, 3, uintptr(msg),
		uintptr(px),
		uintptr(py))
}

func getNoneEmptyCount()int{
	img,err:= ScreenshotWindows()
	if err!=nil{
		log.Println(err)
	}
	rect:= getTopWindowRect()
	mainGrid, _ := BuildLinkBox(img, rect)
	count:=0
	for _, v := range mainGrid {
		for _, p := range v {
			if !p.IsEmpty{
				count++
			}
		}
	}
	return count
}


func getSameBox(n1 *Node,mainGrid [19][11]*Node)(numbers []*Node){
	for _, v := range mainGrid {
		for _, p := range v {
			if p==n1{
				continue
			}
			f,err:=CompareImage(n1.Img,p.Img)
			if err==nil&&f>0.97{
				numbers=append(numbers,p)
			}
		}
	}
	return numbers
}

func Start()(elapsed time.Duration){
	invocation:=time.Now()
	defer func(){
		elapsed = time.Since(invocation)
	}()
	for {
		img,err:= ScreenshotWindows()
		if err!=nil{
			log.Println(err)
			return
		}
		r:=getTopWindowRect()
		if r.Max.X<0||r.Max.Y<0||r.Min.X<0||r.Min.Y<0{
			log.Println(errors.New("Incorrect program location"))
			return
		}
		grid,err := BuildLinkBox(img, r)
		if err != nil {
			log.Println(err)
			return
		}
		removeEngine(grid)
		if getNoneEmptyCount()==0 {
			return 
		}
		time.Sleep(time.Duration(MainConfig.Interval)*time.Second)
	} 
}

func removeEngine(mainGrid [19][11]*Node){
	for _, v := range mainGrid {
		for _, p := range v {
			if !p.IsEmpty{
				numbers:=getSameBox(p,mainGrid)
				for _,node:=range numbers{
					b:=false
					b=matchBox(p,node,mainGrid)
					b=matchBoxOne(p,node,mainGrid)
					b=matchBoxTwo(p,node,mainGrid)
					if b{
						p.Click()
						node.Click()
						p.SetEmpty()
						node.SetEmpty()
					}
				}
			}
		}
	}
}

func matchBox(n1,n2 *Node,mainGrid [19][11]*Node)bool{
	if n1.PointX!=n2.PointX && n1.PointY!=n2.PointY{
		return false
	}
	if n1.PointX==n2.PointX {
		minLen:=min(n1.PointY,n2.PointY)
		maxLen:=max(n1.PointY,n2.PointY)
		for minLen++;minLen<maxLen;minLen++{
			if !mainGrid[n1.PointX][minLen].IsEmpty{
				return false
			}
		}
	}
	if n1.PointY==n2.PointY {
		minLen:=min(n1.PointX,n2.PointX)
		maxLen:=max(n1.PointX,n2.PointX)
		for minLen++;minLen<maxLen;minLen++{
			if !mainGrid[minLen][n1.PointY].IsEmpty{
				return false
			}
		}
	}
	return true
}

func matchBoxOne(n1,n2 *Node,mainGrid [19][11]*Node)bool{
	point1:=point{x:n1.PointX,y:n2.PointY}
	point2:=point{x:n2.PointX,y:n1.PointY}
	if mainGrid[point1.x][point1.y].IsEmpty{
		if matchBox(n1,mainGrid[point1.x][point1.y],mainGrid)&& matchBox(n2,mainGrid[point1.x][point1.y],mainGrid){
			return true
		}
	}
	if mainGrid[point2.x][point2.y].IsEmpty{
		if matchBox(n1,mainGrid[point2.x][point2.y],mainGrid)&& matchBox(n2,mainGrid[point2.x][point2.y],mainGrid){
			return true
		}
	}
	return false
}

func matchBoxTwo(n1,n2 *Node,mainGrid [19][11]*Node)bool{
	for i:=n1.PointX-1;i>=0;i--{
		if mainGrid[i][n1.PointY].IsEmpty{
			if matchBoxOne(mainGrid[i][n1.PointY],n2,mainGrid){
				return true
			}
		}
	}
	for i:=n1.PointX+1;i<19;i++{
		if mainGrid[i][n1.PointY].IsEmpty{
			if matchBoxOne(mainGrid[i][n1.PointY],n2,mainGrid){
				return true
			}
		}
	}

	for i:=n1.PointY-1;i>=0;i--{
		if mainGrid[n1.PointX][i].IsEmpty{
			if matchBoxOne(mainGrid[n1.PointX][i],n2,mainGrid){
				return true
			}
		}
	}

	for i:=n1.PointY+1;i<11;i++{
		if mainGrid[n1.PointX][i].IsEmpty{
			if matchBoxOne(mainGrid[n1.PointX][i],n2,mainGrid){
				return true
			}
		}
	}

	return false
	
}



func max(a, b int) int {
    if a < b {
        return b
    }
    return a
} 

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

