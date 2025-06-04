package game

import (
	"sort"
	"time"

	"game/internal/entity"
)

const (
	ScreenWidth      = 800
	ScreenHeight     = 600
	PlayerSize       = 10
	BulletSize       = 8
	InitialBullets   = 20
	BulletSpeedMin   = 2.0
	BulletSpeedMax   = 5.0
	BulletSpawnRate  = 5  // 1秒あたりの新しい弾の数
	MaxRankingScores = 5  // ランキングに表示するスコア数
	
	// シールド関連の定数
	ShieldDurability = 3  // シールドの耐久値
	ShieldItemSize   = 15 // シールドアイテムのサイズ
	ShieldSpawnRate  = 0.05 // シールドアイテムの出現確率（1フレームあたり）
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
}

// NewGame は新しいゲームインスタンスを作成する
func NewGame() *Game {
	g := &Game{
		Player:        entity.NewPlayer(float64(ScreenWidth)/2, float64(ScreenHeight)/2, PlayerSize),
		Bullets:       make([]*entity.Bullet, 0, InitialBullets),
		ShieldItem:    entity.NewShieldItem(ShieldItemSize),
		GameOver:      false,
		StartTime:     time.Now(),
		CurrentTime:   0,
		Scores:        make([]float64, 0, MaxRankingScores),
		LastBulletAdd: time.Now(),
		
		// UI効果の初期化
		GameOverAlpha: 0,
		GameOverScale: 0.5,
		RankingAppear: 0,
		ScoreAnimations: make([]*entity.ScoreAnimation, 0),
		
		// 難易度の初期化
		Difficulty: 1,
		LastDifficultyIncrease: time.Now(),
	}

	// 初期の弾を生成
	for i := 0; i < InitialBullets; i++ {
		g.addRandomBullet()
	}

	return g
}

// addRandomBullet はランダムな位置と速度で新しい弾を追加する
func (g *Game) addRandomBullet() {
	bullet := entity.NewRandomBullet(ScreenWidth, ScreenHeight, BulletSize, BulletSpeedMin, BulletSpeedMax, g.Difficulty)
	g.Bullets = append(g.Bullets, bullet)
}

// Layout はウィンドウサイズを返す
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
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
	if len(g.Scores) > MaxRankingScores {
		g.Scores = g.Scores[:MaxRankingScores]
	}
	
	// スコアアニメーションを追加
	g.ScoreAnimations = append(g.ScoreAnimations, entity.NewScoreAnimation(
		score,
		ScreenWidth / 2,
		ScreenHeight / 3,
	))
}
