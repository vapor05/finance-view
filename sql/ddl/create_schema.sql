CREATE SCHEMA financeview;

CREATE TABLE financeview.description (
    id SERIAL PRIMARY KEY NOT NULL,
    description TEXT,
    createdate TIMESTAMP,
    updatedate TIMESTAMP
);

CREATE TABLE financeview.expense (
    id SERIAL PRIMARY KEY NOT NULL,
    date DATE,
    description_id INT,
    amount MONEY,
    comment TEXT,
    createdate TIMESTAMP,
    updatedate TIMESTAMP
);

CREATE TABLE financeview.category (
    id SERIAL PRIMARY KEY NOT NULL,
    name TEXT,
    createdate TIMESTAMP,
    updatedate TIMESTAMP
);

CREATE TABLE financeview.expense_category (
    id SERIAL PRIMARY KEY NOT NULL,
    expense_id INT NOT NULL,
    category_id INT NOT NULL,
    createdate TIMESTAMP
);