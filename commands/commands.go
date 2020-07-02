package commands

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/scnewma/ynabrt"
	"github.com/spf13/cobra"
	"go.bmvs.io/ynab"
)

func logError(fn func(cmd *cobra.Command, args []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		if err := fn(cmd, args); err != nil {
			fmt.Printf("ERROR: %v\n", err.Error())
			return err
		}
		return nil
	}
}

func getAPIClient(cmd *cobra.Command) (ynab.ClientServicer, error) {
	accessTokenEnv := os.Getenv("YNAB_ACCESS_TOKEN")
	accessToken := ""

	if flags := cmd.Flags(); flags != nil {
		if flag, err := flags.GetString("access-token"); err == nil {
			accessToken = flag
		}
	}

	if accessToken == "" {
		accessToken = accessTokenEnv

		if accessToken == "" {
			return nil, fmt.Errorf("No access token provided. Either provide the token by setting the YNAB_ACCESS_TOKEN environment variable or the --access-token flag")
		}
	}

	return ynab.NewClient(accessToken), nil
}

func getClient(cmd *cobra.Command) (*ynabrt.Client, error) {
	apiClient, err := getAPIClient(cmd)
	if err != nil {
		return nil, fmt.Errorf("could not construct ynab api client: %v", err)
	}

	client, err := ynabrt.NewClient(apiClient)
	if err != nil {
		return nil, fmt.Errorf("could not construct ynab client: %v", err)
	}

	return client, nil
}

func newTablewriter(out io.Writer) *tablewriter.Table {
	w := tablewriter.NewWriter(out)
	w.SetBorder(false)
	w.SetCenterSeparator("")
	w.SetRowLine(false)
	w.SetColumnSeparator("")
	w.SetHeaderLine(false)
	w.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	w.SetAutoWrapText(false)
	return w
}
