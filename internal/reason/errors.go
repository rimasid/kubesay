// Package reason turns an evidence bundle into a structured diagnosis using a
// single OpenAI-compatible LLM backend (Ollama by default, Groq optionally).
//
// The reasoning types and functions are added in later phases; this file
// defines the sentinel error that cmd/kubesay maps onto a process exit code.
package reason

import "errors"

// ErrRawFallback indicates the reasoner exhausted its retry budget without
// producing a schema-valid diagnosis. The caller prints the raw evidence
// bundle, which is still a usable result, so this maps to a success exit (0).
var ErrRawFallback = errors.New("reasoner fell back to raw evidence bundle")
