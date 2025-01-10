-- DDL
DROP DATABASE IF EXISTS `melifresh_purchase_orders_test_db`;

CREATE DATABASE IF NOT EXISTS `melifresh_purchase_orders_test_db`;

USE `melifresh_purchase_orders_test_db`;

-- table `products`
CREATE TABLE `products`
(
    `id`                               int(11) NOT NULL AUTO_INCREMENT,
    `product_code`                     varchar(25) NOT NULL,
    `description`                      text        NOT NULL,
    `height`                           float       NOT NULL,
    `lenght`                           float       NOT NULL,
    `width`                            float       NOT NULL,
    `weight`                           float       NOT NULL,
    `expiration_rate`                  float       NOT NULL,
    `freezing_rate`                    float       NOT NULL,
    `recommended_freezing_temperature` float       NOT NULL,
    `seller_id`                        int(11) NOT NULL,
    `product_type_id`                  int(11) NOT NULL,
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

CREATE TABLE product_records
(
    id                int(11) NOT NULL AUTO_INCREMENT,
	last_update_date  datetime NOT NULL,
	purchase_price    float NOT NULL,
	sale_price        float NOT NULL,
	product_id        int NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products (id),
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
    FOREIGN KEY (`buyer_id`) REFERENCES buyers (id) ON DELETE SET NULL,
    FOREIGN KEY (`product_record_id`) REFERENCES product_records (id) ON DELETE SET NULL,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;


INSERT INTO products (product_code, description, height, lenght, width, weight, expiration_rate,
                      freezing_rate, recommended_freezing_temperature, seller_id, product_type_id)
VALUES ('P1001', 'Product 1', 10, 5, 8, 2, 0.1, 0.2, -5, 1, 1),
       ('P1002', 'Product 2', 12, 6, 9, 2.5, 0.15, 0.25, -6, 2, 2),
       ('P1003', 'Product 3', 14, 7, 10, 3, 0.2, 0.3, -7, 3, 3),
       ('P1004', 'Product 4', 16, 8, 11, 3.5, 0.25, 0.35, -8, 4, 4),
       ('P1005', 'Product 5', 18, 9, 12, 4, 0.3, 0.4, -9, 5, 5),
       ('P1006', 'Product 6', 20, 10, 13, 4.5, 0.35, 0.45, -10, 6, 6),
       ('P1007', 'Product 7', 22, 11, 14, 5, 0.4, 0.5, -11, 7, 7),
       ('P1008', 'Product 8', 24, 12, 15, 5.5, 0.45, 0.55, -12, 8, 8),
       ('P1009', 'Product 9', 26, 13, 16, 6, 0.5, 0.6, -13, 9, 9),
       ('P1010', 'Product 10', 28, 14, 17, 6.5, 0.55, 0.65, -14, 10, 10);

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