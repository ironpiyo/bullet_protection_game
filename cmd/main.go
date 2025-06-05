package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"game/internal/config"
	"game/internal/game"
	"game/internal/render"
)

// Game はEbitenのゲームインターフェースを実装する
type Game struct {
	gameState *game.Game
}

// Update はゲームの状態を更新する
func (g *Game) Update() error {
	return g.gameState.Update()
}

// Draw はゲームの状態を描画する
func (g *Game) Draw(screen *ebiten.Image) {
	render.Draw(screen, g.gameState)
}

// Layout はウィンドウサイズを返す
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.gameState.Layout(outsideWidth, outsideHeight)
}

func main() {
	// 乱数の初期化
	rand.Seed(time.Now().UnixNano())
	
	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("弾幕避けゲーム")
	
	g := &Game{
		gameState: game.NewGame(),
	}
	
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
