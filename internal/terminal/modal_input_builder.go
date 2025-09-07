package terminal

import "github.com/rivo/tview"

type FormInputBuilder struct {
	*Modal
	inputField *tview.InputField
}

func NewInputBuilder(modal *Modal) *FormInputBuilder {
	inputField := tview.NewInputField()
	modal.form.AddFormItem(inputField)

	return &FormInputBuilder{
		Modal:      modal,
		inputField: inputField,
	}
}

func (b *FormInputBuilder) SetFieldWidth(width int) *FormInputBuilder {
	b.inputField.SetFieldWidth(width)
	return b
}

func (b *FormInputBuilder) SetChangedFunc(f func(text string)) *FormInputBuilder {
	b.inputField.SetChangedFunc(f)
	return b
}

func (b *FormInputBuilder) SetAcceptanceFunc(f func(textToCheck string, lastChar rune) bool) *FormInputBuilder {
	b.inputField.SetAcceptanceFunc(f)
	return b
}

func (b *FormInputBuilder) SetText(value string) *FormInputBuilder {
	b.inputField.SetText(value)
	return b
}

func (b *FormInputBuilder) SetLabel(label string) *FormInputBuilder {
	b.inputField.SetLabel(label)
	return b
}
