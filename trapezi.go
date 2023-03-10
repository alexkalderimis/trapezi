package trapezi

import (
	"fmt"
	"github.com/juju/loggo"
	"io"
	"unicode/utf8"
)

type CellPrinter interface {
	Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error)
}

type Cell interface {
	Render(w io.Writer) (int, error)
	length() int
	align(int)
}

type Alignment int

type CustomPrinter struct {
	Function func(io.Writer, string, ...interface{}) (int, error)
}

func (p *CustomPrinter) Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	return p.Function(w, format, a...)
}

type SimpleCell struct {
	Printer     CellPrinter
	Padding     int
	Value       string
	Align       Alignment
	PaddingRune rune
}

type CompoundCell struct {
	Cells   []Cell
	Divider string
	Prefix  string
	Suffix  string
}

const (
	ALIGN_LEFT  Alignment = iota
	ALIGN_RIGHT Alignment = iota
	EXPAND      Alignment = iota // TODO make this work
)

var (
	logger     = loggo.GetLogger("trapezi")
	FmtPrinter = CustomPrinter{
		Function: fmt.Fprintf,
	}
	defaultPrinter CellPrinter = &FmtPrinter
)

func SetDefaultPrinter(printer CellPrinter) {
	defaultPrinter = printer
}

func (cell SimpleCell) length() int {
	return utf8.RuneCountInString(cell.Value)
}

func (cell *SimpleCell) align(width int) {
	logger.Debugf("align: Value = %s, Padding = %d", cell.Value, width)
	cell.Padding = width - cell.length()
}

func (cell SimpleCell) Render(w io.Writer) (int, error) {

	if cell.Padding == 0 {
		return cell.Printer.Fprintf(w, "%s", cell.Value)
	}
	r := cell.PaddingRune
	if r == 0 {
		r = ' '
	}
	var chars []rune
	for i := 0; i < cell.Padding; i++ {
		chars = append(chars, r)
	}
	paddingValue := string(chars)
	logger.Infof("render: Value = %s, Padding = %d, Padded With = %s", cell.Value, cell.Padding, paddingValue)

	if cell.Align == ALIGN_RIGHT {
		return cell.Printer.Fprintf(w, "%s%s", paddingValue, cell.Value)
	} else {
		return cell.Printer.Fprintf(w, "%s%s", cell.Value, paddingValue)
	}
}

func (cell CompoundCell) Render(w io.Writer) (int, error) {
	var needsDivider bool = false
	var written int = 0
	if cell.Prefix != "" {
		fmt.Fprintf(w, cell.Prefix)
	}
	for _, child := range cell.Cells {
		if needsDivider {
			wrote, err := fmt.Fprint(w, cell.Divider)
			written += wrote
			if err != nil {
				return written, err
			}
		}
		wrote, err := child.Render(w)
		written += wrote
		if err != nil {
			return written, err
		}
		needsDivider = true
	}
	if cell.Suffix != "" {
		fmt.Fprintf(w, cell.Suffix)
	}
	return written, nil
}

func (cell CompoundCell) align(width int) {
	if width > 0 {
		n := len(cell.Cells)
		each := width / n
		leftOver := width % n

		for i, child := range cell.Cells {
			padding := each
			if i+1 >= n {
				padding += leftOver
			}
			logger.Infof("CompoundCell.align: Padding cell to %d", padding)
			child.align(padding)
		}
	} else {
		logger.Infof("CompoundCell.align: No padding to be done")
	}
}

func (cell CompoundCell) length() int {
	var sum int = 0
	for _, child := range cell.Cells {
		sum += child.length()
	}
	if len(cell.Cells) > 1 {
		sum += len(cell.Divider) * (len(cell.Cells) - 1)
	}
	return sum
}

func Join(divider string, cells ...Cell) CompoundCell {
	cell := CompoundCell{
		Divider: divider,
		Cells:   cells,
	}
	return cell
}

func WithPrinter(printer CellPrinter, msg string) *SimpleCell {
	cell := SimpleCell{
		Printer: printer,
		Value:   msg,
		Align:   ALIGN_LEFT,
	}
	return &cell
}

func Text(value string) *SimpleCell {
	cell := SimpleCell{
		Printer: defaultPrinter,
		Value:   value,
		Align:   ALIGN_LEFT,
	}
	return &cell
}

func Empty() *SimpleCell {
	return Text("")
}

func AlignRight(cell *SimpleCell) *SimpleCell {
	cell.Align = ALIGN_RIGHT
	return cell
}

func PaddedWith(padding rune, cell *SimpleCell) *SimpleCell {
	cell.PaddingRune = padding
	return cell
}

func TextR(value string) *SimpleCell {
	return AlignRight(Text(value))
}

func AlignTable(table []CompoundCell) {
	var maxWidth = 0
	for _, row := range table {
		l := row.length()
		if l > maxWidth {
			maxWidth = l
		}
	}
	alignTableTo(table, maxWidth)
}

func alignTableTo(table []CompoundCell, maxWidth int) {
	if len(table) <= 1 {
		logger.Debugf("Table already aligned: rows: %d", len(table))
		return
	}
	width := len(table[0].Cells)
	logger.Debugf("Aligning table: rows: %d, width: %d", len(table), width)
	for i := 0; i < width; i++ {
		column := make([]Cell, len(table), len(table))
		for rowIdx, row := range table {
			column[rowIdx] = row.Cells[i]
		}
		logger.Debugf("Aligning column: %d", i)
		alignColumn(column)
	}
}

func alignColumn(column []Cell) {
	if len(column) <= 1 {
		return
	}
	maxWidth := 0
	for _, cell := range column {
		w := cell.length()
		if w > maxWidth {
			maxWidth = w
		}
	}
	simpleCells := make([]Cell, 0, len(column))
	compoundCells := make([]Cell, 0, len(column))
	for _, cell := range column {
		switch cell.(type) {
		case *SimpleCell:
			simpleCells = append(simpleCells, cell)
		case CompoundCell:
			compoundCells = append(compoundCells, cell)
		}
	}
	alignSimpleCells(simpleCells, maxWidth)
	alignCompoundCells(compoundCells, maxWidth)
}

func alignSimpleCells(column []Cell, maxWidth int) {
	for _, cell := range column {
		cell.align(maxWidth)
	}
}

func Table(lineBreak string, rows []CompoundCell) CompoundCell {
	cells := make([]Cell, len(rows), len(rows))
	for i, row := range rows {
		cells[i] = row
	}
	table := CompoundCell{
		Divider: lineBreak,
		Cells:   cells,
	}
	return table
}

func alignCompoundCells(column []Cell, maxWidth int) {
	subtable := make([]CompoundCell, len(column), len(column))
	for i, cell := range column {
		cc := cell.(CompoundCell)
		subtable[i] = cc
	}
	alignTableTo(subtable, maxWidth)
}
