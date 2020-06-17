package simulation

type data struct {
	w, h    int
	objects [][][]Object
}

func newData(w, h int) *data {
	obs := make([][][]Object, w)
	for i := range obs {
		obs[i] = make([][]Object, h)
		for j := range obs[i] {
			obs[i][j] = make([]Object, 0)
		}
	}

	d := &data{w, h, obs}

	return d
}

func (d *data) addObject(o Object) {
	x, y := o.GetPosition().X, o.GetPosition().Y
	d.objects[x][y] = append(d.objects[x][y], o)
}

func (d *data) objectsByPosition(position Coordinates) ([]Object, int) {
	x, y := position.X, position.Y
	if x > d.w || x < 0 || y > d.h || y < 0 {
		return nil, 0
	}
	return d.objects[position.X][position.Y], len(d.objects[position.X][position.Y])
}

func (d *data) deleteObject(position Coordinates, id string) {
	x, y := position.X, position.Y
	a := &d.objects[x][y]

	n := -1
	for i := range *a {
		if (*a)[i].GetID() == id {
			n = i
			break
		}
	}
	if n == -1 {
		return
	}

	copy((*a)[n:], (*a)[n+1:])
	(*a)[len((*a))-1] = nil
	(*a) = (*a)[:len((*a))-1]
}
