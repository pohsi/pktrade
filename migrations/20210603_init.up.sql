CREATE TABLE player
(
    id         SERIAL 		PRIMARY KEY,
    name       VARCHAR 		NOT NULL,
	password   VARCHAR 		NOT NULL
);

CREATE TABLE card_order
(
    id         SERIAL	  	PRIMARY KEY,
    owner_id   INTEGER      NOT NULL,
	owner_name VARCHAR 		NOT NULL,
	created_at TIMESTAMP 	NOT NULL,
	card_type  INTEGER 		NOT NULL,
	price      FLOAT 		NOT NULL
);

ALTER TABLE card_order ADD FOREIGN KEY (owner_id) REFERENCES player(id);

CREATE TABLE record
(
	id         SERIAL 		PRIMARY KEY,
	from_user  VARCHAR 		NOT NULL,
	to_user    VARCHAR 		NOT NULL,
	created_at TIMESTAMP 	NOT NULL,
	card_type  INTEGER		NOT NULL,
	price      FLOAT 		NOT NULL
);