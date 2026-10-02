package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/akram/redhat-babylon-cli/pkg/api"
	"github.com/akram/redhat-babylon-cli/pkg/output"
)

var startRuntime string

var serviceStartCmd = &cobra.Command{
	Use:   "start <service-name>",
	Short: "Start a stopped service",
	Long: `Start a stopped service with optional custom runtime duration.

Examples:
  babylon service start my-service                # Start with default runtime (6h)
  babylon service start my-service --runtime 12h  # Start with 12-hour runtime
  babylon service start my-service --runtime 2d   # Start with 2-day runtime`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		var runtimeDuration time.Duration
		if startRuntime != "" {
			d, err := parseDurationString(startRuntime)
			if err != nil {
				return fmt.Errorf("invalid --runtime: %w", err)
			}
			runtimeDuration = d
		}

		claim, err := apiClient.StartResourceClaim(namespace, name, runtimeDuration)
		if err != nil {
			return err
		}

		return output.Print(getOutputFormat(), claim, func(w io.Writer) {
			fmt.Fprintf(w, "Service %q start requested.\n", api.DisplayName(claim))
			if claim.Status != nil && claim.Status.Summary != nil {
				fmt.Fprintf(w, "State: %s\n", claim.Status.Summary.State)
			}
		})
	},
}

func init() {
	serviceCmd.AddCommand(serviceStartCmd)
	serviceStartCmd.Flags().StringVar(&startRuntime, "runtime", "", "custom runtime duration (e.g., '12h', '2d'); defaults to service's runtime_default")
}
