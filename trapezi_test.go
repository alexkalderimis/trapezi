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

func TestAlignColumnR(t *testing.T) {
	want := strings.Join([]string{
		"Hello    ",
		"       to",
		"     léft",
		"alignment",
		"with     ",
		"unicode  ",
	}, "|")
	cells := Join("|",
		Text("Hello"),
		AlignRight(Text("to")),
		AlignRight(Text("léft")),
		Text("alignment"),
		Text("with"),
		Text("unicode"),
	)
	alignColumn(cells.Cells)
	assertRenderEql(cells, want, t)
}

func TestAlignTable(t *testing.T) {
	want := strings.Join([]string{
		"This   |  And",
		"column | this",
		"is     |  one",
		"left   |   is",
		"aligned|right",
	}, "/")
	cells := []CompoundCell{
		Join("|", Text("This"), TextR("And")),
		Join("|", Text("column"), TextR("this")),
		Join("|", Text("is"), TextR("one")),
		Join("|", Text("left"), TextR("is")),
		Join("|", Text("aligned"), TextR("right")),
	}
	AlignTable(cells)
	table := Table("/", cells)
	assertRenderEql(table, want, t)
}

func TestAlignNested(t *testing.T) {
	rowSep := " | "
	want := strings.Join([]string{
		"Language | Singular,Plural  ",
		"German   |    Tisch,Tische  ",
		"Greek    |  τραπέζι,τραπέζια",
		"Italian  |   tavola,tavole  ",
	}, "/")
	cells := []CompoundCell{
		Join(rowSep, Text("Language"), Join(",", TextR("Singular"), Text("Plural"))),
		Join(rowSep, Text("German"), Join(",", TextR("Tisch"), Text("Tische"))),
		Join(rowSep, Text("Greek"), Join(",", TextR("τραπέζι"), Text("τραπέζια"))),
		Join(rowSep, Text("Italian"), Join(",", TextR("tavola"), Text("tavole"))),
	}
	AlignTable(cells)
	table := Table("/", cells)
	assertRenderEql(table, want, t)
}

func assertRenderEql(cell Cell, want string, t *testing.T) {
	var buff strings.Builder
	cell.Render(&buff)

	if got := buff.String(); got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}
}
