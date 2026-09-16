package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/akram/redhat-babylon-cli/pkg/api"
	"github.com/akram/redhat-babylon-cli/pkg/output"
)

var (
	retireEndDate string
	retireNow     bool
)

var serviceRetireCmd = &cobra.Command{
	Use:   "retire <service-name>",
	Short: "Retire a service or extend its lifespan",
	Long: `Set the lifespan end date for a service.

Requires either --now (immediate retirement) or --end-date (extend/shorten).

Examples:
  babylon service retire my-service --now                  # Retire immediately
  babylon service retire my-service --end-date 7d          # Extend by 7 days from now
  babylon service retire my-service --end-date 2026-10-01  # Extend to specific date`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !retireNow && retireEndDate == "" {
			return fmt.Errorf("specify --now to retire immediately or --end-date to set a new expiration")
		}
		if retireNow && retireEndDate != "" {
			return fmt.Errorf("--now and --end-date are mutually exclusive")
		}

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
	serviceRetireCmd.Flags().BoolVar(&retireNow, "now", false, "retire immediately (set lifespan end to now)")
	serviceRetireCmd.Flags().StringVar(&retireEndDate, "end-date", "", "lifespan end date (duration like '7d' or date like '2026-10-01')")
}
