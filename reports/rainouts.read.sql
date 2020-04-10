SELECT
	EXTRACT(YEAR FROM posted_on) AS year,
	EXTRACT(MONTH FROM posted_on) AS month,
    TRIM(TO_CHAR(posted_on, 'Month')) AS month_name,
	COUNT(*) AS total_rainout
FROM theme_park.events
WHERE title ILIKE '%rainout%'
GROUP BY year, month, month_name
ORDER BY year DESC, month DESC
