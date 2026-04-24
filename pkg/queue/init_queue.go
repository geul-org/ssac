//ff:func feature=pkg-queue type=loader control=selection
//ff:what 큐 백엔드를 초기화한다 — memory 기본 또는 외부 주입된 Backend
package queue

import (
	"context"
	"errors"
	"sync"
)

var (
	// ErrNotInitialized — Publish/Subscribe/Start called before Init.
	ErrNotInitialized = errors.New("queue: not initialized, call Init first")
	// ErrUnknownBackend — unsupported backend name passed to Init.
	ErrUnknownBackend = errors.New("queue: unknown backend")
)

// singleton state. Exported via the package-level Publish/Subscribe helpers so
// runtime code need not pass a queue handle through every call site.
var (
	mu       sync.RWMutex
	handlers map[string][]func(ctx context.Context, msg []byte) error
	backend  Backend
	cancel   context.CancelFunc
	done     chan struct{}
	inited   bool
)

// Init initializes the queue with the memory backend. For durable backends
// (postgres, redis, etc.) callers use SetBackend(b) after yongol-generated
// code constructs the Backend implementation against the user's sqlc Queries.
//
// Signature accepts backendName for backward-compatible call sites — only
// "memory" is accepted from inside ssac. Any other name (including
// "postgres") returns ErrUnknownBackend; the caller must instead call
// SetBackend(externalImpl) directly.
func Init(ctx context.Context, backendName string) error {
	mu.Lock()
	defer mu.Unlock()

	switch backendName {
	case "memory":
		backend = newMemoryBackend()
	default:
		return errors.New("queue: Init only accepts \"memory\"; use SetBackend for durable backends")
	}

	handlers = make(map[string][]func(ctx context.Context, msg []byte) error)
	inited = true
	_ = ctx // reserved for backends that need ctx at init
	return nil
}

// SetBackend installs an externally constructed Backend (e.g. a yongol-
// generated postgres implementation). Use this from main.go after building
// the backend from the user's sqlc Queries:
//
//	q := db.New(pool)
//	queue.SetBackend(postgresqueue.NewPostgres(q))
//
// Handlers registered via Subscribe before SetBackend survive the swap.
func SetBackend(b Backend) {
	mu.Lock()
	defer mu.Unlock()
	backend = b
	if handlers == nil {
		handlers = make(map[string][]func(ctx context.Context, msg []byte) error)
	}
	inited = true
}
