package entity

import (
	"image/color"
	"math"
)

// Explosion は爆発エフェクトの構造体
type Explosion struct {
	X, Y      float64
	Radius    float64
	MaxRadius float64
	Alpha     float64
	Duration  float64
	Lifetime  float64
	Active    bool
}

// NewExplosion は新しい爆発エフェクトを作成する
func NewExplosion(x, y, radius float64) *Explosion {
	return &Explosion{
		X:         x,
		Y:         y,
		Radius:    0,
		MaxRadius: radius,
		Alpha:     1.0,
		Duration:  0.5, // 爆発の持続時間（秒）
		Lifetime:  0.5,
		Active:    true,
	}
}

// Update は爆発エフェクトを更新する
func (e *Explosion) Update(deltaTime float64) {
	if !e.Active {
		return
	}
	
	e.Lifetime -= deltaTime
	
	// 爆発の半径を拡大
	progress := 1.0 - (e.Lifetime / e.Duration)
	e.Radius = e.MaxRadius * math.Sin(progress * math.Pi)
	
	// 透明度を徐々に下げる
	e.Alpha = 1.0 - progress
	
	if e.Lifetime <= 0 {
		e.Active = false
	}
}

// GetColor は爆発の色を取得する
func (e *Explosion) GetColor() color.RGBA {
	// 爆発の進行に応じて色を変化させる
	r := uint8(255)
	g := uint8(255 * (1.0 - e.Lifetime/e.Duration))
	b := uint8(100 * (1.0 - e.Lifetime/e.Duration))
	a := uint8(e.Alpha * 180)
	
	return color.RGBA{r, g, b, a}
}
