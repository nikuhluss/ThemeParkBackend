SELECT
    EXTRACT(YEAR FROM maintenance_start.year_month) AS year,
    EXTRACT(MONTH FROM maintenance_start.year_month) AS month,
    TRIM(TO_CHAR(maintenance_start.year_month, 'Month')) AS month_name,
    rides.id AS ride_id,
    rides.name AS ride_name,
    maintenance_start.total AS maintenance_start_total,
    COALESCE(maintenance_end.total, 0) AS maintenance_end_total
-- sub-query to get total maintenance started per-month
FROM (
    SELECT
        DATE_TRUNC('month', rides_maintenance.start_datetime) AS year_month,
        rides_maintenance.ride_id AS ride_id,
        COUNT(*) AS total
    FROM theme_park.rides_maintenance
    {{ if isSet "since" }}
    WHERE rides_maintenance.start_datetime >= DATE_TRUNC('month', '{{.since}}'::timestamptz)
    {{ end }}
    GROUP BY year_month, ride_id
) AS maintenance_start
-- join for rides information
JOIN theme_park.rides ON rides.id = maintenance_start.ride_id
-- sub-query to get total maintenance closed per-month
LEFT JOIN (
    SELECT
        DATE_TRUNC('month', rides_maintenance.end_datetime) AS year_month,
        rides_maintenance.ride_id AS ride_id,
        COUNT(*) AS total
    FROM theme_park.rides_maintenance
    WHERE rides_maintenance.end_datetime IS NOT NULL
    GROUP BY year_month, ride_id
) AS maintenance_end
ON maintenance_end.year_month = maintenance_start.year_month AND maintenance_end.ride_id = maintenance_start.ride_id
ORDER BY year DESC, month DESC, ride_name ASC
