-- CREATE courses TABLE
CREATE TABLE IF NOT EXISTS courses (
    id         UUID           NOT NULL DEFAULT uuid_generate_v1mc(),
    name       TEXT           NOT NULL,
    value      NUMERIC(19, 2) NOT NULL,
    created_at TIMESTAMP      NOT NULL DEFAULT NOW(),
    CONSTRAINT courses_name_un UNIQUE (name),
    CONSTRAINT courses_pk PRIMARY KEY (id)
);
