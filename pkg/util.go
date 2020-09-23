package delaunator

import "math"

func dist(ax, ay, bx, by float64) float64 {
    dx := ax - bx
    dy := ay - by

    return dx * dx + dy * dy
}

func circumradius(ax, ay, bx, by, cx, cy float64) float64 {
    dx := bx - ax
    dy := by - ay
    ex := cx - ax
    ey := cy - ay

    bl := dx * dx + dy * dy
    cl := ex * ex + ey * ey
    d := 0.5 / (dx * ey - dy * ex)

    x := (ey * bl - dy * cl) * d
    y := (dx * cl - ex * bl) * d

    return x * x + y * y
}

// return 2d orientation sign if we're confident in it through J. Shewchuk's error bound check
func orientIfSure(px, py, rx, ry, qx, qy float64) float64 {
    l := (ry - py) * (qx - px)
    r := (rx - px) * (qy - py)

    if math.Abs(l - r) > 3.3306690738754716e-16 * math.Abs(l + r) {
        return l - r
    } else {
        return 0
    }
}

// a more robust orientation test that's stable in a given triangle (to fix robustness issues)
func orient(rx, ry, qx, qy, px, py float64) bool {
    o := orientIfSure(px, py, rx, ry, qx, qy)
    if o < 0 {
        return true
    }
orientIfSure(rx, ry, qx, qy, px, py) ||
orientIfSure(qx, qy, px, py, rx, ry)) < 0;
}
