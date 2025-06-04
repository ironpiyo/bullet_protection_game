package game

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	
	"game/internal/entity"
)

// Update はゲームの状態を更新する
func (g *Game) Update() error {
	// ゲームオーバー時のリスタート処理とアニメーション
	if g.GameOver {
		return g.updateGameOver()
	}

	// プレイヤーの位置をマウスカーソルに合わせる
	x, y := ebiten.CursorPosition()
	g.Player.X = float64(x)
	g.Player.Y = float64(y)

	// 経過時間を更新
	g.CurrentTime = time.Since(g.StartTime).Seconds()
	
	// 難易度の更新
	g.updateDifficulty()
	
	// 弾の生成
	g.updateBulletSpawn()
	
	// シールドアイテムの更新
	g.updateShieldItem()
	
	// スコアアニメーションの更新
	g.updateScoreAnimations()

	// 弾の移動と衝突判定
	g.updateBullets()

	return nil
}

// updateGameOver はゲームオーバー時の更新処理
func (g *Game) updateGameOver() error {
	// ゲームオーバーアニメーションの更新
	if g.GameOverAlpha < 0.8 {
		g.GameOverAlpha += 0.02
	}
	
	if g.GameOverScale < 1.2 {
		g.GameOverScale += 0.03
	} else if g.GameOverScale > 1.2 {
		g.GameOverScale = 1.2
	}
	
	// ランキング表示のアニメーション
	if g.RankingAppear < 1.0 {
		g.RankingAppear += 0.03
	}
	
	// リスタート処理
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.Reset()
	}
	
	return nil
}

// updateDifficulty は難易度を更新する
func (g *Game) updateDifficulty() {
	// 6秒ごとに難易度を上げる
	if time.Since(g.LastDifficultyIncrease).Seconds() > 6.0 {
		g.Difficulty++
		g.LastDifficultyIncrease = time.Now()
		
		// デバッグ用に難易度上昇を表示
		log.Printf("難易度上昇: レベル %d", g.Difficulty)
	}
}

// updateBulletSpawn は弾の生成を更新する
func (g *Game) updateBulletSpawn() {
	// 難易度に応じて弾の発生頻度を調整
	bulletSpawnInterval := 1.0 / float64(BulletSpawnRate)
	if time.Since(g.LastBulletAdd).Seconds() > bulletSpawnInterval {
		// 難易度に応じて複数の弾を発射
		bulletsToAdd := g.Difficulty
		for i := 0; i < bulletsToAdd; i++ {
			g.addRandomBullet()
		}
		g.LastBulletAdd = time.Now()
	}
}

// updateShieldItem はシールドアイテムを更新する
func (g *Game) updateShieldItem() {
	// シールドアイテムの生成（ランダムに）
	if !g.ShieldItem.Active && rand.Float64() < ShieldSpawnRate {
		g.ShieldItem.Spawn(ScreenWidth, ScreenHeight)
	}
	
	// シールドアイテムのアニメーション更新
	g.ShieldItem.Update()
	
	// シールドアイテムとプレイヤーの衝突判定
	if g.ShieldItem.CollidesWith(g.Player.X, g.Player.Y, g.Player.Size) {
		// シールドを獲得
		g.Player.AddShield(ShieldDurability)
		g.ShieldItem.Deactivate()
		
		// デバッグ用にシールド獲得を表示
		log.Printf("シールド獲得！ 耐久値: %d", g.Player.Shield)
	}
}

// updateScoreAnimations はスコアアニメーションを更新する
func (g *Game) updateScoreAnimations() {
	newScoreAnims := make([]*entity.ScoreAnimation, 0)
	for _, anim := range g.ScoreAnimations {
		anim.Update(0.016) // 約60FPSを想定
		
		if anim.IsActive() {
			newScoreAnims = append(newScoreAnims, anim)
		}
	}
	g.ScoreAnimations = newScoreAnims
}

// updateBullets は弾の移動と衝突判定を更新する
func (g *Game) updateBullets() {
	newBullets := make([]*entity.Bullet, 0, len(g.Bullets))
	for _, b := range g.Bullets {
		// 弾を移動
		b.Update()
		
		// 画面外に出た弾は削除
		if b.IsOutOfScreen(ScreenWidth, ScreenHeight, 100) {
			continue
		}
		
		// プレイヤーとの衝突判定
		if b.CollidesWith(g.Player.X, g.Player.Y, g.Player.Size) {
			// シールドがある場合
			if g.Player.HasShield() {
				g.Player.ReduceShield()
				log.Printf("シールドが弾を防いだ！ 残り耐久値: %d", g.Player.Shield)
				continue // この弾は消える
			} else {
				// シールドがない場合、ゲームオーバー
				g.GameOver = true
				
				// スコアを記録
				g.AddScore(g.CurrentTime)
				
				break
			}
		}
		
		newBullets = append(newBullets, b)
	}
	
	if !g.GameOver {
		g.Bullets = newBullets
	}
}
