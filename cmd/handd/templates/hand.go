package templates

type HandViewModel struct {
	HandId    string
	TableId   string
	Opponents []OpponentViewModel
	Player    PlayerViewModel
	IsActive bool
}
