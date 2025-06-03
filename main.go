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

// スコアアニメーションの構造体
type ScoreAnimation struct {
	score     float64
	x, y      float64
	scale     float64
	alpha     float64
	lifetime  float64
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
	
	// UI効果用の変数
	gameOverAlpha    float64  // ゲームオーバー画面の透明度
	gameOverScale    float64  // ゲームオーバーテキストのスケール
	rankingAppear    float64  // ランキング表示の進行度
	scoreAnimations  []ScoreAnimation // スコアアニメーション
	
	// 難易度関連の変数
	difficulty       int       // 現在の難易度レベル
	lastDifficultyIncrease time.Time // 最後に難易度を上げた時間
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
		
		// UI効果の初期化
		gameOverAlpha: 0,
		gameOverScale: 0.5,
		rankingAppear: 0,
		scoreAnimations: make([]ScoreAnimation, 0),
		
		// 難易度の初期化
		difficulty: 1,
		lastDifficultyIncrease: time.Now(),
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
	
	// 難易度に応じて弾の速度を調整
	speedMultiplier := 1.0 + float64(g.difficulty-1)*0.1 // 難易度ごとに10%ずつ速くなる
	minSpeed := bulletSpeedMin * speedMultiplier
	maxSpeed := bulletSpeedMax * speedMultiplier
	
	switch side {
	case 0: // 上から
		x = rand.Float64() * screenWidth
		y = -bulletSize
		vx = (rand.Float64()*2 - 1) * maxSpeed
		vy = rand.Float64()*(maxSpeed-minSpeed) + minSpeed
	case 1: // 右から
		x = screenWidth + bulletSize
		y = rand.Float64() * screenHeight
		vx = -(rand.Float64()*(maxSpeed-minSpeed) + minSpeed)
		vy = (rand.Float64()*2 - 1) * maxSpeed
	case 2: // 下から
		x = rand.Float64() * screenWidth
		y = screenHeight + bulletSize
		vx = (rand.Float64()*2 - 1) * maxSpeed
		vy = -(rand.Float64()*(maxSpeed-minSpeed) + minSpeed)
	case 3: // 左から
		x = -bulletSize
		y = rand.Float64() * screenHeight
		vx = rand.Float64()*(maxSpeed-minSpeed) + minSpeed
		vy = (rand.Float64()*2 - 1) * maxSpeed
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
	// ゲームオーバー時のリスタート処理とアニメーション
	if g.gameOver {
		// ゲームオーバーアニメーションの更新
		if g.gameOverAlpha < 0.8 {
			g.gameOverAlpha += 0.02
		}
		
		if g.gameOverScale < 1.2 {
			g.gameOverScale += 0.03
		} else if g.gameOverScale > 1.2 {
			g.gameOverScale = 1.2
		}
		
		// ランキング表示のアニメーション
		if g.rankingAppear < 1.0 {
			g.rankingAppear += 0.03
		}
		
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			// 現在のスコアを保持
			oldScores := g.scores
			
			// ゲームをリセット
			*g = *NewGame()
			
			// スコアを復元
			g.scores = oldScores
			
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
	
	// 6秒ごとに難易度を上げる
	if time.Since(g.lastDifficultyIncrease).Seconds() > 6.0 {
		g.difficulty++
		g.lastDifficultyIncrease = time.Now()
		
		// デバッグ用に難易度上昇を表示
		log.Printf("難易度上昇: レベル %d", g.difficulty)
	}

	// 難易度に応じて弾の発生頻度を調整
	bulletSpawnInterval := 1.0 / float64(bulletSpawnRate)
	if time.Since(g.lastBulletAdd).Seconds() > bulletSpawnInterval {
		// 難易度に応じて複数の弾を発射
		bulletsToAdd := g.difficulty
		for i := 0; i < bulletsToAdd; i++ {
			g.addRandomBullet()
		}
		g.lastBulletAdd = time.Now()
	}

	// スコアアニメーションの更新
	newScoreAnims := make([]ScoreAnimation, 0)
	for _, anim := range g.scoreAnimations {
		anim.lifetime -= 0.016 // 約60FPSを想定
		anim.scale += 0.02
		anim.alpha -= 0.02
		
		if anim.lifetime > 0 && anim.alpha > 0 {
			newScoreAnims = append(newScoreAnims, anim)
		}
	}
	g.scoreAnimations = newScoreAnims

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
			
			// スコアアニメーションを追加
			g.scoreAnimations = append(g.scoreAnimations, ScoreAnimation{
				score:    g.currentTime,
				x:        screenWidth / 2,
				y:        screenHeight / 3,
				scale:    1.0,
				alpha:    1.0,
				lifetime: 2.0,
			})
			
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

	// 経過時間と難易度を表示
	timeText := fmt.Sprintf("Time: %.2f  Difficulty: %d", g.currentTime, g.difficulty)
	ebitenutil.DebugPrintAt(screen, timeText, 20, 20)

	// スコアアニメーションを描画
	for _, anim := range g.scoreAnimations {
		// スケールと透明度に基づいて描画
		scoreText := fmt.Sprintf("%.2f seconds", anim.score)
		
		// 文字サイズを計算（スケールに応じて）
		textWidth := float64(len(scoreText) * 6) * anim.scale
		textX := anim.x - textWidth/2
		
		// テキストを描画（簡易的なスケーリング）
		for i := 0; i < int(anim.scale*2); i++ {
			ebitenutil.DebugPrintAt(screen, scoreText, int(textX), int(anim.y)+i)
		}
	}

	// ゲームオーバー時の表示
	if g.gameOver {
		// 半透明のオーバーレイを描画
		overlayColor := color.RGBA{0, 0, 0, uint8(g.gameOverAlpha * 200)}
		ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), overlayColor)
		
		// ゲームオーバーテキストを描画（拡大/縮小アニメーション付き）
		gameOverText := "GAME OVER"
		textWidth := float64(len(gameOverText) * 8) * g.gameOverScale
		textX := (float64(screenWidth) - textWidth) / 2
		textY := float64(screenHeight)/3 - 20
		
		// テキストを複数回描画してボールド効果を出す
		for i := 0; i < int(g.gameOverScale*3); i++ {
			ebitenutil.DebugPrintAt(screen, gameOverText, int(textX), int(textY)+i)
		}
		
		// リスタート案内
		restartText := "Press SPACE to restart"
		restartX := screenWidth/2 - len(restartText)*3
		restartY := int(textY) + 40
		ebitenutil.DebugPrintAt(screen, restartText, restartX, restartY)
		
		// ランキングを表示（徐々に表示されるアニメーション）
		if g.rankingAppear > 0 {
			rankingTitleY := screenHeight/2 + 30
			ebitenutil.DebugPrintAt(screen, "TOP SCORES:", screenWidth/2-50, rankingTitleY)
			
			// 各スコアを表示（徐々に表示）
			maxScoresToShow := int(float64(len(g.scores)) * g.rankingAppear)
			for i := 0; i < maxScoresToShow && i < len(g.scores); i++ {
				scoreText := fmt.Sprintf("%d. %.2f seconds", i+1, g.scores[i])
				
				// アニメーション効果（少しずつ右から現れる）
				offset := int((1.0 - g.rankingAppear) * 100)
				if offset < 0 {
					offset = 0
				}
				
				ebitenutil.DebugPrintAt(screen, scoreText, screenWidth/2-50+offset, rankingTitleY+20+i*20)
			}
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
