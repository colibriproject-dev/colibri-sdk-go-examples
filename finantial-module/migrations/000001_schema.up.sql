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
CREATE TABLE accounts (
    id           UUID           NOT NULL DEFAULT uuid_generate_v1mc(),
    student_id   UUID           NOT NULL,
    course_id    UUID           NOT NULL,
    installments INT2           NOT NULL,
    value        DECIMAL(19,2)  NOT NULL,
    status       ACCOUNT_STATUS NOT NULL,
    created_at   TIMESTAMP      NOT NULL DEFAULT NOW(),
    CONSTRAINT accounts_pk PRIMARY KEY (id)
);

CREATE TABLE invoices (
    id          UUID          NOT NULL DEFAULT uuid_generate_v1mc(),
    account_id  UUID          NOT NULL,
    installment int2          NOT NULL,
    due_date    DATE          NOT NULL,
    value       DECIMAL(19,2) NOT NULL,
    created_at  TIMESTAMP     NOT NULL DEFAULT NOW(),
    paid_at     TIMESTAMP,
    CONSTRAINT invoices_pk PRIMARY KEY (id),
    CONSTRAINT invoices_accounts_fk FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE ON UPDATE CASCADE
);
