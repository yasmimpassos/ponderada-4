CREATE TABLE IF NOT EXISTS expense_splits (
    expense_id  CHAR(36) NOT NULL,
    user_id     CHAR(36) NOT NULL,
    amount_owed DECIMAL(12, 2) NOT NULL,
    PRIMARY KEY (expense_id, user_id),
    FOREIGN KEY (expense_id) REFERENCES expenses(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
