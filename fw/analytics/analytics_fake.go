package analytics

import "github.com/short-d/app/fw"

var _ Analytics = (*Fake)(nil)

type Fake struct {
}

func (f Fake) Identify(userID string, traits map[string]string) {
}

func (f Fake) Track(eventName string, properties map[string]string, userID string, ctx fw.ExecutionContext) {
}

func (f Fake) Group(userID string, groupID string) {
}

func (f Fake) Alias(prevUserID string, newUserID string) {
}

func NewFake() Fake {
	return Fake{}
}
