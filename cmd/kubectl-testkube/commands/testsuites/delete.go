package testsuites

import (
	"strings"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	"github.com/kubeshop/testkube/pkg/ui"
	"github.com/spf13/cobra"
)

func NewDeleteTestSuiteCmd() *cobra.Command {
	var deleteAll bool
	var selectors []string

	cmd := &cobra.Command{
		Use:     "testsuite <testSuiteName>",
		Aliases: []string{"ts"},
		Short:   "Delete test suite",
		Long:    `Delete test suite by name`,
		Run: func(cmd *cobra.Command, args []string) {
			ui.Logo()
			client, _ := common.GetClient(cmd)
			namespace := cmd.Flag("namespace").Value.String()

			if deleteAll {
				err := client.DeleteTestSuites("")
				ui.ExitOnError("delete all tests from namespace "+namespace, err)
				ui.Success("Succesfully deleted all test suites in namespace", namespace)
			} else if len(args) > 0 {
				name := args[0]
				err := client.DeleteTestSuite(name)
				ui.ExitOnError("delete test suite "+name+" from namespace "+namespace, err)
				ui.Success("Succesfully deleted", name)
			} else if len(selectors) != 0 {
				selector := strings.Join(selectors, ",")
				err := client.DeleteTestSuites(selector)
				ui.ExitOnError("deleting test suites by labels: "+selector, err)
			} else {
				ui.Failf("Pass TestSuite name, --all flag to delete all, labels to delete by labels")
			}
		},
	}

	cmd.Flags().BoolVar(&deleteAll, "all", false, "Delete all tests")
	cmd.Flags().StringSliceVarP(&selectors, "label", "l", nil, "label key value pair: --label key1=value1")

	return cmd
}
