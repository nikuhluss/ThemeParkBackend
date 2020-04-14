SELECT
	EXTRACT(YEAR FROM posted_on) AS year,
	EXTRACT(MONTH FROM posted_on) AS month,
    TRIM(TO_CHAR(posted_on, 'Month')) AS month_name,
	COUNT(*) AS total_rainout
FROM theme_park.events
JOIN theme_park.event_types ON event_types.id = events.event_type_id
WHERE event_types.event_type ILIKE '%rainout%'
{{ if isSet "start" }}
	AND posted_on >= DATE_TRUNC('month', '{{.start}}'::timestamptz)
{{ end }}
{{ if isSet "end" }}
	AND posted_on <= DATE_TRUNC('month', '{{.end}}'::timestamptz + '1 month'::interval)
{{ end }}
GROUP BY year, month, month_name
ORDER BY year DESC, month DESC
