CREATE TABLE verification_tokens (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    token TEXT NOT NULL,
    type TEXT NOT NULL,
    expired_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP
);

CREATE INDEX idx_verification_tokens_user_id ON verification_tokens(user_id);
