/************************************************************************************************100
A simple struct for tracking two points

This might move to another package
***************************************************************************************************/

package ili948x

type Point struct {
	X int16
	Y int16
}

/* take numbers and set them as the point value */
func (p *Point) Set(x, y int16) {
	p.X = x
	p.Y = y
}

/* Move this point to another points location */
func (p *Point) Move(other Point) {
	p.X = other.X
	p.Y = other.Y
}

/* Test if this point has collided with another point */
func (p Point) Collition(other Point) bool {
	return (p.X == other.X) && (p.Y == other.Y)
}
