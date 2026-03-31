DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'alerts'
          AND column_name = 'chat_id'
    ) THEN
        ALTER TABLE alerts ADD COLUMN chat_id BIGINT;
    ELSIF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'alerts'
          AND column_name = 'chat_id'
          AND data_type <> 'bigint'
    ) THEN
        ALTER TABLE alerts
            ALTER COLUMN chat_id TYPE BIGINT
            USING chat_id::BIGINT;
    END IF;
END $$;
