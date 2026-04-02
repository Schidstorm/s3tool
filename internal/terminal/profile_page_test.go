package terminal

import (
	"context"
	"errors"
	"testing"

	"github.com/schidstorm/s3tool/internal/s3lib"
)

type profileTestLoader struct {
	connectors []s3lib.Connector
	err        error
}

func (l profileTestLoader) Load() ([]s3lib.Connector, error) {
	if l.err != nil {
		return nil, l.err
	}
	return l.connectors, nil
}

type profileTestConnector struct {
	name     string
	typeName string
}

func (c profileTestConnector) Name() string { return c.name }
func (c profileTestConnector) Type() string { return c.typeName }
func (c profileTestConnector) CreateClient(context.Context) (s3lib.Client, error) {
	return s3lib.NewMemoryClient(), nil
}

func TestLoadConnectorsAggregates(t *testing.T) {
	loaders := []s3lib.ConnectorLoader{
		profileTestLoader{connectors: []s3lib.Connector{profileTestConnector{name: "a", typeName: "aws"}}},
		profileTestLoader{connectors: []s3lib.Connector{profileTestConnector{name: "b", typeName: "s3tool"}}},
	}

	connectors, err := loadConnectors(loaders)
	if err != nil {
		t.Fatalf("loadConnectors failed: %v", err)
	}
	if len(connectors) != 2 {
		t.Fatalf("expected 2 connectors, got %d", len(connectors))
	}
}

func TestLoadConnectorsSortsByTypeThenName(t *testing.T) {
	loaders := []s3lib.ConnectorLoader{
		profileTestLoader{connectors: []s3lib.Connector{
			profileTestConnector{name: "zeta", typeName: "s3tool"},
			profileTestConnector{name: "beta", typeName: "aws"},
		}},
		profileTestLoader{connectors: []s3lib.Connector{
			profileTestConnector{name: "alpha", typeName: "aws"},
			profileTestConnector{name: "eta", typeName: "s3tool"},
		}},
	}

	connectors, err := loadConnectors(loaders)
	if err != nil {
		t.Fatalf("loadConnectors failed: %v", err)
	}

	got := make([]string, len(connectors))
	for i, c := range connectors {
		got[i] = c.Type() + ":" + c.Name()
	}

	expected := []string{"aws:alpha", "aws:beta", "s3tool:eta", "s3tool:zeta"}
	for i := range expected {
		if got[i] != expected[i] {
			t.Fatalf("unexpected order: got %#v, expected %#v", got, expected)
		}
	}
}

func TestLoadConnectorsReturnsError(t *testing.T) {
	expected := errors.New("boom")
	loaders := []s3lib.ConnectorLoader{
		profileTestLoader{connectors: []s3lib.Connector{profileTestConnector{name: "a", typeName: "aws"}}},
		profileTestLoader{err: expected},
	}

	_, err := loadConnectors(loaders)
	if err == nil {
		t.Fatal("expected loader error, got nil")
	}
}

func TestNewProfilePageBasics(t *testing.T) {
	ctx := NewContext()
	page := NewProfilePage(ctx, nil)

	if page.Title() != "Profiles" {
		t.Fatalf("expected title Profiles, got %q", page.Title())
	}
	if page.Context() == nil {
		t.Fatal("expected non-nil context")
	}
	if len(page.Hotkeys()) != 0 {
		t.Fatalf("expected no hotkeys, got %d", len(page.Hotkeys()))
	}
	if page.multiSelect {
		t.Fatal("expected profile page to disable multiselect")
	}

	columns := page.table.Columns()
	if len(columns) != 2 || columns[0] != "Type" || columns[1] != "Name" {
		t.Fatalf("expected columns [Type Name], got %#v", columns)
	}
}

func TestProfilePageLoadReplacesRows(t *testing.T) {
	ctx := NewContext()
	loaders := []s3lib.ConnectorLoader{
		profileTestLoader{connectors: []s3lib.Connector{
			profileTestConnector{name: "dev", typeName: "aws"},
			profileTestConnector{name: "local", typeName: "s3tool"},
		}},
	}

	page := NewProfilePage(ctx, loaders)
	page.Add(profileTestConnector{name: "stale", typeName: "memory"})

	if err := page.Load(); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	rows := page.table.Rows()
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows after load, got %d", len(rows))
	}

	if rows[0][0] != "aws" || rows[0][1] != "dev" {
		t.Fatalf("unexpected first row %#v", rows[0])
	}
	if rows[1][0] != "s3tool" || rows[1][1] != "local" {
		t.Fatalf("unexpected second row %#v", rows[1])
	}
}
