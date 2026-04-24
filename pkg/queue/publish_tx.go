//ff:func feature=pkg-queue type=util control=selection
//ff:what 트랜잭션 내에서 토픽에 메시지를 발행한다 — Backend.PublishTx 위임 (driver-중립)
package queue

import (
	"context"
	"encoding/json"
)

// PublishTx enqueues payload on topic inside the caller's transaction. The
// tx parameter is driver-neutral (any); the active Backend asserts the
// expected concrete type — typically pgx.Tx for the postgres backend or
// *sql.Tx for legacy database/sql. The memory backend returns
// ErrTxUnsupported.
//
// Atomicity: on Commit the row becomes visible to pollers; on Rollback no
// trace remains. The caller is responsible for the commit/rollback.
func PublishTx(ctx context.Context, tx any, topic string, payload any, opts ...PublishOption) error {
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
	return b.PublishTx(ctx, tx, topic, data, cfg)
}
