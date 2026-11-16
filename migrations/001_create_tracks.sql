CREATE TABLE tracks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    artist TEXT NOT NULL,
    album TEXT,
    year INTEGER,
    rating INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT tracks_rating_check CHECK (rating >= 0 AND rating <= 5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);