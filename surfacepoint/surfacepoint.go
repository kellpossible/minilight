package surfacepoint

import (
	"math/rand"
	"math32"
	"minilight/triangle"
	. "vector"
)

type SurfacePoint struct {
	Tri      *triangle.Triangle
	Position Vector3f
}

func New(tri *triangle.Triangle,
	position Vector3f) *SurfacePoint {
	return &SurfacePoint{tri, position}
}

func (self *SurfacePoint) Emission(toPosition Vector3f,
	outDirection Vector3f,
	isSolidAngle bool) Vector3f {
	ray := toPosition.Sub(self.Position)
	normal := self.Tri.Normal()
	distance2 := ray.Dot(normal)
	cosOut := outDirection.Dot(normal)
	area := self.Tri.Area()

	/* emit front face of surface only */
	var mulval1, mulval2 float32

	mulval1 = math32.Bool2Float(cosOut > 0.0)

	if isSolidAngle {
		/* with infinity clamped-out */
		var divval, powval float32
		powval = math32.Pow10(-6)
		if distance2 >= powval {
			divval = distance2
		} else {
			divval = powval
		}
		mulval2 = (cosOut * area) / divval
	} else {
		mulval2 = 1.0
	}

	solidAngle := mulval1 * mulval2
	return self.Tri.Emitivity.Mulf(solidAngle)
}

func (self *SurfacePoint) Reflection(inDirection,
	inRadiance,
	outDirection Vector3f) Vector3f {
	normal := self.Tri.Normal()
	inDot := inDirection.Dot(normal)
	outDot := outDirection.Dot(normal)
	
	/* directions must be same side of surface (no transmission) */
	isSameSide := !(math32.Xor(inDot < 0.0, outDot < 0.0))
	
	/* ideal diffuse BRDF:
		radiance sacaled by reflectivity, cosine and 1/pi */
	r := inRadiance.Mulv(self.Tri.Reflectivity)
	mulval1 := (math32.Abs(inDot)/math32.Pi)
	mulval2 := math32.Bool2Float(isSameSide)
	return r.Mulf(mulval1 * mulval2)
}

func (self *SurfacePoint) NextDirection(r *rand.Rand,
	inDirection Vector3f, 
	outDirection_o *Vector3f,
	colour *Vector3f) bool {
	reflectivityMean := self.Tri.Reflectivity.Dot(VECTOR_ONE) / 3.0
	
	/* russian-roulette for reflectance 'magnitude' */
	isAlive := r.Float32() < reflectivityMean
	
	if !isAlive {
		return false
	}
	
	/* cosine-weighted importance sample hemisphere */
	_2pr1 := math32.Pi * 2.0 * r.Float32()
	sr2 := math32.Sqrt(r.Float32())
	
	/* make coord frame coefficients (z in normal direction) */
	x := math32.Cos( _2pr1 ) * sr2
	y := math32.Sin( _2pr1 ) * sr2
	z := math32.Sqrt(1.0 - (sr2 * sr2))
	
	/* make coord frame */
	t := self.Tri.Tangent()
	n := self.Tri.Normal()
	var c Vector3f
	/* put normal on inward ray side of surface (preventing transmission) */
	if n.Dot(inDirection) < 0.0 {
		n.Neg()
	}
	c = n.Cross(t)
	
	/* scale frame by coefficients */
	tx := t.Mulf(x)
	cy := c.Mulf(y)
	nz := n.Mulf(z)
	
	/* make direction from sum of scaled components */
	sum := tx.Add(cy)
	*outDirection_o = sum.Add(nz)
	
	/* make color by dividing-out mean from reflectivity */
	*colour = self.Tri.Reflectivity.Mulf(1.0/reflectivityMean)
	
	/* discluding degenerate result direction */
	return isAlive && !outDirection_o.Is_Zero()
}
