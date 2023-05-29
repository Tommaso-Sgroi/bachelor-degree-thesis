// PerceptibleUpTo calcola quanto e' grande (in km) il raggio entro il 
// quale il terremoto e' percettibile sopra una
// determinata soglia Imin (intensita' minima (soglia) di solito 
// impostata a 1.5, in scala MCS), ed il tempo di arrivo del 
// fronte d'onda al perimetro del cerchio che viene
// generato dal raggio calcolato applicato sull'epicentro.
func PerceptibleUpTo(magnitudo float64, depth float64, minIntensity 
float64) (float64, time.Duration) {
	a := 1.51
	b := 1.55
	c := -3.15

	if depth < 5 {
		depth = 5
	}

	Ic := a + b*magnitudo
	R0 := math.Pow(10, (minIntensity-Ic)/c)
	d0 := 0.0
	if R0*R0-depth*depth > 0 {
		d0 = math.Sqrt(R0*R0 - depth*depth)
	}
	ett := d0 / 5
	return d0, time.Duration(ett * 1000 * 1000 * 1000)
}