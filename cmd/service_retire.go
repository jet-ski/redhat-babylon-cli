package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/akram/redhat-babylon-cli/pkg/api"
	"github.com/akram/redhat-babylon-cli/pkg/output"
)

var retireEndDate string

var serviceRetireCmd = &cobra.Command{
	Use:   "retire <service-name>",
	Short: "Retire a service or extend its lifespan",
	Long: `Retire a service by setting its lifespan end date.

Without --end-date, sets lifespan end to now (immediate retirement).
With --end-date, sets lifespan end to the specified date (extend or shorten).

Examples:
  babylon service retire my-service                    # Retire now
  babylon service retire my-service --end-date 7d      # Extend by 7 days from now
  babylon service retire my-service --end-date 2026-10-01  # Extend to specific date`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		var endTime time.Time
		if retireEndDate != "" {
			parsed, err := parseEndDate(retireEndDate, nil)
			if err != nil {
				return fmt.Errorf("invalid --end-date: %w", err)
			}
			endTime = parsed
		}

		claim, err := apiClient.SetLifespanEnd(namespace, name, endTime)
		if err != nil {
			return err
		}

		return output.Print(getOutputFormat(), claim, func(w io.Writer) {
			if endTime.IsZero() || endTime.Before(time.Now()) {
				fmt.Fprintf(w, "Service %q retirement requested.\n", api.DisplayName(claim))
			} else {
				fmt.Fprintf(w, "Service %q lifespan extended to %s.\n", api.DisplayName(claim), endTime.Format(time.RFC3339))
			}
		})
	},
}

func init() {
	serviceCmd.AddCommand(serviceRetireCmd)
	serviceRetireCmd.Flags().StringVar(&retireEndDate, "end-date", "", "lifespan end date (duration like '7d' or date like '2026-10-01'); defaults to now")
}
