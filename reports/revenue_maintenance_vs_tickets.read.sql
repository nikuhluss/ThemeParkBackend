SELECT
    EXTRACT(YEAR FROM tickets_bought.year_month) AS year,
    EXTRACT(MONTH FROM tickets_bought.year_month) AS month,
    TRIM(TO_CHAR(tickets_bought.year_month, 'Month')) AS month_name,
    tickets_bought.total AS ticket_sales,
    maintenance_costs.total AS maintenance_costs,
    tickets_bought.total - maintenance_costs.total AS revenue
FROM(
	SELECT
		DATE_TRUNC('month', tickets.purchased_on) AS year_month,
		sum(tickets.purchase_price) as total
	FROM theme_park.tickets
	GROUP BY year_month
) as tickets_bought
LEFT JOIN(
	SELECT
		DATE_TRUNC('month', rides_maintenance.start_datetime) AS year_month,
		sum(rides_maintenance.cost) AS total
	FROM theme_park.rides_maintenance
	GROUP BY year_month
) AS maintenance_costs
ON maintenance_costs.year_month = tickets_bought.year_month
ORDER BY year DESC, month DESC