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

func TestPrefixesAndSuffixes(t *testing.T) {
	want := strings.Join([]string{
		">>> " + "Hello    ",
		"to       ",
		"léft     ",
		"alignment",
		"with     ",
		"unicode  " + " <<<",
	}, "|")
	cells := Join("|",
		Text("Hello"),
		Text("to"),
		Text("léft"),
		Text("alignment"),
		Text("with"),
		Text("unicode"),
	)
	cells.Prefix = ">>> "
	cells.Suffix = " <<<"
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

func TestPaddingChars(t *testing.T) {
	rowSep := " | "
	want := strings.Join([]string{
		"Name    | Family   | Example",
		"--------|----------|--------",
		"German  | Germanic | Tisch  ",
		"Greek   | Hellenic | τραπέζι",
		"Italian | Romance  | tavola ",
	}, "/")
	cells := []CompoundCell{
		Join(rowSep, Text("Name"), Text("Family"), Text("Example")),
		Join("-|-", PaddedWith('-', Empty()), PaddedWith('-', Empty()), PaddedWith('-', Empty())),
		Join(rowSep, Text("German"), Text("Germanic"), Text("Tisch")),
		Join(rowSep, Text("Greek"), Text("Hellenic"), Text("τραπέζι")),
		Join(rowSep, Text("Italian"), Text("Romance"), Text("tavola")),
	}
	AlignTable(cells)
	table := Table("/", cells)
	assertRenderEql(table, want, t)
}

// TODO
// func TestAlignMixedColumn(t *testing.T) {
//   sepA := " || "
//   sepB := " | "
//   want := strings.Join([]string{
//     "Language || SING              || PL                ",
//     "         || masc | fem | neut || masc  | fem | neut",
//     "German   || der  | die | das  || die               ",
//     "Greek    || ο    | η   | το   || οι          | τα  ",
//     "Italian  || il   | la  | -    || i,gli | le  | -   ",
//   }, "/")
//   mfn := func(m, f, n string) Cell {
//     return Join(sepB, Text(m), Text(f), Text(n))
//   }
//   mf_n := func(mf, n string) Cell {
//     return Join(sepB, Text(mf), Text(n))
//   }
//   cells := []CompoundCell{
//     Join(sepA, Text("Language"), Text("Singular"), Text("Plural")),
//     Join(sepA, Empty(), mfn("masc", "fem", "neut"), mfn("masc", "fem", "neut")),
//     Join(sepA, Text("German"), mfn("der", "die", "das"), Text("die")),
//     Join(sepA, Text("Greek"), mfn("ο", "οι", "το"), mf_n("οι", "τα")),
//     Join(sepA, Text("Italian"), mfn("il", "la", "-"),
//       Join(sepA, Join(",", Text("i"), Text("gli")),
//         Text("le"), Text("-"))),
//   }
//   AlignTable(cells)
//   table := Table("/", cells)
//   assertRenderEql(table, want, t)
// }

func assertRenderEql(cell Cell, want string, t *testing.T) {
	var buff strings.Builder
	cell.Render(&buff)

	if got := buff.String(); got != want {
		t.Errorf("Expected %q, got %q", want, got)
	}
}
