package proxy

import (
	"os"
	"regexp"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/oscarbc96/agbridge/pkg/awsutils"
)

func PrintMappings(handlerMapping map[*regexp.Regexp]Handler) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Path", "Methods", "Stage Variables", "Rest API ID", "Resource ID", "Account ID", "Region", "Identity"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Stage Variables", AutoMerge: true},
		{Name: "Rest API ID", AutoMerge: true},
		{Name: "Account ID", AutoMerge: true},
		{Name: "Region", AutoMerge: true},
		{Name: "Identity", WidthMax: 40, AutoMerge: true},
	})

	for _, handler := range handlerMapping {
		accountID, identity, err := awsutils.GetAccountDetails(handler.Config)
		if err != nil {
			return err
		}

		t.AppendRow(table.Row{
			handler.StagePath,
			handler.Methods,
			handler.StageVariables,
			handler.RestAPIID,
			handler.ResourceID,
			accountID,
			handler.Config.Region,
			identity,
		})
	}

	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateRows = true
	t.Render()

	return nil
}
