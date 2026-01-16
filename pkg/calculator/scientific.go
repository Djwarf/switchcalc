package calculator

import (
	"math"
)

func (e *Engine) toRadians(angle float64) float64 {
	switch e.AngleMode {
	case Degrees:
		return angle * math.Pi / 180
	case Gradians:
		return angle * math.Pi / 200
	default:
		return angle
	}
}

func (e *Engine) fromRadians(rad float64) float64 {
	switch e.AngleMode {
	case Degrees:
		return rad * 180 / math.Pi
	case Gradians:
		return rad * 200 / math.Pi
	default:
		return rad
	}
}

func (e *Engine) Sin() {
	e.CurrentValue = math.Sin(e.toRadians(e.CurrentValue))
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Cos() {
	e.CurrentValue = math.Cos(e.toRadians(e.CurrentValue))
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Tan() {
	e.CurrentValue = math.Tan(e.toRadians(e.CurrentValue))
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Asin() {
	if e.CurrentValue < -1 || e.CurrentValue > 1 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = e.fromRadians(math.Asin(e.CurrentValue))
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Acos() {
	if e.CurrentValue < -1 || e.CurrentValue > 1 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = e.fromRadians(math.Acos(e.CurrentValue))
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Atan() {
	e.CurrentValue = e.fromRadians(math.Atan(e.CurrentValue))
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Sinh() {
	e.CurrentValue = math.Sinh(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Cosh() {
	e.CurrentValue = math.Cosh(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Tanh() {
	e.CurrentValue = math.Tanh(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Asinh() {
	e.CurrentValue = math.Asinh(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Acosh() {
	if e.CurrentValue < 1 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = math.Acosh(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Atanh() {
	if e.CurrentValue <= -1 || e.CurrentValue >= 1 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = math.Atanh(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Log() {
	if e.CurrentValue <= 0 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = math.Log10(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Ln() {
	if e.CurrentValue <= 0 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = math.Log(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Log2() {
	if e.CurrentValue <= 0 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = math.Log2(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Exp() {
	e.CurrentValue = math.Exp(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Exp10() {
	e.CurrentValue = math.Pow(10, e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Exp2() {
	e.CurrentValue = math.Pow(2, e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Sqrt() {
	if e.CurrentValue < 0 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = math.Sqrt(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Cbrt() {
	e.CurrentValue = math.Cbrt(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Square() {
	e.CurrentValue = e.CurrentValue * e.CurrentValue
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Cube() {
	e.CurrentValue = e.CurrentValue * e.CurrentValue * e.CurrentValue
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Reciprocal() {
	if e.CurrentValue == 0 {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	e.CurrentValue = 1 / e.CurrentValue
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Factorial() {
	n := int(e.CurrentValue)
	if n < 0 || e.CurrentValue != float64(n) {
		e.Display = "Error"
		e.NewInput = true
		return
	}
	if n > 170 {
		e.Display = "Overflow"
		e.NewInput = true
		return
	}
	result := 1.0
	for i := 2; i <= n; i++ {
		result *= float64(i)
	}
	e.CurrentValue = result
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Abs() {
	e.CurrentValue = math.Abs(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Floor() {
	e.CurrentValue = math.Floor(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Ceil() {
	e.CurrentValue = math.Ceil(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Round() {
	e.CurrentValue = math.Round(e.CurrentValue)
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) Pi() {
	e.CurrentValue = math.Pi
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) E() {
	e.CurrentValue = math.E
	e.Display = e.formatNumber(e.CurrentValue)
	e.NewInput = true
}

func (e *Engine) SetAngleMode(mode AngleMode) {
	e.AngleMode = mode
}
