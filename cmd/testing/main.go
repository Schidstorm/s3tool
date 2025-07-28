package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)

	table := tview.NewTable().SetSelectable(true, false)
	flex.AddItem(table, 0, 1, true)

	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		fmt.Println("table")
		return event
	})

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		fmt.Println("flex")
		return nil
	})

	tview.NewApplication().SetRoot(flex, true).Run()
}
