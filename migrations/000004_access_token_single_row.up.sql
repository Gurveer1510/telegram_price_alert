TRUNCATE TABLE access_tokens;
ALTER TABLE access_tokens ADD COLUMN IF NOT EXISTS id INT DEFAULT 1;
ALTER TABLE access_tokens ADD CONSTRAINT access_tokens_pkey PRIMARY KEY (id);