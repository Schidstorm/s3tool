package terminal

import "github.com/rivo/tview"

type FormBuilder struct {
	inputfields []*tview.InputField
	buttons     []*tview.Button
}

func NewFormBuilder() *FormBuilder {
	return &FormBuilder{
		buttons: []*tview.Button{},
	}
}

func (b *FormBuilder) Button(label string, selected func(values map[string]string)) *FormBuilder {
	btn := tview.NewButton(label).SetSelectedFunc(func() {
		selected(b.getInputValues())
	})
	b.buttons = append(b.buttons, btn)
	return b
}

func (b *FormBuilder) getInputValues() map[string]string {
	values := make(map[string]string)
	for _, inputField := range b.inputfields {
		values[inputField.GetLabel()] = inputField.GetText()
	}
	return values
}

func (b *FormBuilder) Build() tview.Primitive {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	form := tview.NewForm()
	for _, inputField := range b.inputfields {
		form.AddFormItem(inputField)
	}
	flex.AddItem(form, 0, 1, true)
	buttonsFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	buttonsFlex.AddItem(nil, 0, 1, false)
	for i, btn := range b.buttons {
		if i > 0 {
			buttonsFlex.AddItem(nil, 1, 0, false)
		}
		buttonsFlex.AddItem(btn, 10, 0, false)
	}
	flex.AddItem(buttonsFlex, 1, 0, false)

	return flex
}
