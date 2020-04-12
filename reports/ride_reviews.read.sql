SELECT
    rides.id AS ride_id,
    rides.name as ride_name,
    COUNT(reviews.*) AS review_count,
    AVG(reviews.rating)  AS review_avg
FROM theme_park.rides
LEFT JOIN theme_park.reviews ON reviews.ride_id = rides.id
GROUP BY rides.id
ORDER BY ride_name ASC
