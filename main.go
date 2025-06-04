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
	
	// シールド関連の定数
	shieldDurability = 3  // シールドの耐久値
	shieldItemSize   = 15 // シールドアイテムのサイズ
	shieldSpawnRate  = 0.05 // シールドアイテムの出現確率（1フレームあたり）
)

// プレイヤーの構造体
type Player struct {
	x, y float64
	size float64
	shield int // シールドの耐久値
}

// 弾の構造体
type Bullet struct {
	x, y      float64
	vx, vy    float64
	size      float64
	color     color.RGBA
}

// シールドアイテムの構造体
type ShieldItem struct {
	x, y      float64
	size      float64
	active    bool
	angle     float64  // 回転角度
	glowSize  float64  // 輝きのサイズ
	glowDir   float64  // 輝きの方向（拡大/縮小）
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
	shieldItem    ShieldItem // シールドアイテム
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
			shield: 0, // 初期状態ではシールドなし
		},
		bullets:      make([]Bullet, 0, initialBullets),
		shieldItem: ShieldItem{
			x: -100, // 画面外に配置
			y: -100,
			size: shieldItemSize,
			active: false,
			angle: 0,
			glowSize: 0,
			glowDir: 0.1,
		},
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
	
	// シールドアイテムの生成（ランダムに）
	if !g.shieldItem.active && rand.Float64() < shieldSpawnRate {
		g.spawnShieldItem()
	}
	
	// シールドアイテムのアニメーション更新
	if g.shieldItem.active {
		// 回転させる
		g.shieldItem.angle += 0.05
		
		// 輝きのサイズを変化させる
		g.shieldItem.glowSize += g.shieldItem.glowDir
		if g.shieldItem.glowSize > 3 || g.shieldItem.glowSize < 0 {
			g.shieldItem.glowDir *= -1
		}
	}
	
	// シールドアイテムとプレイヤーの衝突判定
	if g.shieldItem.active {
		dx := g.player.x - g.shieldItem.x
		dy := g.player.y - g.shieldItem.y
		distance := math.Sqrt(dx*dx + dy*dy)
		
		if distance < g.player.size+g.shieldItem.size {
			// シールドを獲得
			g.player.shield = shieldDurability
			g.shieldItem.active = false
			g.shieldItem.x = -100 // 画面外に移動
			g.shieldItem.y = -100
			
			// デバッグ用にシールド獲得を表示
			log.Printf("シールド獲得！ 耐久値: %d", g.player.shield)
		}
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
			// シールドがある場合
			if g.player.shield > 0 {
				g.player.shield--
				log.Printf("シールドが弾を防いだ！ 残り耐久値: %d", g.player.shield)
				continue // この弾は消える
			} else {
				// シールドがない場合、ゲームオーバー
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

	// シールドアイテムを描画
	if g.shieldItem.active {
		// 外側の輝き
		glowColor := color.RGBA{0, 255, 255, 100}
		ebitenutil.DrawCircle(screen, g.shieldItem.x, g.shieldItem.y, 
							 g.shieldItem.size + g.shieldItem.glowSize, glowColor)
		
		// メインの円
		ebitenutil.DrawCircle(screen, g.shieldItem.x, g.shieldItem.y, 
							 g.shieldItem.size, color.RGBA{0, 200, 255, 200})
		
		// 内側の円
		ebitenutil.DrawCircle(screen, g.shieldItem.x, g.shieldItem.y, 
							 g.shieldItem.size * 0.6, color.RGBA{100, 200, 255, 255})
		
		// 回転する星型のエフェクト
		angle := g.shieldItem.angle
		radius := g.shieldItem.size * 0.8
		for i := 0; i < 5; i++ {
			a := angle + float64(i) * math.Pi * 0.4
			x1 := g.shieldItem.x + math.Cos(a) * radius
			y1 := g.shieldItem.y + math.Sin(a) * radius
			x2 := g.shieldItem.x + math.Cos(a+0.2) * (radius * 0.5)
			y2 := g.shieldItem.y + math.Sin(a+0.2) * (radius * 0.5)
			ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{255, 255, 255, 200})
		}
	}

	// 弾を描画
	for _, b := range g.bullets {
		ebitenutil.DrawCircle(screen, b.x, b.y, b.size, b.color)
	}

	// プレイヤーを描画（白い円）
	if !g.gameOver {
		ebitenutil.DrawCircle(screen, g.player.x, g.player.y, g.player.size, color.RGBA{255, 255, 255, 255})
		
		// シールドがある場合、プレイヤーの周りにシールドを描画
		if g.player.shield > 0 {
			// シールドの色は耐久値によって変化
			var shieldColor color.RGBA
			switch g.player.shield {
			case 3:
				shieldColor = color.RGBA{0, 255, 255, 100} // 明るい水色（半透明）
			case 2:
				shieldColor = color.RGBA{100, 200, 255, 100} // やや暗い水色（半透明）
			case 1:
				shieldColor = color.RGBA{150, 150, 255, 100} // 紫がかった色（半透明）
			}
			
			// シールドを描画（プレイヤーと同じサイズだが、半透明）
			ebitenutil.DrawCircle(screen, g.player.x, g.player.y, g.player.size, shieldColor)
			
			// 回転する16角形のエフェクト
			particleCount := 16 // 16角形
			baseAngle := g.currentTime * 1.5 // 時間に基づいて回転（速度調整）
			
			// エフェクトとプレイヤーの間の距離を広げる
			particleDistance := g.player.size * 2.0 // プレイヤーからの距離を2倍に
			
			// 16角形の頂点を描画
			for i := 0; i < particleCount; i++ {
				angle := baseAngle + float64(i) * (2 * math.Pi / float64(particleCount))
				
				// 粒子の位置を計算
				particleX := g.player.x + math.Cos(angle) * particleDistance
				particleY := g.player.y + math.Sin(angle) * particleDistance
				
				// 粒子を描画
				particleSize := 2.5
				particleColor := color.RGBA{255, 255, 255, 200} // 白い光の粒子
				ebitenutil.DrawCircle(screen, particleX, particleY, particleSize, particleColor)
			}
			
			// 16角形の辺を描画（頂点同士を線で結ぶ）
			for i := 0; i < particleCount; i++ {
				angle1 := baseAngle + float64(i) * (2 * math.Pi / float64(particleCount))
				angle2 := baseAngle + float64((i+1)%particleCount) * (2 * math.Pi / float64(particleCount))
				
				x1 := g.player.x + math.Cos(angle1) * particleDistance
				y1 := g.player.y + math.Sin(angle1) * particleDistance
				x2 := g.player.x + math.Cos(angle2) * particleDistance
				y2 := g.player.y + math.Sin(angle2) * particleDistance
				
				// 線の色は耐久値に応じて変える
				lineColor := shieldColor
				lineColor.A = 150 // 線は少し濃く
				
				ebitenutil.DrawLine(screen, x1, y1, x2, y2, lineColor)
			}
			
			// シールドの耐久値を表示
			shieldText := fmt.Sprintf("%d", g.player.shield)
			ebitenutil.DebugPrintAt(screen, shieldText, int(g.player.x)-3, int(g.player.y)-3)
		}
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
// シールドアイテムをランダムな位置に生成する
func (g *Game) spawnShieldItem() {
	// 画面内のランダムな位置に配置
	g.shieldItem.x = rand.Float64() * (screenWidth - 2*shieldItemSize) + shieldItemSize
	g.shieldItem.y = rand.Float64() * (screenHeight - 2*shieldItemSize) + shieldItemSize
	g.shieldItem.active = true
}
