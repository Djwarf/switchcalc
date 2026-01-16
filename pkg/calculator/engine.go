package calculator

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type Operation int

const (
	OpNone Operation = iota
	OpAdd
	OpSubtract
	OpMultiply
	OpDivide
	OpModulo
	OpPower
)

type Engine struct {
	Display        string
	CurrentValue   float64
	StoredValue    float64
	PendingOp      Operation
	NewInput       bool
	Memory         float64
	History        []string
	AngleMode      AngleMode
	NumberBase     NumberBase
}

type AngleMode int

const (
	Degrees AngleMode = iota
	Radians
	Gradians
)

type NumberBase int

const (
	Decimal NumberBase = iota
	Binary
	Octal
	Hexadecimal
)

func NewEngine() *Engine {
	return &Engine{
		Display:    "0",
		NewInput:   true,
		AngleMode:  Degrees,
		NumberBase: Decimal,
		History:    make([]string, 0),
	}
}

func (e *Engine) Clear() {
	e.Display = "0"
	e.CurrentValue = 0
	e.StoredValue = 0
	e.PendingOp = OpNone
	e.NewInput = true
}

func (e *Engine) ClearEntry() {
	e.Display = "0"
	e.CurrentValue = 0
	e.NewInput = true
}

func (e *Engine) InputDigit(digit string) {
	if e.NewInput {
		e.Display = digit
		e.NewInput = false
	} else {
		if e.Display == "0" && digit != "." {
			e.Display = digit
		} else {
			e.Display += digit
		}
	}
	e.CurrentValue, _ = strconv.ParseFloat(e.Display, 64)
}

func (e *Engine) InputDecimal() {
	if e.NewInput {
		e.Display = "0."
		e.NewInput = false
	} else if !strings.Contains(e.Display, ".") {
		e.Display += "."
	}
}

func (e *Engine) InputExponent() {
	if e.NewInput {
		e.Display = "1e"
		e.NewInput = false
	} else if !strings.Contains(strings.ToLower(e.Display), "e") {
		e.Display += "e"
	}
}

func (e *Engine) InputHexDigit(digit string) {
	if e.NumberBase != Hexadecimal {
		return
	}
	if e.NewInput {
		e.Display = digit
		e.NewInput = false
	} else {
		if e.Display == "0" {
			e.Display = digit
		} else {
			e.Display += digit
		}
	}
	val, _ := strconv.ParseInt(e.Display, 16, 64)
	e.CurrentValue = float64(val)
}

func (e *Engine) SetOperation(op Operation) {
	if e.PendingOp != OpNone && !e.NewInput {
		e.Calculate()
	}
	e.StoredValue = e.CurrentValue
	e.PendingOp = op
	e.NewInput = true
}

func (e *Engine) Calculate() float64 {
	if e.PendingOp == OpNone {
		return e.CurrentValue
	}

	var result float64
	switch e.PendingOp {
	case OpAdd:
		result = e.StoredValue + e.CurrentValue
	case OpSubtract:
		result = e.StoredValue - e.CurrentValue
	case OpMultiply:
		result = e.StoredValue * e.CurrentValue
	case OpDivide:
		if e.CurrentValue != 0 {
			result = e.StoredValue / e.CurrentValue
		} else {
			e.Display = "Error"
			e.NewInput = true
			return 0
		}
	case OpModulo:
		if e.CurrentValue != 0 {
			result = math.Mod(e.StoredValue, e.CurrentValue)
		} else {
			e.Display = "Error"
			e.NewInput = true
			return 0
		}
	case OpPower:
		result = math.Pow(e.StoredValue, e.CurrentValue)
	}

	historyEntry := fmt.Sprintf("%s %s %s = %s",
		e.formatNumber(e.StoredValue),
		e.opSymbol(e.PendingOp),
		e.formatNumber(e.CurrentValue),
		e.formatNumber(result))
	e.History = append(e.History, historyEntry)
	if len(e.History) > 50 {
		e.History = e.History[1:]
	}

	e.CurrentValue = result
	e.Display = e.formatNumber(result)
	e.PendingOp = OpNone
	e.NewInput = true
	return result
}

func (e *Engine) opSymbol(op Operation) string {
	switch op {
	case OpAdd:
		return "+"
	case OpSubtract:
		return "−"
	case OpMultiply:
		return "×"
	case OpDivide:
		return "÷"
	case OpModulo:
		return "mod"
	case OpPower:
		return "^"
	}
	return ""
}

func (e *Engine) formatNumber(n float64) string {
	if e.NumberBase != Decimal {
		return e.FormatInBase(int64(n))
	}
	if n == float64(int64(n)) {
		return fmt.Sprintf("%d", int64(n))
	}
	formatted := fmt.Sprintf("%.10f", n)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")
	return formatted
}

func (e *Engine) Negate() {
	e.CurrentValue = -e.CurrentValue
	e.Display = e.formatNumber(e.CurrentValue)
}

func (e *Engine) Percent() {
	if e.PendingOp == OpAdd || e.PendingOp == OpSubtract {
		e.CurrentValue = e.StoredValue * (e.CurrentValue / 100)
	} else {
		e.CurrentValue = e.CurrentValue / 100
	}
	e.Display = e.formatNumber(e.CurrentValue)
}

func (e *Engine) MemoryClear() {
	e.Memory = 0
}

func (e *Engine) MemoryRecall() {
	e.CurrentValue = e.Memory
	e.Display = e.formatNumber(e.Memory)
	e.NewInput = true
}

func (e *Engine) MemoryAdd() {
	e.Memory += e.CurrentValue
}

func (e *Engine) MemorySubtract() {
	e.Memory -= e.CurrentValue
}

func (e *Engine) MemoryStore() {
	e.Memory = e.CurrentValue
}

func (e *Engine) Backspace() {
	if len(e.Display) > 1 {
		e.Display = e.Display[:len(e.Display)-1]
		e.CurrentValue, _ = strconv.ParseFloat(e.Display, 64)
	} else {
		e.Display = "0"
		e.CurrentValue = 0
	}
}

func (e *Engine) SetNumberBase(base NumberBase) {
	intVal := int64(e.CurrentValue)
	e.NumberBase = base
	e.Display = e.FormatInBase(intVal)
}

func (e *Engine) FormatInBase(n int64) string {
	switch e.NumberBase {
	case Binary:
		return strconv.FormatInt(n, 2)
	case Octal:
		return strconv.FormatInt(n, 8)
	case Hexadecimal:
		return strings.ToUpper(strconv.FormatInt(n, 16))
	default:
		return strconv.FormatInt(n, 10)
	}
}

func (e *Engine) ParseCurrentBase(s string) (int64, error) {
	switch e.NumberBase {
	case Binary:
		return strconv.ParseInt(s, 2, 64)
	case Octal:
		return strconv.ParseInt(s, 8, 64)
	case Hexadecimal:
		return strconv.ParseInt(s, 16, 64)
	default:
		return strconv.ParseInt(s, 10, 64)
	}
}
