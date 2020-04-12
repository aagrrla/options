package options

import "time"
import "math"

//Option has basic structure for any options //black scholes requires
type Option struct {
	Underlying        float64
	ExpiryTime        time.Time
	StrikePrice       float64
	Type              string
	ImpliedVolatility float64
	RiskFree          float64
	Dividend          float64
}

//OptionValues has price, delta theta, gamma greeks for specific time with given option
type OptionValues struct {
	Time  time.Time
	Value float64
	Delta float64
	Gamma float64
	Vega  float64
	Theta float64
	Rho   float64
}

//cdf is from stats pkg where normal cdf μ 0 σ 1
func cdf(x float64) float64 {
	return math.Erfc(-(x-0)/(1*math.Sqrt2)) / 2
}
func pdf(x float64) float64 {
	return 0.3989422804014327 * math.Exp(-1*(x-0)*(x-0)/(2*1*1)) //0.39*exp(-1*(x-μ)*(x-μ)/(2*σ*σ))/σ
}

//Calc computes and returns calculated Values
func (op *Option) Calc(atTime time.Time) OptionValues {
	s, k := op.Underlying, op.StrikePrice
	t := float64((op.ExpiryTime).Sub(atTime).Hours()) / (24 * 365)
	r, σ, q := op.RiskFree/100, op.ImpliedVolatility/100, op.Dividend/100
	d1 := (math.Log(s/k) + (r-q+σ*σ/2.0)*t) / (σ * math.Sqrt(t))
	d2 := d1 - σ*math.Sqrt(t)
	kert := k * math.Exp(-1*r*t)

	delta := cdf(d1)
	gamma := pdf(d1) / (s * σ * math.Sqrt(t))
	vega := s * pdf(d1) * math.Sqrt(t)

	var value, theta, rho float64

	switch op.Type {
	case "Call":
		value = s*cdf(d1) - kert*cdf(d2)
		theta = -s*pdf(d1)*σ/(2.0*math.Sqrt(t)) - r*kert*cdf(d2)
		rho = t * kert * cdf(d2)
	case "Put":
		delta = delta - 1
		value = kert*cdf(-d2) - s*cdf(-d1)
		theta = -s*pdf(d1)*σ/(2.0*math.Sqrt(t)) + r*kert*cdf(-d2)
		rho = -t * kert * cdf(-d2)
	}
	return OptionValues{atTime, value, delta, gamma, vega / 100, theta / 365, rho / 100}
}
