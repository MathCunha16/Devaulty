CREATE TABLE card_items (
    card_id TEXT NOT NULL,
    item_type TEXT NOT NULL,
    item_id TEXT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    PRIMARY KEY (card_id, item_type, item_id),
    CONSTRAINT fk_card_items_card FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
);

CREATE INDEX idx_card_items_card ON card_items(card_id);
CREATE INDEX idx_card_items_item ON card_items(item_type, item_id);