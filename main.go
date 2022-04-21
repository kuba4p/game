package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image"
	_ "image/png"
	"log"
)

const (
	screenWidth  = 680
	screenHeight = 240
	jumpHeight = 70		//弹跳高度

)

const (
	frameOX     = 0
	frameOY     = 1050
	frameWidth  = 192
	frameHeight = 300
	frameNum    = 8
)
const(
	walk = iota
	jumpUp
	jumpDown
)
//go:embed resource/character.png
var ch embed.FS

//go:embed resource/barrier.png
var ba embed.FS

var (
	characterImage *ebiten.Image
	barrierImage *ebiten.Image
)

func init() {

	//人物图片加载
	chByte, _ := ch.ReadFile("resource/character.png")
	img, _, err := image.Decode(bytes.NewReader(chByte))
	//img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	characterImage = ebiten.NewImageFromImage(img)
	//障碍物图片加载
	chByteb, _ := ba.ReadFile("resource/barrier.png")
	imgb, _, errb := image.Decode(bytes.NewReader(chByteb))
	//img, _, err := image.Decode(bytes.NewReader(images.Tiles_png))
	if errb != nil {
		log.Fatal(errb)
	}
	barrierImage = ebiten.NewImageFromImage(imgb)
}


// Game implements ebiten.Game interface.
type Game struct{
	count int
	status int
	height int 	//当前高度
	keys []ebiten.Key
	bars []bar	//障碍物
}

func (g *Game) Update() error {
	g.keys = inpututil.AppendPressedKeys(g.keys[:0])
	g.count++

	if(g.count % 100 == 0){
		g.createBar()
	}

	//测试
	if(len(g.bars) > 0){
		index := -1
		for i, _ := range g.bars {
			//障碍物向左移动处理
			g.bars[i].x -= 5
			if index == -1 && g.bars[i].x < -50{
				index = i
			}
		}

		if index == len(g.bars) - 1{
			g.bars = g.bars[:0]
		}else{
			g.bars = g.bars[index + 1:]
		}

	}

	if g.status == jumpUp && g.height == 0 {
		g.status = jumpDown
	}

	if g.status == jumpDown && g.height == jumpHeight {
		g.status = walk
	}

	if g.status == jumpUp {
		g.height--
	}

	if g.status == jumpDown {
		g.height += 2
	}
	// Write your game's logical update.
	return nil
}

func (g *Game) createBar() {
	//	创建一个障碍物
	ba := bar{
		x: 1080*2 + 50,
		y: 240,
	}
	g.bars = append(g.bars, ba)
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(screenWidth/2 + frameWidth - 200, screenHeight/2 + frameHeight + float64(g.height * 5))
	//op.GeoM.Translate(float64(500), float64(200))
	op.GeoM.Scale(0.2, 0.2)

	if g.status == jumpUp || g.status == jumpDown{
		screen.DrawImage(characterImage.SubImage(image.Rect(192, 0, 192 + frameWidth,  frameHeight)).(*ebiten.Image), op)
		//log.Println("跳")
	}else if g.status == walk {

		//if g.count % 300 == 0{
		//	g.status = jumpUp
		//}

		i := (g.count / 6) % frameNum
		//i := 2
		sx, sy := frameOX+i*frameWidth, frameOY
		//sx := frame
		screen.DrawImage(characterImage.SubImage(image.Rect(sx, sy, sx+frameWidth, frameOY+frameHeight)).(*ebiten.Image), op)
	}



	//画障碍物
	if len(g.bars) > 0{
		for _, b := range g.bars {
			op = &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(b.x), float64(b.y))
			//op.GeoM.Translate(float64(500), float64(200))
			op.GeoM.Scale(0.3, 0.5)
			screen.DrawImage(barrierImage.SubImage(image.Rect(225, 260, 225+64, 260+128)).(*ebiten.Image), op)
		}

	}


	//判断按键
	//如果状态是走路，则可能进行其他状态转换
	if(g.status == walk){
		for _, p := range g.keys {
			//如果是空格键
			if p == ebiten.KeySpace && g.status != jumpUp && g.status != jumpDown{
				log.Println(p)
				g.status = jumpUp
			}

		}
	}



	ebitenutil.DebugPrint(screen, fmt.Sprintf("SCORE: %d\nTPS: %0.2f\nFPS: %0.2f", g.count, ebiten.CurrentTPS(),ebiten.CurrentFPS()))

}



func (g *Game) Layout(outsideWidth, outsideHeight int) (int,  int) {
	return screenWidth, screenHeight
}

func main() {
	game := &Game{
		status: walk,
		height: 50,
	}
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(1080, 480)
	ebiten.SetWindowTitle("game")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
