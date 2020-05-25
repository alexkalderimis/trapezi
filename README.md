Trapezi - Easy automatically aligned text tables
================================================

> τα κλειδιά είναι πάνω στο τραπέζι

This library is for creating text based tables that have automatic alignment
and support coloured printing.

# Usage

```go
import "io"
import "fmt"
import "gitlab.com/alexkalderimis/trapezi"

func usageExample() {
  table := []CompoundCell{
    trapezi.Join(" | ", trapezi.Text("Name"), trapezi.Text("Value")),
    trapezi.Join(" | ", trapezi.Text("foo"), trapezi.TextR("5")),
    trapezi.Join(" | ", trapezi.Text("bar"), trapezi.TextR("17")),
  }
  trapezi.AlignTable(table)

  for _, row := range table {
    row.Render(io.Stdout)
    fmt.Println("")
  }
}
```

The core abstraction of this package is the `Cell`, with is either
a `SimpleCell` containing a single text value, or a `CompoundCell`, containing
other cells separated by a divider. A table row is just a compound cell, but any
cell can contain other nested compound cells. Cells can be either left or right aligned.

Once the entire table is created, it can be aligned, and once aligned, rendered to a `io.Writer`.


