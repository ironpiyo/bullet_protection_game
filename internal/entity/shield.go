package entity

import (
	"math/rand"
)

// ShieldItem はシールドアイテムの構造体
type ShieldItem struct {
	X, Y      float64
	Size      float64
	Active    bool
	Angle     float64  // 回転角度
	GlowSize  float64  // 輝きのサイズ
	GlowDir   float64  // 輝きの方向（拡大/縮小）
}

// NewShieldItem は新しいシールドアイテムを作成する
func NewShieldItem(size float64) *ShieldItem {
	return &ShieldItem{
		X:        -100, // 画面外に配置
		Y:        -100,
		Size:     size,
		Active:   false,
		Angle:    0,
		GlowSize: 0,
		GlowDir:  0.1,
	}
}

// Spawn はシールドアイテムをランダムな位置に生成する
func (s *ShieldItem) Spawn(screenWidth, screenHeight float64) {
	// 画面内のランダムな位置に配置
	s.X = rand.Float64() * (screenWidth - 2*s.Size) + s.Size
	s.Y = rand.Float64() * (screenHeight - 2*s.Size) + s.Size
	s.Active = true
}

// Update はシールドアイテムのアニメーションを更新する
func (s *ShieldItem) Update() {
	if s.Active {
		// 回転させる
		s.Angle += 0.05
		
		// 輝きのサイズを変化させる
		s.GlowSize += s.GlowDir
		if s.GlowSize > 3 || s.GlowSize < 0 {
			s.GlowDir *= -1
		}
	}
}

// CollidesWith はシールドアイテムが指定された座標と衝突するかどうかを判定する
func (s *ShieldItem) CollidesWith(x, y, size float64) bool {
	if !s.Active {
		return false
	}
	
	dx := s.X - x
	dy := s.Y - y
	distance := dx*dx + dy*dy
	return distance < (s.Size+size)*(s.Size+size)
}

// Deactivate はシールドアイテムを非アクティブにする
func (s *ShieldItem) Deactivate() {
	s.Active = false
	s.X = -100 // 画面外に移動
	s.Y = -100
}
