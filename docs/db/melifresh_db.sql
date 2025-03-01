-- DDL
DROP DATABASE IF EXISTS `melifresh`;

CREATE DATABASE `melifresh`;

USE `melifresh`;

-- table `localities`
CREATE TABLE `localities`
(
    `id` int(11) NOT NULL,
    `name` varchar(255) NOT NULL,
    `province_name` varchar(255) NOT NULL,
    `country_name` varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `sellers`
CREATE TABLE `sellers`
(
    `id`           int(11) NOT NULL AUTO_INCREMENT,
    `cid`          int(11) NOT NULL,
    `company_name` varchar(255) NOT NULL,
    `address`      varchar(255) NOT NULL,
    `telephone`    varchar(15)  NOT NULL,
    `locality_id`  int(11) NOT NULL,
    FOREIGN KEY (`locality_id`) REFERENCES localities (id) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `warehouses`
CREATE TABLE `warehouses`
(
    `id`                  int(11) NOT NULL AUTO_INCREMENT,
    `warehouse_code`      varchar(25)  NOT NULL,
    `address`             varchar(255) NOT NULL,
    `telephone`           varchar(15)  NOT NULL,
    `minimum_capacity`    int          NOT NULL,
    `minimum_temperature` float        NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `product_type`
CREATE TABLE `product_type`
(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `description` varchar(255) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `sections`
CREATE TABLE `sections`
(
    `id`                  int(11) NOT NULL AUTO_INCREMENT,
    `section_number`      int(11) NOT NULL,
    `current_temperature` decimal(19, 2) NOT NULL,
    `minimum_temperature` decimal(19, 2) NOT NULL,
    `current_capacity`    int   NOT NULL,
    `minimum_capacity`    int   NOT NULL,
    `maximum_capacity`    int   NOT NULL,
    `warehouse_id`        int(11) NOT NULL,
    `product_type_id`     int(11) NOT NULL,
    FOREIGN KEY (`warehouse_id`) REFERENCES warehouses(id) ON DELETE CASCADE,
    FOREIGN KEY (`product_type_id`) REFERENCES product_type(id) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `products`
CREATE TABLE `products`
(
    id                               int(11) NOT NULL AUTO_INCREMENT,
    product_code                     varchar(25) NOT NULL,
    description                      varchar(25) NOT NULL,
    height                           float       NOT NULL,
    length                           float       NOT NULL,
    net_weight                        float       NOT NULL,
    expiration_rate                  float       NOT NULL,
    freezing_rate                    float       NOT NULL,
    recommended_freezing_temperature float       NOT NULL,
    width                            float       NOT NULL,
    seller_id                        int(11) NOT NULL,
    product_type_id                  int(11) NOT NULL,
    FOREIGN KEY (`seller_id`) REFERENCES sellers(id) ON DELETE CASCADE,
    FOREIGN KEY (`product_type_id`) REFERENCES product_type(id) ON DELETE CASCADE,
    PRIMARY KEY (id)

) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `employees`
CREATE TABLE `employees`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT,
    `card_number_id` varchar(25) NOT NULL,
    `first_name`     varchar(50) NOT NULL,
    `last_name`      varchar(50) NOT NULL,
    `warehouse_id`   int(11) NOT NULL,
    FOREIGN KEY (`warehouse_id`) REFERENCES warehouses(id) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `buyers`
CREATE TABLE `buyers`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT,
    `card_number_id` varchar(25) NOT NULL,
    `first_name`     varchar(50) NOT NULL,
    `last_name`      varchar(50) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `product_records`
CREATE TABLE product_records
(
    id                int(11) NOT NULL AUTO_INCREMENT,
	last_update_date  datetime  NOT NULL,
	purchase_price    float NOT NULL,
	sale_price        float NOT NULL,
	product_id        int NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE,
	PRIMARY KEY (id)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `purchase_orders`
CREATE TABLE `purchase_orders`
(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `order_number` varchar(255)  NOT NULL,
    `order_date` date NOT NULL,
    `tracking_code` varchar(255) NOT NULL,
    `buyer_id` int(11) NULL,
    `product_record_id` int(11) NULL,
    FOREIGN KEY (`buyer_id`) REFERENCES buyers (id) ON DELETE CASCADE,
    FOREIGN KEY (`product_record_id`) REFERENCES product_records (id) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `carries`
CREATE TABLE `carries`
(
    `id`             int(11) NOT NULL AUTO_INCREMENT,
	`cid`            varchar(10) UNIQUE NOT NULL,
	`company_name`   varchar(100) NOT NULL,
	`address`        varchar(100) NOT NULL,
	`phone_number`   varchar(20) NOT NULL,
	`locality_id`    int(11) NOT NULL,
    FOREIGN KEY (`locality_id`) REFERENCES localities (id) ON DELETE CASCADE,
	PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

-- table `product_batches`
CREATE TABLE `product_batches` (
    `id`                 INT(11) NOT NULL AUTO_INCREMENT,
    `batch_number`      INT(11) NOT NULL,
    `current_quantity`   INT(11) NOT NULL,
    `current_temperature` FLOAT NOT NULL,
    `due_date`          DATE NOT NULL,     
    `initial_quantity`  INT(11) NOT NULL,
    `manufacturing_date` DATE NOT NULL,                   
    `manufacturing_hour` INT(11) NOT NULL,                  
    `minumum_temperature` FLOAT NOT NULL,
    `product_id`        INT(11) NOT NULL,
    `section_id`        INT(11) NOT NULL,
    FOREIGN KEY (`product_id`) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (`section_id`) REFERENCES sections(id) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- table `inbound_orders`
CREATE TABLE `inbound_orders`
(
    `id`               int(11) NOT NULL AUTO_INCREMENT,
    `order_date`       date NOT NULL,
    `order_number`     varchar(255) UNIQUE NOT NULL,
    `employee_id`      int(11) NOT NULL,
    `product_batch_id` int(11) NOT NULL,
    `warehouse_id`     int(11) NOT NULL,
    FOREIGN KEY (`employee_id`) REFERENCES employees (id) ON DELETE CASCADE,
    FOREIGN KEY (`product_batch_id`) REFERENCES product_batches (id) ON DELETE CASCADE,
    FOREIGN KEY (`warehouse_id`) REFERENCES warehouses (id) ON DELETE CASCADE,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;


-- DML
INSERT INTO localities (id, name, province_name, country_name)
VALUES (1, 'New York City', 'New York', 'United States'),
       (2, 'Los Angeles', 'California', 'United States'),
       (3, 'Chicago', 'Illinois', 'United States'),
       (4, 'Houston', 'Texas', 'United States'),
       (5, 'Phoenix', 'Arizona', 'United States'),
       (6, 'Philadelphia', 'Pennsylvania', 'United States'),
       (7, 'San Antonio', 'Texas', 'United States'),
       (8, 'San Diego', 'California', 'United States'),
       (9, 'Dallas', 'Texas', 'United States'),
       (10, 'San Jose', 'California', 'United States');

INSERT INTO sellers (cid, company_name, address, telephone, locality_id)
VALUES (1, 'Company A', '123 Main St', '123-456-7890', 1),
       (2, 'Company B', '456 Elm St', '123-456-7891', 2),
       (3, 'Company C', '789 Oak St', '123-456-7892', 3),
       (4, 'Company D', '101 Pine St', '123-456-7893', 4),
       (5, 'Company E', '102 Maple St', '123-456-7894', 5),
       (6, 'Company F', '103 Cedar St', '123-456-7895', 6),
       (7, 'Company G', '104 Birch St', '123-456-7896', 7),
       (8, 'Company H', '105 Willow St', '123-456-7897', 8),
       (9, 'Company I', '106 Cherry St', '123-456-7898', 9),
       (10, 'Company J', '107 Walnut St', '123-456-7899', 10);

INSERT INTO warehouses (warehouse_code, address, telephone, minimum_capacity, minimum_temperature)
VALUES ('WH01', '200 Warehouse Rd', '234-567-8901', 100, 0),
       ('WH02', '201 Warehouse Ln', '234-567-8902', 150, -5),
       ('WH03', '202 Storage Blvd', '234-567-8903', 120, 2),
       ('WH04', '203 Distribution Ave', '234-567-8904', 200, -2),
       ('WH05', '204 Inventory St', '234-567-8905', 180, 0),
       ('WH06', '205 Logistics Way', '234-567-8906', 160, -3),
       ('WH07', '206 Depot Dr', '234-567-8907', 140, 1),
       ('WH08', '207 Supply Ct', '234-567-8908', 170, -4),
       ('WH09', '208 Goods Rd', '234-567-8909', 130, 3),
       ('WH10', '209 Freight St', '234-567-8910', 190, -1);

INSERT INTO product_type (description)
VALUES  ('Dairy'),
        ('Meat'),
        ('Vegetables'),
        ('Fruits'),
        ('Bakery'),
        ('Seafood'),
        ('Beverages'),
        ('Snacks'),
        ('Condiments'),
        ('Frozen Foods');

INSERT INTO sections (section_number, current_temperature, minimum_temperature, current_capacity,
                      minimum_capacity, maximum_capacity, warehouse_id, product_type_id)
VALUES (1, 0, -5, 50, 20, 100, 1, 1),
       (2, -2, -6, 60, 30, 110, 2, 2),
       (3, 1, -4, 70, 40, 120, 3, 3),
       (4, -3, -7, 80, 50, 130, 4, 4),
       (5, 2, -5, 90, 60, 140, 5, 5),
       (6, -4, -8, 100, 70, 150, 6, 6),
       (7, 3, -6, 110, 80, 160, 7, 7),
       (8, -5, -9, 120, 90, 170, 8, 8),
       (9, 4, -7, 130, 100, 180, 9, 9),
       (10, -6, -10, 140, 110, 190, 10, 10);

INSERT INTO products (product_code, description, height, length, width, net_weight, expiration_rate,
                      freezing_rate, recommended_freezing_temperature, seller_id, product_type_id)
VALUES 
('P1001', 'Product 1', 10, 5, 8, 2, 0.1, 0.2, -5, 1, 1),
('P1002', 'Product 2', 12, 6, 9, 2.5, 0.15, 0.25, -6, 2, 2),
('P1003', 'Product 3', 14, 7, 10, 3, 0.2, 0.3, -7, 3, 3),
('P1004', 'Product 4', 16, 8, 11, 3.5, 0.25, 0.35, -8, 4, 4),
('P1005', 'Product 5', 18, 9, 12, 4, 0.3, 0.4, -9, 5, 5),
('P1006', 'Product 6', 20, 10, 13, 4.5, 0.35, 0.45, -10, 6, 6),
('P1007', 'Product 7', 22, 11, 14, 5, 0.4, 0.5, -11, 7, 7),
('P1008', 'Product 8', 24, 12, 15, 5.5, 0.45, 0.55, -12, 8, 8),
('P1009', 'Product 9', 26, 13, 16, 6, 0.5, 0.6, -13, 9, 9),
('P1010', 'Product 10', 28, 14, 17, 6.5, 0.55, 0.65, -14, 10, 10);

INSERT INTO employees (card_number_id, first_name, last_name, warehouse_id)
VALUES ('E1001', 'John', 'Doe', 1),
       ('E1002', 'Jane', 'Smith', 2),
       ('E1003', 'Michael', 'Johnson', 3),
       ('E1004', 'Emily', 'Davis', 4),
       ('E1005', 'David', 'Miller', 5),
       ('E1006', 'Sarah', 'Wilson', 6),
       ('E1007', 'Robert', 'Moore', 7),
       ('E1008', 'Jennifer', 'Taylor', 8),
       ('E1009', 'William', 'Anderson', 9),
       ('E1010', 'Jessica', 'Thomas', 10);

INSERT INTO buyers (card_number_id, first_name, last_name)
VALUES ('B1001', 'Alice', 'Brown'),
       ('B1002', 'Mark', 'Jones'),
       ('B1003', 'Linda', 'Garcia'),
       ('B1004', 'Brian', 'Williams'),
       ('B1005', 'Susan', 'Martinez'),
       ('B1006', 'Richard', 'Lee'),
       ('B1007', 'Karen', 'Harris'),
       ('B1008', 'Steven', 'Clark'),
       ('B1009', 'Betty', 'Lopez'),
       ('B1010', 'Edward', 'Gonzalez');

INSERT INTO carries (cid, company_name, address, phone_number, locality_id)
VALUES  (1, 'Meli Fresh Logistics', '123 Fresh St', '555-1001', 1),
        (2, 'Quick Delivery Services', '456 Fast Ave', '555-1002', 2),
        (3, 'Fresh Express', '789 Speed Blvd', '555-1003', 3),
        (4, 'Swift Transport Co.', '101 Pine St', '555-1004', 4),
        (5, 'Rapid Freight Solutions', '202 Oak Dr', '555-1005', 5);

INSERT INTO product_records (id, last_update_date, purchase_price, sale_price, product_id)
VALUES (1, '2025-01-01 10:00:00', 50.00, 70.00, 1),
(2, '2025-01-02 11:30:00', 30.00, 45.00, 2),
(3, '2025-01-03 14:45:00', 100.00, 150.00, 3),
(4, '2025-01-04 09:15:00', 20.00, 35.00, 4),
(5, '2025-01-05 16:00:00', 75.00, 110.00, 5);

INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id)
VALUES  ('PO1001', '2021-01-01', 'T1001', 1, 1),
        ('PO1002', '2021-01-02', 'T1002', 2, 2),
        ('PO1003', '2021-01-03', 'T1003', 3, 3),
        ('PO1004', '2021-01-04', 'T1004', 4, 4),
        ('PO1005', '2021-01-05', 'T1005', 5, 5);

INSERT INTO product_batches (batch_number, current_quantity, current_temperature, due_date, initial_quantity, manufacturing_date, manufacturing_hour, minumum_temperature, product_id, section_id)
VALUES  (1, 100, 20.0, '2022-01-08', 150, '2022-01-01', 10, -5.0, 1, 1),
        (2, 200, 18.5, '2022-02-04', 250, '2022-01-02', 11, -4.0, 2, 1),
        (3, 150, 15.0, '2022-03-01', 180, '2022-01-03', 12, -3.0, 1, 2),
        (4, 300, 22.0, '2022-03-15', 350, '2022-01-04', 9, -6.0, 2, 2),
        (5, 250, 25.0, '2022-04-10', 300, '2022-01-05', 8, -2.0, 1, 3),
        (6, 400, 30.0, '2022-05-05', 450, '2022-01-06', 7, -1.0, 3, 2),
        (7, 500, 5.5,  '2022-06-01', 600, '2022-01-07', 6, -5.0, 3, 3),
        (8, 600, 10.2, '2022-06-15', 550, '2022-01-08', 5, -4.1, 4, 4),
        (9, 350, 12.3, '2022-07-01', 400, '2022-01-09', 4, -3.2, 5, 1),
        (10, 450, 16.4, '2022-07-15', 250, '2022-01-10', 3, -7.5, 5, 3);

INSERT INTO inbound_orders (order_date, order_number, employee_id, product_batch_id, warehouse_id)
VALUES ('2025-01-01', 'ORD001', 1, 1, 1),
       ('2025-01-02', 'ORD002', 2, 2, 2),
       ('2025-01-03', 'ORD003', 3, 3, 3),
       ('2025-01-04', 'ORD004', 4, 4, 4),
       ('2025-01-05', 'ORD005', 5, 5, 5),
       ('2025-01-06', 'ORD006', 6, 6, 6),
       ('2025-01-07', 'ORD007', 7, 7, 7),
       ('2025-01-08', 'ORD008', 8, 8, 8),
       ('2025-01-09', 'ORD009', 9, 9, 9),
       ('2025-01-10', 'ORD010', 10, 10, 10);