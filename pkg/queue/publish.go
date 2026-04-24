//ff:func feature=pkg-queue type=util control=selection
//ff:what 토픽에 메시지를 발행한다 — Backend 에 위임
package queue

import (
	"context"
	"encoding/json"
)

// Publish serializes payload to JSON and delegates to the active Backend.
// Returns ErrNotInitialized if neither Init nor SetBackend has run.
func Publish(ctx context.Context, topic string, payload any, opts ...PublishOption) error {
	mu.RLock()
	if !inited {
		mu.RUnlock()
		return ErrNotInitialized
	}
	b := backend
	mu.RUnlock()

	cfg := applyPublishOpts(opts)

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return b.Publish(ctx, topic, data, cfg)
}
