SELECT
    EXTRACT(YEAR FROM COALESCE(tickets_sold.year_month, maintenance_costs.year_month)) AS year,
    EXTRACT(MONTH FROM COALESCE(tickets_sold.year_month, maintenance_costs.year_month)) AS month,
    TRIM(TO_CHAR(COALESCE(tickets_sold.year_month, maintenance_costs.year_month), 'Month')) AS month_name,
    COALESCE(tickets_sold.total, 0) AS ticket_sales,
    COALESCE(maintenance_costs.total, 0) AS maintenance_costs,
    COALESCE(tickets_sold.total, 0) - COALESCE(maintenance_costs.total, 0) AS revenue
FROM(
	SELECT
		DATE_TRUNC('month', tickets.purchased_on) AS year_month,
		sum(tickets.purchase_price) as total
	FROM theme_park.tickets
	GROUP BY year_month
) as tickets_sold
FULL OUTER JOIN (
	SELECT
		DATE_TRUNC('month', rides_maintenance.end_datetime) AS year_month,
		sum(rides_maintenance.cost) AS total
	FROM theme_park.rides_maintenance
	WHERE rides_maintenance.end_datetime IS NOT NULL
	GROUP BY year_month
) AS maintenance_costs
ON maintenance_costs.year_month = tickets_sold.year_month
WHERE 1=1
{{ if isSet "start" }}
	AND COALESCE(tickets_sold.year_month, maintenance_costs.year_month) >= DATE_TRUNC('month', '{{.start}}'::timestamptz)
{{ end }}
{{ if isSet "end" }}
	AND COALESCE(tickets_sold.year_month, maintenance_costs.year_month) <= DATE_TRUNC('month', '{{.end}}'::timestamptz)
{{ end }}
ORDER BY year DESC, month DESC
