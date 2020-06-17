package simulation

type food struct {
	id       string
	position Coordinates
}

func (f *food) GetID() string {
	return f.id
}
func (f *food) GetPosition() Coordinates {
	return f.position
}
func (f *food) GetType() ObjectType {
	return FOOD
}
func (f *food) process() {
}

func (f *food) IsDead() bool { return false }
