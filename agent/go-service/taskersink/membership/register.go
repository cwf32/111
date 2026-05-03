package membership

import "github.com/MaaXYZ/maa-framework-go/v4"

var _ maa.TaskerEventSink = &MembershipChecker{}

// Register registers the membership checker as a tasker sink.
func Register() {
	maa.AgentServerAddTaskerSink(&MembershipChecker{})
}
