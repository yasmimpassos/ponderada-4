CREATE TABLE IF NOT EXISTS settlements (
    id         CHAR(36)       NOT NULL PRIMARY KEY,
    payer_id   CHAR(36)       NOT NULL,
    payee_id   CHAR(36)       NOT NULL,
    amount     DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (payer_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (payee_id) REFERENCES users(id) ON DELETE CASCADE
);
