SELECT
	EXTRACT(YEAR FROM t.year_month) AS year,
	EXTRACT(MONTH FROM t.year_month) AS month,
    TRIM(TO_CHAR(t.year_month, 'Month')) AS month_name,
    t.ride_id AS ride_id,
    rides.name AS ride_name,
    t.times_ridden AS times_ridden
FROM (
	SELECT
        DATE_TRUNC('month', scan_datetime) AS year_month,
        tickets_on_rides.ride_id AS ride_id,
        COUNT(*) AS times_ridden,
        RANK() OVER (PARTITION BY date_trunc('month', scan_datetime) ORDER BY COUNT(*) DESC) AS rnk
	FROM theme_park.tickets_on_rides
	GROUP BY date_trunc('month', scan_datetime), ride_id
) AS t
JOIN theme_park.rides ON rides.id = t.ride_id
WHERE t.rnk = 1
ORDER BY year DESC, month DESC
