package terminal

import (
	"fmt"
	"strings"

	"github.com/10gen/realm-cli/internal/utils/flags"
)

// set of supported terminal flags
const (
	FlagAutoConfirm      = "yes"
	FlagAutoConfirmShort = "y"
	FlagAutoConfirmUsage = "set to automatically proceed through command confirmations"

	FlagDisableColors      = "disable-colors"
	FlagDisableColorsUsage = "disable output styling"

	FlagOutputFormat      = "output-format"
	FlagOutputFormatShort = "f"
	FlagOutputFormatUsage = "set the output format, available options: [json]"

	FlagOutputTarget      = "output-target"
	FlagOutputTargetShort = "o"
	FlagOutputTargetUsage = "write output to the specified filepath"
)

// DelimiterInline is the preferred inline delimiter when presenting information
const DelimiterInline = " - "

// OutputFormat is the terminal output format
type OutputFormat string

// String returns the output format display
func (of OutputFormat) String() string { return string(of) }

// Type returns the OutputFormat type
func (of OutputFormat) Type() string { return flags.TypeString }

// Set validates and sets the output format value
func (of *OutputFormat) Set(val string) error {
	outputFormat := OutputFormat(val)

	if !isValidOutputFormat(outputFormat) {
		allOutputFormats := []string{OutputFormatJSON.String()}
		return fmt.Errorf("unsupported value, use one of [%s] instead", strings.Join(allOutputFormats, ", "))
	}

	*of = outputFormat
	return nil
}

// set of supported terminal output formats
const (
	OutputFormatText OutputFormat = "" // zero-valued to be flag's default
	OutputFormatJSON OutputFormat = "json"
)

func isValidOutputFormat(outputFormat OutputFormat) bool {
	switch outputFormat {
	case
		OutputFormatJSON,
		OutputFormatText:
		return true
	}
	return false
}
