package s3lib

import (
	"context"
	"errors"
	"testing"
)

func TestMemoryPaginatorError(t *testing.T) {
	p := &memoryPaginator[int]{err: errors.New("fail")}

	if p.HasMorePages() {
		t.Fatal("expected HasMorePages false when paginator has error")
	}
	if _, err := p.NextPage(context.Background()); err == nil {
		t.Fatal("expected NextPage to return error")
	}
}

func TestMemoryPaginatorReadOnce(t *testing.T) {
	p := &memoryPaginator[string]{items: []string{"a", "b"}}

	if !p.HasMorePages() {
		t.Fatal("expected HasMorePages true before first read")
	}

	items, err := p.NextPage(context.Background())
	if err != nil {
		t.Fatalf("first NextPage failed: %v", err)
	}
	if len(items) != 2 || items[0] != "a" || items[1] != "b" {
		t.Fatalf("unexpected items: %#v", items)
	}

	if p.HasMorePages() {
		t.Fatal("expected HasMorePages false after first read")
	}

	items, err = p.NextPage(context.Background())
	if err != nil {
		t.Fatalf("second NextPage failed: %v", err)
	}
	if items != nil {
		t.Fatalf("expected nil items on second read, got %#v", items)
	}
}

func TestMemoryPaginatorEmptyItems(t *testing.T) {
	p := &memoryPaginator[int]{items: []int{}}
	if p.HasMorePages() {
		t.Fatal("expected no pages for empty item list")
	}
}
