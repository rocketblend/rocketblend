package command

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/rocketpack"
	"github.com/spf13/cobra"
)

// newDescribeCommand creates a new cobra command that fetches the definition of a package.
// It retrieves the definition based on the reference provided as an argument and formats the output based on the 'output' flag.
func (srv *Service) newDescribeCommand() *cobra.Command {
	var output string

	c := &cobra.Command{
		Use:   "describe [reference]",
		Short: "Fetches a package definition",
		Long:  `Fetches the definition of a package by its reference. The output can be formatted by specifying the 'output' flag.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref, err := srv.parseReference(args[0])
			if err != nil {
				return err
			}

			pack, err := srv.getPackageDescription(ref)
			if err != nil {
				return err
			}

			display, err := srv.getPackageOutput(pack, output)
			if err != nil {
				return err
			}

			cmd.Print(display)

			return nil
		},
	}

	c.Flags().StringVarP(&output, "output", "o", "table", "output format (table, json)")

	return c
}

// getPackageDescription returns the description of a package by reference.
func (srv *Service) getPackageDescription(ref *reference.Reference) (*rocketpack.RocketPack, error) {
	rocketblend, err := srv.factory.CreateRocketBlendService()
	if err != nil {
		return nil, fmt.Errorf("failed to create rocketblend service: %w", err)
	}

	pack, err := rocketblend.DescribePackByReference(*ref)
	if err != nil {
		return nil, fmt.Errorf("failed to describe package: %w", err)
	}

	return pack, nil
}

// getPackageOutput returns a string representation of the package based on the output format.
func (srv *Service) getPackageOutput(pack *rocketpack.RocketPack, output string) (string, error) {
	var displayString string
	var err error

	switch output {
	case "json":
		packJSON, err := json.Marshal(pack)
		if err != nil {
			return "", fmt.Errorf("failed to marshal package to JSON: %w", err)
		}
		displayString = string(packJSON)
	case "table":
		displayString, err = srv.renderRocketPackTable(pack)
		if err != nil {
			return "", fmt.Errorf("failed to print package table: %w", err)
		}
	default:
		return "", fmt.Errorf("unsupported output format: %s", output)
	}

	return displayString, nil
}

// renderRocketPackTable renders a table representation of a RocketPack.
func (srv *Service) renderRocketPackTable(pack *rocketpack.RocketPack) (string, error) {
	if pack.Addon == nil && pack.Build == nil {
		return "", fmt.Errorf("no addon or build present in the package")
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetCaption(true, "Package Description")

	if pack.Addon != nil {
		table.SetHeader([]string{"Addon Name", "Addon Version", "Addon File Source", "Addon URL Source"})

		addonName := pack.Addon.Name
		addonVersion := pack.Addon.Version.String()
		addonFileSource, addonURLSource := "", ""
		if pack.Addon.Source != nil {
			addonFileSource = pack.Addon.Source.File
			addonURLSource = pack.Addon.Source.URL
		}

		table.Append([]string{addonName, addonVersion, addonFileSource, addonURLSource})
	}

	if pack.Build != nil {
		table.SetHeader([]string{"Build Version", "Build Args"})

		buildVersion := pack.Build.Version.String()
		buildArgs := pack.Build.Args

		table.Append([]string{buildVersion, buildArgs})
	}

	table.Render()
	return tableString.String(), nil
}
