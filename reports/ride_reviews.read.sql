SELECT rides.id, rides.name, COUNT(reviews.*), AVG(reviews.rating)
    FROM theme_park.rides
    LEFT JOIN theme_park.reviews ON reviews.ride_id = rides.id
    GROUP BY rides.id
