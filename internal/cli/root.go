// Package cli builds the kubesay command tree: the root command and its
// configuration, plus the list and investigate subcommands (added in later
// phases). cmd/kubesay stays a thin entrypoint that delegates here.
package cli

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config keys. Centralizing them keeps flag binding, environment binding, and
// documented defaults in sync, and gives the subcommands one source of truth.
const (
	keyBackendActive  = "backend.active"
	keyOllamaEndpoint = "backend.ollama.endpoint"
	keyOllamaModel    = "backend.ollama.model"
	keyGroqEndpoint   = "backend.groq.endpoint"
	keyGroqModel      = "backend.groq.model"
	keyGroqAPIKey     = "backend.groq.api_key"
	keyRetryLimit     = "reasoner.retry_limit"
	keyOutputFormat   = "output.format"
	keyNamespace      = "namespace"
)

// commandFactory builds a subcommand given the shared Viper instance so the
// subcommand can read resolved configuration. list.go and investigate.go
// register their factories via registerCommand in an init function, and
// NewRootCmd wires them onto the root without importing each file explicitly.
type commandFactory func(v *viper.Viper) *cobra.Command

var commandFactories []commandFactory

// registerCommand records a subcommand factory to be attached to the root.
// Subcommands call this from an init function in their own file.
func registerCommand(f commandFactory) {
	commandFactories = append(commandFactories, f)
}

// NewRootCmd builds the kubesay root command: persistent flags bound to Viper,
// documented defaults and environment bindings for every configuration key, and
// a slog handler for operational diagnostics. Running kubesay with no
// subcommand prints help and exits 0.
func NewRootCmd() *cobra.Command {
	v := viper.New()

	var (
		cfgFile string
		verbose bool
	)

	cmd := &cobra.Command{
		Use:   "kubesay",
		Short: "Explain why a Kubernetes workload is failing (read-only)",
		Long: "KubeSay gathers the read-only evidence a human operator would collect about a\n" +
			"failing Kubernetes workload and uses a local or free LLM to explain the likely\n" +
			"root cause, separating observed facts from hypotheses. It never writes to the\n" +
			"cluster.",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return initConfig(v, cfgFile, verbose)
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			// No subcommand given: show help and exit 0.
			return cmd.Help()
		},
	}

	flags := cmd.PersistentFlags()
	flags.StringP(keyNamespace, "n", "", "Kubernetes namespace (defaults to the kubeconfig context namespace)")
	flags.String("backend", "", "active LLM backend: ollama or groq")
	flags.StringVar(&cfgFile, "config", "", "path to a config file")
	flags.StringP("output", "o", "", "output format: table or json")
	flags.BoolVar(&verbose, "verbose", false, "enable debug-level diagnostics on stderr")

	bindFlags(v, flags)
	bindEnv(v)
	setDefaults(v)

	for _, factory := range commandFactories {
		cmd.AddCommand(factory(v))
	}

	return cmd
}

// bindFlags wires the persistent flags onto their configuration keys so that an
// explicitly set flag takes precedence over the environment, the config file,
// and the defaults.
func bindFlags(v *viper.Viper, flags *pflag.FlagSet) {
	// Ignore bind errors: they can only occur if a flag name is misspelled,
	// which is a programming error caught immediately in testing.
	_ = v.BindPFlag(keyNamespace, flags.Lookup(keyNamespace))
	_ = v.BindPFlag(keyBackendActive, flags.Lookup("backend"))
	_ = v.BindPFlag(keyOutputFormat, flags.Lookup("output"))
}

// bindEnv binds each configuration key to its documented environment variable.
// The names do not follow a single prefix pattern (for example backend.active
// reads KUBESAY_BACKEND and backend.groq.api_key reads GROQ_API_KEY), so every
// key is bound explicitly rather than through a key replacer.
func bindEnv(v *viper.Viper) {
	// BindEnv never errors when given an explicit variable name.
	_ = v.BindEnv(keyBackendActive, "KUBESAY_BACKEND")
	_ = v.BindEnv(keyOllamaEndpoint, "KUBESAY_OLLAMA_ENDPOINT")
	_ = v.BindEnv(keyOllamaModel, "KUBESAY_OLLAMA_MODEL")
	_ = v.BindEnv(keyGroqEndpoint, "KUBESAY_GROQ_ENDPOINT")
	_ = v.BindEnv(keyGroqModel, "KUBESAY_GROQ_MODEL")
	_ = v.BindEnv(keyGroqAPIKey, "GROQ_API_KEY")
	_ = v.BindEnv(keyRetryLimit, "KUBESAY_RETRY_LIMIT")
	_ = v.BindEnv(keyOutputFormat, "KUBESAY_OUTPUT")
}

// setDefaults registers the documented default for every configuration key.
// These sit at the bottom of the precedence order (flags > env > file >
// defaults), so any unset key resolves to the value documented in the README
// and .env.example.
func setDefaults(v *viper.Viper) {
	v.SetDefault(keyBackendActive, "ollama")
	v.SetDefault(keyOllamaEndpoint, "http://localhost:11434/v1")
	v.SetDefault(keyOllamaModel, "qwen2.5")
	v.SetDefault(keyGroqEndpoint, "https://api.groq.com/openai/v1")
	v.SetDefault(keyGroqModel, "llama-3.1-8b-instant")
	v.SetDefault(keyGroqAPIKey, "")
	v.SetDefault(keyRetryLimit, 3)
	v.SetDefault(keyOutputFormat, "table")
}

// initConfig reads an optional config file and configures the slog handler.
// It runs after flags are parsed, so the --config path and --verbose level are
// known. Configuration precedence is flags > env > config file > defaults.
func initConfig(v *viper.Viper, cfgFile string, verbose bool) error {
	if cfgFile != "" {
		v.SetConfigFile(cfgFile)
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("reading config file %q: %w", cfgFile, err)
		}
	}
	configureLogger(verbose)
	return nil
}

// configureLogger installs a text slog handler on stderr as the default logger.
// Operational diagnostics flow through slog and stay separate from user-facing
// output; --verbose raises the level from Info to Debug.
func configureLogger(verbose bool) {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}
