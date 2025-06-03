package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

// Game はゲームの状態を管理する構造体
type Game struct{}

// Update はゲームの状態を更新する
func (g *Game) Update() error {
	return nil
}

// Draw は画面に描画する
func (g *Game) Draw(screen *ebiten.Image) {
	// 背景を水色で塗りつぶす
	screen.Fill(color.RGBA{135, 206, 235, 255})
	
	// Hello World! を表示
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

// Layout はウィンドウサイズを返す
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello World")
	
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
