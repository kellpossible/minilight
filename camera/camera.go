package camera

import (
	"math/rand"
	. "math32"
	mimage "minilight/image"
	mscene "minilight/scene"
	. "vector"
)

type Camera struct {
	/* eye definition */
	viewPosition Vector3f
	viewAngle    float32

	/* view frame */
	viewDirection, right, up Vector3f
}

func New(viewPosition Vector3f,
	viewDirection Vector3f,
	viewAngle float32) *Camera {
	Y := Vector3f{0.0, 1.0, 0.0}
	Z := Vector3f{0.0, 0.0, 1.0}

	self := &Camera{}

	self.viewPosition = viewPosition
	self.viewDirection = viewDirection.UnitizeCopy()
	/* if degenerate, default to Z */
	if self.viewDirection.Is_Zero() {
		self.viewDirection = Z
	}

	/* clamp and convert to radians */
	self.viewAngle = Min(Max(10.0, viewAngle), 160.0) * (Pi / 180.0)

	/* make other directions of view coord frame */
	/* make trial 'right', using viewDirection and assuming 'up' is Y */
	self.right = Y
	self.right = self.right.Cross(self.viewDirection).UnitizeCopy()

	/* check if 'right' is not valid
	   -- i.e. viewDirection was co-linear with 'up' */
	if self.right.Is_Zero() {
		var z float32
		/* 'up' is Z if viewDirection is down, otherwise -Z */
		if self.viewDirection.Z != 0.0 {
			z = 1.0
		} else {
			z = -1.0
		}
		self.up = Vector3f{0.0, 0.0, z}
		/* remake 'right' */
		self.right = self.up.Cross(self.viewDirection).UnitizeCopy()
	} else {
		/* use 'right', and make 'up' properly orthogonal */
		self.up = self.viewDirection.Cross(self.right).UnitizeCopy()
	}

	return self
}

func (self *Camera) GetFrame(scene *mscene.Scene,
	r *rand.Rand,
	image *mimage.Image) {
	
}
