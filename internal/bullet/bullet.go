package bullet

const (
	ShotByPlayer = "shot_by_player"
	ShotByEnemy  = "shot_by_enemy"
)

type Bullet struct {
	X, Y   int
	GoingY int
	ShotBy string
}