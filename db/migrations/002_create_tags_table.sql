-- this is just a piece of shit.
-- TODO: change the schema to:
-- CREATE TABLE IF NOT EXISTS tags (
--     name TEXT NOT NULL,
--     post_id INTEGER NOT NULL,
--     PRIMARY KEY (post_id, name),
--     FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
-- );
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE
);
