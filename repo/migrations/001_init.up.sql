CREATE TABLE IF NOT EXISTS notifications (
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMPTZ NOT NULL,
    modified_at     TIMESTAMPTZ NOT NULL,
    expiration_date TIMESTAMPTZ NOT NULL,
    message         TEXT,
    error           TEXT,
    user_uid        TEXT,
    message_type    TEXT,
    link            TEXT,
    status          TEXT CHECK (status IN ('NEW','COMPLETE')),
    subject         TEXT,
    created_by      TEXT
);

CREATE INDEX idx_notifications_user_uid ON notifications(user_uid);
CREATE INDEX idx_notifications_status   ON notifications(status);