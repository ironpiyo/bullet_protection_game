package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth      = 800
	screenHeight     = 600
	playerSize       = 10
	bulletSize       = 8
	initialBullets   = 20
	bulletSpeedMin   = 2.0
	bulletSpeedMax   = 5.0
	bulletSpawnRate  = 5  // 1秒あたりの新しい弾の数
	maxRankingScores = 5  // ランキングに表示するスコア数
)

// プレイヤーの構造体
type Player struct {
	x, y float64
	size float64
}

// 弾の構造体
type Bullet struct {
	x, y      float64
	vx, vy    float64
	size      float64
	color     color.RGBA
}

// ゲームの状態を管理する構造体
type Game struct {
	player        Player
	bullets       []Bullet
	gameOver      bool
	startTime     time.Time
	currentTime   float64
	scores        []float64
	lastBulletAdd time.Time
}

// NewGame は新しいゲームインスタンスを作成する
func NewGame() *Game {
	g := &Game{
		player: Player{
			x:    float64(screenWidth) / 2,
			y:    float64(screenHeight) / 2,
			size: playerSize,
		},
		bullets:      make([]Bullet, 0, initialBullets),
		gameOver:     false,
		startTime:    time.Now(),
		currentTime:  0,
		scores:       make([]float64, 0, maxRankingScores),
		lastBulletAdd: time.Now(),
	}

	// 初期の弾を生成
	for i := 0; i < initialBullets; i++ {
		g.addRandomBullet()
	}

	return g
}

// addRandomBullet はランダムな位置と速度で新しい弾を追加する
func (g *Game) addRandomBullet() {
	// 画面の端から弾を発射
	var x, y float64
	var vx, vy float64
	
	side := rand.Intn(4) // 0: 上, 1: 右, 2: 下, 3: 左
	
	switch side {
	case 0: // 上から
		x = rand.Float64() * screenWidth
		y = -bulletSize
		vx = (rand.Float64()*2 - 1) * bulletSpeedMax
		vy = rand.Float64()*(bulletSpeedMax-bulletSpeedMin) + bulletSpeedMin
	case 1: // 右から
		x = screenWidth + bulletSize
		y = rand.Float64() * screenHeight
		vx = -(rand.Float64()*(bulletSpeedMax-bulletSpeedMin) + bulletSpeedMin)
		vy = (rand.Float64()*2 - 1) * bulletSpeedMax
	case 2: // 下から
		x = rand.Float64() * screenWidth
		y = screenHeight + bulletSize
		vx = (rand.Float64()*2 - 1) * bulletSpeedMax
		vy = -(rand.Float64()*(bulletSpeedMax-bulletSpeedMin) + bulletSpeedMin)
	case 3: // 左から
		x = -bulletSize
		y = rand.Float64() * screenHeight
		vx = rand.Float64()*(bulletSpeedMax-bulletSpeedMin) + bulletSpeedMin
		vy = (rand.Float64()*2 - 1) * bulletSpeedMax
	}
	
	// ランダムな色を生成
	r := uint8(rand.Intn(200) + 55)
	gb := uint8(rand.Intn(200) + 55)
	b := uint8(rand.Intn(200) + 55)
	
	bullet := Bullet{
		x:    x,
		y:    y,
		vx:   vx,
		vy:   vy,
		size: bulletSize,
		color: color.RGBA{r, gb, b, 255},
	}
	
	g.bullets = append(g.bullets, bullet)
}

// Update はゲームの状態を更新する
func (g *Game) Update() error {
	// ゲームオーバー時のリスタート処理
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			*g = *NewGame()
			return nil
		}
		return nil
	}

	// プレイヤーの位置をマウスカーソルに合わせる
	x, y := ebiten.CursorPosition()
	g.player.x = float64(x)
	g.player.y = float64(y)

	// 経過時間を更新
	g.currentTime = time.Since(g.startTime).Seconds()

	// 一定間隔で新しい弾を追加
	if time.Since(g.lastBulletAdd).Seconds() > 1.0/float64(bulletSpawnRate) {
		g.addRandomBullet()
		g.lastBulletAdd = time.Now()
	}

	// 弾の移動と画面外判定
	newBullets := make([]Bullet, 0, len(g.bullets))
	for _, b := range g.bullets {
		// 弾を移動
		b.x += b.vx
		b.y += b.vy
		
		// 画面外に出た弾は削除
		if b.x < -100 || b.x > screenWidth+100 || b.y < -100 || b.y > screenHeight+100 {
			continue
		}
		
		// プレイヤーとの衝突判定
		dx := g.player.x - b.x
		dy := g.player.y - b.y
		distance := math.Sqrt(dx*dx + dy*dy)
		
		if distance < g.player.size+b.size {
			// 衝突した場合、ゲームオーバー
			g.gameOver = true
			
			// スコアを記録
			g.scores = append(g.scores, g.currentTime)
			
			// スコアを降順にソート
			sort.Slice(g.scores, func(i, j int) bool {
				return g.scores[i] > g.scores[j]
			})
			
			// 上位スコアだけを保持
			if len(g.scores) > maxRankingScores {
				g.scores = g.scores[:maxRankingScores]
			}
			
			break
		}
		
		newBullets = append(newBullets, b)
	}
	
	if !g.gameOver {
		g.bullets = newBullets
	}

	return nil
}

// Draw は画面に描画する
func (g *Game) Draw(screen *ebiten.Image) {
	// 背景を黒で塗りつぶす
	screen.Fill(color.RGBA{20, 20, 40, 255})

	// 弾を描画
	for _, b := range g.bullets {
		ebitenutil.DrawCircle(screen, b.x, b.y, b.size, b.color)
	}

	// プレイヤーを描画（白い円）
	if !g.gameOver {
		ebitenutil.DrawCircle(screen, g.player.x, g.player.y, g.player.size, color.RGBA{255, 255, 255, 255})
	}

	// 経過時間を表示
	timeText := fmt.Sprintf("Time: %.2f", g.currentTime)
	ebitenutil.DebugPrintAt(screen, timeText, 20, 20)

	// ゲームオーバー時の表示
	if g.gameOver {
		gameOverText := "GAME OVER - Press SPACE to restart"
		ebitenutil.DebugPrintAt(screen, gameOverText, screenWidth/2-len(gameOverText)*3, screenHeight/2)

		// ランキングを表示
		ebitenutil.DebugPrintAt(screen, "TOP SCORES:", screenWidth/2-50, screenHeight/2+30)
		for i, score := range g.scores {
			scoreText := fmt.Sprintf("%d. %.2f seconds", i+1, score)
			ebitenutil.DebugPrintAt(screen, scoreText, screenWidth/2-50, screenHeight/2+50+i*20)
		}
	}
}

// Layout はウィンドウサイズを返す
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	// 乱数の初期化
	rand.Seed(time.Now().UnixNano())
	
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("弾幕避けゲーム")
	
	game := NewGame()
	
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
