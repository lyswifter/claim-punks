package fsm

type mutator interface {
	apply(state *TaskInfo)
}

// globalMutator is an event which can apply in every state
type globalMutator interface {
	// applyGlobal applies the event to the state. If if returns true,
	//  event processing should be interrupted
	applyGlobal(state *TaskInfo) bool
}

type Ignorable interface {
	Ignore()
}

type TaskStart struct {
	ID     TaskNumber
	IP     string
	Wallet string
}

func (evt TaskStart) apply(state *TaskInfo) {
	state.ID = evt.ID
	state.IP = evt.IP
	state.Wallet = evt.Wallet
}

type TaskSuccessed struct {
	TokenID string
}

func (evt TaskSuccessed) apply(state *TaskInfo) {
	state.TokenID = evt.TokenID
}

type TaskFailed struct{ error }

func (evt TaskFailed) apply(state *TaskInfo) {
	state.TokenID = ""
}

type TaskFinished struct {
	TokenID string
}

func (evt TaskFinished) apply(state *TaskInfo) {
	state.TokenID = evt.TokenID
}
