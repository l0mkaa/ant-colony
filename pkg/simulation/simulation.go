package simulation

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// ObjectType ...
type ObjectType int

// Типы объектов
const (
	_ ObjectType = iota
	ANT
	FOOD
	PHEROMONEFOOD
	PHEROMONEHOME
	HOME
)

// Object ...
type Object interface {
	// Уникальный идентификатор объекта симуляции
	GetID() string
	// Текущая позиция объекта
	GetPosition() Coordinates
	// Тип обекта (ANT, FOOD ...)
	GetType() ObjectType
	// изменение состояния объекта на итерации симуляции
	process()
	// Жив ли объект
	IsDead() bool
}

// Simulation ...
type Simulation struct {
	vars          Vars
	width, height int
	Home          Coordinates
	data          *data
	steps         int
}

// Vars Переменные симуляции
// AntCount количество муравьев в симуляции
// Lifespan время жизни муравья
// Sight насколько далеко видит муравей
// FoodPheremoneDecay, HomePheremoneDecay время распада феромонов
type Vars struct {
	AntCount           int
	Lifespan           int
	Sight              int
	FoodPheremoneDecay float64
	HomePheremoneDecay float64
}

var defaultSimulationVars = Vars{100, 1000, 3, 0.9, 0.9}

// NewSimulation создает новую симуляцию.
// Если не указана какая либо переменная, то все переменные берутся по дефолту
func NewSimulation(w, h int, vars Vars) *Simulation {
	if w < 100 {
		w = 100
	}
	if h < 100 {
		h = 100
	}

	data := newData(w, h)
	if vars.AntCount == 0 || vars.Lifespan == 0 || vars.Sight == 0 || vars.FoodPheremoneDecay == 0 || vars.HomePheremoneDecay == 0 {
		vars = defaultSimulationVars
	}
	s := &Simulation{
		vars:  vars, //SimulationVars{1000, 3, 0.9, 0.9},
		width: w, height: h,
		data:  data,
		steps: 0,
	}

	homeCoord := Coordinates{X: int(math.Floor(float64(w) / 2)), Y: int(math.Floor(float64(h) / 2))}
	for x := homeCoord.X; x <= homeCoord.X+5; x++ {
		for y := homeCoord.Y; y <= homeCoord.Y+5; y++ {
			u := uuid()
			s.data.addObject(&home{u, Coordinates{x, y}})
		}
	}
	s.Home = homeCoord

	foodCoord := Coordinates{X: homeCoord.X, Y: homeCoord.Y - 30}
	s.AddFood(foodCoord)

	for i := 0; i < s.vars.AntCount; i++ {
		u := uuid()
		s.data.addObject(newAnt(u, s.Home, s.vars.Sight, s.vars.Lifespan,
			s.isCell, s.cellContents, s.takeFood, s.addPheromone, s.shouldRespawned))
	}

	return s
}

// AddFood добавить еду по определенным координатам
func (s *Simulation) AddFood(position Coordinates) {
	for x := position.X; x <= position.X+5; x++ {
		for y := position.Y; y <= position.Y+5; y++ {
			u := uuid()
			s.data.addObject(&food{u, Coordinates{x, y}})
		}
	}
}

// Run запускает симуляцию.
// На вход требует канал для остановки цикла симуляции.
// Возвращает канал (int) счетчик итерации симуляции
// и канал для считывания объектов на текущем шаге.
func (s *Simulation) Run(abort <-chan bool) (chan int, chan [][][]Object) {

	obs := make(chan [][][]Object)
	step := make(chan int)

	go func(step chan int, obs chan [][][]Object, abort <-chan bool) {
		defer close(step)
		defer close(obs)
		var mu sync.Mutex
		i := 0
		for {
			select {
			case abort := <-abort:
				if abort {
					return
				}
			}
			i++
			s.step()
			mu.Lock()
			go func() {
				obs <- s.data.objects
				mu.Unlock()
			}()
			step <- i
			time.Sleep(time.Millisecond * 100)
		}
	}(step, obs, abort)
	return step, obs
}
func (s *Simulation) step() {
	for _, row := range s.data.objects {
		for _, c := range row {
			l := len(c)

			for i := 0; i < l; i++ {
				switch c[i].GetType() {
				case PHEROMONEFOOD:
					if c[i].(*PheromoneFood).power < 0.001 {
						s.data.deleteObject(c[i].GetPosition(), c[i].GetID())
						i--
						l--
					}
				case PHEROMONEHOME:
					if c[i].(*PheromoneHome).power < 0.001 {
						s.data.deleteObject(c[i].GetPosition(), c[i].GetID())
						i--
						l--
					}
				}
				if i >= 0 {
					c[i].process()
				}
			}

		}
	}
	s.steps++
}

func (s *Simulation) shouldRespawned() (bool, Coordinates) {
	return randFloat(0, 1000) < 5, s.Home
}

func (s *Simulation) isCell(position Coordinates) bool {
	if position.X < 0 || position.X >= s.width {
		return false
	}
	if position.Y < 0 || position.Y >= s.height {
		return false
	}
	return true
}

func (s *Simulation) cellContents(position Coordinates) ([]Object, int) {
	return s.data.objectsByPosition(position)
}

func (s *Simulation) takeFood(position Coordinates, foodID string) {
	s.data.deleteObject(position, foodID)
}

func (s *Simulation) addPheromone(position Coordinates, t PheromoneType) {
	createPheromone := func(t PheromoneType) {
		u := uuid()
		if t == PHEROMONEFOODTYPE {
			s.data.addObject(newpheromoneFood(u, position, s.vars.FoodPheremoneDecay))
			return
		}
		s.data.addObject(newPheromoneHome(u, position, s.vars.HomePheremoneDecay))
		return
	}

	obs, _ := s.data.objectsByPosition(position)

	for _, v := range obs {
		switch v.GetType() {
		case PHEROMONEFOOD:
			if t == PHEROMONEFOODTYPE {
				v.(*PheromoneFood).power++
				return
			}
		case PHEROMONEHOME:
			if t == PHEROMONEHOMETYPE {
				v.(*PheromoneHome).power++
				return
			}
		}
	}
	createPheromone(t)
}

func distance(p1, p2 Coordinates) float64 {
	first := math.Pow(float64(p2.X-p1.X), 2)
	second := math.Pow(float64(p2.Y-p1.Y), 2)
	return math.Sqrt(first + second)
}

func uuid() (s string) {
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		return uuid()
	}
	s = fmt.Sprintf("%x", b)
	return
}

func randFloat(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}
