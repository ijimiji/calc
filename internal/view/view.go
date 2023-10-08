package view

import (
	"calc/internal/calculator"
	"log"
	"slices"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func New(calculator *calculator.Calculator) *View {
	v := &View{
		calc: calculator,
	}
	v.app = app.New()
	v.window = v.app.NewWindow("calculator")
	v.input = widget.NewEntry()
	v.result = widget.NewLabel("")
	v.result.TextStyle = fyne.TextStyle{Monospace: true}
	v.result.Alignment = fyne.TextAlignTrailing

	v.window.SetContent(
		v.numpad(),
	)
	v.window.Canvas().AddShortcut(&fyne.ShortcutCopy{}, func(shorcut fyne.Shortcut) {
		v.window.Clipboard().SetContent(v.input.Text)
		return
	})
	v.window.Canvas().AddShortcut(&fyne.ShortcutPaste{}, func(shorcut fyne.Shortcut) {
		v.input.SetText(v.window.Clipboard().Content())
		return
	})

	return v
}

type View struct {
	app       fyne.App
	window    fyne.Window
	calc      *calculator.Calculator
	stack     []string
	input     *widget.Entry
	form      *widget.Form
	result    *widget.Label
	lastInput string
	lock      bool
}

func (v *View) Show() {
	v.window.ShowAndRun()
}

func (v *View) numpad() *fyne.Container {
	equals := v.button("=", func() {
		res, err := v.calc.Eval(v.input.Text)
		log.Printf("%s = %s", v.input.Text, res)
		if err != nil {
			v.showError(err.Error())
			return
		}
		v.result.SetText(res)
	})
	equals.Importance = widget.HighImportance

	return container.NewGridWithColumns(1,
		v.input,
		v.result,
		container.NewGridWithColumns(4,
			v.digit(7),
			v.digit(8),
			v.digit(9),
			v.char("*")),
		container.NewGridWithColumns(4,
			v.digit(4),
			v.digit(5),
			v.digit(6),
			v.char("-")),
		container.NewGridWithColumns(4,
			v.digit(1),
			v.digit(2),
			v.digit(3),
			v.char("+")),
		container.NewGridWithColumns(4,
			v.button("C", v.clear),
			v.digit(0),
			v.char("."),
			v.char("/"),
		),
		container.NewGridWithColumns(4,
			v.char("round1"),
			v.char("round2"),
			v.char("round3"),
			equals,
		),
		widget.NewLabelWithStyle("Ларин Егор Сергеевич, 4 курс, 4 группа, 2023", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true}),
	)
}

func (v *View) digit(digit int) *widget.Button {
	str := strconv.Itoa(digit)
	return v.button(str, func() {
		v.update(str)
	})
}
func (v *View) char(ch string) *widget.Button {
	button := v.button(ch, func() {
		v.update(ch)
	})
	return button
}

func (v *View) button(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	return button
}

func (v *View) update(text string) {
	v.result.SetText("")
	pos := v.input.CursorColumn
	current := v.input.Text
	toUpdate := current[:pos] + text + current[pos:]
	log.Println(toUpdate)
	v.input.SetText(toUpdate)
	v.input.CursorColumn++
}

func (v *View) clear() {
	v.stack = v.stack[:0]
	v.input.SetText("")
}

func (v *View) showError(text string) {
	errDialog := widget.NewLabel(text)

	errorWindow := v.app.NewWindow("Ошибка")
	errorWindow.SetContent(container.NewVBox(
		errDialog,
		widget.NewButton("Окей", func() {
			errorWindow.Close()
		}),
	))
	errorWindow.Show()
}

func (v *View) isOp(str string) bool {
	return slices.Contains([]string{"+", "-", "/", "*"}, str)
}
