package scene

import (
	misect "minilight/intersection"
	. "minilight/ray"
	sindex "minilight/spatialindex"
	"minilight/triangle"
	memit "minilight/emitter"
	. "vector"
	"math/rand"
)

type SceneParams struct {
	triangles                                  []triangle.Triangle
	skyEmission, groundReflection, eyePosition Vector3f
}

type Scene struct {
	/* objects */
	triangles []triangle.Triangle
	emitters  []memit.Emitter
	index     *sindex.SpatialIndex

	/* background */
	skyEmission, groundReflection Vector3f
}

func New(sp SceneParams) *Scene {
	self := &Scene{triangles: sp.triangles,
		skyEmission:      sp.skyEmission,
		groundReflection: sp.groundReflection}

	/* find emitting objects, until maximum reached */
	self.emitters = make([]memit.Emitter, 0, len(self.triangles))
	for _, v := range self.triangles {
		/* has non-zero emission and area */
		is_emitter := (!v.Emitivity.Is_Zero()) && (v.Area() > 0.0)
		if is_emitter {
			/* append to emitters storage */
			self.emitters = append(self.emitters, memit.Emitter{&v})
		}
	}

	/* condition background sky and ground values */
	self.skyEmission = self.skyEmission.Clamped(VECTOR_ZERO, self.skyEmission)
	self.groundReflection = self.groundReflection.Clamped(VECTOR_ZERO, VECTOR_ONE)

	/* make index of objects */
	self.index = sindex.New(sp.eyePosition, sp.triangles)
	return self
}

func (self *Scene) Intersection(ray *Ray) *misect.Intersection {
	return self.index.Intersection(ray)
}

func (self *Scene) Emitter(r *rand.Rand) *memit.Emitter {
	if len(self.emitters) == 0 {
		return nil
	}
	i := r.Int()%len(self.emitters)
	return &self.emitters[i]
}

func (self *Scene) DefaultEmission(backDirection Vector3f) Vector3f {
	/* sky for downward ray, ground for upward ray */
	if backDirection.Y < 0.0 {
		return self.skyEmission
	}
	return self.skyEmission.Mulv(self.groundReflection)
}
