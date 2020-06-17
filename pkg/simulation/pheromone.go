package simulation

type PheromoneType int

// ss
const (
	_ PheromoneType = iota
	PHEROMONEFOODTYPE
	PHEROMONEHOMETYPE
)

type pheromone struct {
	id       string
	position Coordinates
	// pheromoneType  PheromoneType
	power float64
	decay float64
}

func (p *pheromone) GetPower() float64 {
	return p.power
}

func (p *pheromone) GetPosition() Coordinates {
	return p.position
}

func (p *pheromone) GetID() string {
	return p.id
}

func (p *pheromone) process() {
	p.power *= p.decay
}

func (p *pheromone) IsDead() bool { return false }

// -------------------------------------------------------------

type PheromoneFood struct {
	pheromone
}

func newpheromoneFood(id string, position Coordinates, decay float64) *PheromoneFood {
	return &PheromoneFood{pheromone{id, position, 1, decay}}
}

func (pf *PheromoneFood) GetType() ObjectType {
	return PHEROMONEFOOD
}

//--------------------------------------------------------------

type PheromoneHome struct {
	pheromone
}

func newPheromoneHome(id string, position Coordinates, decay float64) *PheromoneHome {
	return &PheromoneHome{pheromone{id, position, 1, decay}}
}

func (ph *PheromoneHome) GetType() ObjectType {
	return PHEROMONEHOME
}
