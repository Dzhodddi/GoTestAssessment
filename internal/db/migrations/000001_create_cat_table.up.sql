CREATE TABLE IF NOT EXISTS spycat(
    id bigserial PRIMARY KEY,
    name varchar(255) NOT NULL,
    years int NOT NULL,
    breed text NOT NULL,
    salary int NOT NULL
);