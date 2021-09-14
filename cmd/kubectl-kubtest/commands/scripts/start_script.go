package scripts

import (
	"fmt"
	"os"
	"time"

	"github.com/kubeshop/kubtest/pkg/ui"
	"github.com/spf13/cobra"
)

const WatchInterval = 2 * time.Second

var watch bool
var params map[string]string

func NewStartScriptCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"run"},
		Short:   "Starts new script",
		Long:    `Starts new script based on Script Custom Resource name, returns results to console`,
		Run: func(cmd *cobra.Command, args []string) {
			ui.Logo()

			if len(args) == 0 {
				ui.ExitOnError("Invalid arguments", fmt.Errorf("please pass script name to run"))
			}

			scriptID := args[0]

			client, namespace := GetClient(cmd)
			namespacedName := fmt.Sprintf("%s/%s", namespace, scriptID)

			scriptExecution, err := client.ExecuteScript(scriptID, namespace, name, params)
			ui.ExitOnError("starting script execution "+namespacedName, err)

			execution := scriptExecution.Execution

			switch true {

			case execution.IsQueued():
				ui.Warn("Script queued for execution")

			case execution.IsPending():
				ui.Warn("Script execution started")

			case execution.IsSuccesful():
				fmt.Println(execution.Result.RawOutput)
				duration := execution.EndTime.Sub(execution.StartTime)
				ui.Success("Script execution completed with sucess in " + duration.String())

			case execution.IsFailed():
				fmt.Println(execution.Result.ErrorMessage)
				ui.Errf("Script execution failed")

			}

			ui.ShellCommand(
				"Use following command to get script execution details",
				"kubectl kubtest scripts execution "+scriptExecution.Id,
			)
			ui.ShellCommand(
				"or watch script execution until complete",
				"kubectl kubtest scripts watch "+scriptExecution.Id,
			)
			ui.NL()
			if watch {
				for range time.Tick(time.Second) {
					scriptExecution, err := client.GetExecution("-", scriptExecution.Id)
					ui.ExitOnError("getting API for script completion", err)
					render := GetRenderer(cmd)
					err = render.Render(scriptExecution, os.Stdout)
					ui.ExitOnError("rendering", err)
					if scriptExecution.Execution.IsCompleted() {
						return
					}
				}
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "execution name, if empty will be autogenerated")
	cmd.Flags().StringToStringVarP(&params, "param", "p", map[string]string{}, "execution envs passed to executor")
	cmd.Flags().BoolVarP(&watch, "watch", "f", false, "watch for changes after start")

	return cmd
}
