package delaunator

import (
    "fmt"
    "math"
)

type triangulation struct {
    coords    []float64
    triangles []uint
    halfEdges []int
    hull      []int
}

type point struct {
    x float64
    y float64
}

type Point2d interface {
    X() float64
    Y() float64
}

func (p point) X() float64 {
    return p.x
}

func (p point) Y() float64 {
    return p.y
}

func NewTriangulationFromPoints(points []Point2d) (triangulation, error) {
    if len(points) < 3 {
        return triangulation{}, fmt.Errorf("cannot triangulate fewer than 3 points")
    }

    coords := make([]float64, len(points)*2)
    for i, p := range points {
        coords[i*2] = p.X()
        coords[i*2+1] = p.Y()
    }

    return NewTriangulationFromCoords(coords)
}

func NewTriangulationFromCoords(coords []float64) (triangulation, error) {
    if len(coords) < 6 {
        return triangulation{}, fmt.Errorf("cannot trianguate fewer than 3 points")
    }

    coordsCopy := make([]float64, len(coords))
    copy(coordsCopy, coords)

    // TODO: This is from mapbox/delaunator, but why are we dividing the length by 2, then multiplying it by 2?
    maxTriangles := 2*(len(coordsCopy)>>1) - 5

    tri := triangulation{
        coords:    coordsCopy,
        triangles: make([]uint, maxTriangles*3),
        halfEdges: make([]int, maxTriangles*3),
    }

    err := tri.Update()

    return tri, err
}

// Update recalculates this triangulation
func (t triangulation) Update() error {
    coords := t.coords

    minX := math.Inf(1)
    minY := minX
    maxX := math.Inf(-1)
    maxY := maxX

    n := len(coords) >> 1
    ids := make([]int, n)

    for i := 0; i < n; i++ {
        x := coords[2*i]
        y := coords[2*i+1]
        if x < minX {
            minX = x
        }
        if y < minY {
            minY = y
        }
        if x > maxX {
            maxX = x
        }
        if y > maxY {
            maxY = y
        }
        ids[i] = i
    }

    cx := (minX + maxX) / 2.0
    cy := (minY + maxY) / 2.0

    minDist := math.Inf(1)

    // pick a seed point close to the center
    var i0, i1, i2 int
    for i := 0; i < n; i++ {
        d := dist(cx, cy, coords[2*i], coords[2*i+1])
        if d < minDist {
            i0 = i
            minDist = d
        }
    }
    i0x, i0y := coords[2*i0], coords[2*i0+1]

    minDist = math.Inf(1)

    // find the point closest to the seed
    for i := 0; i < n; i++ {
        if i == i0 {
            continue
        }
        d := dist(i0x, i0y, coords[2*i], coords[2*i+1])
        if d < minDist {
            i1 = i
            minDist = d
        }
    }
    i1x, i1y := coords[2*i1], coords[2*i1+1]

    minRadius := math.Inf(1)

    // find the point which makes the smallest circumcircle with the first two
    for i := 0; i < n; i++ {
        if i == i0 || i == i1 {
            continue
        }
        r := circumradius(i0x, i0y, i1x, i1y, coords[2*i], coords[2*i+1])
        if r < minRadius {
            i2 = i
            minRadius = r
        }
    }
    i2x, i2y := coords[i2*2], coords[i2*2+1]

    dists := make([]float64, n)

    if math.IsInf(minRadius, 0) {
        // Somehow, all of our points are collinear.  Sort them by dx or dy and return them as a hull
        for i := 0; i < n; i++ {
            dists[i] = math.Max(coords[i*2]-coords[0], coords[i*2+1]-coords[1])
        }

        sortByValues(ids, dists)

        hull := make([]int, n)
        j := 0
        d0 := math.Inf(-1)
        for i := 0; i < n; i++ {
            if dists[i] > d0 {
                hull[j] = i
                j++
                d0 = dists[i]
            }
        }

        t.triangles = make([]uint, 0)
        t.hull = hull[0:j]
        t.halfEdges = make([]int, 0)

        // We're done
        return nil
    }


}

