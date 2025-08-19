-- CREATE students TABLE
CREATE TABLE IF NOT EXISTS students (
    id         UUID      NOT NULL DEFAULT uuid_generate_v1mc(),
    name       TEXT      NOT NULL,
    email      TEXT      NOT NULL,
    birthday   DATE      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT students_pk PRIMARY KEY (id),
    CONSTRAINT students_email_un UNIQUE (email)
);
