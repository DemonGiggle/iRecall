package main

import (
	"flag"
	"strings"
	"testing"
)

func TestUsageTextIncludesCurrentFlagsAndExamples(t *testing.T) {
	t.Parallel()

	fs := flag.NewFlagSet("irecall", flag.ContinueOnError)
	fs.Bool("debug", false, "enable debug logging")
	fs.String("data-path", "", "store database, config, and logs under this root path")
	fs.Bool("version", false, "print version and exit")

	text := usageText(fs, "irecall")

	for _, want := range []string{
		"Usage:",
		"-debug",
		"-data-path",
		"-version",
		"manual quote sharing via exported JSON",
		"asks for your display name",
		"/tmp/irecall-alice",
	} {
		if !strings.Contains(text, want) {
			t.Fatalf("usage text missing %q:\n%s", want, text)
		}
	}
}
