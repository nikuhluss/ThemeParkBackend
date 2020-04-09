SELECT rides.id, rides.name, COUNT(reviews.*) AS review_count, AVG(reviews.rating)  AS review_avg
    FROM theme_park.rides
    LEFT JOIN theme_park.reviews ON reviews.ride_id = rides.id
    GROUP BY rides.id
