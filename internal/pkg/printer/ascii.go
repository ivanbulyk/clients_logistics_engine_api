package printer

import (
    "fmt"
    "strings"
)

// ASCIITablePrinter to STDOUT
type ASCIITablePrinter struct {
    headers []string
    rows    [][]string
}

// NewASCIITablePrinter simple table printer
func NewASCIITablePrinter() *ASCIITablePrinter {
    return &ASCIITablePrinter{
        headers: make([]string, 0),
        rows:    make([][]string, 0),
    }
}

// AddHeader to table
func (t *ASCIITablePrinter) AddHeader(headers []string) {
    t.headers = headers
}

// AddRow to table
func (t *ASCIITablePrinter) AddRow(values []string) {
    t.rows = append(t.rows, values)
}

func (t *ASCIITablePrinter) String() string {
    var builder strings.Builder

    // Calculate column widths
    columnWidths := make([]int, len(t.headers))
    for i, header := range t.headers {
        columnWidths[i] = len(header)
    }
    for _, row := range t.rows {
        for i, value := range row {
            if len(value) > columnWidths[i] {
                columnWidths[i] = len(value)
            }
        }
    }

    // Print headers
    for i, header := range t.headers {
        builder.WriteString(fmt.Sprintf("| %-*s ", columnWidths[i], header))
    }
    builder.WriteString("|\n")

    // Print separator
    separator := strings.Repeat("-", sum(columnWidths)+3*len(columnWidths)+1)
    builder.WriteString(separator + "\n")

    // Print rows
    for _, row := range t.rows {
        for i, value := range row {
            builder.WriteString(fmt.Sprintf("| %-*s ", columnWidths[i], value))
        }
        builder.WriteString("|\n")
    }

    // Print bottom separator
    builder.WriteString(separator + "\n")

    return builder.String()
}

func sum(nums []int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}
