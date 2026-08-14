CREATE TABLE app_settings (
    key TEXT NOT NULL PRIMARY KEY,
    value TEXT NOT NULL
);

INSERT INTO app_settings (key, value) VALUES ('crypto_version', '1');
