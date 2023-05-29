// La funzione squaring prende un punto come centro di un 
// cerchio con raggio = `distanza`.
// Restituisce i limiti del quadrato che inscrive tale cerchio.
func Squaring(point Point2D, distance float64) BoundingBox {
	offset := (distance * 1000 / metersPerUnit) * latLngChange

	return BoundingBox{
		A: Point2D{
			Latitude:  point.Latitude - Latitude(offset),
			Longitude: point.Longitude - Longitude(offset),
		},
		B: Point2D{
			Latitude:  point.Latitude + Latitude(offset),
			Longitude: point.Longitude + Longitude(offset),
		},
	}
}
