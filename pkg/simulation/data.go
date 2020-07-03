package simulation

import (
	"container/list"
	"sync"
)

type data struct {
	w, h int
	// феромоны, еда и т.д
	staticObjects *staticObjects
	// муравьи
	dynamicObjects []Object
}

type staticObjects struct {
	obs [][]objectsList
}
type objectsList struct {
	mu sync.RWMutex
	l  list.List
}

func newStaticObjects(w, h int) *staticObjects {
	obs := make([][]objectsList, w)
	for i := range obs {
		obs[i] = make([]objectsList, h)
	}
	return &staticObjects{obs}
}

func (s *staticObjects) addObject(x, y int, o Object) {
	s.obs[x][y].mu.Lock()
	s.obs[x][y].l.PushBack(o)
	s.obs[x][y].mu.Unlock()
}

func (s *staticObjects) objectsByPosition(x, y int) []Object {
	var obs []Object
	s.obs[x][y].mu.Lock()
	defer s.obs[x][y].mu.Unlock()
	for e := s.obs[x][y].l.Front(); e != nil; e = e.Next() {
		obs = append(obs, e.Value.(Object))
	}
	return obs
}

func (s *staticObjects) deleteObject(x, y int, id string) {
	s.obs[x][y].mu.Lock()
	defer s.obs[x][y].mu.Unlock()
	for e := s.obs[x][y].l.Front(); e != nil; e = e.Next() {
		if e.Value.(Object).GetID() == id {
			s.obs[x][y].l.Remove(e)
			return
		}
	}
}

func (s *staticObjects) asSlice() []Object {
	var obs []Object
	for x := range s.obs {
		for y := range s.obs[x] {
			for e := s.obs[x][y].l.Front(); e != nil; e = e.Next() {
				obs = append(obs, e.Value.(Object))
			}
		}
	}
	return obs
}

func (s *staticObjects) _range(f func(o Object) bool) {
	for x := range s.obs {
		for y := range s.obs[x] {
			for e := s.obs[x][y].l.Front(); e != nil; e = e.Next() {
				if !f(e.Value.(Object)) {
					break
				}
			}
		}
	}
}

func newData(w, h int) *data {
	d := &data{w: w, h: h,
		staticObjects:  newStaticObjects(w, h),
		dynamicObjects: make([]Object, 0)}
	return d
}

func (d *data) addObject(o Object) {
	if o.GetType() == ANT {
		d.dynamicObjects = append(d.dynamicObjects, o)
		return
	}
	pos := o.GetPosition()
	// obs, _ := d.staticObjects.LoadOrStore(pos, &sync.Map{})
	//	obs.(*sync.Map).Store(o.GetID(), o)
	d.staticObjects.addObject(pos.X, pos.Y, o)
}

func (d *data) objectsByPosition(position Coordinates) []Object /*map[string]Object*/ {
	x, y := position.X, position.Y
	if x > d.w || x < 0 || y > d.h || y < 0 {
		return nil
	}

	/*m, ok := d.staticObjects.Load(position)
	if !ok || m == nil {
		return nil
	}
	obs := make([]Object, 0, 3)
	m.(*sync.Map).Range(func(key, value interface{}) bool {
		obs = append(obs, value.(Object))
		return true
	})
	return obs*/
	return d.staticObjects.objectsByPosition(x, y)
}

func (d *data) deleteObject(position Coordinates, id string) {
	/*obs, ok := d.staticObjects.Load(position)
	if ok {
		obs.(*sync.Map).Delete(id)
	}*/
	d.staticObjects.deleteObject(position.X, position.Y, id)
}
