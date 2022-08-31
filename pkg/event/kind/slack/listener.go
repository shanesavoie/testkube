package slack

import (
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/event/kind/common"
	"github.com/kubeshop/testkube/pkg/log"
	"github.com/kubeshop/testkube/pkg/slacknotifier"
	"go.uber.org/zap"
)

var _ common.Listener = &SlackListener{}

func NewSlackListener(selector string, events []testkube.TestkubeEventType) *SlackListener {
	return &SlackListener{
		Log:      log.DefaultLogger,
		selector: selector,
		events:   events,
	}
}

type SlackListener struct {
	Log      *zap.SugaredLogger
	events   []testkube.TestkubeEventType
	selector string
}

func (l *SlackListener) Selector() string {
	return l.selector
}

func (l *SlackListener) Events() []testkube.TestkubeEventType {
	return l.events
}
func (l *SlackListener) Metadata() map[string]string {
	return map[string]string{}
}

func (l *SlackListener) Notify(event testkube.TestkubeEvent) (result testkube.TestkubeEventResult) {
	err := slacknotifier.SendEvent(event.Type_, *event.Execution)
	if err != nil {
		return testkube.NewFailedTestkubeEventResult(event.Id, err)
	}

	return testkube.NewSuccessTestkubeEventResult(event.Id, "event sent to slack")
}

func (l *SlackListener) Kind() string {
	return "slack"
}
