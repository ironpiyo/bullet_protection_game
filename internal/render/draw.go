package render

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	"game/internal/config"
	"game/internal/entity"
	"game/internal/game"
)

// Draw はゲームの状態を描画する
func Draw(screen *ebiten.Image, g *game.Game) {
	// 背景を黒で塗りつぶす
	screen.Fill(color.RGBA{20, 20, 40, 255})

	// シールドアイテムを描画
	drawShieldItem(screen, g.ShieldItem)

	// 弾を描画
	for _, b := range g.Bullets {
		drawBullet(screen, b)
	}

	// 爆発エフェクトを描画
	if g.Explosion != nil && g.Explosion.Active {
		drawExplosion(screen, g.Explosion)
	}

	// プレイヤーを描画
	if !g.GameOver {
		drawPlayer(screen, g.Player, g.CurrentTime)
	}

	// 経過時間と難易度を表示
	timeText := fmt.Sprintf("Time: %.2f  Difficulty: %d", g.CurrentTime, g.Difficulty)
	ebitenutil.DebugPrintAt(screen, timeText, 20, 20)

	// 爆発スキルのクールダウン表示
	drawBombCooldown(screen, g.Player)

	// スコアアニメーションを描画
	for _, anim := range g.ScoreAnimations {
		drawScoreAnimation(screen, anim)
	}

	// ゲームオーバー時の表示
	if g.GameOver {
		drawGameOver(screen, g)
	}
}

// drawPlayer はプレイヤーを描画する
func drawPlayer(screen *ebiten.Image, player *entity.Player, currentTime float64) {
	// プレイヤーを描画（白い円）
	ebitenutil.DrawCircle(screen, player.X, player.Y, player.Size, color.RGBA{255, 255, 255, 255})
	
	// シールドがある場合、プレイヤーの周りにシールドを描画
	if player.HasShield() {
		drawPlayerShield(screen, player, currentTime)
	}
}

// drawPlayerShield はプレイヤーのシールドを描画する
func drawPlayerShield(screen *ebiten.Image, player *entity.Player, currentTime float64) {
	// シールドの色は耐久値によって変化
	var shieldColor color.RGBA
	switch player.Shield {
	case 3:
		shieldColor = color.RGBA{0, 255, 255, 100} // 明るい水色（半透明）
	case 2:
		shieldColor = color.RGBA{100, 200, 255, 100} // やや暗い水色（半透明）
	case 1:
		shieldColor = color.RGBA{150, 150, 255, 100} // 紫がかった色（半透明）
	}
	
	// シールドを描画（プレイヤーと同じサイズだが、半透明）
	ebitenutil.DrawCircle(screen, player.X, player.Y, player.Size, shieldColor)
	
	// 回転する16角形のエフェクト
	particleCount := 16 // 16角形
	baseAngle := currentTime * 1.5 // 時間に基づいて回転（速度調整）
	
	// エフェクトとプレイヤーの間の距離を広げる
	particleDistance := player.Size * 2.0 // プレイヤーからの距離を2倍に
	
	// 16角形の頂点を描画
	for i := 0; i < particleCount; i++ {
		angle := baseAngle + float64(i) * (2 * math.Pi / float64(particleCount))
		
		// 粒子の位置を計算
		particleX := player.X + math.Cos(angle) * particleDistance
		particleY := player.Y + math.Sin(angle) * particleDistance
		
		// 粒子を描画
		particleSize := 2.5
		particleColor := color.RGBA{255, 255, 255, 200} // 白い光の粒子
		ebitenutil.DrawCircle(screen, particleX, particleY, particleSize, particleColor)
	}
	
	// 16角形の辺を描画（頂点同士を線で結ぶ）
	for i := 0; i < particleCount; i++ {
		angle1 := baseAngle + float64(i) * (2 * math.Pi / float64(particleCount))
		angle2 := baseAngle + float64((i+1)%particleCount) * (2 * math.Pi / float64(particleCount))
		
		x1 := player.X + math.Cos(angle1) * particleDistance
		y1 := player.Y + math.Sin(angle1) * particleDistance
		x2 := player.X + math.Cos(angle2) * particleDistance
		y2 := player.Y + math.Sin(angle2) * particleDistance
		
		// 線の色は耐久値に応じて変える
		lineColor := shieldColor
		lineColor.A = 150 // 線は少し濃く
		
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, lineColor)
	}
	
	// シールドの耐久値を表示
	shieldText := fmt.Sprintf("%d", player.Shield)
	ebitenutil.DebugPrintAt(screen, shieldText, int(player.X)-3, int(player.Y)-3)
}

// drawBullet は弾を描画する
func drawBullet(screen *ebiten.Image, bullet *entity.Bullet) {
	ebitenutil.DrawCircle(screen, bullet.X, bullet.Y, bullet.Size, bullet.Color)
}

// drawShieldItem はシールドアイテムを描画する
func drawShieldItem(screen *ebiten.Image, shieldItem *entity.ShieldItem) {
	if !shieldItem.Active {
		return
	}
	
	// 外側の輝き
	glowColor := color.RGBA{0, 255, 255, 100}
	ebitenutil.DrawCircle(screen, shieldItem.X, shieldItem.Y, 
						 shieldItem.Size + shieldItem.GlowSize, glowColor)
	
	// メインの円
	ebitenutil.DrawCircle(screen, shieldItem.X, shieldItem.Y, 
						 shieldItem.Size, color.RGBA{0, 200, 255, 200})
	
	// 内側の円
	ebitenutil.DrawCircle(screen, shieldItem.X, shieldItem.Y, 
						 shieldItem.Size * 0.6, color.RGBA{100, 200, 255, 255})
	
	// 回転する星型のエフェクト
	angle := shieldItem.Angle
	radius := shieldItem.Size * 0.8
	for i := 0; i < 5; i++ {
		a := angle + float64(i) * math.Pi * 0.4
		x1 := shieldItem.X + math.Cos(a) * radius
		y1 := shieldItem.Y + math.Sin(a) * radius
		x2 := shieldItem.X + math.Cos(a+0.2) * (radius * 0.5)
		y2 := shieldItem.Y + math.Sin(a+0.2) * (radius * 0.5)
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.RGBA{255, 255, 255, 200})
	}
}

// drawExplosion は爆発エフェクトを描画する
func drawExplosion(screen *ebiten.Image, explosion *entity.Explosion) {
	// 爆発の円を描画
	ebitenutil.DrawCircle(screen, explosion.X, explosion.Y, explosion.Radius, explosion.GetColor())
	
	// 爆発の波紋を描画
	waveColor := explosion.GetColor()
	waveColor.A = uint8(float64(waveColor.A) * 0.5)
	ebitenutil.DrawCircle(screen, explosion.X, explosion.Y, explosion.Radius * 0.8, waveColor)
	
	// 中心の明るい部分
	centerColor := color.RGBA{255, 255, 200, uint8(explosion.Alpha * 255)}
	ebitenutil.DrawCircle(screen, explosion.X, explosion.Y, explosion.Radius * 0.3, centerColor)
}

// drawBombCooldown はボムのクールダウンを表示する
func drawBombCooldown(screen *ebiten.Image, player *entity.Player) {
	// クールダウン表示の位置
	x, y := 20, 40
	width := 100.0
	height := 10.0
	
	// 背景バー
	ebitenutil.DrawRect(screen, float64(x), float64(y), width, height, color.RGBA{50, 50, 50, 200})
	
	// クールダウン進行バー
	if !player.BombAvailable {
		progress := 1.0 - (player.BombCooldown / player.BombCooldownMax)
		ebitenutil.DrawRect(screen, float64(x), float64(y), width * progress, height, color.RGBA{0, 200, 255, 200})
	} else {
		// 使用可能時は満タン
		ebitenutil.DrawRect(screen, float64(x), float64(y), width, height, color.RGBA{0, 255, 255, 200})
	}
	
	// テキスト表示
	bombText := "BOMB [X]"
	ebitenutil.DebugPrintAt(screen, bombText, x, y-5)
}

// drawScoreAnimation はスコアアニメーションを描画する
func drawScoreAnimation(screen *ebiten.Image, anim *entity.ScoreAnimation) {
	// スケールと透明度に基づいて描画
	scoreText := fmt.Sprintf("%.2f seconds", anim.Score)
	
	// 文字サイズを計算（スケールに応じて）
	textWidth := float64(len(scoreText) * 6) * anim.Scale
	textX := anim.X - textWidth/2
	
	// テキストを描画（簡易的なスケーリング）
	for i := 0; i < int(anim.Scale*2); i++ {
		ebitenutil.DebugPrintAt(screen, scoreText, int(textX), int(anim.Y)+i)
	}
}

// drawGameOver はゲームオーバー画面を描画する
func drawGameOver(screen *ebiten.Image, g *game.Game) {
	// 半透明のオーバーレイを描画
	overlayColor := color.RGBA{0, 0, 0, uint8(g.GameOverAlpha * 200)}
	ebitenutil.DrawRect(screen, 0, 0, float64(config.ScreenWidth), float64(config.ScreenHeight), overlayColor)
	
	// ゲームオーバーテキストを描画（拡大/縮小アニメーション付き）
	gameOverText := "GAME OVER"
	textWidth := float64(len(gameOverText) * 8) * g.GameOverScale
	textX := (float64(config.ScreenWidth) - textWidth) / 2
	textY := float64(config.ScreenHeight)/3 - 20
	
	// テキストを複数回描画してボールド効果を出す
	for i := 0; i < int(g.GameOverScale*3); i++ {
		ebitenutil.DebugPrintAt(screen, gameOverText, int(textX), int(textY)+i)
	}
	
	// リスタート案内
	restartText := "Press SPACE to restart"
	restartX := config.ScreenWidth/2 - len(restartText)*3
	restartY := int(textY) + 40
	ebitenutil.DebugPrintAt(screen, restartText, restartX, restartY)
	
	// ランキングを表示（徐々に表示されるアニメーション）
	if g.RankingAppear > 0 {
		rankingTitleY := config.ScreenHeight/2 + 30
		ebitenutil.DebugPrintAt(screen, "TOP SCORES:", config.ScreenWidth/2-50, rankingTitleY)
		
		// 各スコアを表示（徐々に表示）
		maxScoresToShow := int(float64(len(g.Scores)) * g.RankingAppear)
		for i := 0; i < maxScoresToShow && i < len(g.Scores); i++ {
			scoreText := fmt.Sprintf("%d. %.2f seconds", i+1, g.Scores[i])
			
			// アニメーション効果（少しずつ右から現れる）
			offset := int((1.0 - g.RankingAppear) * 100)
			if offset < 0 {
				offset = 0
			}
			
			ebitenutil.DebugPrintAt(screen, scoreText, config.ScreenWidth/2-50+offset, rankingTitleY+20+i*20)
		}
	}
}
