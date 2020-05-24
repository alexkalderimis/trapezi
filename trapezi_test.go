package trapezi

import (
	"strings"
	"testing"
)

func TestRenderSimple(t *testing.T) {
	want := "Hello"
	cell := Text("Hello")

	assertRenderEql(cell, want, t)
}

func TestRenderJoin(t *testing.T) {
	want := "Hello World"
	cell := Join(" ", Text("Hello"), Text("World"))

	assertRenderEql(cell, want, t)
}

func TestAlignColumn(t *testing.T) {
	want := strings.Join([]string{
		"Hello    ",
		"to       ",
		"léft     ",
		"alignment",
		"with     ",
		"unicode  ",
	}, "|")
	cells := Join("|",
		Text("Hello"),
		Text("to"),
		Text("léft"),
		Text("alignment"),
		Text("with"),
		Text("unicode"),
	)
	alignColumn(cells.Cells)
	assertRenderEql(cells, want, t)
}

func assertRenderEql(cell Cell, want string, t *testing.T) {
	var buff strings.Builder
	cell.Render(&buff)

	if got := buff.String(); got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}
}
