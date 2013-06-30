package spatialindex

import (
	"math32"
	misect "minilight/intersection"
	. "minilight/ray"
	mtri "minilight/triangle"
	. "vector"
)

/**
 * A minimal spatial index for ray tracing.<br/><br/>
 *
 * Suitable for a scale of 1 metre == 1 numerical unit, and with a resolution
 * of 1 millimetre. (Implementation uses fixed tolerances, and single precision
 * FP.)
 *
 * Constant.<br/><br/>
 *
 * @implementation
 * A crude State pattern: typed by isBranch field to be either a branch
 * or leaf cell.<br/><br/>
 *
 * Octree: axis-aligned, cubical. Subcells are numbered thusly:
 * <pre>      110---111
 *            /|    /|
 *         010---011 |
 *    y z   | 100-|-101
 *    |/    |/    | /
 *    .-x  000---001      </pre><br/><br/>
 *
 * Each cell stores its bound (fatter data, but simpler code).<br/><br/>
 *
 * Calculations for building and tracing are absolute rather than incremental --
 * so quite numerically solid. Uses tolerances in: bounding triangles (in
 * TriangleBound), and checking intersection is inside cell (both effective
 * for axis-aligned items). Also, depth is constrained to an absolute subcell
 * size (easy way to handle overlapping items).
 *
 * @invariants
 * * aBound[0-2] <= aBound[3-5]
 * * bound encompasses the cell's contents
 * if isBranch
 * * apArray elements are SpatialIndex pointers or zeros
 * * length (of apArray) is 8
 * else
 * * apArray elements are non-zero Triangle pointers
 */

/* constants */
/* accommodates scene including sun and earth, down to cm cells
   (use 47 for mm) */
const MAX_LEVELS = 44

/* 8 seemed reasonably optimal in casual testing */
const MAX_ITEMS = 8

/* implementation */
type SpatialIndex struct {
	isBranch bool
	aBound   []float32
	apArray  [][]mtri.Triangle
	length   int
}

func New(eyePosition Vector3f,
	items []mtri.Triangle) *SpatialIndex {
	self := &SpatialIndex{}
	self.aBound = make([]float32, 6)

	/* set overall bound (and convert to collection of pointers) */
	itemPs := make([]*mtri.Triangle, len(items))
	{
		var i, j int
		/* accommodate eye position (makes tracing algorithm simpler) */
		for i = 6; i > 0; i-- {
			self.aBound[i] = eyePosition.XYZ(i % 3)
		}

		/* accommodate all items */
		for i = 0; i < len(items); i++ {
			itemPs[i] = &items[i]
			itemBound := items[i].Bound()

			/* accommodate item */
			for j = 0; j < 6; j++ {
				if math32.Xor(self.aBound[j] > itemBound[j], j > 2) {
					self.aBound[j] = itemBound[j]
				}
			}

		}
		// make cubical
		var maxSize float32 = 0.0
		v1 := SliceToVector3f(self.aBound[3:6])
		v2 := SliceToVector3f(self.aBound[0:3])
		v3 := v1.Sub(v2)
		maxSize = math32.MaxSlice(v3.ToSlice())
		v4 := FloatToVector3f(maxSize)
		v5 := v2.Add(v4)
		v6 := v1.Clamped(v5, VECTOR_MAX)
		v6slice := v6.ToSlice()
		self.aBound = append(self.aBound[0:3], v6slice...)
	}
	//TODO: octree construct method
	return self
}



func (self *SpatialIndex) Intersection(ray *Ray) *misect.Intersection {
	return &misect.Intersection{}
}
