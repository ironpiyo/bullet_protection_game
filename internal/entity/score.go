package entity

// ScoreAnimation はスコアアニメーションの構造体
type ScoreAnimation struct {
	Score     float64
	X, Y      float64
	Scale     float64
	Alpha     float64
	Lifetime  float64
}

// NewScoreAnimation は新しいスコアアニメーションを作成する
func NewScoreAnimation(score, x, y float64) *ScoreAnimation {
	return &ScoreAnimation{
		Score:    score,
		X:        x,
		Y:        y,
		Scale:    1.0,
		Alpha:    1.0,
		Lifetime: 2.0,
	}
}

// Update はスコアアニメーションを更新する
func (s *ScoreAnimation) Update(deltaTime float64) {
	s.Lifetime -= deltaTime
	s.Scale += 0.02
	s.Alpha -= 0.02
}

// IsActive はアニメーションがまだアクティブかどうかを返す
func (s *ScoreAnimation) IsActive() bool {
	return s.Lifetime > 0 && s.Alpha > 0
}
