package backend

import (
	"context"
	"slices"

	"github.com/duxlabs-dev/dux-go/ipc"
)

const APIVersion = 1

type Backend struct {
	version string
	bus     *ipc.EventBus
	server  *ipc.Server
	done    chan error
	stop    context.CancelFunc
	bar     bar
}

func Start(ctx context.Context, appID, version string) (*Backend, error) {
	ctx, stop := context.WithCancel(ctx)
	b := &Backend{
		version: version,
		bus:     ipc.NewEventBus(),
		done:    make(chan error, 1),
		stop:    stop,
	}
	b.server = ipc.NewServer(ipc.Config{
		AppName:      appID,
		APIVersion:   APIVersion,
		Capabilities: []string{"subscribe", "bar", "sysinfo"},
		Bus:          b.bus,
		OnSubscribe:  b.replayState,
	}, b.routes().ServeIPC)

	if err := b.server.Listen(); err != nil {
		stop()
		return nil, err
	}
	go func() { b.done <- b.server.Serve(ctx) }()
	go b.publishSystemInfo(ctx)
	return b, nil
}

func (b *Backend) SocketPath() string { return b.server.SocketPath() }

func (b *Backend) Done() <-chan error { return b.done }

func (b *Backend) Close() {
	b.stop()
	b.server.Close()
}

func (b *Backend) replayState(topics []string, _ *ipc.Subscriber) {
	if slices.Contains(topics, "bar") {
		b.publishBar(b.bar.isVisible())
	}
	if slices.Contains(topics, "sysinfo") {
		b.sampleSystemInfo()
	}
}

func (b *Backend) routes() *ipc.Mux {
	mux := ipc.NewMux()
	mux.Handle("version", b.handleVersion)
	mux.Handle("bar.get", b.handleBarGet)
	mux.Handle("bar.toggle", b.handleBarToggle)
	mux.Handle("sysinfo.get", b.handleSystemInfoGet)
	return mux
}

func (b *Backend) handleVersion(_ context.Context, w *ipc.ConnWriter, req ipc.Request, _ *ipc.Subscriber) {
	w.Respond(req.ID, map[string]any{"version": b.version, "apiVersion": APIVersion})
}
