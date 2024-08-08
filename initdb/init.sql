CREATE TABLE IF NOT EXISTS games (
    game_id TEXT PRIMARY KEY,
    player_id1 TEXT,
    player_id2 TEXT,
    game_state TEXT,
    live BOOLEAN DEFAULT TRUE
);