CREATE TABLE IF NOT EXISTS subs (
                                    email VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    frequency VARCHAR(10) NOT NULL,
    confirmation_token VARCHAR(255) UNIQUE,
    unsubscribe_token VARCHAR(255) UNIQUE,
    confirmed BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (email, city, frequency)
    );

CREATE INDEX IF NOT EXISTS idx_subs_confirmation_token ON subs (confirmation_token);
CREATE INDEX IF NOT EXISTS idx_subs_unsubscribe_token ON subs (unsubscribe_token);