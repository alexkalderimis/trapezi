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
	Printer CellPrinter
	Padding int
	Value   string
	Align   Alignment
}

type CompoundCell struct {
	Cells   []Cell
	Divider string
}

const (
	ALIGN_LEFT  Alignment = iota
	ALIGN_RIGHT Alignment = iota
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
	logger.Infof("render: Value = %s, Padding = %d", cell.Value, cell.Padding)
	if cell.Align == ALIGN_RIGHT {
		return cell.Printer.Fprintf(w, "%*s%s", -cell.Padding, "", cell.Value)
	} else {
		return cell.Printer.Fprintf(w, "%s%*s", cell.Value, -cell.Padding, "")
	}
}

func (cell CompoundCell) Render(w io.Writer) (int, error) {
	var needsDivider bool = false
	var written int = 0
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

func AlignTable(table []CompoundCell) {
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
	switch column[0].(type) {
	case *SimpleCell:
		alignSimpleCells(column)
	case CompoundCell:
		alignCompoundCells(column)
	}
}

func alignSimpleCells(column []Cell) {
	maxWidth := 0
	for _, cell := range column {
		w := cell.length()
		if w > maxWidth {
			maxWidth = w
		}
	}
	for _, cell := range column {
		cell.align(maxWidth)
	}
}

func alignCompoundCells(column []Cell) {
	subtable := make([]CompoundCell, len(column), len(column))
	for i, cell := range column {
		cc := cell.(CompoundCell)
		subtable[i] = cc
	}
	AlignTable(subtable)
}
