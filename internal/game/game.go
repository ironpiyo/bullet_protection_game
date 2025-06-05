package game

import (
	"sort"
	"time"

	"game/internal/config"
	"game/internal/entity"
)

// Game はゲームの状態を管理する構造体
type Game struct {
	Player        *entity.Player
	Bullets       []*entity.Bullet
	ShieldItem    *entity.ShieldItem
	GameOver      bool
	StartTime     time.Time
	CurrentTime   float64
	Scores        []float64
	LastBulletAdd time.Time
	
	// UI効果用の変数
	GameOverAlpha    float64
	GameOverScale    float64
	RankingAppear    float64
	ScoreAnimations  []*entity.ScoreAnimation
	
	// 難易度関連の変数
	Difficulty       int
	LastDifficultyIncrease time.Time
	
	// 爆発関連
	Explosion     *entity.Explosion
}

// NewGame は新しいゲームインスタンスを作成する
func NewGame() *Game {
	g := &Game{
		Player:        entity.NewPlayer(float64(config.ScreenWidth)/2, float64(config.ScreenHeight)/2, config.PlayerSize),
		Bullets:       make([]*entity.Bullet, 0, config.InitialBullets),
		ShieldItem:    entity.NewShieldItem(config.ShieldItemSize),
		GameOver:      false,
		StartTime:     time.Now(),
		CurrentTime:   0,
		Scores:        make([]float64, 0, config.MaxRankingScores),
		LastBulletAdd: time.Now(),
		
		// UI効果の初期化
		GameOverAlpha: 0,
		GameOverScale: 0.5,
		RankingAppear: 0,
		ScoreAnimations: make([]*entity.ScoreAnimation, 0),
		
		// 難易度の初期化
		Difficulty: 1,
		LastDifficultyIncrease: time.Now(),
		
		// 爆発は初期状態ではnil
		Explosion: nil,
	}

	// 初期の弾を生成
	for i := 0; i < config.InitialBullets; i++ {
		g.addRandomBullet()
	}

	return g
}

// addRandomBullet はランダムな位置と速度で新しい弾を追加する
func (g *Game) addRandomBullet() {
	bullet := entity.NewRandomBullet(config.ScreenWidth, config.ScreenHeight, config.BulletSize, config.BulletSpeedMin, config.BulletSpeedMax, g.Difficulty)
	g.Bullets = append(g.Bullets, bullet)
}

// Layout はウィンドウサイズを返す
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return config.ScreenWidth, config.ScreenHeight
}

// Reset はゲームをリセットする（スコアは保持）
func (g *Game) Reset() {
	oldScores := g.Scores
	*g = *NewGame()
	g.Scores = oldScores
}

// AddScore はスコアを追加する
func (g *Game) AddScore(score float64) {
	g.Scores = append(g.Scores, score)
	
	// スコアを降順にソート
	sort.Slice(g.Scores, func(i, j int) bool {
		return g.Scores[i] > g.Scores[j]
	})
	
	// 上位スコアだけを保持
	if len(g.Scores) > config.MaxRankingScores {
		g.Scores = g.Scores[:config.MaxRankingScores]
	}
	
	// スコアアニメーションを追加
	g.ScoreAnimations = append(g.ScoreAnimations, entity.NewScoreAnimation(
		score,
		config.ScreenWidth / 2,
		config.ScreenHeight / 3,
	))
}
