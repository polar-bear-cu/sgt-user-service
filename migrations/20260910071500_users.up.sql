CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        email TEXT NOT NULL UNIQUE,
        google_sub TEXT NOT NULL UNIQUE,
        display_name TEXT NOT NULL DEFAULT '',
        picture_url TEXT NOT NULL DEFAULT ''
    );