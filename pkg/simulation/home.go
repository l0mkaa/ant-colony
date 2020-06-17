package simulation

type home struct {
	id       string
	position Coordinates
}

func (h *home) GetID() string {
	return h.id
}
func (h *home) GetPosition() Coordinates {
	return h.position
}
func (h *home) GetType() ObjectType {
	return HOME
}
func (h *home) process() {
}

func (h *home) IsDead() bool { return false }
