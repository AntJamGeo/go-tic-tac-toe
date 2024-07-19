package player

type Player struct {
	id      string
	name    string
	channel chan *map[string]string
}

func NewPlayer(id string, name string) *Player {
	return &Player{
		id:      id,
		name:    name,
		channel: make(chan *map[string]string),
	}
}
