package entity

// Player はプレイヤーの構造体
type Player struct {
	X, Y   float64
	Size   float64
	Shield int // シールドの耐久値
}

// NewPlayer は新しいプレイヤーを作成する
func NewPlayer(x, y, size float64) *Player {
	return &Player{
		X:      x,
		Y:      y,
		Size:   size,
		Shield: 0, // 初期状態ではシールドなし
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
