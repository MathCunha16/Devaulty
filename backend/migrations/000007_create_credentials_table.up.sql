CREATE TABLE credentials (
    id TEXT NOT NULL PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    secret_type TEXT NOT NULL,
    payload_encrypted BLOB NOT NULL,
    encryption_iv BLOB NOT NULL,
    encryption_auth_tag BLOB NOT NULL,
    notes TEXT,
    related_url TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT fk_credentials_project FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX idx_credentials_project ON credentials(project_id);
