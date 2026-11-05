// Command kubesay is the entrypoint for the KubeSay CLI. It builds the root
// command, executes it, and maps classified errors onto process exit codes.
package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/rimasid/kubesay/internal/cli"
	"github.com/rimasid/kubesay/internal/k8s"
	"github.com/rimasid/kubesay/internal/reason"
)

func main() {
	os.Exit(run())
}

// run executes the root command and returns the process exit code.
func run() int {
	if err := cli.NewRootCmd().Execute(); err != nil {
		return exitCode(err)
	}
	return 0
}

// exitCode maps a command error onto a deterministic process exit code.
//
//	nil / ErrRawFallback -> 0 (the raw evidence bundle is still a usable result)
//	ErrKubeconfig        -> 1
//	ErrNoFailingPod      -> 1
//	any other error      -> 1
//
// Errors that map to a non-zero code are printed to stderr; the raw-bundle
// fallback is not, because its output has already been shown to the user.
func exitCode(err error) int {
	switch {
	case err == nil:
		return 0
	case errors.Is(err, reason.ErrRawFallback):
		return 0
	case errors.Is(err, k8s.ErrKubeconfig), errors.Is(err, k8s.ErrNoFailingPod):
		fmt.Fprintln(os.Stderr, "kubesay: "+err.Error())
		return 1
	default:
		fmt.Fprintln(os.Stderr, "kubesay: "+err.Error())
		return 1
	}
}
