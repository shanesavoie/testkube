package commands

import (
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/oauth"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/telemetry"
	"github.com/kubeshop/testkube/cmd/kubectl-testkube/config"
	"github.com/kubeshop/testkube/pkg/ui"
	"github.com/spf13/cobra"
)

func NewStatusCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status [feature|resource]",
		Short: "Show status of feature or resource",
		Run: func(cmd *cobra.Command, args []string) {
			ui.NL()
			ui.Print(ui.IconRocket + "  Getting status on the testkube CLI")

			cfg, err := config.Load()
			ui.ExitOnError("   Loading config file failed", err)

			if cfg.TelemetryEnabled {
				ui.PrintEnabled("Telemetry on CLI", "enabled")
			} else {
				ui.PrintDisabled("Telemetry on CLI", "disabled")
			}

			if cfg.OAuth2Data.Enabled {
				ui.PrintEnabled("OAuth", "enabled")
			} else {
				ui.PrintDisabled("Oauth", "disabled")
			}
			ui.NL()
		},
	}

	cmd.AddCommand(telemetry.NewStatusTelemetryCmd())
	cmd.AddCommand(oauth.NewStatusOAuthCmd())

	return cmd
}
