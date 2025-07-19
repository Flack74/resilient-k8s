-- Alter the parameters column to accept text instead of JSONB
ALTER TABLE experiments ALTER COLUMN parameters TYPE TEXT;