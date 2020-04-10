SELECT
    users.id AS user_id,
    users.email AS user_email,
    user_details.first_name AS user_first_name,
    user_details.last_name AS user_last_name,
    COUNT(tickets_on_rides.*) AS scan_count
FROM theme_park.tickets_on_rides
JOIN theme_park.tickets ON tickets.id = tickets_on_rides.ticket_id
JOIN theme_park.users ON users.id = tickets.user_id
LEFT JOIN theme_park.user_details ON user_details.user_id = users.id
GROUP BY (users.id, users.email, user_details.first_name, user_details.last_name)
ORDER BY scan_count DESC
