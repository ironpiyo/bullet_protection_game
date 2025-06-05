package config

// 画面サイズ
const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

// ゲームオブジェクトのサイズ
const (
	PlayerSize    = 10
	BulletSize    = 8
	ShieldItemSize = 15
)

// ゲームロジック関連
const (
	InitialBullets   = 20
	BulletSpeedMin   = 2.0
	BulletSpeedMax   = 5.0
	BulletSpawnRate  = 5  // 1秒あたりの新しい弾の数
	MaxRankingScores = 5  // ランキングに表示するスコア数
	
	// シールド関連の定数
	ShieldDurability = 3     // シールドの耐久値
	ShieldSpawnRate  = 0.05  // シールドアイテムの出現確率（1フレームあたり）
)

// 時間関連
const (
	DeltaTime = 0.016  // 1フレームあたりの時間（約60FPS）
)
