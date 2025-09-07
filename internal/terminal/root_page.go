package terminal

import (
	"errors"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RootPage struct {
	*tview.Flex
	pages       *tview.Pages
	profileInfo *ProfileInfoBox
	hotkeyInfo  *HotkeyInfoBox

	pageStask      []*Page
	openModalNames []string
}

func NewRootPage() *RootPage {
	header := tview.NewFlex()
	header.SetDirection(tview.FlexColumn)

	profileInfo := NewProfileInfoBox()
	header.AddItem(profileInfo, 0, 1, false)

	hotkeyInfo := NewHotkeyInfoBox()
	hotkeyInfo.Update(nil)
	header.AddItem(hotkeyInfo, 40, 0, false)

	content := tview.NewPages()
	content.SetBorder(true)
	content.SetBorderStyle(DefaultTheme.PageBorder)
	content.SetBorderPadding(0, 0, 1, 1)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(header, 5, 1, false)
	flex.AddItem(content, 0, 1, true)

	a := &RootPage{
		profileInfo: profileInfo,
		hotkeyInfo:  hotkeyInfo,
		pages:       content,
		Flex:        flex,
	}

	flex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			if len(a.openModalNames) > 0 {
				a.closeModal(a.openModalNames[len(a.openModalNames)-1])
				return nil
			}
		}
		return event
	})

	return a
}

func (a *RootPage) Modal(p ModalBuilder) {
	name := "modal_" + strconv.FormatInt(time.Now().UnixNano(), 16)
	a.openModalNames = append(a.openModalNames, name)

	content := p(func() {
		a.closeModal(name)
	})

	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(content, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 1, true).
		AddItem(nil, 0, 1, false)

	a.pages.AddPage(name, modal, true, true)
}

func (a *RootPage) closeModal(name string) {
	var topModalIndex int = -1
	for i := len(a.openModalNames) - 1; i >= 0; i-- {
		if a.openModalNames[i] == name {
			topModalIndex = i
			break
		}
	}
	if topModalIndex == -1 {
		return
	}

	if a.pages.HasPage(name) {
		a.pages.RemovePage(name)
		a.openModalNames = append(a.openModalNames[:topModalIndex], a.openModalNames[topModalIndex+1:]...)
	}
}

func (a *RootPage) OpenPage(pageContent PageContent) {
	if pageContent == nil {
		return
	}

	page := NewPage(pageContent)
	page.SetCloseHandler(func() {
		a.closePage()
	})
	a.pageStask = append(a.pageStask, page)
	a.openPage(page)
}

func (a *RootPage) openPage(page *Page) {
	a.pages.AddPage(page.Title(), page, true, true)
	a.pages.SetTitle(" " + page.Title() + " ")
	a.pages.SetTitleColor(DefaultTheme.PageTitlePrimary)
	a.pages.SwitchToPage(page.Title())

	err := page.content.Load()
	if err != nil {
		a.SetError(err)
	}

	a.hotkeyInfo.Update(page.content)
	a.profileInfo.UpdateContext(page.content.Context())
}

func (a *RootPage) closePage() {
	if len(a.pageStask) == 1 {
		return
	}

	a.pageStask = a.pageStask[:len(a.pageStask)-1]

	a.openPage(a.pageStask[len(a.pageStask)-1])
}

func (a *RootPage) SetError(err error) {
	title, message := errorText(err)

	a.Modal(func(close func()) tview.Primitive {
		return NewModal().
			SetTextStyle(DefaultTheme.ErrorMessage).
			SetText(message).
			AddButtons([]string{"OK"}).
			SetDoneFunc(func(buttonLabel string, formValues map[string]string) {
				close()
			}).SetTitle(" Error: " + title)
	})
}

func (a *RootPage) UpdateContext(c Context) {
	a.profileInfo.UpdateContext(c)
}

func errorText(err error) (title, message string) {
	if err == nil {
		return "", ""
	}

	if smithyErr, ok := findError[*smithy.OperationError](err); ok {
		if err, ok := findError[*http.ResponseError](err); ok {
			return "Http Error", err.Err.Error()
		}
		if err, ok := findError[*smithy.GenericAPIError](err); ok {
			return err.Code, err.Message
		}
		if err, ok := findError[*types.BucketAlreadyExists](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.BucketAlreadyOwnedByYou](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.EncryptionTypeMismatch](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.IdempotencyParameterMismatch](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.InvalidObjectState](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.InvalidRequest](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.InvalidWriteOffset](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.NoSuchBucket](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.NoSuchKey](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.NoSuchUpload](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.NotFound](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.ObjectAlreadyInActiveTierError](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.ObjectNotInActiveTierError](err); ok {
			return smithyErr.Operation(), err.Error()
		}
		if err, ok := findError[*types.TooManyParts](err); ok {
			return smithyErr.Operation(), err.Error()
		}

		return smithyErr.Operation(), smithyErr.Error()
	}

	return "Unknown Error", err.Error()
}

func findError[T any](e error) (T, bool) {
	var errorList []error
	currentErr := e
	for {
		if currentErr == nil {
			break
		}
		errorList = append(errorList, currentErr)
		currentErr = errors.Unwrap(currentErr)
	}

	for _, err := range errorList {
		if e, ok := err.(T); ok {
			return e, true
		}
	}

	var zero T
	return zero, false
}
