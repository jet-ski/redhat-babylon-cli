package cmd

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"

	"github.com/akram/redhat-babylon-cli/pkg/api"
	"github.com/akram/redhat-babylon-cli/pkg/output"
)

var stopAt string

var serviceStopCmd = &cobra.Command{
	Use:   "stop <service-name>",
	Short: "Stop a running service",
	Long: `Stop a running service immediately or schedule a stop for later.

Examples:
  babylon service stop my-service                     # Stop immediately
  babylon service stop my-service --at 20:30          # Stop today at 20:30 (local time)
  babylon service stop my-service --at 2026-10-02T20:30:00+04:00  # Stop at specific time
  babylon service stop my-service --at 6h             # Stop in 6 hours`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		var stopTime time.Time
		if stopAt != "" {
			parsed, err := parseStopAt(stopAt)
			if err != nil {
				return fmt.Errorf("invalid --at: %w", err)
			}
			if parsed.Before(time.Now()) {
				return fmt.Errorf("--at is in the past (%s); omit --at to stop immediately", parsed.Format(time.RFC3339))
			}
			stopTime = parsed
		}

		claim, err := apiClient.StopResourceClaim(namespace, name, stopTime)
		if err != nil {
			return err
		}

		return output.Print(getOutputFormat(), claim, func(w io.Writer) {
			if stopTime.IsZero() {
				fmt.Fprintf(w, "Service %q stop requested.\n", api.DisplayName(claim))
			} else {
				fmt.Fprintf(w, "Service %q scheduled to stop at %s.\n", api.DisplayName(claim), stopTime.Format(time.RFC3339))
			}
			if claim.Status != nil && claim.Status.Summary != nil {
				fmt.Fprintf(w, "State: %s\n", claim.Status.Summary.State)
			}
		})
	},
}

func parseStopAt(input string) (time.Time, error) {
	// Try as duration first (e.g., "6h", "2d")
	if d, err := parseDurationString(input); err == nil {
		return time.Now().UTC().Add(d), nil
	}

	// Try as time-only (e.g., "20:30") — assume today, local time
	if t, err := time.Parse("15:04", input); err == nil {
		now := time.Now()
		scheduled := time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
		if scheduled.Before(now) {
			scheduled = scheduled.Add(24 * time.Hour)
		}
		return scheduled.UTC(), nil
	}

	// Try standard date/time formats
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02T15:04",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, input); err == nil {
			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse %q (try '20:30', '6h', '2d', or '2026-10-02T20:30:00+04:00')", input)
}

func init() {
	serviceCmd.AddCommand(serviceStopCmd)
	serviceStopCmd.Flags().StringVar(&stopAt, "at", "", "schedule stop at a specific time (e.g., '20:30', '6h', '2026-10-02T20:30:00+04:00')")
}
