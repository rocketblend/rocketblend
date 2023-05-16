package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rocketblend/rocketblend/pkg/rocketblend"
	"github.com/spf13/cobra"
)

// newResolveCommand creates a new cobra.Command that outputs resolved information about the project.
// This information includes dependencies and paths for the project on the local machine.
func (srv *Service) newResolveCommand() *cobra.Command {
	var output string

	c := &cobra.Command{
		Use:   "resolve",
		Short: "Resolves and outputs project details",
		Long:  `Fetches and prints the resolved dependencies and paths for the project in the local machine in the specified output format (JSON or table)`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			blend, err := srv.findBlendFile(srv.flags.workingDirectory)
			if err != nil {
				return fmt.Errorf("failed to find blend file: %w", err)
			}

			blendOutput, err := srv.getBlendFileOutput(blend, output)
			if err != nil {
				return err
			}

			cmd.Println(blendOutput)

			return nil
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "json", "output format (table, json)")

	return c
}

// getBlendFileOutput generates the output for the blend file in the specified format.
// This function supports two output formats: 'json' and 'table'.
func (srv *Service) getBlendFileOutput(blend *rocketblend.BlendFile, output string) (string, error) {
	switch output {
	case "json":
		jsonOutput, err := json.Marshal(blend)
		if err != nil {
			return "", fmt.Errorf("failed to resolve config: %w", err)
		}
		return string(jsonOutput), nil
	case "table":
		return renderBlendFileTable(blend), nil
	default:
		return "", fmt.Errorf("invalid output format: %s", output)
	}
}

// renderBlendFileTable generates a table string representation of the blend file.
// This function is used when the output format for the 'resolve' command is 'table'.
func renderBlendFileTable(blend *rocketblend.BlendFile) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	// Set table header
	table.SetHeader([]string{"Type", "Name", "Version", "Path", "ARGS"})

	// Set table caption
	table.SetCaption(true, "Resolved Dependencies")

	// Add Build to the table
	if blend.Build != nil {
		table.Append([]string{"Build", "", "", blend.Build.Path, blend.Build.ARGS})
		for _, addon := range *blend.Build.Addons {
			table.Append([]string{"Addon (Build)", addon.Name, addon.Version.String(), addon.Path, ""})
		}
	}

	// Add Addons to the table
	if blend.Addons != nil {
		for _, addon := range *blend.Addons {
			table.Append([]string{"Addon", addon.Name, addon.Version.String(), addon.Path, ""})
		}
	}

	table.Render()
	return tableString.String()
}
