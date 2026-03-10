CREATE TABLE reservations (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id),
    room_id    BIGINT       NOT NULL REFERENCES rooms(id),
    start_at   TIMESTAMPTZ  NOT NULL,
    end_at     TIMESTAMPTZ  NOT NULL,
    status     VARCHAR(20)  NOT NULL DEFAULT 'confirmed',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_time CHECK (start_at < end_at)
);

CREATE INDEX idx_reservations_room_time ON reservations (room_id, start_at, end_at);
CREATE INDEX idx_reservations_user      ON reservations (user_id);
CREATE INDEX idx_reservations_start_at ON reservations (start_at);
CREATE INDEX idx_reservations_created  ON reservations (created_at);
