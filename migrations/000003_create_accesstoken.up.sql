CREATE TABLE IF NOT EXISTS access_tokens(
    access_token VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
)