package terminal

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Modal struct {
	*tview.Box

	// The frame embedded in the modal.
	frame *tview.Frame

	// The form embedded in the modal's frame.
	form *tview.Form

	// The message text (original, not word-wrapped).
	text string

	// The text color.
	textStyle tcell.Style

	title string

	// The optional callback for when the user clicked one of the buttons. It
	// receives the index of the clicked button and the button's label.
	done func(buttonLabel string, formValues map[string]string)
}

// NewModal returns a new modal message window.
func NewModal() *Modal {
	m := &Modal{
		Box:       tview.NewBox().SetBorder(true).SetBackgroundColor(tview.Styles.ContrastBackgroundColor),
		textStyle: tcell.StyleDefault.Foreground(tview.Styles.PrimaryTextColor),
	}
	m.form = tview.NewForm().
		SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tview.Styles.PrimitiveBackgroundColor).
		SetButtonTextColor(tview.Styles.PrimaryTextColor)
	m.form.SetBackgroundColor(tview.Styles.ContrastBackgroundColor).SetBorderPadding(0, 0, 0, 0)
	m.form.SetCancelFunc(func() {
		if m.done != nil {
			m.done("", nil)
		}
	})
	m.frame = tview.NewFrame(m.form).SetBorders(0, 0, 1, 0, 0, 0)
	m.frame.SetBackgroundColor(tview.Styles.ContrastBackgroundColor).
		SetBorderPadding(1, 1, 1, 1)
	return m
}

// SetTextStyle sets the color of the message text.
func (m *Modal) SetTextStyle(style tcell.Style) *Modal {
	m.textStyle = style
	_, bg, _ := style.Decompose()
	m.form.SetBackgroundColor(bg)
	m.frame.SetBackgroundColor(bg)
	m.Box.SetBackgroundColor(bg)
	return m
}

// SetButtonBackgroundColor sets the background color of the buttons.
func (m *Modal) SetButtonBackgroundColor(color tcell.Color) *Modal {
	m.form.SetButtonBackgroundColor(color)
	return m
}

// SetButtonTextColor sets the color of the button texts.
func (m *Modal) SetButtonTextColor(color tcell.Color) *Modal {
	m.form.SetButtonTextColor(color)
	return m
}

// SetButtonStyle sets the style of the buttons when they are not focused.
func (m *Modal) SetButtonStyle(style tcell.Style) *Modal {
	m.form.SetButtonStyle(style)
	return m
}

// SetButtonActivatedStyle sets the style of the buttons when they are focused.
func (m *Modal) SetButtonActivatedStyle(style tcell.Style) *Modal {
	m.form.SetButtonActivatedStyle(style)
	return m
}

// SetDoneFunc sets a handler which is called when one of the buttons was
// pressed. It receives the index of the button as well as its label text. The
// handler is also called when the user presses the Escape key. The index will
// then be negative and the label text an empty string.
func (m *Modal) SetDoneFunc(handler func(buttonLabel string, formValues map[string]string)) *Modal {
	m.done = handler
	return m
}

// SetText sets the message text of the window. The text may contain line
// breaks but style tag states will not transfer to following lines. Note that
// words are wrapped, too, based on the final size of the window.
func (m *Modal) SetText(text string) *Modal {
	m.text = text
	return m
}

func (m *Modal) SetTitle(title string) *Modal {
	m.title = title
	m.Box.SetTitle(" " + title + " ")
	m.Box.SetTitleAlign(tview.AlignLeft)
	return m
}

// AddButtons adds buttons to the window. There must be at least one button and
// a "done" handler so the window can be closed again.
func (m *Modal) AddButtons(labels []string) *Modal {
	for index, label := range labels {
		func(i int, l string) {
			m.form.AddButton(label, func() {
				if m.done != nil {
					inputValues := m.getInputValues()
					m.done(l, inputValues)
				}
			})
			button := m.form.GetButton(m.form.GetButtonCount() - 1)
			button.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
				switch event.Key() {
				case tcell.KeyDown, tcell.KeyRight:
					return tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
				case tcell.KeyUp, tcell.KeyLeft:
					return tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
				}
				return event
			})
		}(index, label)
	}
	return m
}

func (m *Modal) getInputValues() map[string]string {
	values := make(map[string]string)
	for i := 0; i < m.form.GetFormItemCount(); i++ {
		inputField, ok := m.form.GetFormItem(i).(*tview.InputField)
		if !ok {
			continue
		}
		values[inputField.GetLabel()] = inputField.GetText()
	}
	return values
}

// ClearButtons removes all buttons from the window.
func (m *Modal) ClearButtons() *Modal {
	m.form.ClearButtons()
	return m
}

// SetFocus shifts the focus to the button with the given index.
func (m *Modal) SetFocus(index int) *Modal {
	m.form.SetFocus(index)
	return m
}

// Focus is called when this primitive receives focus.
func (m *Modal) Focus(delegate func(p tview.Primitive)) {
	delegate(m.form)
}

// HasFocus returns whether or not this primitive has focus.
func (m *Modal) HasFocus() bool {
	return m.form.HasFocus()
}

// Draw draws this primitive onto the screen.
func (m *Modal) Draw(screen tcell.Screen) {

	// Calculate the width of this modal.
	buttonsWidth := 0
	for i := 0; i < m.form.GetButtonCount(); i++ {
		button := m.form.GetButton(i)
		buttonsWidth += tview.TaggedStringWidth(button.GetLabel()) + 4 + 2
	}
	buttonsWidth -= 2
	screenWidth, screenHeight := screen.Size()
	width := max(screenWidth/3, buttonsWidth)
	// width is now without the box border.

	// Reset the text and find out how wide it is.
	m.frame.Clear()
	lines := tview.WordWrap(m.text, width)
	if len(lines) == 1 && lines[0] == "" {
		lines = nil
	}
	fg, _, _ := m.textStyle.Decompose()
	for _, line := range lines {
		m.frame.AddText(line, true, tview.AlignCenter, fg)
	}

	// Set the modal's position and size.
	padding := 4
	buttonSize := 0
	if m.form.GetButtonCount() > 0 {
		buttonSize = 1
	}
	if lines != nil {
		padding += 1 // extra padding if there is text
	}
	height := len(lines) + padding + buttonSize + m.form.GetFormItemCount()*2
	width += 4
	x := (screenWidth - width) / 2
	y := (screenHeight - height) / 2
	m.SetRect(x, y, width, height)

	// Set color of the form's input fields.
	for i := 0; i < m.form.GetFormItemCount(); i++ {
		inputField, ok := m.form.GetFormItem(i).(*tview.InputField)
		if !ok {
			continue
		}
		_, bg, _ := m.textStyle.Decompose()
		inputField.SetFieldBackgroundColor(bg)
		inputField.SetFieldTextColor(tview.Styles.PrimaryTextColor)
	}

	// Draw the frame.
	m.Box.DrawForSubclass(screen, m)
	x, y, width, height = m.GetInnerRect()
	m.frame.SetRect(x, y, width, height)
	m.frame.Draw(screen)
}

// MouseHandler returns the mouse handler for this primitive.
func (m *Modal) MouseHandler() func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
	return m.WrapMouseHandler(func(action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive)) (consumed bool, capture tview.Primitive) {
		// Pass mouse events on to the form.
		consumed, capture = m.form.MouseHandler()(action, event, setFocus)
		if !consumed && action == tview.MouseLeftDown && m.InRect(event.Position()) {
			setFocus(m)
			consumed = true
		}
		return
	})
}

// InputHandler returns the handler for this primitive.
func (m *Modal) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return m.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if m.frame.HasFocus() {
			if handler := m.frame.InputHandler(); handler != nil {
				handler(event, setFocus)
				return
			}
		}
	})
}

// GetForm returns the form embedded in the modal.
func (m *Modal) GetForm() *tview.Form {
	return m.form
}

func (m *Modal) AddInput() *FormInputBuilder {
	return NewInputBuilder(m)
}
