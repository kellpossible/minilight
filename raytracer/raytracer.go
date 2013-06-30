package raytracer

import (
	mscene "minilight/scene"
	. "vector"
)

type RayTracer struct {
	scene *mscene.Scene
}

func New(scene *mscene.Scene) *RayTracer {
	return &RayTracer{scene}
}

func (self *RayTracer) Radiance(rayOrigin *Vector3f,
	rayDirection *Vector3f) {
	
}
