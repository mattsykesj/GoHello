DROP TABLE user_foods CASCADE;

DROP TABLE users;

DROP TABLE foods;

CREATE TABLE foods (
	food_id bigserial not null PRIMARY KEY,
	name varchar(256) not null, 
	protein real not null,
	carbohydrate real not null,
	fat real not null,
	calories real not null
);

CREATE TABLE users (
	user_id bigserial not null PRIMARY KEY,
	user_name varchar(20) not null UNIQUE,
	target_calories real not null,
	target_protein real not null,
	target_carbohydrate real not null,
	target_fat real not null
);

CREATE TABLE user_foods (
	user_foods_id bigserial not null PRIMARY KEY,
	user_id bigserial not null REFERENCES users (user_id),
	food_id bigserial not null REFERENCES foods (food_id)
);

INSERT INTO users (user_name, target_calories, target_protein, target_carbohydrate, target_fat)
VALUES('matt', 1931, 120, 166, 87)

 