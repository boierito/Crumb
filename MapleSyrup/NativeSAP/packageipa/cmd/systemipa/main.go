// Host entry point for the same preparation, planning and export code used by
// Crumb. This command supplies no emulated system decryption or license keys.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/majd/ipatool/v2/wafflebridge/packageipa"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: systemipa prepare|unpack|archive-signed|plan|export [options]")
		os.Exit(2)
	}
	flags := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
	source := flags.String("source", "", "installed helper bundle or target directory")
	ipa := flags.String("ipa", "", "original purchased IPA (prepare only)")
	overrides := flags.String("overrides", "", "private native file outputs (export only)")
	output := flags.String("output", "", "new private IPA path")
	ids := flags.String("identifiers-file", "", "private JSON list of known identifiers; never included in output")
	_ = flags.Parse(os.Args[2:])
	var result any
	var err error
	switch os.Args[1] {
	case "prepare":
		result, err = packageipa.PrepareSystemHelper(*source, *ipa, *output)
	case "unpack":
		result, err = packageipa.UnpackSystemSigning(*source, *output)
	case "archive-signed":
		result, err = packageipa.ArchiveSystemSigning(*source, *ipa, *output)
	case "plan":
		result, err = packageipa.PlanSystemTarget(*source)
	case "export":
		var identifiers []string
		if *ids != "" {
			var data []byte
			data, err = os.ReadFile(*ids)
			if err == nil {
				err = json.Unmarshal(data, &identifiers)
			}
		}
		if err == nil {
			result, err = packageipa.ExportSystemResult(*source, *overrides, *output, identifiers)
		}
	default:
		fmt.Fprintln(os.Stderr, "unsupported operation")
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "system IPA operation rejected; no successful result produced")
		os.Exit(1)
	}
	_ = json.NewEncoder(os.Stdout).Encode(result)
}
