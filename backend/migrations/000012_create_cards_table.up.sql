CREATE TABLE cards (
    id TEXT NOT NULL PRIMARY KEY,
    column_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    position INTEGER NOT NULL,
    priority TEXT,
    due_date DATETIME NULL,
    linked_item_type TEXT,
    linked_item_id TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT fk_cards_column FOREIGN KEY (column_id) REFERENCES board_columns(id) ON DELETE CASCADE
);

CREATE INDEX idx_cards_column_position ON cards(column_id, position);