package runner

import (
	"github.com/kubeshop/kubtest/pkg/api/kubtest"
)

// Runner interface to abstract runners implementations
type Runner interface {
	// Run takes Execution data and returns execution result
	Run(execution kubtest.Execution) kubtest.ExecutionResult
}