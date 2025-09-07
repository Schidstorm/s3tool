package terminal

import "github.com/rivo/tview"

func ConfirmModal(message string, onConfirm func()) ModalBuilder {
	return func(close func()) tview.Primitive {
		modal := NewModal().
			SetText(message).
			AddButtons([]string{"Cancel", "Confirm"}).
			SetDoneFunc(func(buttonLabel string, values map[string]string) {
				close()
				if buttonLabel == "Confirm" {
					onConfirm()
				}
			})
		modal.SetTitleAlign(tview.AlignLeft)
		return modal
	}
}
