-- ================================================================
-- Team 14 - Theme Park
--
-- Conventions:
--'ID' in column names is upper-case.
-- Everything else in column names is lower-case.
-- ID(s) are varchar(32).
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
    ID varchar(32) NOT NULL,
    username varchar(32) NOT NULL,
    email varchar(64) NOT NULL,
    password_salt varchar(32) NOT NULL,
    password_hash varchar(32) NOT NULL,
    registered_on timestamp NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE (username),
    UNIQUE (email)
);

CREATE TABLE genders (
    ID varchar(32) NOT NULL,
    gender varchar(16) NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE (gender)
);

CREATE TABLE user_details (
    user_ID varchar(32) NOT NULL,
    gender_ID varchar(32),
    first_name varchar(32),
    last_name varchar(32),
    date_of_birth date,
    phone varchar(16),
    address varchar(32),
    PRIMARY KEY (user_ID),
    FOREIGN KEY (user_ID) REFERENCES users (ID)
);

CREATE TABLE customers (
    user_ID varchar(32) NOT NULL,
    PRIMARY KEY (user_ID),
    FOREIGN KEY (user_ID) REFERENCES users (ID)
);

-- Pictures
-- --------------------------------
-- Section that focuses on picture collections.

CREATE TABLE pictures (
    ID varchar(32) NOT NULL,
    format varchar(16) NOT NULL,
    blob bytea NOT NULL,
    PRIMARY KEY (ID)
);

CREATE TABLE picture_collections (
    ID varchar(32) NOT NULL,
    PRIMARY KEY (ID)
);

CREATE TABLE pictures_in_collection (
    picture_ID varchar(32) NOT NULL,
    collection_ID varchar(32) NOT NULL,
    picture_sequence integer,
    PRIMARY KEY (picture_ID, collection_ID),
    FOREIGN KEY (picture_ID) REFERENCES pictures (ID),
    FOREIGN KEY (collection_ID) REFERENCES picture_collections (ID)
);

-- Shop, Transactions
-- --------------------------------
-- Section that focuses on the shops.

CREATE TABLE shop_types (
    ID varchar(32) NOT NULL,
    shop_type varchar(32) NOT NULL,
    PRIMARY KEY (ID)
);

CREATE TABLE shops (
    ID varchar(32) NOT NULL,
    shop_type_ID varchar(32) NOT NULL,
    picture_collection_ID varchar(32),
    name varchar(32) NOT NULL,
    description varchar(128),
    PRIMARY KEY (ID),
    FOREIGN KEY (shop_type_ID) REFERENCES shop_types (ID),
    FOREIGN KEY (picture_collection_ID) REFERENCES picture_collections (ID)
);

CREATE TABLE items_types (
    ID varchar(32) NOT NULL,
    picture_collection_ID varchar(32),
    name varchar(32) NOT NULL,
    description varchar(128),
    PRIMARY KEY (ID),
    FOREIGN KEY (picture_collection_ID) REFERENCES picture_collections (ID)
);

CREATE TABLE items_in_shop (
    ID varchar(32) NOT NULL,
    item_ID varchar(32) NOT NULL,
    shop_ID varchar(32) NOT NULL,
    price numeric(10, 2) NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE (ID, item_ID, shop_ID),
    FOREIGN KEY (item_ID) REFERENCES items_types (ID),
    FOREIGN KEY (shop_ID) REFERENCES shops (ID),
    CHECK (price >= 0)
);

CREATE TABLE transactions (
    ID varchar(32) NOT NULL,
    customer_ID varchar(32) NOT NULL,
    item_in_shop_ID varchar(32) NOT NULL,
    transaction_datetime timestamp NOT NULL,
    purchase_reference varchar(64) NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE (ID, customer_ID, item_in_shop_ID),
    FOREIGN KEY (customer_ID) REFERENCES customers (user_ID),
    FOREIGN KEY (item_in_shop_ID) REFERENCES items_in_shop (ID)
);

-- Rides, and Tickets
-- --------------------------------
-- Section that focuses on rides, ticket purchase, and ticket usage.

CREATE TABLE rides (
    ID varchar(32) NOT NULL,
    picture_collection_ID varchar(32),
    name varchar(32) NOT NULL,
    description varchar(128),
    min_age integer,
    min_height integer,
    longitude integer,
    latitude integer,
    PRIMARY KEY (ID),
    UNIQUE (name),
    CHECK (longitude >= - 180 AND longitude <= 180),
    CHECK (latitude >= - 90 AND latitude <= 90),
    CHECK (min_age >= 0),
    CHECK (min_height >= 0)
);

CREATE TABLE reviews (
    ID varchar(32) NOT NULL,
    ride_ID varchar(32) NOT NULL,
    customer_ID varchar(32) NOT NULL,
    rating integer NOT NULL,
    title text,
    content text,
    PRIMARY KEY (ID),
    UNIQUE (ID, ride_ID, customer_ID),
    FOREIGN KEY (ride_ID) REFERENCES rides (ID),
    FOREIGN KEY (customer_ID) REFERENCES customers (user_ID),
    CHECK (rating >= 1 AND rating <= 5)
);

CREATE TABLE tickets (
    ID varchar(32) NOT NULL,
    user_ID varchar(32) NOT NULL,
    is_kid boolean DEFAULT TRUE NOT NULL,
    puchase_price numeric(10, 2) NOT NULL,
    puchased_on timestamp NOT NULL,
    purchase_reference varchar(64) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (user_ID) REFERENCES customers (user_ID)
);

CREATE TABLE tickets_on_rides (
    ID varchar(32) NOT NULL,
    ride_ID varchar(32) NOT NULL,
    ticket_ID varchar(32) NOT NULL,
    scan_datetime timestamp NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (ride_ID) REFERENCES rides (ID),
    FOREIGN KEY (ticket_ID) REFERENCES tickets (ID)
);

-- Employees
-- --------------------------------
-- Section that focuses on employees and their schedule on rides.

CREATE TABLE roles (
    ID varchar(32) NOT NULL,
    role varchar(16) NOT NULL,
    hourly_rate numeric(10, 2) NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE (ROLE),
    CHECK (hourly_rate >= 0)
);

CREATE TABLE employees (
    ID varchar(32) NOT NULL,
    user_ID varchar(32) NOT NULL,
    role_ID varchar(32) NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (user_ID) REFERENCES users (ID),
    FOREIGN KEY (role_ID) REFERENCES roles (ID)
);

CREATE TABLE employees_on_rides (
    ID varchar(32) NOT NULL,
    ride_ID varchar(32) NOT NULL,
    employee_ID varchar(32) NOT NULL,
    start_datetime timestamp NOT NULL,
    end_datetime timestamp,
    PRIMARY KEY (ID),
    FOREIGN KEY (ride_ID) REFERENCES rides (ID),
    FOREIGN KEY (employee_ID) REFERENCES employees (ID),
    CHECK (start_datetime <= end_datetime)
);

-- Maintenance
-- --------------------------------
-- Section that focuses on rides maintenance.

CREATE TABLE maintenance_types (
    ID varchar(32) NOT NULL,
    maintenance_type varchar(16) NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE (maintenance_type)
);

CREATE TABLE rides_maintenance (
    ID varchar(32) NOT NULL,
    ride_ID varchar(32) NOT NULL,
    maintenance_type_ID varchar(32) NOT NULL,
    description varchar(128),
    cost numeric(10, 2),
    start_datetime timestamp NOT NULL,
    end_datetime timestamp,
    PRIMARY KEY (ID, ride_ID),
    UNIQUE (ID),
    FOREIGN KEY (ride_ID) REFERENCES rides (ID),
    FOREIGN KEY (maintenance_type_ID) REFERENCES maintenance_types (ID),
    CHECK (start_datetime <= end_datetime)
);

CREATE TABLE employees_on_maintenance (
    ID varchar(32) NOT NULL,
    maintenance_ID varchar(32) NOT NULL,
    employee_ID varchar(32) NOT NULL,
    start_datetime timestamp NOT NULL,
    end_datetime timestamp,
    PRIMARY KEY (ID, maintenance_ID, employee_ID),
    FOREIGN KEY (maintenance_ID) REFERENCES rides_maintenance (ID),
    FOREIGN KEY (employee_ID) REFERENCES employees (ID),
    CHECK (start_datetime <= end_datetime)
);

-- Events
-- --------------------------------
-- Section that focuses on the event board (rainouts, specials, etc).

CREATE TABLE event_types (
    ID varchar(32) NOT NULL,
    event_type varchar(32) NOT NULL,
    PRIMARY KEY (ID),
    UNIQUE (event_type)
);

CREATE TABLE events (
    ID varchar(32) NOT NULL,
    employee_ID varchar(32) NOT NULL,
    event_type_ID varchar(32) NOT NULL,
    title varchar(32) NOT NULL,
    description varchar(128) NOT NULL,
    posted_on timestamp NOT NULL,
    PRIMARY KEY (ID),
    FOREIGN KEY (employee_ID) REFERENCES employees (ID),
    FOREIGN KEY (event_type_ID) REFERENCES event_types (ID)
);
