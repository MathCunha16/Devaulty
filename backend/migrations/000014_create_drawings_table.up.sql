CREATE TABLE drawings (
    id TEXT NOT NULL PRIMARY KEY,
    project_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT,
    thumbnail TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT fk_project
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
);

CREATE INDEX idx_drawings_projects ON drawings(project_id);
