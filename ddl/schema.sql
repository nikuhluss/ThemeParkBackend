-- ================================================================
-- Team 14 - Theme Park
--
-- Conventions:
--'id' in column names is lower-case.
-- Everything else in column names is lower-case.
-- id(s) are varchar(64).
-- description(s) are varchar(128).
-- all contraints (including primary key) at the end of CREATE TABLE statement.
-- prices/costs are of type numeric(10, 2).
-- *keywords* are upper-case.
-- *datatypes* are lower-case.
--
-- formatted using: https://github.com/darold/pgFormatter
-- ================================================================

-- Users, and Customers
-- --------------------------------
-- Section that focuses on users, details, and customers.

SET search_path TO theme_park;

CREATE TABLE users (
    id varchar(64) NOT NULL,
    username varchar(32) NOT NULL,
    email varchar(64) NOT NULL,
    password_salt varchar(32) NOT NULL,
    password_hash varchar(32) NOT NULL,
    registered_on timestamp NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (username),
    UNIQUE (email)
);

CREATE TABLE genders (
    id varchar(64) NOT NULL,
    gender varchar(16) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (gender)
);

CREATE TABLE user_details (
    user_id varchar(64) NOT NULL,
    gender_id varchar(64),
    first_name varchar(32),
    last_name varchar(32),
    date_of_birth date,
    phone varchar(16),
    address varchar(32),
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE customers (
    user_id varchar(64) NOT NULL,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Pictures
-- --------------------------------
-- Section that focuses on picture collections.

CREATE TABLE pictures (
    id varchar(64) NOT NULL,
    format varchar(16) NOT NULL,
    blob bytea NOT NULL,
    PRIMARY KEY (id),
    CHECK (format <> '')
);

CREATE TABLE picture_collections (
    id varchar(64) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE pictures_in_collection (
    picture_id varchar(64) NOT NULL,
    collection_id varchar(64) NOT NULL,
    picture_sequence integer,
    PRIMARY KEY (picture_id, collection_id),
    FOREIGN KEY (picture_id) REFERENCES pictures (id),
    FOREIGN KEY (collection_id) REFERENCES picture_collections (id)
);

-- Shop, Transactions
-- --------------------------------
-- Section that focuses on the shops.

CREATE TABLE shop_types (
    id varchar(64) NOT NULL,
    shop_type varchar(32) NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE shops (
    id varchar(64) NOT NULL,
    shop_type_id varchar(64) NOT NULL,
    picture_collection_id varchar(64),
    name varchar(32) NOT NULL,
    description varchar(128),
    PRIMARY KEY (id),
    FOREIGN KEY (shop_type_id) REFERENCES shop_types (id),
    FOREIGN KEY (picture_collection_id) REFERENCES picture_collections (id)
);

CREATE TABLE items_types (
    id varchar(64) NOT NULL,
    picture_collection_id varchar(64),
    name varchar(32) NOT NULL,
    description varchar(128),
    PRIMARY KEY (id),
    FOREIGN KEY (picture_collection_id) REFERENCES picture_collections (id)
);

CREATE TABLE items_in_shop (
    id varchar(64) NOT NULL,
    item_id varchar(64) NOT NULL,
    shop_id varchar(64) NOT NULL,
    price numeric(10, 2) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (id, item_id, shop_id),
    FOREIGN KEY (item_id) REFERENCES items_types (id),
    FOREIGN KEY (shop_id) REFERENCES shops (id),
    CHECK (price >= 0)
);

CREATE TABLE transactions (
    id varchar(64) NOT NULL,
    customer_id varchar(64) NOT NULL,
    item_in_shop_id varchar(64) NOT NULL,
    transaction_datetime timestamp NOT NULL,
    purchase_reference varchar(64) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (id, customer_id, item_in_shop_id),
    FOREIGN KEY (customer_id) REFERENCES customers (user_id),
    FOREIGN KEY (item_in_shop_id) REFERENCES items_in_shop (id)
);

-- Rides, and Tickets
-- --------------------------------
-- Section that focuses on rides, ticket purchase, and ticket usage.

CREATE TABLE rides (
    id varchar(64) NOT NULL,
    picture_collection_id varchar(64),
    name varchar(32) NOT NULL,
    description varchar(128),
    min_age integer,
    min_height integer,
    longitude real,
    latitude real,
    PRIMARY KEY (id),
    UNIQUE (name),
    CHECK (name <> ''),
    CHECK (longitude >= - 180 AND longitude <= 180),
    CHECK (latitude >= - 90 AND latitude <= 90),
    CHECK (min_age >= 0),
    CHECK (min_height >= 0)
);

CREATE TABLE reviews (
    id varchar(64) NOT NULL,
    ride_id varchar(64) NOT NULL,
    customer_id varchar(64) NOT NULL,
    rating integer NOT NULL,
    title text,
    content text,
    posted_on timestamp NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (id, ride_id, customer_id),
    FOREIGN KEY (ride_id) REFERENCES rides (id),
    FOREIGN KEY (customer_id) REFERENCES customers (user_id),
    CHECK (rating >= 1 AND rating <= 5)
);

CREATE TABLE tickets (
    id varchar(64) NOT NULL,
    user_id varchar(64) NOT NULL,
    is_kid boolean DEFAULT TRUE NOT NULL,
    puchase_price numeric(10, 2) NOT NULL,
    puchased_on timestamp NOT NULL,
    purchase_reference varchar(64) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES customers (user_id)
);

CREATE TABLE tickets_on_rides (
    id varchar(64) NOT NULL,
    ride_id varchar(64) NOT NULL,
    ticket_id varchar(64) NOT NULL,
    scan_datetime timestamp NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (ride_id) REFERENCES rides (id),
    FOREIGN KEY (ticket_id) REFERENCES tickets (id)
);

-- Employees
-- --------------------------------
-- Section that focuses on employees and their schedule on rides.

CREATE TABLE roles (
    id varchar(64) NOT NULL,
    role varchar(16) NOT NULL,
    hourly_rate numeric(10, 2) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (ROLE),
    CHECK (hourly_rate >= 0)
);

CREATE TABLE employees (
    id varchar(64) NOT NULL,
    user_id varchar(64) NOT NULL,
    role_id varchar(64) NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TABLE employees_on_rides (
    id varchar(64) NOT NULL,
    ride_id varchar(64) NOT NULL,
    employee_id varchar(64) NOT NULL,
    start_datetime timestamp NOT NULL,
    end_datetime timestamp,
    PRIMARY KEY (id),
    FOREIGN KEY (ride_id) REFERENCES rides (id),
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    CHECK (start_datetime <= end_datetime)
);

-- Maintenance
-- --------------------------------
-- Section that focuses on rides maintenance.

CREATE TABLE maintenance_types (
    id varchar(64) NOT NULL,
    maintenance_type varchar(16) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (maintenance_type)
);

CREATE TABLE rides_maintenance (
    id varchar(64) NOT NULL,
    ride_id varchar(64) NOT NULL,
    maintenance_type_id varchar(64) NOT NULL,
    description varchar(128),
    cost numeric(10, 2),
    start_datetime timestamp NOT NULL,
    end_datetime timestamp,
    PRIMARY KEY (id, ride_id),
    UNIQUE (id),
    FOREIGN KEY (ride_id) REFERENCES rides (id),
    FOREIGN KEY (maintenance_type_id) REFERENCES maintenance_types (id),
    CHECK (start_datetime <= end_datetime)
);

CREATE TABLE employees_on_maintenance (
    id varchar(64) NOT NULL,
    maintenance_id varchar(64) NOT NULL,
    employee_id varchar(64) NOT NULL,
    start_datetime timestamp NOT NULL,
    end_datetime timestamp,
    PRIMARY KEY (id, maintenance_id, employee_id),
    FOREIGN KEY (maintenance_id) REFERENCES rides_maintenance (id),
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    CHECK (start_datetime <= end_datetime)
);

-- Events
-- --------------------------------
-- Section that focuses on the event board (rainouts, specials, etc).

CREATE TABLE event_types (
    id varchar(64) NOT NULL,
    event_type varchar(32) NOT NULL,
    PRIMARY KEY (id),
    UNIQUE (event_type)
);

CREATE TABLE events (
    id varchar(64) NOT NULL,
    employee_id varchar(64) NOT NULL,
    event_type_id varchar(64) NOT NULL,
    title varchar(32) NOT NULL,
    description varchar(128) NOT NULL,
    posted_on timestamp NOT NULL,
    PRIMARY KEY (id),
    FOREIGN KEY (employee_id) REFERENCES employees (id),
    FOREIGN KEY (event_type_id) REFERENCES event_types (id)
);
