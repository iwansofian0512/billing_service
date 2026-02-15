CREATE TYPE loan_status AS ENUM ('inprogress', 'completed');

CREATE TABLE IF NOT EXISTS loans (
    id SERIAL PRIMARY KEY,
    borrower_id INT NOT NULL,
    principal_amount NUMERIC(15, 2) NOT NULL,
    total_interest NUMERIC(15, 2) NOT NULL,
    total_payable NUMERIC(15, 2) NOT NULL,
    outstanding_amount NUMERIC(15, 2) NOT NULL,
    duration_weeks INT NOT NULL,
    weekly_payment_amount NUMERIC(15, 2) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    status loan_status NOT NULL DEFAULT 'inprogress',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE billing_status AS ENUM ('pending', 'paid');

CREATE TABLE IF NOT EXISTS billing_schedules (
    id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES loans(id) ON DELETE CASCADE,
    week_number INT NOT NULL,
    due_date DATE NOT NULL,
    amount_due NUMERIC(15, 2) NOT NULL,
    amount_paid NUMERIC(15, 2) DEFAULT 0,
    status billing_status DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_billing_schedules_loan_id ON billing_schedules(loan_id);
CREATE INDEX idx_billing_schedules_status ON billing_schedules(status);

CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    loan_id INT REFERENCES loans(id) ON DELETE CASCADE,
    billing_schedule_id INT REFERENCES billing_schedules(id) ON DELETE CASCADE,
    amount NUMERIC(15, 2) NOT NULL,
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS borrowers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE loans
    ADD CONSTRAINT fk_loans_borrower
    FOREIGN KEY (borrower_id) REFERENCES borrowers(id) ON DELETE RESTRICT;
