-- +goose Up
-- +goose StatementBegin
CREATE TABLE invites (
    id UUID PRIMARY KEY,
    group_id UUID NOT NULL,
    sender_id UUID NOT NULL,
    invited_user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_sender FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_invited_user FOREIGN KEY (invited_user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unique_invite UNIQUE (group_id, invited_user_id)
);

CREATE INDEX idx_invites_group_id ON invites(group_id);
CREATE INDEX idx_invites_sender_id ON invites(sender_id);
CREATE INDEX idx_invites_invited_user_id ON invites(invited_user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS invites;
DROP INDEX IF EXISTS idx_invites_group_id;
DROP INDEX IF EXISTS idx_invites_sender_id;
DROP INDEX IF EXISTS idx_invites_invited_user_id;
-- +goose StatementEnd
