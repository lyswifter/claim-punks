package fsm

import (
	"log"
	"reflect"

	"github.com/filecoin-project/go-statemachine"
	"golang.org/x/xerrors"
)

func (ic *ICPunks) Plan(events []statemachine.Event, user interface{}) (interface{}, uint64, error) {
	next, processed, err := ic.plan(events, user.(*TaskInfo))
	if err != nil || next == nil {
		return nil, processed, err
	}

	return func(ctx statemachine.Context, si TaskInfo) error {
		err := next(ctx, si)
		if err != nil {
			log.Fatalf("unhandled sector error (%s): %+v", si.Wallet, err)
			return nil
		}

		return nil
	}, processed, nil // TODO: This processed event count is not very correct
}

var fsmPlanners = map[TaskState]func(events []statemachine.Event, state *TaskInfo) (uint64, error){
	UndefinedSectorState: planOne(
		on(TaskStart{}, Start),
	),
	Start: planOne(
		on(TaskStart{}, Processing),
	),
	Processing: planOne(
		on(TaskSuccessed{}, Finished),
		on(TaskFailed{}, Failed),
	),
	// Successed: planOne(
	// 	on(TaskSuccessed{}, Finished),
	// ),
	// Failed: planOne(
	// 	on(TaskStart{}, Start),
	// ),
}

func (ic *ICPunks) plan(events []statemachine.Event, state *TaskInfo) (func(statemachine.Context, TaskInfo) error, uint64, error) {
	/////
	// First process all events

	log.Printf("State: %s", state.State)

	p := fsmPlanners[state.State]
	if p == nil {
		if len(events) == 1 {
			if _, ok := events[0].User.(globalMutator); ok {
				p = planOne() // in case we're in a really weird state, allow restart / update state / remove
			}
		}

		if p == nil {
			return nil, 0, xerrors.Errorf("planner for state %s not found", state.State)
		}
	}

	processed, err := p(events, state)
	if err != nil {
		return nil, 0, xerrors.Errorf("running planner for state %s failed: %w", state.State, err)
	}

	switch state.State {
	// Happy path
	case Empty:
		fallthrough
	case Start:
		fallthrough
	case Processing:
		return ic.hanleProcessing, processed, nil
	case Successed:
		return ic.handleOk, processed, nil
	case Failed:
		return ic.handleFailed, processed, nil
	}

	return nil, processed, nil
}

func on(mut mutator, next TaskState) func() (mutator, func(*TaskInfo) (bool, error)) {
	return func() (mutator, func(*TaskInfo) (bool, error)) {
		return mut, func(state *TaskInfo) (bool, error) {
			state.State = next
			return false, nil
		}
	}
}

func planOne(ts ...func() (mut mutator, next func(*TaskInfo) (more bool, err error))) func(events []statemachine.Event, state *TaskInfo) (uint64, error) {
	return func(events []statemachine.Event, state *TaskInfo) (uint64, error) {
	eloop:
		for i, event := range events {
			if gm, ok := event.User.(globalMutator); ok {
				gm.applyGlobal(state)
				return uint64(i + 1), nil
			}

			for _, t := range ts {
				mut, next := t()

				if reflect.TypeOf(event.User) != reflect.TypeOf(mut) {
					continue
				}

				if err, iserr := event.User.(error); iserr {
					log.Printf("wallet %s got error event %T: %+v", state.Wallet, event.User, err)
				}

				event.User.(mutator).apply(state)
				more, err := next(state)
				if err != nil || !more {
					return uint64(i + 1), err
				}

				continue eloop
			}

			_, ok := event.User.(Ignorable)
			if ok {
				continue
			}

			return uint64(i + 1), xerrors.Errorf("planner for state %s received unexpected event %T (%+v)", state.State, event.User, event)
		}

		return uint64(len(events)), nil
	}
}
