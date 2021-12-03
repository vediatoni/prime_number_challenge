CREATE TABLE prime_numbers
(
    number_tested                    INTEGER,
    is_prime                         BOOLEAN,
    validation_time                  BIGINT,
    auto_incrementing_seq            SERIAL PRIMARY KEY,
    time_needed_to_validate_microsec INTEGER
);
