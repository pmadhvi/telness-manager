CREATE TABLE IF NOT EXISTS subscription(
    msidn VARCHAR(50) NOT NULL,
    activate_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    PRIMARY KEY (msidn)
);