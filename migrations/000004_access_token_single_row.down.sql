ALTER TABLE access_tokens DROP CONSTRAINT IF EXISTS access_tokens_pkey;
ALTER TABLE access_tokens DROP COLUMN IF EXISTS id;