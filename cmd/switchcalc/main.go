package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"

	"switchcalc/pkg/calculator"
)

type CalculatorMode int

const (
	ModeStandard CalculatorMode = iota
	ModeScientific
	ModeProgrammer
	ModeDateTime
)

type App struct {
	window        *gtk.ApplicationWindow
	engine        *calculator.Engine
	dateCalc      *calculator.DateTimeCalc
	mode          CalculatorMode
	display       *gtk.Label
	expressionLbl *gtk.Label
	historyList   *gtk.ListBox
	mainStack     *gtk.Stack
	modeButtons   map[CalculatorMode]*gtk.ToggleButton

	// Programmer mode widgets
	baseLabels    map[calculator.NumberBase]*gtk.Label
	bitDisplay    *gtk.Label
	hexButtons    []*gtk.Button
	bitWidth      calculator.BitWidth
	shiftAmount   int
	bitWidthLabel *gtk.Label
	shiftLabel    *gtk.Label
	bitPosEntry   *gtk.Entry

	// Date calculator widgets
	startDateEntry   *gtk.Entry
	endDateEntry     *gtk.Entry
	dateResultLbl    *gtk.Label
	addSubEntry      *gtk.Entry
	addSubResult     *gtk.Label
	birthDateEntry   *gtk.Entry
	ageResultLbl     *gtk.Label
	timestampEntry   *gtk.Entry
	timestampResult  *gtk.Label

	// Scientific mode
	angleModeLbl *gtk.Label
}

func main() {
	app := gtk.NewApplication("com.switchcalc.app", 0)
	app.ConnectActivate(func() {
		activate(app)
	})

	if code := app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	calcApp := &App{
		engine:      calculator.NewEngine(),
		dateCalc:    calculator.NewDateTimeCalc(),
		mode:        ModeStandard,
		modeButtons: make(map[CalculatorMode]*gtk.ToggleButton),
		baseLabels:  make(map[calculator.NumberBase]*gtk.Label),
		bitWidth:    calculator.Bits32,
		shiftAmount: 1,
	}

	calcApp.window = gtk.NewApplicationWindow(app)
	calcApp.window.SetTitle("SwitchCalc")
	calcApp.window.SetDefaultSize(400, 600)
	calcApp.window.SetResizable(true)

	mainBox := gtk.NewBox(gtk.OrientationVertical, 0)
	mainBox.AddCSSClass("main-container")

	// Mode selector
	modeBox := calcApp.createModeSelector()
	mainBox.Append(modeBox)

	// Display area
	displayBox := calcApp.createDisplayArea()
	mainBox.Append(displayBox)

	// Stack for different calculator modes
	calcApp.mainStack = gtk.NewStack()
	calcApp.mainStack.SetTransitionType(gtk.StackTransitionTypeCrossfade)
	calcApp.mainStack.SetTransitionDuration(200)
	calcApp.mainStack.SetVExpand(true)

	// Add mode pages
	standardPage := calcApp.createStandardKeypad()
	calcApp.mainStack.AddNamed(standardPage, "standard")

	scientificPage := calcApp.createScientificKeypad()
	calcApp.mainStack.AddNamed(scientificPage, "scientific")

	programmerPage := calcApp.createProgrammerKeypad()
	calcApp.mainStack.AddNamed(programmerPage, "programmer")

	dateTimePage := calcApp.createDateTimePage()
	calcApp.mainStack.AddNamed(dateTimePage, "datetime")

	mainBox.Append(calcApp.mainStack)

	// Apply CSS
	calcApp.applyCSS()

	// Setup keyboard handling
	calcApp.setupKeyboardHandling()

	calcApp.window.SetChild(mainBox)
	calcApp.window.Show()
}

func (a *App) createModeSelector() *gtk.Box {
	box := gtk.NewBox(gtk.OrientationHorizontal, 4)
	box.SetMarginTop(8)
	box.SetMarginBottom(8)
	box.SetMarginStart(8)
	box.SetMarginEnd(8)
	box.SetHomogeneous(true)

	modes := []struct {
		mode  CalculatorMode
		label string
		name  string
	}{
		{ModeStandard, "Standard", "standard"},
		{ModeScientific, "Scientific", "scientific"},
		{ModeProgrammer, "Programmer", "programmer"},
		{ModeDateTime, "Date", "datetime"},
	}

	for _, m := range modes {
		btn := gtk.NewToggleButton()
		btn.SetLabel(m.label)
		btn.AddCSSClass("mode-button")
		if m.mode == ModeStandard {
			btn.SetActive(true)
		}

		mode := m.mode
		name := m.name
		btn.ConnectClicked(func() {
			a.setMode(mode, name)
		})

		a.modeButtons[m.mode] = btn
		box.Append(btn)
	}

	return box
}

func (a *App) setMode(mode CalculatorMode, stackName string) {
	a.mode = mode
	a.mainStack.SetVisibleChildName(stackName)

	for m, btn := range a.modeButtons {
		btn.SetActive(m == mode)
	}

	if mode == ModeProgrammer {
		a.engine.SetNumberBase(calculator.Decimal)
		a.updateProgrammerDisplay()
	}
}

func (a *App) createDisplayArea() *gtk.Box {
	box := gtk.NewBox(gtk.OrientationVertical, 4)
	box.AddCSSClass("display-area")
	box.SetMarginStart(16)
	box.SetMarginEnd(16)
	box.SetMarginTop(8)
	box.SetMarginBottom(8)

	a.expressionLbl = gtk.NewLabel("")
	a.expressionLbl.AddCSSClass("expression-label")
	a.expressionLbl.SetXAlign(1)
	box.Append(a.expressionLbl)

	a.display = gtk.NewLabel("0")
	a.display.AddCSSClass("main-display")
	a.display.SetXAlign(1)
	a.display.SetSelectable(true)
	box.Append(a.display)

	return box
}

func (a *App) createStandardKeypad() *gtk.Box {
	box := gtk.NewBox(gtk.OrientationVertical, 4)
	box.SetMarginStart(8)
	box.SetMarginEnd(8)
	box.SetMarginBottom(8)
	box.SetVExpand(true)

	// Memory buttons row
	memRow := gtk.NewBox(gtk.OrientationHorizontal, 4)
	memRow.SetHomogeneous(true)
	memButtons := []struct {
		label string
		fn    func()
	}{
		{"MC", a.engine.MemoryClear},
		{"MR", func() { a.engine.MemoryRecall(); a.updateDisplay() }},
		{"M+", a.engine.MemoryAdd},
		{"M-", a.engine.MemorySubtract},
		{"MS", a.engine.MemoryStore},
	}
	for _, mb := range memButtons {
		btn := a.createButton(mb.label, "memory-button", mb.fn)
		memRow.Append(btn)
	}
	box.Append(memRow)

	// Main keypad
	keys := [][]struct {
		label string
		class string
		fn    func()
	}{
		{
			{"%", "function-button", func() { a.engine.Percent(); a.updateDisplay() }},
			{"CE", "function-button", func() { a.engine.ClearEntry(); a.updateDisplay() }},
			{"C", "function-button", func() { a.engine.Clear(); a.updateDisplay() }},
			{"⌫", "function-button", func() { a.engine.Backspace(); a.updateDisplay() }},
		},
		{
			{"1/x", "function-button", func() { a.engine.Reciprocal(); a.updateDisplay() }},
			{"x²", "function-button", func() { a.engine.Square(); a.updateDisplay() }},
			{"√", "function-button", func() { a.engine.Sqrt(); a.updateDisplay() }},
			{"÷", "operator-button", func() { a.engine.SetOperation(calculator.OpDivide); a.updateExpression() }},
		},
		{
			{"7", "number-button", func() { a.engine.InputDigit("7"); a.updateDisplay() }},
			{"8", "number-button", func() { a.engine.InputDigit("8"); a.updateDisplay() }},
			{"9", "number-button", func() { a.engine.InputDigit("9"); a.updateDisplay() }},
			{"×", "operator-button", func() { a.engine.SetOperation(calculator.OpMultiply); a.updateExpression() }},
		},
		{
			{"4", "number-button", func() { a.engine.InputDigit("4"); a.updateDisplay() }},
			{"5", "number-button", func() { a.engine.InputDigit("5"); a.updateDisplay() }},
			{"6", "number-button", func() { a.engine.InputDigit("6"); a.updateDisplay() }},
			{"−", "operator-button", func() { a.engine.SetOperation(calculator.OpSubtract); a.updateExpression() }},
		},
		{
			{"1", "number-button", func() { a.engine.InputDigit("1"); a.updateDisplay() }},
			{"2", "number-button", func() { a.engine.InputDigit("2"); a.updateDisplay() }},
			{"3", "number-button", func() { a.engine.InputDigit("3"); a.updateDisplay() }},
			{"+", "operator-button", func() { a.engine.SetOperation(calculator.OpAdd); a.updateExpression() }},
		},
		{
			{"±", "function-button", func() { a.engine.Negate(); a.updateDisplay() }},
			{"0", "number-button", func() { a.engine.InputDigit("0"); a.updateDisplay() }},
			{".", "number-button", func() { a.engine.InputDecimal(); a.updateDisplay() }},
			{"=", "equals-button", func() { a.engine.Calculate(); a.updateDisplay(); a.expressionLbl.SetText("") }},
		},
	}

	for _, row := range keys {
		rowBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
		rowBox.SetHomogeneous(true)
		rowBox.SetVExpand(true)
		for _, key := range row {
			btn := a.createButton(key.label, key.class, key.fn)
			rowBox.Append(btn)
		}
		box.Append(rowBox)
	}

	return box
}

func (a *App) createScientificKeypad() *gtk.Box {
	box := gtk.NewBox(gtk.OrientationVertical, 4)
	box.SetMarginStart(8)
	box.SetMarginEnd(8)
	box.SetMarginBottom(8)
	box.SetVExpand(true)

	// Angle mode selector
	angleModeBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
	angleModeBox.SetHomogeneous(true)

	degBtn := gtk.NewToggleButton()
	degBtn.SetLabel("DEG")
	degBtn.SetActive(true)
	degBtn.AddCSSClass("angle-button")

	radBtn := gtk.NewToggleButton()
	radBtn.SetLabel("RAD")
	radBtn.AddCSSClass("angle-button")

	gradBtn := gtk.NewToggleButton()
	gradBtn.SetLabel("GRAD")
	gradBtn.AddCSSClass("angle-button")

	degBtn.ConnectClicked(func() {
		a.engine.SetAngleMode(calculator.Degrees)
		degBtn.SetActive(true)
		radBtn.SetActive(false)
		gradBtn.SetActive(false)
	})
	radBtn.ConnectClicked(func() {
		a.engine.SetAngleMode(calculator.Radians)
		degBtn.SetActive(false)
		radBtn.SetActive(true)
		gradBtn.SetActive(false)
	})
	gradBtn.ConnectClicked(func() {
		a.engine.SetAngleMode(calculator.Gradians)
		degBtn.SetActive(false)
		radBtn.SetActive(false)
		gradBtn.SetActive(true)
	})

	angleModeBox.Append(degBtn)
	angleModeBox.Append(radBtn)
	angleModeBox.Append(gradBtn)
	box.Append(angleModeBox)

	// Scientific function keys
	sciKeys := [][]struct {
		label string
		class string
		fn    func()
	}{
		{
			{"sin", "sci-button", func() { a.engine.Sin(); a.updateDisplay() }},
			{"cos", "sci-button", func() { a.engine.Cos(); a.updateDisplay() }},
			{"tan", "sci-button", func() { a.engine.Tan(); a.updateDisplay() }},
			{"log", "sci-button", func() { a.engine.Log(); a.updateDisplay() }},
			{"ln", "sci-button", func() { a.engine.Ln(); a.updateDisplay() }},
		},
		{
			{"asin", "sci-button", func() { a.engine.Asin(); a.updateDisplay() }},
			{"acos", "sci-button", func() { a.engine.Acos(); a.updateDisplay() }},
			{"atan", "sci-button", func() { a.engine.Atan(); a.updateDisplay() }},
			{"10ˣ", "sci-button", func() { a.engine.Exp10(); a.updateDisplay() }},
			{"eˣ", "sci-button", func() { a.engine.Exp(); a.updateDisplay() }},
		},
		{
			{"sinh", "sci-button", func() { a.engine.Sinh(); a.updateDisplay() }},
			{"cosh", "sci-button", func() { a.engine.Cosh(); a.updateDisplay() }},
			{"tanh", "sci-button", func() { a.engine.Tanh(); a.updateDisplay() }},
			{"x²", "sci-button", func() { a.engine.Square(); a.updateDisplay() }},
			{"x³", "sci-button", func() { a.engine.Cube(); a.updateDisplay() }},
		},
		{
			{"asinh", "sci-button", func() { a.engine.Asinh(); a.updateDisplay() }},
			{"acosh", "sci-button", func() { a.engine.Acosh(); a.updateDisplay() }},
			{"atanh", "sci-button", func() { a.engine.Atanh(); a.updateDisplay() }},
			{"log₂", "sci-button", func() { a.engine.Log2(); a.updateDisplay() }},
			{"2ˣ", "sci-button", func() { a.engine.Exp2(); a.updateDisplay() }},
		},
		{
			{"π", "sci-button", func() { a.engine.Pi(); a.updateDisplay() }},
			{"e", "sci-button", func() { a.engine.E(); a.updateDisplay() }},
			{"n!", "sci-button", func() { a.engine.Factorial(); a.updateDisplay() }},
			{"√", "sci-button", func() { a.engine.Sqrt(); a.updateDisplay() }},
			{"∛", "sci-button", func() { a.engine.Cbrt(); a.updateDisplay() }},
		},
		{
			{"xʸ", "sci-button", func() { a.engine.SetOperation(calculator.OpPower); a.updateExpression() }},
			{"mod", "sci-button", func() { a.engine.SetOperation(calculator.OpModulo); a.updateExpression() }},
			{"|x|", "sci-button", func() { a.engine.Abs(); a.updateDisplay() }},
			{"⌊x⌋", "sci-button", func() { a.engine.Floor(); a.updateDisplay() }},
			{"⌈x⌉", "sci-button", func() { a.engine.Ceil(); a.updateDisplay() }},
		},
		{
			{"round", "sci-button", func() { a.engine.Round(); a.updateDisplay() }},
			{"1/x", "sci-button", func() { a.engine.Reciprocal(); a.updateDisplay() }},
			{"%", "sci-button", func() { a.engine.Percent(); a.updateDisplay() }},
			{"±", "sci-button", func() { a.engine.Negate(); a.updateDisplay() }},
			{"EXP", "sci-button", func() { a.engine.InputExponent(); a.updateDisplay() }},
		},
	}

	for _, row := range sciKeys {
		rowBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
		rowBox.SetHomogeneous(true)
		for _, key := range row {
			btn := a.createButton(key.label, key.class, key.fn)
			rowBox.Append(btn)
		}
		box.Append(rowBox)
	}

	// Standard number pad
	numKeys := [][]struct {
		label string
		class string
		fn    func()
	}{
		{
			{"C", "function-button", func() { a.engine.Clear(); a.updateDisplay() }},
			{"CE", "function-button", func() { a.engine.ClearEntry(); a.updateDisplay() }},
			{"⌫", "function-button", func() { a.engine.Backspace(); a.updateDisplay() }},
			{"÷", "operator-button", func() { a.engine.SetOperation(calculator.OpDivide); a.updateExpression() }},
		},
		{
			{"7", "number-button", func() { a.engine.InputDigit("7"); a.updateDisplay() }},
			{"8", "number-button", func() { a.engine.InputDigit("8"); a.updateDisplay() }},
			{"9", "number-button", func() { a.engine.InputDigit("9"); a.updateDisplay() }},
			{"×", "operator-button", func() { a.engine.SetOperation(calculator.OpMultiply); a.updateExpression() }},
		},
		{
			{"4", "number-button", func() { a.engine.InputDigit("4"); a.updateDisplay() }},
			{"5", "number-button", func() { a.engine.InputDigit("5"); a.updateDisplay() }},
			{"6", "number-button", func() { a.engine.InputDigit("6"); a.updateDisplay() }},
			{"−", "operator-button", func() { a.engine.SetOperation(calculator.OpSubtract); a.updateExpression() }},
		},
		{
			{"1", "number-button", func() { a.engine.InputDigit("1"); a.updateDisplay() }},
			{"2", "number-button", func() { a.engine.InputDigit("2"); a.updateDisplay() }},
			{"3", "number-button", func() { a.engine.InputDigit("3"); a.updateDisplay() }},
			{"+", "operator-button", func() { a.engine.SetOperation(calculator.OpAdd); a.updateExpression() }},
		},
		{
			{"±", "function-button", func() { a.engine.Negate(); a.updateDisplay() }},
			{"0", "number-button", func() { a.engine.InputDigit("0"); a.updateDisplay() }},
			{".", "number-button", func() { a.engine.InputDecimal(); a.updateDisplay() }},
			{"=", "equals-button", func() { a.engine.Calculate(); a.updateDisplay(); a.expressionLbl.SetText("") }},
		},
	}

	for _, row := range numKeys {
		rowBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
		rowBox.SetHomogeneous(true)
		rowBox.SetVExpand(true)
		for _, key := range row {
			btn := a.createButton(key.label, key.class, key.fn)
			rowBox.Append(btn)
		}
		box.Append(rowBox)
	}

	return box
}

func (a *App) createProgrammerKeypad() *gtk.Box {
	box := gtk.NewBox(gtk.OrientationVertical, 4)
	box.SetMarginStart(8)
	box.SetMarginEnd(8)
	box.SetMarginBottom(8)
	box.SetVExpand(true)

	// Bit width and shift amount controls
	controlsRow := gtk.NewBox(gtk.OrientationHorizontal, 8)
	controlsRow.SetMarginBottom(4)

	// Bit width selector
	bitWidthBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
	bitWidthLabel := gtk.NewLabel("Width:")
	bitWidthLabel.AddCSSClass("dim-label")
	bitWidthBox.Append(bitWidthLabel)

	bitWidths := []struct {
		label string
		width calculator.BitWidth
	}{
		{"8", calculator.Bits8},
		{"16", calculator.Bits16},
		{"32", calculator.Bits32},
		{"64", calculator.Bits64},
	}

	for _, bw := range bitWidths {
		btn := gtk.NewToggleButton()
		btn.SetLabel(bw.label)
		btn.AddCSSClass("bit-width-button")
		if bw.width == calculator.Bits32 {
			btn.SetActive(true)
		}
		width := bw.width
		btn.ConnectClicked(func() {
			a.bitWidth = width
			a.updateProgrammerDisplay()
		})
		bitWidthBox.Append(btn)
	}
	controlsRow.Append(bitWidthBox)

	// Shift amount controls
	shiftBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
	shiftBox.SetHExpand(true)
	shiftBox.SetHAlign(gtk.AlignEnd)

	shiftLabelPre := gtk.NewLabel("Shift:")
	shiftLabelPre.AddCSSClass("dim-label")
	shiftBox.Append(shiftLabelPre)

	shiftMinus := gtk.NewButton()
	shiftMinus.SetLabel("-")
	shiftMinus.AddCSSClass("shift-ctrl-button")
	shiftMinus.ConnectClicked(func() {
		if a.shiftAmount > 1 {
			a.shiftAmount--
			a.shiftLabel.SetText(fmt.Sprintf("%d", a.shiftAmount))
		}
	})
	shiftBox.Append(shiftMinus)

	a.shiftLabel = gtk.NewLabel("1")
	a.shiftLabel.AddCSSClass("shift-amount")
	a.shiftLabel.SetWidthChars(2)
	shiftBox.Append(a.shiftLabel)

	shiftPlus := gtk.NewButton()
	shiftPlus.SetLabel("+")
	shiftPlus.AddCSSClass("shift-ctrl-button")
	shiftPlus.ConnectClicked(func() {
		if a.shiftAmount < 63 {
			a.shiftAmount++
			a.shiftLabel.SetText(fmt.Sprintf("%d", a.shiftAmount))
		}
	})
	shiftBox.Append(shiftPlus)

	controlsRow.Append(shiftBox)
	box.Append(controlsRow)

	// Base display area
	baseBox := gtk.NewBox(gtk.OrientationVertical, 2)
	baseBox.AddCSSClass("base-display")
	baseBox.SetMarginBottom(4)

	bases := []struct {
		name string
		base calculator.NumberBase
	}{
		{"HEX", calculator.Hexadecimal},
		{"DEC", calculator.Decimal},
		{"OCT", calculator.Octal},
		{"BIN", calculator.Binary},
	}

	for _, b := range bases {
		row := gtk.NewBox(gtk.OrientationHorizontal, 8)
		label := gtk.NewLabel(b.name)
		label.AddCSSClass("base-label")
		label.SetWidthChars(4)

		valueLabel := gtk.NewLabel("0")
		valueLabel.AddCSSClass("base-value")
		valueLabel.SetXAlign(0)
		valueLabel.SetHExpand(true)
		valueLabel.SetSelectable(true)

		a.baseLabels[b.base] = valueLabel

		base := b.base
		clickCtrl := gtk.NewGestureClick()
		clickCtrl.ConnectPressed(func(n int, x, y float64) {
			a.engine.SetNumberBase(base)
			a.updateProgrammerDisplay()
		})
		row.AddController(clickCtrl)

		row.Append(label)
		row.Append(valueLabel)
		baseBox.Append(row)
	}
	box.Append(baseBox)

	// Bit display
	a.bitDisplay = gtk.NewLabel("0000 0000 0000 0000 0000 0000 0000 0000")
	a.bitDisplay.AddCSSClass("bit-display")
	a.bitDisplay.SetSelectable(true)
	box.Append(a.bitDisplay)

	// Hex digits row
	hexRow := gtk.NewBox(gtk.OrientationHorizontal, 4)
	hexRow.SetHomogeneous(true)

	hexDigits := []string{"A", "B", "C", "D", "E", "F"}
	for _, h := range hexDigits {
		digit := h
		btn := a.createButton(h, "hex-button", func() {
			a.engine.InputHexDigit(digit)
			a.updateProgrammerDisplay()
		})
		a.hexButtons = append(a.hexButtons, btn)
		hexRow.Append(btn)
	}
	box.Append(hexRow)

	// Bitwise operations - Row 1 (basic)
	bitwiseRow1 := gtk.NewBox(gtk.OrientationHorizontal, 4)
	bitwiseRow1.SetHomogeneous(true)
	bitwiseKeys1 := []struct {
		label string
		fn    func()
	}{
		{"AND", func() { a.engine.SetBitwiseOperation(calculator.BitOpAnd); a.updateExpression() }},
		{"OR", func() { a.engine.SetBitwiseOperation(calculator.BitOpOr); a.updateExpression() }},
		{"XOR", func() { a.engine.SetBitwiseOperation(calculator.BitOpXor); a.updateExpression() }},
		{"NOT", func() { a.engine.Not(); a.updateProgrammerDisplay() }},
	}
	for _, k := range bitwiseKeys1 {
		btn := a.createButton(k.label, "bitwise-button", k.fn)
		bitwiseRow1.Append(btn)
	}
	box.Append(bitwiseRow1)

	// Bitwise operations - Row 2 (NAND, NOR, shifts)
	bitwiseRow2 := gtk.NewBox(gtk.OrientationHorizontal, 4)
	bitwiseRow2.SetHomogeneous(true)
	bitwiseKeys2 := []struct {
		label string
		fn    func()
	}{
		{"NAND", func() { a.engine.SetBitwiseOperation(calculator.BitOpNand); a.updateExpression() }},
		{"NOR", func() { a.engine.SetBitwiseOperation(calculator.BitOpNor); a.updateExpression() }},
		{"<<", func() { a.engine.LeftShift(uint(a.shiftAmount)); a.updateProgrammerDisplay() }},
		{">>", func() { a.engine.RightShift(uint(a.shiftAmount)); a.updateProgrammerDisplay() }},
	}
	for _, k := range bitwiseKeys2 {
		btn := a.createButton(k.label, "bitwise-button", k.fn)
		bitwiseRow2.Append(btn)
	}
	box.Append(bitwiseRow2)

	// Bitwise operations - Row 3 (rotate, count)
	bitwiseRow3 := gtk.NewBox(gtk.OrientationHorizontal, 4)
	bitwiseRow3.SetHomogeneous(true)
	bitwiseKeys3 := []struct {
		label string
		fn    func()
	}{
		{"RoL", func() { a.engine.RotateLeft(uint(a.shiftAmount), a.bitWidth); a.updateProgrammerDisplay() }},
		{"RoR", func() { a.engine.RotateRight(uint(a.shiftAmount), a.bitWidth); a.updateProgrammerDisplay() }},
		{"Cnt", func() { a.engine.CountBits(); a.updateProgrammerDisplay() }},
		{"2's", func() { a.engine.TwosComplement(a.bitWidth); a.updateProgrammerDisplay() }},
	}
	for _, k := range bitwiseKeys3 {
		btn := a.createButton(k.label, "bitwise-button", k.fn)
		bitwiseRow3.Append(btn)
	}
	box.Append(bitwiseRow3)

	// Bitwise operations - Row 4 (advanced)
	bitwiseRow4 := gtk.NewBox(gtk.OrientationHorizontal, 4)
	bitwiseRow4.SetHomogeneous(true)
	bitwiseKeys4 := []struct {
		label string
		fn    func()
	}{
		{"LZ", func() { a.engine.LeadingZeros(a.bitWidth); a.updateProgrammerDisplay() }},
		{"TZ", func() { a.engine.TrailingZeros(); a.updateProgrammerDisplay() }},
		{"Swap", func() { a.engine.ByteSwap(a.bitWidth); a.updateProgrammerDisplay() }},
		{"Tog", func() { a.toggleBitAtPosition(); a.updateProgrammerDisplay() }},
	}
	for _, k := range bitwiseKeys4 {
		btn := a.createButton(k.label, "bitwise-button", k.fn)
		bitwiseRow4.Append(btn)
	}
	box.Append(bitwiseRow4)

	// Number pad
	numKeys := [][]struct {
		label string
		class string
		fn    func()
	}{
		{
			{"C", "function-button", func() { a.engine.Clear(); a.updateProgrammerDisplay() }},
			{"CE", "function-button", func() { a.engine.ClearEntry(); a.updateProgrammerDisplay() }},
			{"⌫", "function-button", func() { a.engine.Backspace(); a.updateProgrammerDisplay() }},
			{"÷", "operator-button", func() { a.engine.SetOperation(calculator.OpDivide); a.updateExpression() }},
		},
		{
			{"7", "number-button", func() { a.engine.InputDigit("7"); a.updateProgrammerDisplay() }},
			{"8", "number-button", func() { a.engine.InputDigit("8"); a.updateProgrammerDisplay() }},
			{"9", "number-button", func() { a.engine.InputDigit("9"); a.updateProgrammerDisplay() }},
			{"×", "operator-button", func() { a.engine.SetOperation(calculator.OpMultiply); a.updateExpression() }},
		},
		{
			{"4", "number-button", func() { a.engine.InputDigit("4"); a.updateProgrammerDisplay() }},
			{"5", "number-button", func() { a.engine.InputDigit("5"); a.updateProgrammerDisplay() }},
			{"6", "number-button", func() { a.engine.InputDigit("6"); a.updateProgrammerDisplay() }},
			{"−", "operator-button", func() { a.engine.SetOperation(calculator.OpSubtract); a.updateExpression() }},
		},
		{
			{"1", "number-button", func() { a.engine.InputDigit("1"); a.updateProgrammerDisplay() }},
			{"2", "number-button", func() { a.engine.InputDigit("2"); a.updateProgrammerDisplay() }},
			{"3", "number-button", func() { a.engine.InputDigit("3"); a.updateProgrammerDisplay() }},
			{"+", "operator-button", func() { a.engine.SetOperation(calculator.OpAdd); a.updateExpression() }},
		},
		{
			{"±", "function-button", func() { a.engine.Negate(); a.updateProgrammerDisplay() }},
			{"0", "number-button", func() { a.engine.InputDigit("0"); a.updateProgrammerDisplay() }},
			{"mod", "function-button", func() { a.engine.SetOperation(calculator.OpModulo); a.updateExpression() }},
			{"=", "equals-button", func() { a.calculateProgrammer() }},
		},
	}

	for _, row := range numKeys {
		rowBox := gtk.NewBox(gtk.OrientationHorizontal, 4)
		rowBox.SetHomogeneous(true)
		rowBox.SetVExpand(true)
		for _, key := range row {
			btn := a.createButton(key.label, key.class, key.fn)
			rowBox.Append(btn)
		}
		box.Append(rowBox)
	}

	return box
}

func (a *App) toggleBitAtPosition() {
	// Toggle bit at position equal to shift amount (0-indexed)
	pos := a.shiftAmount - 1
	if pos < 0 {
		pos = 0
	}
	a.engine.ToggleBit(uint(pos))
}

func (a *App) createDateTimePage() *gtk.Box {
	// Create scrollable container for all date content
	scrollWin := gtk.NewScrolledWindow()
	scrollWin.SetVExpand(true)
	scrollWin.SetPolicy(gtk.PolicyNever, gtk.PolicyAutomatic)

	box := gtk.NewBox(gtk.OrientationVertical, 12)
	box.SetMarginStart(16)
	box.SetMarginEnd(16)
	box.SetMarginTop(16)
	box.SetMarginBottom(16)

	// Date Difference Calculator
	diffFrame := gtk.NewFrame("Date Difference")
	diffBox := gtk.NewBox(gtk.OrientationVertical, 8)
	diffBox.SetMarginStart(12)
	diffBox.SetMarginEnd(12)
	diffBox.SetMarginTop(12)
	diffBox.SetMarginBottom(12)

	// Start date
	startBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	startLabel := gtk.NewLabel("From:")
	startLabel.SetWidthChars(6)
	a.startDateEntry = gtk.NewEntry()
	a.startDateEntry.SetPlaceholderText("DD/MM/YYYY")
	a.startDateEntry.SetText(time.Now().Format("02/01/2006"))
	startBox.Append(startLabel)
	startBox.Append(a.startDateEntry)
	diffBox.Append(startBox)

	// End date
	endBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	endLabel := gtk.NewLabel("To:")
	endLabel.SetWidthChars(6)
	a.endDateEntry = gtk.NewEntry()
	a.endDateEntry.SetPlaceholderText("DD/MM/YYYY")
	a.endDateEntry.SetText(time.Now().Format("02/01/2006"))
	endBox.Append(endLabel)
	endBox.Append(a.endDateEntry)
	diffBox.Append(endBox)

	// Calculate button
	calcDiffBtn := gtk.NewButton()
	calcDiffBtn.SetLabel("Calculate Difference")
	calcDiffBtn.AddCSSClass("suggested-action")
	calcDiffBtn.ConnectClicked(func() {
		a.calculateDateDifference()
	})
	diffBox.Append(calcDiffBtn)

	// Result
	a.dateResultLbl = gtk.NewLabel("")
	a.dateResultLbl.AddCSSClass("date-result")
	a.dateResultLbl.SetWrap(true)
	a.dateResultLbl.SetSelectable(true)
	diffBox.Append(a.dateResultLbl)

	diffFrame.SetChild(diffBox)
	box.Append(diffFrame)

	// Add/Subtract from Date
	addFrame := gtk.NewFrame("Add/Subtract from Date")
	addBox := gtk.NewBox(gtk.OrientationVertical, 8)
	addBox.SetMarginStart(12)
	addBox.SetMarginEnd(12)
	addBox.SetMarginTop(12)
	addBox.SetMarginBottom(12)

	// Instructions
	instrLabel := gtk.NewLabel("Enter: +/- years, months, days (e.g., +1y 2m -5d)")
	instrLabel.AddCSSClass("dim-label")
	addBox.Append(instrLabel)

	// Input
	a.addSubEntry = gtk.NewEntry()
	a.addSubEntry.SetPlaceholderText("+1y 2m 3d")
	addBox.Append(a.addSubEntry)

	// Calculate button
	addSubBtn := gtk.NewButton()
	addSubBtn.SetLabel("Calculate")
	addSubBtn.AddCSSClass("suggested-action")
	addSubBtn.ConnectClicked(func() {
		a.calculateAddSubtract()
	})
	addBox.Append(addSubBtn)

	// Result
	a.addSubResult = gtk.NewLabel("")
	a.addSubResult.AddCSSClass("date-result")
	a.addSubResult.SetSelectable(true)
	addBox.Append(a.addSubResult)

	addFrame.SetChild(addBox)
	box.Append(addFrame)

	// Age Calculator
	ageFrame := gtk.NewFrame("Age Calculator")
	ageBox := gtk.NewBox(gtk.OrientationVertical, 8)
	ageBox.SetMarginStart(12)
	ageBox.SetMarginEnd(12)
	ageBox.SetMarginTop(12)
	ageBox.SetMarginBottom(12)

	// Birth date
	birthBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	birthLabel := gtk.NewLabel("Birth:")
	birthLabel.SetWidthChars(6)
	a.birthDateEntry = gtk.NewEntry()
	a.birthDateEntry.SetPlaceholderText("DD/MM/YYYY")
	birthBox.Append(birthLabel)
	birthBox.Append(a.birthDateEntry)
	ageBox.Append(birthBox)

	// Calculate age button
	calcAgeBtn := gtk.NewButton()
	calcAgeBtn.SetLabel("Calculate Age")
	calcAgeBtn.AddCSSClass("suggested-action")
	calcAgeBtn.ConnectClicked(func() {
		a.calculateAge()
	})
	ageBox.Append(calcAgeBtn)

	// Age result
	a.ageResultLbl = gtk.NewLabel("")
	a.ageResultLbl.AddCSSClass("date-result")
	a.ageResultLbl.SetWrap(true)
	a.ageResultLbl.SetSelectable(true)
	ageBox.Append(a.ageResultLbl)

	ageFrame.SetChild(ageBox)
	box.Append(ageFrame)

	// Unix Timestamp Converter
	timestampFrame := gtk.NewFrame("Unix Timestamp")
	timestampBox := gtk.NewBox(gtk.OrientationVertical, 8)
	timestampBox.SetMarginStart(12)
	timestampBox.SetMarginEnd(12)
	timestampBox.SetMarginTop(12)
	timestampBox.SetMarginBottom(12)

	// Timestamp input
	tsInputBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	tsLabel := gtk.NewLabel("Unix:")
	tsLabel.SetWidthChars(6)
	a.timestampEntry = gtk.NewEntry()
	a.timestampEntry.SetPlaceholderText("e.g., 1700000000")
	tsInputBox.Append(tsLabel)
	tsInputBox.Append(a.timestampEntry)
	timestampBox.Append(tsInputBox)

	// Button row
	tsBtnBox := gtk.NewBox(gtk.OrientationHorizontal, 8)
	tsBtnBox.SetHomogeneous(true)

	// Convert timestamp to date
	tsToDateBtn := gtk.NewButton()
	tsToDateBtn.SetLabel("To Date")
	tsToDateBtn.AddCSSClass("suggested-action")
	tsToDateBtn.ConnectClicked(func() {
		a.timestampToDate()
	})
	tsBtnBox.Append(tsToDateBtn)

	// Get current timestamp
	nowTsBtn := gtk.NewButton()
	nowTsBtn.SetLabel("Now")
	nowTsBtn.ConnectClicked(func() {
		a.timestampEntry.SetText(fmt.Sprintf("%d", time.Now().Unix()))
		a.timestampToDate()
	})
	tsBtnBox.Append(nowTsBtn)

	// Date to timestamp
	dateToTsBtn := gtk.NewButton()
	dateToTsBtn.SetLabel("From Start Date")
	dateToTsBtn.ConnectClicked(func() {
		a.dateToTimestamp()
	})
	tsBtnBox.Append(dateToTsBtn)

	timestampBox.Append(tsBtnBox)

	// Timestamp result
	a.timestampResult = gtk.NewLabel("")
	a.timestampResult.AddCSSClass("date-result")
	a.timestampResult.SetWrap(true)
	a.timestampResult.SetSelectable(true)
	timestampBox.Append(a.timestampResult)

	timestampFrame.SetChild(timestampBox)
	box.Append(timestampFrame)

	// Quick info section
	infoFrame := gtk.NewFrame("Today's Info")
	infoBox := gtk.NewBox(gtk.OrientationVertical, 4)
	infoBox.SetMarginStart(12)
	infoBox.SetMarginEnd(12)
	infoBox.SetMarginTop(12)
	infoBox.SetMarginBottom(12)

	now := time.Now()
	_, week := now.ISOWeek()

	a.dateCalc.StartDate = now
	leapYear := "No"
	if a.dateCalc.IsLeapYear() {
		leapYear = "Yes"
	}

	infoItems := []string{
		fmt.Sprintf("Today: %s", now.Format("Monday, 02 January 2006")),
		fmt.Sprintf("Week number: %d", week),
		fmt.Sprintf("Day of year: %d / %d", now.YearDay(), func() int {
			if a.dateCalc.IsLeapYear() {
				return 366
			}
			return 365
		}()),
		fmt.Sprintf("Days until end of year: %d", daysUntilEndOfYear(now)),
		fmt.Sprintf("Days until end of month: %d", a.dateCalc.DaysUntilEndOfMonth()),
		fmt.Sprintf("Leap year: %s", leapYear),
		fmt.Sprintf("Unix timestamp: %d", now.Unix()),
	}

	for _, item := range infoItems {
		lbl := gtk.NewLabel(item)
		lbl.SetXAlign(0)
		infoBox.Append(lbl)
	}

	infoFrame.SetChild(infoBox)
	box.Append(infoFrame)

	scrollWin.SetChild(box)

	// Wrap in container box
	containerBox := gtk.NewBox(gtk.OrientationVertical, 0)
	containerBox.Append(scrollWin)

	return containerBox
}

func daysUntilEndOfYear(t time.Time) int {
	endOfYear := time.Date(t.Year(), 12, 31, 23, 59, 59, 0, t.Location())
	return int(endOfYear.Sub(t).Hours() / 24)
}

func (a *App) calculateAge() {
	birthStr := a.birthDateEntry.Text()

	birthDate, err := time.Parse("02/01/2006", birthStr)
	if err != nil {
		a.ageResultLbl.SetText("Invalid birth date format (use DD/MM/YYYY)")
		return
	}

	now := time.Now()
	if birthDate.After(now) {
		a.ageResultLbl.SetText("Birth date cannot be in the future")
		return
	}

	years, months, days := a.dateCalc.GetAge(birthDate)

	// Calculate next birthday
	nextBirthday := time.Date(now.Year(), birthDate.Month(), birthDate.Day(), 0, 0, 0, 0, time.Local)
	if nextBirthday.Before(now) || nextBirthday.Equal(now) {
		nextBirthday = nextBirthday.AddDate(1, 0, 0)
	}
	daysUntilBirthday := int(nextBirthday.Sub(now).Hours() / 24)

	// Total days alive
	totalDays := int(now.Sub(birthDate).Hours() / 24)

	result := fmt.Sprintf("Age: %d years, %d months, %d days\n\n"+
		"Total days alive: %d\n"+
		"Next birthday: %s\n"+
		"Days until birthday: %d",
		years, months, days,
		totalDays,
		nextBirthday.Format("Monday, 02 January 2006"),
		daysUntilBirthday)

	a.ageResultLbl.SetText(result)
}

func (a *App) timestampToDate() {
	tsStr := a.timestampEntry.Text()

	timestamp, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		a.timestampResult.SetText("Invalid timestamp (must be a number)")
		return
	}

	date := time.Unix(timestamp, 0)

	result := fmt.Sprintf("Date: %s\n"+
		"UTC: %s\n"+
		"ISO 8601: %s",
		date.Format("Monday, 02 January 2006 15:04:05 MST"),
		date.UTC().Format("2006-01-02 15:04:05 UTC"),
		date.Format(time.RFC3339))

	a.timestampResult.SetText(result)
}

func (a *App) dateToTimestamp() {
	startStr := a.startDateEntry.Text()

	startDate, err := time.Parse("02/01/2006", startStr)
	if err != nil {
		a.timestampResult.SetText("Invalid start date format")
		return
	}

	timestamp := startDate.Unix()
	a.timestampEntry.SetText(fmt.Sprintf("%d", timestamp))

	result := fmt.Sprintf("Timestamp: %d\n"+
		"Date: %s",
		timestamp,
		startDate.Format("Monday, 02 January 2006"))

	a.timestampResult.SetText(result)
}

func (a *App) calculateDateDifference() {
	startStr := a.startDateEntry.Text()
	endStr := a.endDateEntry.Text()

	startDate, err := time.Parse("02/01/2006", startStr)
	if err != nil {
		a.dateResultLbl.SetText("Invalid start date format")
		return
	}

	endDate, err := time.Parse("02/01/2006", endStr)
	if err != nil {
		a.dateResultLbl.SetText("Invalid end date format")
		return
	}

	a.dateCalc.StartDate = startDate
	a.dateCalc.EndDate = endDate
	diff := a.dateCalc.CalculateDifference()

	result := fmt.Sprintf("Difference: %s\n\nOr:\n• %d total days\n• %d weeks and %d days\n• %d total hours\n• Working days (excl. weekends): %d",
		calculator.FormatDifference(diff),
		diff.TotalDays,
		diff.TotalWeeks, diff.TotalDays%7,
		diff.TotalHours,
		a.dateCalc.GetWorkingDays(true),
	)

	a.dateResultLbl.SetText(result)
}

func (a *App) calculateAddSubtract() {
	input := a.addSubEntry.Text()
	startStr := a.startDateEntry.Text()

	startDate, err := time.Parse("02/01/2006", startStr)
	if err != nil {
		a.addSubResult.SetText("Invalid start date")
		return
	}

	years, months, days := parseTimeDelta(input)
	result := startDate.AddDate(years, months, days)

	a.addSubResult.SetText(fmt.Sprintf("Result: %s (%s)",
		result.Format("02/01/2006"),
		result.Format("Monday")))
}

func parseTimeDelta(input string) (years, months, days int) {
	input = strings.ToLower(strings.TrimSpace(input))
	parts := strings.Fields(input)

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if len(part) < 2 {
			continue
		}

		unit := part[len(part)-1]
		numStr := part[:len(part)-1]

		num, err := strconv.Atoi(numStr)
		if err != nil {
			continue
		}

		switch unit {
		case 'y':
			years = num
		case 'm':
			months = num
		case 'd':
			days = num
		case 'w':
			days += num * 7
		}
	}
	return
}

func (a *App) createButton(label, cssClass string, onClick func()) *gtk.Button {
	btn := gtk.NewButton()
	btn.SetLabel(label)
	btn.AddCSSClass("calc-button")
	btn.AddCSSClass(cssClass)
	btn.SetVExpand(true)
	btn.SetHExpand(true)
	btn.ConnectClicked(onClick)
	return btn
}

func (a *App) updateDisplay() {
	a.display.SetText(a.engine.Display)
}

func (a *App) updateExpression() {
	if a.engine.PendingOp != calculator.OpNone {
		a.expressionLbl.SetText(fmt.Sprintf("%s %s", a.engine.Display, a.opSymbol(a.engine.PendingOp)))
	}
	a.updateDisplay()
}

func (a *App) opSymbol(op calculator.Operation) string {
	switch op {
	case calculator.OpAdd:
		return "+"
	case calculator.OpSubtract:
		return "−"
	case calculator.OpMultiply:
		return "×"
	case calculator.OpDivide:
		return "÷"
	case calculator.OpModulo:
		return "mod"
	case calculator.OpPower:
		return "^"
	}
	return ""
}

func (a *App) updateProgrammerDisplay() {
	a.updateDisplay()

	val := int64(a.engine.CurrentValue)

	// Update all base displays
	for base, label := range a.baseLabels {
		oldBase := a.engine.NumberBase
		a.engine.NumberBase = base
		label.SetText(a.engine.FormatInBase(val))
		a.engine.NumberBase = oldBase
	}

	// Update bit display
	a.bitDisplay.SetText(a.engine.GetBinaryString(calculator.Bits32))

	// Enable/disable hex buttons based on current base
	hexEnabled := a.engine.NumberBase == calculator.Hexadecimal
	for _, btn := range a.hexButtons {
		btn.SetSensitive(hexEnabled)
	}
}

func (a *App) calculateProgrammer() {
	if int(a.engine.PendingOp) >= 100 {
		a.engine.CalculateBitwise()
	} else {
		a.engine.Calculate()
	}
	a.updateProgrammerDisplay()
	a.expressionLbl.SetText("")
}

func (a *App) setupKeyboardHandling() {
	keyCtrl := gtk.NewEventControllerKey()
	keyCtrl.ConnectKeyPressed(func(keyval, keycode uint, state gdk.ModifierType) bool {
		// Check for Ctrl modifier for scientific shortcuts
		ctrlPressed := state&gdk.ControlMask != 0

		// Handle programmer mode hex input (A-F)
		if a.mode == ModeProgrammer && a.engine.NumberBase == calculator.Hexadecimal {
			switch keyval {
			case gdk.KEY_a, gdk.KEY_A:
				a.engine.InputHexDigit("A")
				a.updateProgrammerDisplay()
				return true
			case gdk.KEY_b, gdk.KEY_B:
				if !ctrlPressed {
					a.engine.InputHexDigit("B")
					a.updateProgrammerDisplay()
					return true
				}
			case gdk.KEY_c, gdk.KEY_C:
				if !ctrlPressed {
					a.engine.InputHexDigit("C")
					a.updateProgrammerDisplay()
					return true
				}
			case gdk.KEY_d, gdk.KEY_D:
				a.engine.InputHexDigit("D")
				a.updateProgrammerDisplay()
				return true
			case gdk.KEY_e, gdk.KEY_E:
				a.engine.InputHexDigit("E")
				a.updateProgrammerDisplay()
				return true
			case gdk.KEY_f, gdk.KEY_F:
				a.engine.InputHexDigit("F")
				a.updateProgrammerDisplay()
				return true
			}
		}

		// Scientific mode shortcuts (Ctrl+key)
		if a.mode == ModeScientific && ctrlPressed {
			switch keyval {
			case gdk.KEY_s, gdk.KEY_S:
				a.engine.Sin()
				a.updateDisplay()
				return true
			case gdk.KEY_c, gdk.KEY_C:
				a.engine.Cos()
				a.updateDisplay()
				return true
			case gdk.KEY_t, gdk.KEY_T:
				a.engine.Tan()
				a.updateDisplay()
				return true
			case gdk.KEY_l, gdk.KEY_L:
				a.engine.Log()
				a.updateDisplay()
				return true
			case gdk.KEY_n, gdk.KEY_N:
				a.engine.Ln()
				a.updateDisplay()
				return true
			case gdk.KEY_r, gdk.KEY_R:
				a.engine.Sqrt()
				a.updateDisplay()
				return true
			case gdk.KEY_p, gdk.KEY_P:
				a.engine.Pi()
				a.updateDisplay()
				return true
			case gdk.KEY_e, gdk.KEY_E:
				a.engine.E()
				a.updateDisplay()
				return true
			}
		}

		// Programmer mode shortcuts
		if a.mode == ModeProgrammer {
			switch keyval {
			case gdk.KEY_ampersand:
				a.engine.SetBitwiseOperation(calculator.BitOpAnd)
				a.updateExpression()
				return true
			case gdk.KEY_bar:
				a.engine.SetBitwiseOperation(calculator.BitOpOr)
				a.updateExpression()
				return true
			case gdk.KEY_asciicircum:
				a.engine.SetBitwiseOperation(calculator.BitOpXor)
				a.updateExpression()
				return true
			case gdk.KEY_asciitilde:
				a.engine.Not()
				a.updateProgrammerDisplay()
				return true
			case gdk.KEY_less:
				a.engine.LeftShift(uint(a.shiftAmount))
				a.updateProgrammerDisplay()
				return true
			case gdk.KEY_greater:
				a.engine.RightShift(uint(a.shiftAmount))
				a.updateProgrammerDisplay()
				return true
			}
		}

		switch keyval {
		case gdk.KEY_0, gdk.KEY_KP_0:
			a.engine.InputDigit("0")
		case gdk.KEY_1, gdk.KEY_KP_1:
			a.engine.InputDigit("1")
		case gdk.KEY_2, gdk.KEY_KP_2:
			a.engine.InputDigit("2")
		case gdk.KEY_3, gdk.KEY_KP_3:
			a.engine.InputDigit("3")
		case gdk.KEY_4, gdk.KEY_KP_4:
			a.engine.InputDigit("4")
		case gdk.KEY_5, gdk.KEY_KP_5:
			a.engine.InputDigit("5")
		case gdk.KEY_6, gdk.KEY_KP_6:
			a.engine.InputDigit("6")
		case gdk.KEY_7, gdk.KEY_KP_7:
			a.engine.InputDigit("7")
		case gdk.KEY_8, gdk.KEY_KP_8:
			a.engine.InputDigit("8")
		case gdk.KEY_9, gdk.KEY_KP_9:
			a.engine.InputDigit("9")
		case gdk.KEY_period, gdk.KEY_comma, gdk.KEY_KP_Decimal:
			a.engine.InputDecimal()
		case gdk.KEY_plus, gdk.KEY_KP_Add:
			a.engine.SetOperation(calculator.OpAdd)
			a.updateExpression()
			return true
		case gdk.KEY_minus, gdk.KEY_KP_Subtract:
			a.engine.SetOperation(calculator.OpSubtract)
			a.updateExpression()
			return true
		case gdk.KEY_asterisk, gdk.KEY_KP_Multiply:
			a.engine.SetOperation(calculator.OpMultiply)
			a.updateExpression()
			return true
		case gdk.KEY_slash, gdk.KEY_KP_Divide:
			a.engine.SetOperation(calculator.OpDivide)
			a.updateExpression()
			return true
		case gdk.KEY_percent:
			a.engine.Percent()
			a.updateDisplay()
			return true
		case gdk.KEY_Return, gdk.KEY_KP_Enter, gdk.KEY_equal:
			if a.mode == ModeProgrammer {
				a.calculateProgrammer()
			} else {
				a.engine.Calculate()
				a.expressionLbl.SetText("")
			}
		case gdk.KEY_Escape:
			a.engine.Clear()
		case gdk.KEY_BackSpace:
			a.engine.Backspace()
		case gdk.KEY_Delete:
			a.engine.ClearEntry()
		default:
			return false
		}

		if a.mode == ModeProgrammer {
			a.updateProgrammerDisplay()
		} else {
			a.updateDisplay()
		}
		return true
	})
	a.window.AddController(keyCtrl)
}

func (a *App) applyCSS() {
	css := `
/* SwitchCalc - Glassmorphism with SwitchSides Branding
   Primary: Deep Burgundy #2e0000
   Font: Crimson Text (serif)
   Background: Newsprint warm tones with glass effects
*/

/* Base window - dark grey/black */
.main-container {
	background: linear-gradient(145deg, #1a1a1a 0%, #0f0f0f 50%, #000000 100%);
}

window {
	background: linear-gradient(145deg, #1a1a1a 0%, #0f0f0f 50%, #000000 100%);
}

/* Universal font - Crimson Text */
* {
	font-family: "Crimson Text", "Georgia", "Times New Roman", serif;
}

/* Glass display panel */
.display-area {
	background: alpha(#FAFAF8, 0.03);
	border: 1px solid alpha(#FAFAF8, 0.08);
	border-radius: 16px;
	padding: 16px;
	margin: 8px;
	box-shadow: 0 8px 32px alpha(black, 0.3), inset 0 1px 0 alpha(#FAFAF8, 0.05);
}

.expression-label {
	font-size: 1rem;
	font-style: italic;
	color: alpha(#FAFAF8, 0.5);
	min-height: 20px;
}

.main-display {
	font-size: 2.5rem;
	font-weight: 600;
	font-family: "Crimson Text", Georgia, serif;
	color: #FAFAF8;
	min-height: 48px;
	text-shadow: 0 2px 4px alpha(black, 0.4);
	letter-spacing: -0.01em;
}

/* Mode selector - glass buttons with burgundy accent */
.mode-button {
	background: alpha(#FAFAF8, 0.02);
	border: 1px solid alpha(#FAFAF8, 0.06);
	border-radius: 8px;
	padding: 10px 14px;
	font-weight: 600;
	font-size: 13px;
	color: alpha(#FAFAF8, 0.75);
	transition: all 200ms ease;
}

.mode-button:hover {
	background: alpha(#FAFAF8, 0.06);
	border-color: alpha(#FAFAF8, 0.12);
}

.mode-button:checked {
	background: linear-gradient(135deg, alpha(#3e0000, 0.5), alpha(#2e0000, 0.5));
	border-color: alpha(#8b5555, 0.4);
	color: #FAFAF8;
	box-shadow: 0 4px 15px alpha(#2e0000, 0.3);
}

/* Base calculator button - frosted glass */
.calc-button {
	background: alpha(#FAFAF8, 0.02);
	border: 1px solid alpha(#FAFAF8, 0.05);
	border-radius: 10px;
	font-size: 1.1rem;
	font-weight: 600;
	min-height: 44px;
	margin: 2px;
	padding: 8px 4px;
	color: #FAFAF8;
	transition: all 150ms ease;
	box-shadow: 0 2px 8px alpha(black, 0.15);
}

.calc-button:hover {
	background: alpha(#FAFAF8, 0.06);
	border-color: alpha(#FAFAF8, 0.1);
	box-shadow: 0 4px 12px alpha(black, 0.2);
}

.calc-button:active {
	background: alpha(#FAFAF8, 0.1);
}

/* Number buttons - warm glass */
.number-button {
	background: alpha(#FAFAF8, 0.03);
	border: 1px solid alpha(#FAFAF8, 0.06);
	color: #FAFAF8;
}

.number-button:hover {
	background: alpha(#FAFAF8, 0.08);
	border-color: alpha(#FAFAF8, 0.12);
}

/* Operator buttons - burgundy accent */
.operator-button {
	background: linear-gradient(135deg, alpha(#3e0000, 0.4), alpha(#2e0000, 0.4));
	border: 1px solid alpha(#8b5555, 0.3);
	color: #FAFAF8;
}

.operator-button:hover {
	background: linear-gradient(135deg, alpha(#4e2821, 0.55), alpha(#3e0000, 0.55));
	border-color: alpha(#a87070, 0.4);
	box-shadow: 0 4px 15px alpha(#2e0000, 0.25);
}

/* Equals button - rich burgundy */
.equals-button {
	background: linear-gradient(135deg, alpha(#3e0000, 0.6), alpha(#2e0000, 0.6), alpha(#1e0000, 0.6));
	border: 1px solid alpha(#8b5555, 0.35);
	color: #FAFAF8;
	font-weight: 700;
	box-shadow: 0 4px 15px alpha(#2e0000, 0.35);
}

.equals-button:hover {
	background: linear-gradient(135deg, alpha(#4e2821, 0.7), alpha(#3e0000, 0.7), alpha(#2e0000, 0.7));
	box-shadow: 0 6px 20px alpha(#2e0000, 0.45);
}

/* Function buttons - subtle warm glass */
.function-button {
	background: alpha(#FAFAF8, 0.015);
	border: 1px solid alpha(#FAFAF8, 0.04);
	color: alpha(#FAFAF8, 0.85);
	font-size: 16px;
}

.function-button:hover {
	background: alpha(#FAFAF8, 0.05);
	border-color: alpha(#FAFAF8, 0.08);
}

/* Memory buttons - minimal elegant */
.memory-button {
	font-size: 12px;
	min-height: 36px;
	background: transparent;
	border: 1px solid alpha(#FAFAF8, 0.03);
	color: alpha(#FAFAF8, 0.6);
	font-style: italic;
}

.memory-button:hover {
	background: alpha(#FAFAF8, 0.04);
	color: alpha(#FAFAF8, 0.9);
}

/* Scientific buttons - burgundy medium tint */
.sci-button {
	font-size: 0.85rem;
	min-height: 36px;
	padding: 4px 2px;
	background: alpha(#532b2b, 0.2);
	border: 1px solid alpha(#8b5555, 0.25);
	color: #ddbfbf;
}

.sci-button:hover {
	background: alpha(#532b2b, 0.35);
	border-color: alpha(#a87070, 0.35);
	box-shadow: 0 4px 12px alpha(#2e0000, 0.2);
}

/* Angle mode buttons */
.angle-button {
	font-size: 12px;
	padding: 6px 10px;
	background: alpha(#FAFAF8, 0.02);
	border: 1px solid alpha(#FAFAF8, 0.05);
	color: alpha(#FAFAF8, 0.7);
}

.angle-button:checked {
	background: alpha(#3e0000, 0.35);
	border-color: alpha(#8b5555, 0.3);
	color: #ddbfbf;
}

/* Hex buttons - warm amber (warning color from brand) */
.hex-button {
	font-size: 0.9rem;
	min-height: 36px;
	padding: 4px 2px;
	background: alpha(#92400E, 0.2);
	border: 1px solid alpha(#F59E0B, 0.25);
	color: #FCD34D;
}

.hex-button:hover {
	background: alpha(#92400E, 0.35);
	border-color: alpha(#F59E0B, 0.4);
	box-shadow: 0 4px 12px alpha(#92400E, 0.2);
}

/* Bitwise buttons - error red from brand */
.bitwise-button {
	font-size: 0.75rem;
	min-height: 36px;
	padding: 4px 2px;
	background: alpha(#991B1B, 0.2);
	border: 1px solid alpha(#EF4444, 0.25);
	color: #F87171;
}

.bitwise-button:hover {
	background: alpha(#991B1B, 0.35);
	border-color: alpha(#EF4444, 0.4);
	box-shadow: 0 4px 12px alpha(#991B1B, 0.2);
}

/* Programmer base display */
.base-display {
	background: alpha(#FAFAF8, 0.02);
	border: 1px solid alpha(#FAFAF8, 0.05);
	border-radius: 12px;
	padding: 12px;
}

.base-label {
	font-family: "SF Mono", "Consolas", monospace;
	font-weight: 600;
	font-size: 11px;
	color: alpha(#FAFAF8, 0.45);
}

.base-value {
	font-family: "SF Mono", "Consolas", monospace;
	font-size: 13px;
	color: #FAFAF8;
}

/* Bit display */
.bit-display {
	font-family: "SF Mono", "Consolas", monospace;
	font-size: 11px;
	background: alpha(#FAFAF8, 0.015);
	border: 1px solid alpha(#FAFAF8, 0.04);
	border-radius: 8px;
	padding: 10px;
	margin: 6px 0;
	color: alpha(#FAFAF8, 0.7);
}

/* Date calculator result */
.date-result {
	font-family: "Crimson Text", Georgia, serif;
	font-size: 14px;
	background: alpha(#FAFAF8, 0.02);
	border: 1px solid alpha(#FAFAF8, 0.05);
	border-radius: 12px;
	padding: 14px;
	margin-top: 10px;
	color: #FAFAF8;
	line-height: 1.6;
}

/* Frame styling for date page */
frame {
	background: alpha(#FAFAF8, 0.02);
	border: 1px solid alpha(#FAFAF8, 0.05);
	border-radius: 12px;
}

frame > border {
	border: none;
}

frame > label {
	color: #c59999;
	font-weight: 600;
	font-size: 14px;
	padding: 0 8px;
}

/* Entry fields - warm glass */
entry {
	background: alpha(#FAFAF8, 0.03);
	border: 1px solid alpha(#FAFAF8, 0.06);
	border-radius: 8px;
	padding: 8px 12px;
	color: #FAFAF8;
	caret-color: #c59999;
	font-family: "Crimson Text", Georgia, serif;
}

entry:focus {
	border-color: alpha(#3e0000, 0.5);
	box-shadow: 0 0 0 3px alpha(#2e0000, 0.15);
}

entry placeholder {
	color: alpha(#FAFAF8, 0.35);
	font-style: italic;
}

/* Suggested action button - burgundy */
.suggested-action {
	background: linear-gradient(135deg, alpha(#3e0000, 0.5), alpha(#2e0000, 0.5));
	border: 1px solid alpha(#8b5555, 0.35);
	border-radius: 8px;
	color: #FAFAF8;
	font-weight: 600;
	padding: 10px 16px;
	box-shadow: 0 4px 12px alpha(#2e0000, 0.3);
}

.suggested-action:hover {
	background: linear-gradient(135deg, alpha(#4e2821, 0.6), alpha(#3e0000, 0.6));
	box-shadow: 0 6px 16px alpha(#2e0000, 0.4);
}

/* Dim label */
.dim-label {
	font-size: 12px;
	color: alpha(#FAFAF8, 0.45);
	font-style: italic;
}

/* Scrollbar styling */
scrollbar {
	background: transparent;
}

scrollbar slider {
	background: alpha(#8b5555, 0.4);
	border-radius: 4px;
	min-width: 6px;
}

scrollbar slider:hover {
	background: alpha(#a87070, 0.5);
}

/* Selection color - burgundy */
*:selected {
	background: alpha(#3e0000, 0.5);
}

/* Labels in date section */
label {
	color: alpha(#FAFAF8, 0.85);
}

/* Bit width selector buttons */
.bit-width-button {
	font-size: 11px;
	padding: 4px 8px;
	min-height: 24px;
	background: alpha(#FAFAF8, 0.02);
	border: 1px solid alpha(#FAFAF8, 0.05);
	border-radius: 4px;
	color: alpha(#FAFAF8, 0.6);
	font-family: "SF Mono", "Consolas", monospace;
}

.bit-width-button:checked {
	background: alpha(#3e0000, 0.4);
	border-color: alpha(#8b5555, 0.35);
	color: #ddbfbf;
}

.bit-width-button:hover {
	background: alpha(#FAFAF8, 0.05);
}

/* Shift control buttons */
.shift-ctrl-button {
	font-size: 14px;
	min-width: 28px;
	min-height: 24px;
	padding: 2px 8px;
	background: alpha(#FAFAF8, 0.03);
	border: 1px solid alpha(#FAFAF8, 0.06);
	border-radius: 4px;
	color: alpha(#FAFAF8, 0.8);
}

.shift-ctrl-button:hover {
	background: alpha(#FAFAF8, 0.08);
}

/* Shift amount label */
.shift-amount {
	font-family: "SF Mono", "Consolas", monospace;
	font-size: 13px;
	font-weight: 600;
	color: #FAFAF8;
	padding: 0 4px;
}
`

	provider := gtk.NewCSSProvider()
	provider.LoadFromString(css)

	display := gdk.DisplayGetDefault()
	gtk.StyleContextAddProviderForDisplay(display, provider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func init() {
	log.SetFlags(log.Lshortfile)
}
