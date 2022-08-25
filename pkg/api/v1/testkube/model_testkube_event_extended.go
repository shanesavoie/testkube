package testkube

import "k8s.io/apimachinery/pkg/labels"

func (e TestkubeEvent) Log() []any {

	var executionId, executionName, eventType, labels string
	if e.Execution != nil {
		executionId = e.Execution.Id
		executionName = e.Execution.Name
		for k, v := range e.Execution.Labels {
			labels += k + "=" + v + " "
		}
	}

	if e.Type_ != nil {
		eventType = e.Type_.String()
	}

	return []any{
		"uri", e.Uri,
		"type", eventType,
		"executionId", executionId,
		"executionName", executionName,
		"labels", labels,
	}
}

func (e TestkubeEvent) Valid(selector string) (valid bool) {
	if e.Execution == nil {
		return false
	}

	valid = selector == ""
	if !valid {
		selector, err := labels.Parse(selector)
		if err != nil {
			return false
		}

		valid = selector.Matches(labels.Set(e.Execution.Labels))
	}

	return

}
