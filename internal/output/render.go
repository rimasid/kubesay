// Package output renders typed values for the terminal, either as an aligned
// table (the default) or as JSON.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/rimasid/kubesay/internal/k8s"
)

// Format selects how values are rendered.
type Format string

const (
	// FormatTable renders human-readable aligned columns (the default).
	FormatTable Format = "table"
	// FormatJSON renders indented JSON.
	FormatJSON Format = "json"
)

// ParseFormat validates s (typically the `output.format` config value) and
// returns the matching Format.
func ParseFormat(s string) (Format, error) {
	switch f := Format(s); f {
	case FormatTable, FormatJSON:
		return f, nil
	default:
		return "", fmt.Errorf("unsupported output format %q (use \"table\" or \"json\")", s)
	}
}

// RenderPods writes pods to w in the given format ("table" or "json").
//
// It takes an io.Writer instead of printing to os.Stdout directly so the
// output can be captured in tests (e.g. with a bytes.Buffer). The caller
// passes os.Stdout in the real program.
func RenderPods(w io.Writer, format string, pods []k8s.PodSummary) error {
	f, err := ParseFormat(format)
	if err != nil {
		return err
	}

	if f == FormatJSON {
		// A nil slice would marshal to "null"; force an empty list to "[]".
		if pods == nil {
			pods = []k8s.PodSummary{}
		}
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(pods)
	}

	if len(pods) == 0 {
		_, err := fmt.Fprintln(w, "No pods found.")
		return err
	}

	// tabwriter buffers tab-separated cells and pads them so each column lines
	// up. Nothing is written to w until Flush is called.
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	fmt.Fprintln(tw, "NAMESPACE\tNAME\tREADY\tSTATUS\tRESTARTS\tAGE")
	for _, p := range pods {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%d\t%s\n",
			p.Namespace, p.Name, p.Ready, p.Phase, p.Restarts, p.Age)
	}
	return tw.Flush()
}
