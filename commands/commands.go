package commands

import (
	"fmt"
	"os"

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

func getClient(cmd *cobra.Command) (*ynabrt.Client, error) {
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

	client, err := ynabrt.NewClient(ynab.NewClient(accessToken))
	if err != nil {
		return nil, fmt.Errorf("could not construct ynab client: %v", err)
	}
	return client, nil
}
