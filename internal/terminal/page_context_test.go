package terminal

import (
	"errors"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/internal/s3lib"
)

type pageTestContent struct {
	*tview.Box
	title      string
	searchTerm string
	ctx        Context
	hotkeys    map[tcell.EventKey]Hotkey
	loadErr    error
}

func newPageTestContent(title string, ctx Context) *pageTestContent {
	return &pageTestContent{
		Box:     tview.NewBox(),
		title:   title,
		ctx:     ctx,
		hotkeys: map[tcell.EventKey]Hotkey{},
	}
}

func (p *pageTestContent) Title() string { return p.title }

func (p *pageTestContent) Hotkeys() map[tcell.EventKey]Hotkey { return p.hotkeys }

func (p *pageTestContent) SetSearch(term string) { p.searchTerm = term }

func (p *pageTestContent) Context() Context { return p.ctx }

func (p *pageTestContent) Load() error { return p.loadErr }

func TestPageTitleAndCloseHandler(t *testing.T) {
	content := newPageTestContent("Sample", NewContext())
	page := NewPage(content)

	if page.Title() != "Sample" {
		t.Fatalf("expected title Sample, got %q", page.Title())
	}

	closed := false
	page.SetCloseHandler(func() { closed = true })
	page.handleClose()
	if !closed {
		t.Fatal("expected close handler to be called")
	}
}

func TestPageSearchActivationAndDeactivation(t *testing.T) {
	content := newPageTestContent("Sample", NewContext())
	page := NewPage(content)

	if page.isSearchActive {
		t.Fatal("expected search inactive initially")
	}

	page.activateSearch()
	if !page.isSearchActive {
		t.Fatal("expected search active after activation")
	}

	input, ok := page.searchFlex.GetItem(0).(*tview.InputField)
	if !ok {
		t.Fatalf("expected first item to be input field, got %T", page.searchFlex.GetItem(0))
	}
	input.SetText("needle")
	if content.searchTerm != "needle" {
		t.Fatalf("expected search term needle, got %q", content.searchTerm)
	}

	page.deactivateSearch()
	if page.isSearchActive {
		t.Fatal("expected search inactive after deactivation")
	}
	if _, ok := page.searchFlex.GetItem(0).(*tview.InputField); ok {
		t.Fatal("expected input field to be removed after deactivation")
	}
}

func TestPageEventKey(t *testing.T) {
	ek := EventKey(tcell.KeyRune, 'x', tcell.ModCtrl)
	if ek.Rune() != 'x' || ek.Modifiers() != tcell.ModCtrl {
		t.Fatalf("unexpected event key: key=%v rune=%v mod=%v", ek.Key(), ek.Rune(), ek.Modifiers())
	}
}

func TestContextCallbacksAndDefaults(t *testing.T) {
	ctx := NewContext()

	// Defaults should be no-op and no panic.
	ctx.Modal(nil)
	ctx.SetError(errors.New("noop"))
	ctx.OpenPage(newPageTestContent("noop", ctx))
	if ctx.SuspendApp(func() {}) {
		t.Fatal("expected default SuspendApp to return false")
	}

	var modalCalled bool
	var gotErr error
	var openCalled bool
	var suspendCalled bool

	ctx = ctx.
		WithModalFunc(func(build ModalBuilder) { modalCalled = true }).
		WithErrorFunc(func(err error) { gotErr = err }).
		WithOpenPageFunc(func(page PageContent) { openCalled = page != nil }).
		WithSuspendAppFunc(func(f func()) bool {
			suspendCalled = true
			f()
			return true
		})

	ctx.Modal(nil)
	err := errors.New("boom")
	ctx.SetError(err)
	ctx.OpenPage(newPageTestContent("open", ctx))
	if !ctx.SuspendApp(func() {}) {
		t.Fatal("expected configured SuspendApp to return true")
	}

	if !modalCalled {
		t.Fatal("expected modal callback to be called")
	}
	if gotErr != err {
		t.Fatalf("expected propagated error %v, got %v", err, gotErr)
	}
	if !openCalled {
		t.Fatal("expected open page callback to be called")
	}
	if !suspendCalled {
		t.Fatal("expected suspend callback to be called")
	}
}

func TestContextWithValues(t *testing.T) {
	client := s3lib.NewMemoryClient()
	ctx := NewContext().
		WithClient(client).
		WithBucket("bucket").
		WithObjectKey("key.txt")

	if ctx.S3Client() != client {
		t.Fatal("expected client to be set")
	}
	if ctx.Bucket() != "bucket" {
		t.Fatalf("expected bucket to be bucket, got %q", ctx.Bucket())
	}
	if ctx.ObjectKey() != "key.txt" {
		t.Fatalf("expected object key key.txt, got %q", ctx.ObjectKey())
	}
}
