package calculator

import (
	"fmt"
)

type BitWidth int

const (
	Bits8  BitWidth = 8
	Bits16 BitWidth = 16
	Bits32 BitWidth = 32
	Bits64 BitWidth = 64
)

func (e *Engine) And(other int64) {
	val := int64(e.CurrentValue)
	result := val & other
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) Or(other int64) {
	val := int64(e.CurrentValue)
	result := val | other
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) Xor(other int64) {
	val := int64(e.CurrentValue)
	result := val ^ other
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) Not() {
	val := int64(e.CurrentValue)
	result := ^val
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) Nand(other int64) {
	val := int64(e.CurrentValue)
	result := ^(val & other)
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) Nor(other int64) {
	val := int64(e.CurrentValue)
	result := ^(val | other)
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) LeftShift(bits uint) {
	val := int64(e.CurrentValue)
	result := val << bits
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) RightShift(bits uint) {
	val := int64(e.CurrentValue)
	result := val >> bits
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) RotateLeft(bits uint, width BitWidth) {
	val := uint64(e.CurrentValue)
	mask := uint64((1 << width) - 1)
	val &= mask
	bits = bits % uint(width)
	result := ((val << bits) | (val >> (uint(width) - bits))) & mask
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(int64(result))
	e.NewInput = true
}

func (e *Engine) RotateRight(bits uint, width BitWidth) {
	val := uint64(e.CurrentValue)
	mask := uint64((1 << width) - 1)
	val &= mask
	bits = bits % uint(width)
	result := ((val >> bits) | (val << (uint(width) - bits))) & mask
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(int64(result))
	e.NewInput = true
}

func (e *Engine) GetBit(position uint) int {
	val := int64(e.CurrentValue)
	return int((val >> position) & 1)
}

func (e *Engine) SetBit(position uint) {
	val := int64(e.CurrentValue)
	result := val | (1 << position)
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) ClearBit(position uint) {
	val := int64(e.CurrentValue)
	result := val &^ (1 << position)
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) ToggleBit(position uint) {
	val := int64(e.CurrentValue)
	result := val ^ (1 << position)
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) CountBits() {
	val := uint64(e.CurrentValue)
	count := 0
	for val != 0 {
		count += int(val & 1)
		val >>= 1
	}
	e.CurrentValue = float64(count)
	e.Display = e.FormatInBase(int64(count))
	e.NewInput = true
}

func (e *Engine) LeadingZeros(width BitWidth) {
	val := uint64(e.CurrentValue)
	mask := uint64((1 << width) - 1)
	val &= mask
	count := 0
	for i := int(width) - 1; i >= 0; i-- {
		if (val>>uint(i))&1 == 0 {
			count++
		} else {
			break
		}
	}
	e.CurrentValue = float64(count)
	e.Display = e.FormatInBase(int64(count))
	e.NewInput = true
}

func (e *Engine) TrailingZeros() {
	val := uint64(e.CurrentValue)
	if val == 0 {
		e.CurrentValue = 64
		e.Display = e.FormatInBase(64)
		e.NewInput = true
		return
	}
	count := 0
	for (val & 1) == 0 {
		count++
		val >>= 1
	}
	e.CurrentValue = float64(count)
	e.Display = e.FormatInBase(int64(count))
	e.NewInput = true
}

func (e *Engine) ByteSwap(width BitWidth) {
	val := uint64(e.CurrentValue)
	var result uint64

	switch width {
	case Bits16:
		result = ((val & 0xFF) << 8) | ((val & 0xFF00) >> 8)
	case Bits32:
		result = ((val & 0xFF) << 24) |
			((val & 0xFF00) << 8) |
			((val & 0xFF0000) >> 8) |
			((val & 0xFF000000) >> 24)
	case Bits64:
		result = ((val & 0xFF) << 56) |
			((val & 0xFF00) << 40) |
			((val & 0xFF0000) << 24) |
			((val & 0xFF000000) << 8) |
			((val & 0xFF00000000) >> 8) |
			((val & 0xFF0000000000) >> 24) |
			((val & 0xFF000000000000) >> 40) |
			((val & 0xFF00000000000000) >> 56)
	default:
		result = val
	}
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(int64(result))
	e.NewInput = true
}

func (e *Engine) TwosComplement(width BitWidth) {
	val := int64(e.CurrentValue)
	mask := int64((1 << width) - 1)
	result := (^val + 1) & mask
	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.NewInput = true
}

func (e *Engine) GetBinaryString(width BitWidth) string {
	val := uint64(e.CurrentValue)
	result := ""
	for i := int(width) - 1; i >= 0; i-- {
		if (val>>uint(i))&1 == 1 {
			result += "1"
		} else {
			result += "0"
		}
		if i > 0 && i%4 == 0 {
			result += " "
		}
	}
	return result
}

func (e *Engine) GetAllBases() map[string]string {
	val := int64(e.CurrentValue)
	return map[string]string{
		"DEC": fmt.Sprintf("%d", val),
		"HEX": fmt.Sprintf("%X", val),
		"OCT": fmt.Sprintf("%o", val),
		"BIN": e.GetBinaryString(Bits64),
	}
}

type BitwiseOperation int

const (
	BitOpAnd BitwiseOperation = iota
	BitOpOr
	BitOpXor
	BitOpNand
	BitOpNor
	BitOpLeftShift
	BitOpRightShift
)

func (e *Engine) SetBitwiseOperation(op BitwiseOperation) {
	if e.PendingOp != OpNone && !e.NewInput {
		e.Calculate()
	}
	e.StoredValue = e.CurrentValue
	e.PendingOp = Operation(100 + int(op))
	e.NewInput = true
}

func (e *Engine) CalculateBitwise() int64 {
	op := BitwiseOperation(int(e.PendingOp) - 100)
	stored := int64(e.StoredValue)
	current := int64(e.CurrentValue)
	var result int64

	switch op {
	case BitOpAnd:
		result = stored & current
	case BitOpOr:
		result = stored | current
	case BitOpXor:
		result = stored ^ current
	case BitOpNand:
		result = ^(stored & current)
	case BitOpNor:
		result = ^(stored | current)
	case BitOpLeftShift:
		result = stored << uint(current)
	case BitOpRightShift:
		result = stored >> uint(current)
	default:
		return current
	}

	e.CurrentValue = float64(result)
	e.Display = e.FormatInBase(result)
	e.PendingOp = OpNone
	e.NewInput = true
	return result
}
