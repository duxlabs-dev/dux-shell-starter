package backend

import (
	"bufio"
	"context"
	"encoding/json"
	"net"
	"testing"

	"github.com/duxlabs-dev/dux-go/ipc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func startBackend(t *testing.T) *Backend {
	t.Helper()
	t.Setenv("XDG_RUNTIME_DIR", t.TempDir())
	b, err := Start(context.Background(), "myshell", "test")
	require.NoError(t, err)
	t.Cleanup(b.Close)
	return b
}

func call(t *testing.T, b *Backend, method string) map[string]any {
	t.Helper()
	client, err := ipc.Dial(b.SocketPath())
	require.NoError(t, err)
	defer client.Close()

	resp, err := client.Call(ipc.Request{ID: 1, Method: method})
	require.NoError(t, err)
	require.Empty(t, resp.Error)
	result, ok := (*resp.Result).(map[string]any)
	require.True(t, ok)
	return result
}

type barEvent struct {
	Event string `json:"event"`
	Data  struct {
		Visible bool `json:"visible"`
	} `json:"data"`
}

func nextBarEvent(t *testing.T, lines *bufio.Scanner) bool {
	t.Helper()
	require.True(t, lines.Scan(), "bar event")
	var event barEvent
	require.NoError(t, json.Unmarshal(lines.Bytes(), &event))
	require.Equal(t, "bar", event.Event)
	return event.Data.Visible
}

func TestBarToggleFlipsVisibility(t *testing.T) {
	b := startBackend(t)

	assert.Equal(t, true, call(t, b, "bar.get")["visible"])
	assert.Equal(t, false, call(t, b, "bar.toggle")["visible"])
	assert.Equal(t, false, call(t, b, "bar.get")["visible"])
	assert.Equal(t, true, call(t, b, "bar.toggle")["visible"])
}

func TestSubscriberGetsBarStateThenChanges(t *testing.T) {
	b := startBackend(t)
	call(t, b, "bar.toggle")

	conn, err := net.Dial("unix", b.SocketPath())
	require.NoError(t, err)
	defer conn.Close()
	lines := bufio.NewScanner(conn)
	require.True(t, lines.Scan(), "capabilities line")

	_, err = conn.Write([]byte(`{"id":1,"method":"subscribe","params":{"topics":["bar"]}}` + "\n"))
	require.NoError(t, err)
	require.True(t, lines.Scan(), "subscribe response")
	assert.False(t, nextBarEvent(t, lines))

	call(t, b, "bar.toggle")
	assert.True(t, nextBarEvent(t, lines))
}
