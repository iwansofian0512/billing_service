INSERT INTO borrowers (id, name, email, is_active)
VALUES
    (1, 'iwan', 'iwan@example.com', TRUE),
    (2, 'sofian', 'sofian@example.com', TRUE),
    (3, 'wawan', 'wawan@example.com', TRUE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO loans (id, borrower_id, principal_amount, total_interest, total_payable, outstanding_amount, duration_weeks, weekly_payment_amount, is_active, status)
VALUES
    (1, 1, 5000000, 500000, 5500000, 4950000, 50, 110000, TRUE, 'inprogress'),
    (2, 2, 5000000, 500000, 5500000, 0,       50, 110000, FALSE, 'completed'),
    (3, 3, 5000000, 500000, 5500000, 4400000, 50, 110000, TRUE, 'inprogress')
ON CONFLICT (id) DO NOTHING;

INSERT INTO billing_schedules (loan_id, week_number, due_date, amount_due, amount_paid, status)
SELECT
    1,
    g,
    CURRENT_DATE + (g * 7) * INTERVAL '1 day',
    110000,
    CASE WHEN g <= 5 THEN 110000 ELSE 0 END,
    CASE WHEN g <= 5 THEN 'paid'::billing_status ELSE 'pending'::billing_status END
FROM generate_series(1, 50) AS g;

INSERT INTO billing_schedules (loan_id, week_number, due_date, amount_due, amount_paid, status)
SELECT
    2,
    g,
    CURRENT_DATE + (g * 7) * INTERVAL '1 day',
    110000,
    110000,
    'paid'::billing_status
FROM generate_series(1, 50) AS g;

INSERT INTO billing_schedules (loan_id, week_number, due_date, amount_due, amount_paid, status)
SELECT
    3,
    g,
    DATE '2025-12-01' + (g - 1) * INTERVAL '7 days',
    110000,
    CASE WHEN g <= 9 THEN 110000 ELSE 0 END,
    CASE WHEN g <= 9 THEN 'paid'::billing_status ELSE 'pending'::billing_status END
FROM generate_series(1, 50) AS g;

INSERT INTO payments (loan_id, amount, payment_date)
SELECT
    1,
    110000,
    CURRENT_TIMESTAMP - (5 - g) * INTERVAL '7 days'
FROM generate_series(1, 5) AS g;

INSERT INTO payments (loan_id, amount, payment_date)
SELECT
    2,
    110000,
    CURRENT_TIMESTAMP - (50 - g) * INTERVAL '7 days'
FROM generate_series(1, 50) AS g;

INSERT INTO payments (loan_id, amount, payment_date)
SELECT
    3,
    110000,
    TIMESTAMP '2025-12-01' + (g - 1) * INTERVAL '7 days'
FROM generate_series(1, 9) AS g;

SELECT setval('borrowers_id_seq', (SELECT COALESCE(MAX(id), 1) FROM borrowers));
SELECT setval('loans_id_seq', (SELECT COALESCE(MAX(id), 1) FROM loans));
SELECT setval('billing_schedules_id_seq', (SELECT COALESCE(MAX(id), 1) FROM billing_schedules));
SELECT setval('payments_id_seq', (SELECT COALESCE(MAX(id), 1) FROM payments));
