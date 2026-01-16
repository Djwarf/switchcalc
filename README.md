# SwitchCalc

A complete GTK4 calculator application for Linux with Standard, Scientific, Programmer, and Date calculation modes.

## Features

### Standard Mode
- Basic arithmetic operations (add, subtract, multiply, divide)
- Percentage calculations
- Memory functions (MC, MR, M+, M-, MS)
- Keyboard support for all operations

### Scientific Mode
- Trigonometric functions (sin, cos, tan, asin, acos, atan)
- Hyperbolic functions (sinh, cosh, tanh, asinh, acosh, atanh)
- Logarithmic functions (log, ln, log2)
- Exponential functions (exp, 10^x, 2^x)
- Powers and roots (square, cube, sqrt, cbrt, x^y)
- Constants (pi, e)
- Factorial, absolute value, floor, ceil, round
- Angle mode selector (Degrees, Radians, Gradians)
- Keyboard shortcuts: Ctrl+S (sin), Ctrl+C (cos), Ctrl+T (tan), Ctrl+L (log), Ctrl+N (ln), Ctrl+R (sqrt), Ctrl+P (pi), Ctrl+E (e)

### Programmer Mode
- Number base conversion (Decimal, Binary, Octal, Hexadecimal)
- Live display of value in all bases
- Bit width selection (8, 16, 32, 64-bit)
- Adjustable shift amount (1-63 bits)
- Bitwise operations:
  - AND, OR, XOR, NOT, NAND, NOR
  - Left shift, Right shift
  - Rotate left, Rotate right
  - Bit count (popcount)
  - Leading zeros, Trailing zeros
  - Byte swap
  - Two's complement
  - Toggle bit at position
- Keyboard shortcuts: & (AND), | (OR), ^ (XOR), ~ (NOT), < (left shift), > (right shift)
- Hex input via A-F keys

### Date Mode
- Date difference calculator (years, months, days, weeks, hours)
- Add/subtract time from dates (years, months, days, weeks)
- Age calculator with next birthday countdown
- Unix timestamp converter (to/from date)
- Today's info panel (week number, day of year, leap year status)
- Working days calculation (excluding weekends)

## Installation

### Arch Linux (AUR)

```bash
yay -S switchcalc
```

Or manually:

```bash
git clone https://aur.archlinux.org/switchcalc.git
cd switchcalc
makepkg -si
```

### Building from Source

Requirements:
- Go 1.21 or later
- GTK4
- GLib2

```bash
git clone https://github.com/yourusername/switchcalc.git
cd switchcalc
go build -o switchcalc ./cmd/switchcalc
```

## Usage

Launch from your application menu or run:

```bash
switchcalc
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| 0-9 | Enter digits |
| . | Decimal point |
| + | Add |
| - | Subtract |
| * | Multiply |
| / | Divide |
| % | Percent |
| Enter/= | Calculate |
| Escape | Clear all |
| Backspace | Delete last digit |
| Delete | Clear entry |

### Scientific Mode (Ctrl+key)
| Key | Action |
|-----|--------|
| Ctrl+S | Sine |
| Ctrl+C | Cosine |
| Ctrl+T | Tangent |
| Ctrl+L | Log (base 10) |
| Ctrl+N | Natural log |
| Ctrl+R | Square root |
| Ctrl+P | Pi |
| Ctrl+E | Euler's number |

### Programmer Mode
| Key | Action |
|-----|--------|
| A-F | Hex digits (in hex mode) |
| & | Bitwise AND |
| \| | Bitwise OR |
| ^ | Bitwise XOR |
| ~ | Bitwise NOT |
| < | Left shift |
| > | Right shift |

## License

MIT License

## Contributing

Contributions are welcome. Please open an issue or submit a pull request.
