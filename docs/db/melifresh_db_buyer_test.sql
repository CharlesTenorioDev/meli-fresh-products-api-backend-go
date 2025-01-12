-- DDL
DROP DATABASE IF EXISTS `melifresh_db_buyer_test`;

CREATE DATABASE `melifresh_db_buyer_test`;

USE `melifresh_db_buyer_test`;

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

-- table `product_records`
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

-- DML
INSERT INTO buyers (card_number_id, first_name, last_name)
VALUES ('B1001', 'Alice', 'Brown'),
       ('B1002', 'Mark', 'Jones'),
       ('B1003', 'Linda', 'Garcia');

INSERT INTO products (product_code, description, height, lenght, width, weight, expiration_rate,
                      freezing_rate, recommended_freezing_temperature, seller_id, product_type_id)
VALUES ('P1001', 'Product 1', 10, 5, 8, 2, 0.1, 0.2, -5, 1, 1),
       ('P1002', 'Product 2', 12, 6, 9, 2.5, 0.15, 0.25, -6, 2, 2);

INSERT INTO product_records (id, last_update_date, purchase_price, sale_price, product_id)
VALUES (1, '2025-01-01 10:00:00', 50.00, 70.00, 1),
(2, '2025-01-02 11:30:00', 30.00, 45.00, 2);

INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id, product_record_id)
VALUES  ('PO1001', '2021-01-01', 'T1001', 1, 1),
        ('PO1002', '2021-01-02', 'T1002', 2, 2);