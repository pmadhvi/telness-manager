CREATE TABLE IF NOT EXISTS subscription(
    msisdn VARCHAR(12) NOT NULL UNIQUE,
    activate_at TIMESTAMP NOT NULL,
    sub_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    modified_at TIMESTAMP NOT NULL,
    PRIMARY KEY (msisdn)
);