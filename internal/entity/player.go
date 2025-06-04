package entity

// Player はプレイヤーの構造体
type Player struct {
	X, Y   float64
	Size   float64
	Shield int // シールドの耐久値
	
	// 爆発スキル関連
	BombAvailable    bool    // 爆発スキルが使用可能かどうか
	BombCooldown     float64 // クールダウン残り時間
	BombCooldownMax  float64 // クールダウン最大時間
	BombRadius       float64 // 爆発の半径
}

// NewPlayer は新しいプレイヤーを作成する
func NewPlayer(x, y, size float64) *Player {
	return &Player{
		X:               x,
		Y:               y,
		Size:            size,
		Shield:          0, // 初期状態ではシールドなし
		BombAvailable:   true,
		BombCooldown:    0,
		BombCooldownMax: 10.0, // 10秒のクールダウン
		BombRadius:      150.0, // 爆発の半径
	}
}

// AddShield はプレイヤーにシールドを追加する
func (p *Player) AddShield(durability int) {
	p.Shield = durability
}

// ReduceShield はシールドの耐久値を減らす
func (p *Player) ReduceShield() {
	if p.Shield > 0 {
		p.Shield--
	}
}

// HasShield はシールドを持っているかどうかを返す
func (p *Player) HasShield() bool {
	return p.Shield > 0
}

// UseBomb は爆発スキルを使用する
func (p *Player) UseBomb() bool {
	if p.BombAvailable {
		p.BombAvailable = false
		p.BombCooldown = p.BombCooldownMax
		return true
	}
	return false
}

// UpdateBombCooldown はクールダウンを更新する
func (p *Player) UpdateBombCooldown(deltaTime float64) {
	if !p.BombAvailable {
		p.BombCooldown -= deltaTime
		if p.BombCooldown <= 0 {
			p.BombAvailable = true
			p.BombCooldown = 0
		}
	}
}
