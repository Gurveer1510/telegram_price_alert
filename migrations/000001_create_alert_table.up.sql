CREATE TABLE IF NOT EXISTS alerts(
    id SERIAL PRIMARY KEY,
    instrument_token INTEGER,
    instrument_name VARCHAR(255),
    chat_id BIGINT,
    exchange VARCHAR(255),
    trigger_price DECIMAL,
    condition VARCHAR(255),
    created_at TIMESTAMP
)
