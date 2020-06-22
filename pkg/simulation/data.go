package simulation

import "sync"

type data struct {
	w, h int
	// objects [][][]Object
	// феромоны, еда и т.д
	staticObjects sync.Map //map[Coordinates]map[string]Object
	// муравьи
	dynamicObjects []Object
}

func newData(w, h int) *data {
	d := &data{w: w, h: h, dynamicObjects: make([]Object, 0)}
	return d
}

func (d *data) addObject(o Object) {
	if o.GetType() == ANT {
		d.dynamicObjects = append(d.dynamicObjects, o)
		return
	}
	pos := o.GetPosition()
	obs, _ := d.staticObjects.LoadOrStore(pos, &sync.Map{})
	obs.(*sync.Map).Store(o.GetID(), o)
}

func (d *data) objectsByPosition(position Coordinates) []Object /*map[string]Object*/ {
	x, y := position.X, position.Y
	if x > d.w || x < 0 || y > d.h || y < 0 {
		return nil
	}
	m, ok := d.staticObjects.Load(position)
	if !ok || m == nil {
		return nil
	}
	obs := make([]Object, 0, 3)
	m.(*sync.Map).Range(func(key, value interface{}) bool {
		obs = append(obs, value.(Object))
		return true
	})
	return obs

}

func (d *data) deleteObject(position Coordinates, id string) {
	obs, ok := d.staticObjects.Load(position)
	if ok {
		obs.(*sync.Map).Delete(id)
	}
}
