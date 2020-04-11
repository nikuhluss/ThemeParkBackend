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

-- emit_maintenance_status_event creates a new event with the passed status
-- for the ride_id in the row.
CREATE OR REPLACE FUNCTION emit_maintenance_status_event ()
    RETURNS TRIGGER AS
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

    -- get required values
    _ride_name = (SELECT name FROM rides WHERE id = NEW.ride_id);
    _maintenance_type = (SELECT maintenance_type FROM maintenance_types WHERE id = NEW.maintenance_type_id);

    -- create variables for the event insert
    -- see: https://stackoverflow.com/a/21327318
    _event_id = (SELECT md5(random()::text || clock_timestamp()::text)::uuid);
    _event_type_id = (SELECT id FROM event_types WHERE event_type = 'System');
    _title = CONCAT('Maintenance ', _status);
    _description = CONCAT(_maintenance_type, ' - ', 'For ride "', _ride_name, '"');
    _posted_on = NOW();

    -- insert event
    INSERT INTO events (id, event_type_id, title, description, posted_on)
        VALUES (_event_id, _event_type_id, _title, _description, _posted_on);

    RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;


-- emit_bad_review_posted_event a new event when a bad review
-- is posted.
CREATE OR REPLACE FUNCTION emit_bad_review_posted_event ()
    RETURNS TRIGGER AS
$BODY$
DECLARE
    _user_email text;
    _ride_name text;

    _event_id text;
    _event_type_id text;
    _title text;
    _description text;
    _posted_on timestamp;
BEGIN
    -- get required values
    _user_email = (SELECT email FROM users WHERE id = NEW.customer_id);
    _ride_name = (SELECT name FROM rides WHERE id = NEW.ride_id);

    -- create variables for the event insert
    -- see: https://stackoverflow.com/a/21327318
    _event_id = (SELECT md5(random()::text || clock_timestamp()::text)::uuid);
    _event_type_id = (SELECT id FROM event_types WHERE event_type = 'System');
    _title = 'Bad review posted';
    _description = CONCAT('Rating of ', NEW.rating, ' by user ', _user_email, ' - ', 'For ride "', _ride_name, '"');
    _posted_on = NOW();

    -- insert event
    INSERT INTO events (id, event_type_id, title, description, posted_on)
        VALUES (_event_id, _event_type_id, _title, _description, _posted_on);

    RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

-- emit_bad_reviews_event_when_rating_average_below creates a new event
-- when the rating average for the ride given by NEW.ride_id is below
-- the given threshold.
CREATE OR REPLACE FUNCTION emit_bad_reviews_event_when_rating_average_below ()
    RETURNS TRIGGER AS
$BODY$
DECLARE
    _threshold real;
    _review_count integer;
    _rating_avg real;

    _ride_name text;

    _event_id text;
    _event_type_id text;
    _title text;
    _description text;
    _posted_on timestamp;
BEGIN
    -- threshold parameter
    _threshold = 3.0;
    IF TG_NARGS = 1 THEN
        _threshold = TG_ARGV[0];
    END IF;

    -- early return if review count is 0
    _review_count = (SELECT COUNT(*) FROM reviews WHERE ride_id = NEW.ride_id);
    IF _review_count <= 0 THEN
        return NEW;
    END IF;

    -- early return if reviews average is equal or greater than threshold
    _rating_avg = (SELECT AVG(rating) FROM reviews WHERE ride_id = NEW.ride_id);
    IF _rating_avg >= _threshold THEN
        return NEW;
    END IF;

    -- get required values
    _ride_name = (SELECT name FROM rides WHERE id = NEW.ride_id);

    -- create variables for the event insert
    -- see: https://stackoverflow.com/a/21327318
    _event_id = (SELECT md5(random()::text || clock_timestamp()::text)::uuid);
    _event_type_id = (SELECT id FROM event_types WHERE event_type = 'System');
    _title = 'Ratings below threshold';
    _description = CONCAT('Ratings below ', _threshold, ' - ', 'For ride "', _ride_name, '"');
    _posted_on = NOW();

    -- insert event
    INSERT INTO events (id, event_type_id, title, description, posted_on)
        VALUES (_event_id, _event_type_id, _title, _description, _posted_on);

    RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;

-- Triggers
-- --------------------------------

-- maintenance started

DROP TRIGGER IF EXISTS ride_maintenance_started_event ON rides_maintenance;

CREATE TRIGGER ride_maintenance_started_event
    AFTER INSERT ON rides_maintenance
    FOR EACH ROW
    WHEN (NEW.end_datetime IS NULL)
    EXECUTE FUNCTION emit_maintenance_status_event ('started');

-- maintenance closed

DROP TRIGGER IF EXISTS ride_maintenance_closed_event ON rides_maintenance;

CREATE TRIGGER ride_maintenance_closed_event
    AFTER UPDATE ON rides_maintenance
    FOR EACH ROW
    WHEN (OLD.end_datetime IS NULL AND NEW.end_datetime IS NOT NULL)
    EXECUTE FUNCTION emit_maintenance_status_event ('closed');

-- maintenance reopened

DROP TRIGGER IF EXISTS ride_maintenance_reopened_event ON rides_maintenance;

CREATE TRIGGER ride_maintenance_reopened_event
    AFTER UPDATE ON rides_maintenance
    FOR EACH ROW
    WHEN (OLD.end_datetime IS NOT NULL AND NEW.end_datetime IS NULL)
    EXECUTE FUNCTION emit_maintenance_status_event ('reopened');

-- bad review posted on ride

DROP TRIGGER IF EXISTS ride_bad_review_posted_event ON reviews;

CREATE TRIGGER ride_bad_review_posted_event
    AFTER INSERT OR UPDATE ON reviews
    FOR EACH ROW
    WHEN (NEW.rating < 3.0)
    EXECUTE FUNCTION emit_bad_review_posted_event();

-- bad reviews on ride

DROP TRIGGER IF EXISTS ride_bad_reviews_event ON reviews;

CREATE TRIGGER ride_bad_reviews_event
    AFTER INSERT OR UPDATE ON reviews
    FOR EACH ROW
    -- NOTE: WHEN inside function since subqueries are not supported by postgres
    EXECUTE FUNCTION emit_bad_reviews_event_when_rating_average_below(3.0);
