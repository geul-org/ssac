//ff:func feature=pkg-queue type=util control=sequence
//ff:what 큐 메시지 처리를 시작한다 — Backend 에 Start 시그니처 있을 때 위임
package queue

import "context"

// Starter is implemented by Backends that own a polling loop (e.g. the
// generated postgres backend). memory has no loop — Start blocks until ctx
// is cancelled.
type Starter interface {
	Start(ctx context.Context) error
}

// Start begins processing queued messages. It blocks until the context is
// cancelled. For backends that do not implement Starter this is a no-op that
// blocks on ctx.Done (memory semantics).
func Start(ctx context.Context) error {
	mu.RLock()
	b := backend
	mu.RUnlock()

	innerCtx, c := context.WithCancel(ctx)
	mu.Lock()
	cancel = c
	done = make(chan struct{})
	mu.Unlock()

	defer close(done)

	if s, ok := b.(Starter); ok {
		return s.Start(innerCtx)
	}
	<-innerCtx.Done()
	return nil
}
