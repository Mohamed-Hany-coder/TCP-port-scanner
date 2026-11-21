package output

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/mohamedhany/tcp-port-scanner-go/internal/scanner"
)

func SortResults(results []scanner.Result) {
	sort.Slice(results, func(i, j int) bool { return results[i].Port < results[j].Port })
}

func WriteHuman(w io.Writer, results []scanner.Result) error {
	SortResults(results)
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	if _, err := fmt.Fprintln(tw, "PORT\tSTATUS"); err != nil {
		return err
	}
	open, closed, errors := 0, 0, 0
	for _, result := range results {
		if _, err := fmt.Fprintf(tw, "%d\t%s\n", result.Port, result.Status); err != nil {
			return err
		}
		switch result.Status {
		case scanner.StatusOpen:
			open++
		case scanner.StatusClosed:
			closed++
		default:
			errors++
		}
	}
	if err := tw.Flush(); err != nil {
		return err
	}
	_, err := fmt.Fprintf(w, "\nSummary: %d open, %d closed, %d errors\n", open, closed, errors)
	return err
}

func WriteJSON(w io.Writer, results []scanner.Result) error {
	SortResults(results)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(results)
}
