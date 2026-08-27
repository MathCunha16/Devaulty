CREATE TABLE board_columns(
    id TEXT NOT NULL PRIMARY KEY,
    board_id TEXT NOT NULL,
    name TEXT NOT NULL,
    position INTEGER NOT NULL,
    wip_limit INTEGER,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT fk_board_column_board FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);

CREATE INDEX idx_board_columns_board_position ON board_columns(board_id, position);