-- CREATE UUID EXTENSION
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- CREATE TYPES
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'ACCOUNT_STATUS') THEN
		CREATE TYPE ACCOUNT_STATUS AS ENUM ('ADIMPLENTE', 'INADIMPLENTE');
    END IF;
END;
$$ LANGUAGE plpgsql;

-- CREATE SCHEMA
CREATE TABLE IF NOT EXISTS students (
    id         UUID      NOT NULL DEFAULT uuid_generate_v1mc(),
    name       TEXT      NOT NULL,
    email      TEXT      NOT NULL,
    birthday   DATE      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT students_pk PRIMARY KEY (id),
    CONSTRAINT students_email_un UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS courses (
    id         UUID           NOT NULL DEFAULT uuid_generate_v1mc(),
    name       TEXT           NOT NULL,
    value      NUMERIC(19, 2) NOT NULL,
    created_at TIMESTAMP      NOT NULL DEFAULT NOW(),
    CONSTRAINT courses_name_un UNIQUE (name),
    CONSTRAINT courses_pk PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS enrollments (
    student_id   UUID           NOT NULL,
    course_id    UUID           NOT NULL,
    installments INT2           NOT NULL,
    status       ACCOUNT_STATUS NOT NULL,
    created_at   TIMESTAMP      NOT NULL DEFAULT NOW(),
    CONSTRAINT enrollments_pk PRIMARY KEY (student_id, course_id),
    CONSTRAINT enrollments_students_fk FOREIGN KEY (student_id) REFERENCES students (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT enrollments_courses_fk FOREIGN KEY (course_id) REFERENCES courses (id) ON DELETE CASCADE ON UPDATE CASCADE
);
