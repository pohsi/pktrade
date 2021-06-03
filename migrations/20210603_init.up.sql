CREATE TABLE player
(
    id         SERIAL 		PRIMARY KEY,
    name       VARCHAR 		NOT NULL,
	password   VARCHAR 		NOT NULL
);

CREATE TABLE card_purchase
(
    id         SERIAL	  	PRIMARY KEY,
    owner_id   INTEGER      NOT NULL,
	owner_name VARCHAR 		NOT NULL,
	created_at TIMESTAMP 	NOT NULL,
	card_type  INTEGER 		NOT NULL,
	price      FLOAT 		NOT NULL
);

ALTER TABLE card_purchase ADD FOREIGN KEY (owner_id) REFERENCES player(id);

CREATE TABLE card_sell (LIKE card_purchase INCLUDING ALL);
ALTER TABLE card_sell ALTER id DROP DEFAULT;
CREATE SEQUENCE card_sell_id_seq;
INSERT INTO card_sell SELECT * FROM card_purchase;
SELECT setval('card_sell_id_seq', (SELECT max(id) FROM card_sell), true);
ALTER TABLE card_sell ALTER id SET DEFAULT nextval('card_sell_id_seq');
ALTER SEQUENCE card_sell_id_seq OWNED BY card_sell.id;

CREATE TABLE record
(
	id         SERIAL 		PRIMARY KEY,
	from_user  VARCHAR 		NOT NULL,
	to_user    VARCHAR 		NOT NULL,
	created_at TIMESTAMP 	NOT NULL,
	card_type  INTEGER		NOT NULL,
	price      FLOAT 		NOT NULL
);