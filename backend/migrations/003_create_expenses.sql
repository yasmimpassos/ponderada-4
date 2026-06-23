CREATE TABLE IF NOT EXISTS expenses (
    id           CHAR(36) PRIMARY KEY,
    group_id     CHAR(36) NOT NULL,
    paid_by      CHAR(36) NOT NULL,
    amount       DECIMAL(12, 2) NOT NULL,
    description  TEXT,
    receipt_url  TEXT,
    expense_date DATE NOT NULL,
    created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (group_id) REFERENCES user_groups(id) ON DELETE CASCADE,
    FOREIGN KEY (paid_by) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
