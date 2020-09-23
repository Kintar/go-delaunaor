package delaunator

import "sort"

type ByDistance struct {
    coords []float64
    center [2]float64
}

func (b ByDistance) Len() int {
    return len(b.coords) / 2
}

func (b ByDistance) Less(i, j int) bool {
    ax := b.coords[i * 2]
    ay := b.coords[i * 2 + 1]
    bx := b.coords[j * 2]
    by := b.coords[j * 2 + 1]
    da := dist(b.center[0], b.center[1], ax, ay)
    db := dist(b.center[0], b.center[1], bx, by)
    return da < db
}

func (b ByDistance) Swap(i, j int) {
    b.coords[i * 2], b.coords[i * 2 + 1], b.coords[j * 2], b.coords[j * 2 + 1] =
        b.coords[j * 2], b.coords[j * 2 + 1], b.coords[i * 2], b.coords[i * 2 + 1]
}

type IndexSort struct {
    indexes []int
    values []float64
}

func (is IndexSort) Len() int {
    return len(is.indexes)
}

func (is IndexSort) Less(i, j int) bool {
    return is.values[i] < is.values[j]
}

func (is IndexSort) Swap(i, j int) {
    is.values[i], is.values[j], is.indexes[i], is.indexes[j] =
        is.values[j], is.values[i], is.indexes[j], is.indexes[i]
}

func sortByValues(indexes []int, values []float64) {
    sort.Sort(IndexSort{
        indexes: indexes,
        values:  values,
    })
}
