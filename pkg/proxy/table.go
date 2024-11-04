package proxy

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/oscarbc96/agbridge/pkg/awsutils"
)

type AWSConfigInfo struct {
	Profile   string
	AccountID string
	Region    string
}

func PrintMappings(handlerMapping map[string]Handler) error {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Path", "Methods", "Rest API ID", "Resource ID", "Account ID", "Region", "Identity"})
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: "Rest API ID", AutoMerge: true},
		{Name: "Account ID", AutoMerge: true},
		{Name: "Region", AutoMerge: true},
		{Name: "Identity", WidthMax: 40, AutoMerge: true},
	})

	for path, handler := range handlerMapping {
		accountID, identity, err := awsutils.GetAWSAccountDetails(handler.Config)
		if err != nil {
			return err
		}

		t.AppendRow(table.Row{
			path,
			handler.Methods,
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
