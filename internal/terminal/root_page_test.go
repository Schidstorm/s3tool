package terminal

import (
	"errors"
	"net/http"
	"testing"

	awserrhttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"github.com/rivo/tview"
)

func TestRootPageOpenAndClosePageStack(t *testing.T) {
	root := NewRootPage()
	ctx := NewContext()

	root.OpenPage(nil)
	if len(root.pageStask) != 0 {
		t.Fatalf("expected empty stack for nil page, got %d", len(root.pageStask))
	}

	p1 := newPageTestContent("P1", ctx)
	p2 := newPageTestContent("P2", ctx)
	root.OpenPage(p1)
	root.OpenPage(p2)

	if len(root.pageStask) != 2 {
		t.Fatalf("expected stack size 2, got %d", len(root.pageStask))
	}

	root.closePage()
	if len(root.pageStask) != 1 {
		t.Fatalf("expected stack size 1 after close, got %d", len(root.pageStask))
	}

	root.closePage()
	if len(root.pageStask) != 1 {
		t.Fatalf("expected stack to remain at 1, got %d", len(root.pageStask))
	}
}

func TestRootPageModalOpenAndClose(t *testing.T) {
	root := NewRootPage()

	var closeFn func()
	root.Modal(func(close func()) tview.Primitive {
		closeFn = close
		return tview.NewBox()
	})

	if len(root.openModalNames) != 1 {
		t.Fatalf("expected one modal open, got %d", len(root.openModalNames))
	}
	name := root.openModalNames[0]
	if !root.pages.HasPage(name) {
		t.Fatalf("expected page %s to exist", name)
	}

	closeFn()
	if len(root.openModalNames) != 0 {
		t.Fatalf("expected no open modals after close, got %d", len(root.openModalNames))
	}
	if root.pages.HasPage(name) {
		t.Fatalf("expected page %s removed", name)
	}
}

func TestRootPageCloseModalUnknownNoop(t *testing.T) {
	root := NewRootPage()
	root.closeModal("does-not-exist")
}

func TestRootPageOpenPageLoadErrorShowsModal(t *testing.T) {
	root := NewRootPage()
	ctx := NewContext()

	p := newPageTestContent("ErrPage", ctx)
	p.loadErr = errors.New("load failed")
	root.OpenPage(p)

	if len(root.openModalNames) == 0 {
		t.Fatal("expected error modal to be opened")
	}
}

func TestErrorTextGenericAPIError(t *testing.T) {
	err := &smithy.OperationError{
		OperationName: "ListBuckets",
		Err: &smithy.GenericAPIError{
			Code:    "NoSuchResource",
			Message: "missing",
		},
	}

	title, message := errorText(err)
	if title != "NoSuchResource" || message != "missing" {
		t.Fatalf("unexpected title/message: %q / %q", title, message)
	}
}

func TestErrorTextHTTPResponseError(t *testing.T) {
	httpErr := &awserrhttp.ResponseError{
		ResponseError: &smithyhttp.ResponseError{
			Response: &smithyhttp.Response{Response: &http.Response{StatusCode: 500}},
			Err:      errors.New("boom"),
		},
	}

	err := &smithy.OperationError{OperationName: "GetObject", Err: httpErr}
	title, message := errorText(err)
	if title != "Http Error" || message != "boom" {
		t.Fatalf("unexpected title/message: %q / %q", title, message)
	}
}

func TestErrorTextFallbacks(t *testing.T) {
	title, message := errorText(errors.New("plain"))
	if title != "Unknown Error" || message != "plain" {
		t.Fatalf("unexpected fallback title/message: %q / %q", title, message)
	}

	opErr := &smithy.OperationError{OperationName: "DeleteObject", Err: errors.New("inner")}
	title, message = errorText(opErr)
	if title != "DeleteObject" || message == "" {
		t.Fatalf("unexpected operation fallback: %q / %q", title, message)
	}
}

func TestFindError(t *testing.T) {
	inner := &smithy.GenericAPIError{Code: "C", Message: "M"}
	wrapped := &smithy.OperationError{OperationName: "Op", Err: inner}

	found, ok := findError[*smithy.GenericAPIError](wrapped)
	if !ok || found != inner {
		t.Fatal("expected to find generic api error in chain")
	}

	_, ok = findError[*awserrhttp.ResponseError](wrapped)
	if ok {
		t.Fatal("did not expect to find aws response error")
	}
}
