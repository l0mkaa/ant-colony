package simulation

import (
	"math"
)

var directions = []Coordinates{
	//N
	{0, -1},
	//NE
	{1, -1},
	//E
	{1, 0},
	//SE
	{1, 1},
	//S
	{0, 1},
	//SW
	{-1, 1},
	//W
	{-1, 0},
	//NW
	{-1, -1},
}

// Ant ..
type Ant struct {
	id string
	// Позиция
	position Coordinates
	// Дистанция взгляда
	sight int
	// Время жизни
	lifeSpan int
	// Количество шагов симуляции
	steps int
	// Угол куда смотрит
	angle int
	// Несет ли еду
	CarryingFood bool
	// Принадлежит ли ячейка куда смотрит муравей симуляции
	isCell func(Coordinates) bool
	// Содержимое ячейки куда смотрит муравей
	cellContents func(Coordinates) ([]Object, int)
	// Муравей взял еду
	takeFood func(position Coordinates, foodID string)
	// Муравей оставил феромон
	addPheromone func(Coordinates, PheromoneType)
	dead         bool
	// Надо ли муравью идти на респаун
	shouldRespawned func() (bool, Coordinates)
}

func newAnt(
	id string,
	position Coordinates,
	sight int,
	lifeSpan int,
	isCell func(Coordinates) bool,
	cellContents func(Coordinates) ([]Object, int),
	takeFood func(position Coordinates, foodID string),
	addPheromone func(Coordinates, PheromoneType),
	shouldRespawned func() (bool, Coordinates),
) *Ant {
	return &Ant{
		id:              id,
		position:        position,
		sight:           sight,
		lifeSpan:        lifeSpan,
		steps:           0,
		angle:           0,
		CarryingFood:    false,
		isCell:          isCell,
		cellContents:    cellContents,
		takeFood:        takeFood,
		addPheromone:    addPheromone,
		shouldRespawned: shouldRespawned, dead: false,
	}
}

func (a *Ant) GetID() string {
	return a.id
}

func (a *Ant) GetPosition() Coordinates {
	return a.position
}
func (a *Ant) GetType() ObjectType {
	return ANT
}
func (a *Ant) IsDead() bool { return a.dead }
func (a *Ant) stepInc() {
	a.steps++
}

func (a *Ant) process() {
	defer a.stepInc()

	if a.IsDead() {
		r, pos := a.shouldRespawned()
		if r {
			a.randomizeDirection()
			a.dead = false
			a.steps = 0
			a.position = pos
		}
		return
	}

	if float64(a.steps) > float64(a.lifeSpan)*randFloat(1, 2) {
		a.dead = true
		return
	}

	startPosition := a.position
	sensed := a.sensed()
	fwd := sensed[1]

	if ok := a.isCell(Coordinates{a.position.X + fwd.X, a.position.Y + fwd.Y}); !ok {
		a.randomizeDirection()
		return
	}
	fwdCellContents, _ := a.cellContents(Coordinates{a.position.X + fwd.X, a.position.Y + fwd.Y})
	containsType := func(ctc []Object, t ObjectType) (int, bool) {
		for i, v := range ctc {
			if v.GetType() == t {
				return i, true
			}
		}
		return -1, false
	}

	if a.CarryingFood {
		if _, ok := containsType(fwdCellContents, HOME); ok {
			a.CarryingFood = false
			a.steps = 0

			a.turnRight()
			a.turnRight()
			a.turnRight()
			a.turnRight()

			a.forageForFood()

		} else {
			a.lookForHome()
		}
	} else {
		if index, ok := containsType(fwdCellContents, FOOD); ok {
			a.CarryingFood = true
			a.turnRight()
			a.turnRight()
			a.turnRight()
			a.turnRight()

			a.steps = 0

			a.takeFood(fwdCellContents[index].GetPosition(), fwdCellContents[index].GetID())

			a.lookForHome()

		} else {
			a.forageForFood()
		}
	}

	if !a.IsDead() && a.position.X != startPosition.Y && a.position.Y != startPosition.Y {
		if a.CarryingFood {
			a.addPheromone(a.position, PHEROMONEFOODTYPE)
		} else {
			a.addPheromone(a.position, PHEROMONEHOMETYPE)
		}
	}
}

func (a *Ant) turnLeft() {
	a.angle--
	if a.angle < 0 {
		a.angle = len(directions) - 1
	}
}
func (a *Ant) turnRight() {
	a.angle++
	a.angle = a.angle % len(directions)
}
func (a *Ant) forward() Coordinates {
	forward := directions[a.angle]
	return forward
}

func (a *Ant) randomizeDirection() {
	max := float64(len(directions))
	min := 0.0
	a.angle = int(math.Floor(randFloat(min, max)))
}

func (a *Ant) walkRandomly() {
	forward := a.forward()
	action := int(math.Floor(randFloat(0, 6)))
	//Slightly more likely to move forwards than to turn
	if action < 4 {
		a.position.X += forward.X
		a.position.Y += forward.Y
	} else if action == 4 {
		a.turnLeft()
	} else if action == 5 {
		a.turnRight()
	}
}

func (a *Ant) sensed() []Coordinates {
	forward := a.forward()

	i := 0
	for _, direction := range directions {
		if direction == forward {
			break
		}
		i++
	}

	d := 0
	if i > 0 {
		d = i - 1
	} else {
		d = len(directions) - 1
	}
	forwardLeft := directions[d]

	forwardRight := directions[(i+1)%len(directions)]

	return []Coordinates{forwardLeft, forward, forwardRight}
}

func (a *Ant) seek(isFood bool) {
	sensed := a.sensed()
	forwardLeft := sensed[0]
	forward := sensed[1]
	forwardRight := sensed[2]

	maxScore := 0.0
	bestDirection := forward

	scores := []float64{}

	for i := range sensed {
		direction := sensed[i]
		score := a.getScoreForDirection(a.position, direction, isFood)
		scores = append(scores, score)
		if score > maxScore {
			maxScore = score
			bestDirection = direction
		}
	}
	if float64(maxScore) < 0.01 {
		a.walkRandomly()
		return
	}

	if bestDirection == forwardRight {
		a.turnRight()
		return
	} else if bestDirection == forwardLeft {
		a.turnLeft()
		return
	}

	a.position.X += forward.X
	a.position.Y += forward.Y
}

func (a *Ant) scoreForCell(position Coordinates, isFood bool) float64 {
	obs, count := a.cellContents(position)
	if count == 0 {
		return 0
	}
	containsObjectByType := func(obs []Object, t ObjectType) (int, bool) {
		for i, o := range obs {
			if o.GetType() == t {
				return i, true
			}
		}
		return -1, false
	}
	if isFood {

		_, ok := containsObjectByType(obs, FOOD)
		if ok {
			return 100
		}
		i, ok := containsObjectByType(obs, PHEROMONEFOOD)
		if !ok {
			return 0
		}
		return float64(obs[i].(*PheromoneFood).power)
	}
	_, ok := containsObjectByType(obs, HOME)
	if ok {
		return 100
	}
	i, ok := containsObjectByType(obs, PHEROMONEHOME)
	if !ok {
		return 0
	}
	return float64(obs[i].(*PheromoneHome).power)

}

func (a *Ant) getScoreForDirection(position Coordinates, direction Coordinates, isFood bool) float64 {
	r := a.sight

	x0 := position.X + direction.X*r
	y0 := position.Y + direction.Y*r
	score := 0.0

	for x := x0 - r/2; x <= x0+(r/2); x++ {
		for y := y0 - (r / 2); y <= y0+(r/2); y++ {
			wScore := 0.0
			if ok := a.isCell(Coordinates{x, y}); ok {
				wScore = a.scoreForCell(Coordinates{x, y}, isFood)
			}
			wScore /= (distance(Coordinates{x0, y0}, Coordinates{x, y}) + 1)
			score += wScore
		}
	}

	fwd := Coordinates{position.X + direction.X, position.Y + direction.Y}
	if a.isCell(fwd) {
		score += a.scoreForCell(fwd, isFood)
	}

	return score
}

func (a *Ant) forageForFood() {
	a.seek(true)
}
func (a *Ant) lookForHome() {
	a.seek(false)
}
