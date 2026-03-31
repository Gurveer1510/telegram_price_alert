DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'alerts'
          AND column_name = 'chat_id'
          AND data_type = 'bigint'
    ) THEN
        ALTER TABLE alerts
            ALTER COLUMN chat_id TYPE INTEGER
            USING chat_id::INTEGER;
    END IF;
END $$;
