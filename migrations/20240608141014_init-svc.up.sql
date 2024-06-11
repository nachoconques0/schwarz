BEGIN;

CREATE SCHEMA IF NOT EXISTS schwarz;

CREATE OR REPLACE FUNCTION schwarz.set_updated_at ()
    RETURNS TRIGGER STABLE
    AS $plpgsql$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$plpgsql$
LANGUAGE plpgsql;

COMMIT;
