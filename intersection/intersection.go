package intersection

import . "vector"
import . "minilight/ray"

/*TODO: interface intersector, for things like the triangle, 
or things like acceleration structures,
all implementing the intersection interface*/

type Intersection struct {
	IsHit   bool
	Distance float32
	Location Vector3f
}

func NewLocationCalc(is_hit bool,
	distance float32,
	ray *Ray) *Intersection {
	/* Creates new Intersection, and also calculates the location
	from the supplied ray*/
	ray_vec := ray.Direction.Mulf(distance)
	location := ray_vec.Add(ray.Origin)
	return &Intersection{is_hit, distance, location}
}

func NewNoHit() *Intersection {
	return &Intersection{IsHit: false}
}

func New(is_hit bool,
	distance float32,
	location Vector3f) *Intersection {
	return &Intersection{is_hit, distance, location}
}
