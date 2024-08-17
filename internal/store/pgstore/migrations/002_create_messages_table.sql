-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS messages(
    id uuid primary key not null default gen_random_uuid(),
    room_id uuid not null,
    message varchar(255) not null,
    reaction_count int not null default 0,
    answered bool not null default false,

    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now(),

    FOREIGN KEY (room_id) REFERENCES rooms(id)
);
---- create above / drop below ----
DROP TABLE IF EXISTS messages;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
