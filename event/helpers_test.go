package event_test

import (
	"artk.dev/event"
	"artk.dev/testbarrier"
	"context"
	"reflect"
	"sync"
	"testing"
	"time"
)

func assertIsDeepCopy[Event any](t *testing.T, originalEvent Event) {
	const numObservers = 2
	var receivedEvents [numObservers]Event
	var wg sync.WaitGroup
	wg.Add(numObservers)

	t.Log("Given that there are two observers,")
	mux := event.NewMux[Event]()
	for i := range receivedEvents {
		mux.WillNotify(func(_ context.Context, e Event) error {
			defer wg.Done()
			receivedEvents[i] = e
			return nil
		})
	}

	t.Log("When Mux observes an event,")
	err := mux.Observe(context.TODO(), originalEvent)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	testbarrier.WaitForGroup(t, &wg, 5*time.Second)

	t.Log("Then the events sent to the observers are deep copies.")
	if same(receivedEvents[0], receivedEvents[1]) {
		t.Error("The received event are shallow copies")
	}
	if same(originalEvent, receivedEvents[0]) {
		t.Error("The first event is a shallow copy")
	}
	if same(originalEvent, receivedEvents[1]) {
		t.Error("The second event is a shallow copy")
	}
}

func exampleEvent() Event {
	return Event{
		ID:   expectedID,
		Name: expectedName,
	}
}

func same(x, y any) bool {
	vx := reflect.ValueOf(x)
	vy := reflect.ValueOf(y)
	if !vx.IsValid() {
		return !vy.IsValid()
	}

	return vx.Pointer() == vy.Pointer()
}

const expectedID = 1234
const expectedName = "Test Event"

type key struct{}
