DO $$DECLARE i INTEGER;
BEGIN
       FOR i IN 1 .. 10000 LOOP
              INSERT INTO player(name, password)
              VALUES(CONCAT('name', i), 'pass');
       END LOOP;
END$$;

INSERT INTO card_purchase(owner_id, owner_name, created_at, card_type, price)
VALUES (31, 'user31', '2021-06-03 05:22:37'::TIMESTAMP, 1, 1.27),
       (78, 'user78', '2021-06-02 11:45:46'::TIMESTAMP, 2, 9.9),
       (76, 'user46', '2021-05-03 22:09:12'::TIMESTAMP, 2, 2.85);


INSERT INTO card_sell(owner_id, owner_name, created_at, card_type, price)
VALUES (338, 'user338', '2021-06-03 05:22:37'::TIMESTAMP, 4, 9.6),
       (15,  'user15',  '2021-06-02 11:45:46'::TIMESTAMP, 2, 5.57),
       (442, 'user442', '2021-05-03 22:09:12'::TIMESTAMP, 3, 1.45),
       (167, 'user167', '2021-05-03 22:09:12'::TIMESTAMP, 1, 3.00);

INSERT INTO record(from_user, to_user, created_at, card_type, price)
VALUES ('user61'  ,'user37',  '2021-05-11 14:41:17'::TIMESTAMP, 1, 6.1),
       ('user192' ,'user765', '2021-06-01 21:35:05'::TIMESTAMP, 2, 3.32),
       ('user37'  ,'user579', '2021-06-04 01:55:05'::TIMESTAMP, 4, 8.05),
       ('user3'   ,'user16',  '2021-05-23 01:12:31'::TIMESTAMP, 2, 4.73);



