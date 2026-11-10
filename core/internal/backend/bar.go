package backend

import (
	"context"
	"sync"

	"github.com/duxlabs-dev/dux-go/ipc"
)

type bar struct {
	mu     sync.Mutex
	hidden bool
}

func (b *bar) isVisible() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return !b.hidden
}

func (b *bar) toggle() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.hidden = !b.hidden
	return !b.hidden
}

func (b *Backend) publishBar(visible bool) map[string]any {
	state := map[string]any{"visible": visible}
	b.bus.Publish("bar", state)
	return state
}

func (b *Backend) handleBarGet(_ context.Context, w *ipc.ConnWriter, req ipc.Request, _ *ipc.Subscriber) {
	w.Respond(req.ID, map[string]any{"visible": b.bar.isVisible()})
}

func (b *Backend) handleBarToggle(_ context.Context, w *ipc.ConnWriter, req ipc.Request, _ *ipc.Subscriber) {
	w.Respond(req.ID, b.publishBar(b.bar.toggle()))
}
