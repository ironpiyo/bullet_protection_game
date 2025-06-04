package entity

import (
	"image/color"
	"math"
	"math/rand"
)

// Bullet は弾の構造体
type Bullet struct {
	X, Y    float64
	VX, VY  float64
	Size    float64
	Color   color.RGBA
}

// NewRandomBullet は画面の端から発射されるランダムな弾を作成する
func NewRandomBullet(screenWidth, screenHeight, bulletSize, minSpeed, maxSpeed float64, difficulty int) *Bullet {
	var x, y float64
	var vx, vy float64
	
	side := rand.Intn(4) // 0: 上, 1: 右, 2: 下, 3: 左
	
	// 難易度に応じて弾の速度を調整
	speedMultiplier := 1.0 + float64(difficulty-1)*0.1 // 難易度ごとに10%ずつ速くなる
	minSpeed = minSpeed * speedMultiplier
	maxSpeed = maxSpeed * speedMultiplier
	
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
	g := uint8(rand.Intn(200) + 55)
	b := uint8(rand.Intn(200) + 55)
	
	return &Bullet{
		X:    x,
		Y:    y,
		VX:   vx,
		VY:   vy,
		Size: bulletSize,
		Color: color.RGBA{r, g, b, 255},
	}
}

// Update は弾の位置を更新する
func (b *Bullet) Update() {
	b.X += b.VX
	b.Y += b.VY
}

// IsOutOfScreen は弾が画面外に出たかどうかを判定する
func (b *Bullet) IsOutOfScreen(screenWidth, screenHeight float64, margin float64) bool {
	return b.X < -margin || b.X > screenWidth+margin || b.Y < -margin || b.Y > screenHeight+margin
}

// CollidesWith は弾が指定された座標と衝突するかどうかを判定する
func (b *Bullet) CollidesWith(x, y, size float64) bool {
	dx := b.X - x
	dy := b.Y - y
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < b.Size+size
}
