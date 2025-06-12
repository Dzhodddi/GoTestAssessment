CREATE TABLE IF NOT EXISTS missions (
    id bigserial PRIMARY KEY,
    cat_id integer REFERENCES spycat(id),
    completed BOOLEAN NOT NULL DEFAULT FALSE
);