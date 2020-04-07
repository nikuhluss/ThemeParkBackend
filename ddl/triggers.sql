-- ================================================================
-- Team 14 - Theme Park
-- Functions and triggers
--
-- formatted using: https://github.com/darold/pgFormatter
-- ================================================================

-- Functions
-- --------------------------------
-- Functions that execute after a trigger successfully matches.

SET search_path TO theme_park;

CREATE OR REPLACE FUNCTION emit_maintenance_status_event()
    RETURNS trigger AS
$BODY$
DECLARE
_status text;
_ride_name text;
_maintenance_type text;
_event_id text;
_event_type_id text;
_title text;
_description text;
_posted_on timestamp;
BEGIN

    -- status parameter
    _status = 'unknown';
    IF TG_NARGS = 1 THEN
        _status = TG_ARGV[0];
    END IF;

    -- create variables for the insert
    -- see: https://stackoverflow.com/a/21327318
    _event_id = (SELECT md5(random()::text || clock_timestamp()::text)::uuid);
    _event_type_id = (SELECT id FROM event_types WHERE event_type = 'System');
    _title = CONCAT('Maintenance ', _status);
    _ride_name = (SELECT name FROM rides WHERE id = NEW.ride_id);
    _maintenance_type = (SELECT maintenance_type FROM maintenance_types WHERE id = NEW.maintenance_type_id);
    _description = CONCAT(_maintenance_type, ' - ', 'On ride "', _ride_name, '"');
    _posted_on = NOW();

    -- insert event
    INSERT INTO events (id, event_type_id, title, description, posted_on)
        VALUES (_event_id, _event_type_id, _title, _description, _posted_on);

    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

-- Triggers
-- --------------------------------

-- maintenance started

DROP TRIGGER IF EXISTS ride_maintenance_started_event
    ON rides_maintenance;

CREATE TRIGGER ride_maintenance_started_event
    AFTER INSERT
    ON rides_maintenance
    FOR EACH ROW
    WHEN (NEW.end_datetime IS NULL)
    EXECUTE FUNCTION emit_maintenance_status_event('started');

-- maintenance closed

DROP TRIGGER IF EXISTS ride_maintenance_closed_event
    ON rides_maintenance;

CREATE TRIGGER ride_maintenance_closed_event
    AFTER UPDATE
    ON rides_maintenance
    FOR EACH ROW
    WHEN (OLD.end_datetime IS NULL AND NEW.end_datetime IS NOT NULL)
    EXECUTE FUNCTION emit_maintenance_status_event('closed');
