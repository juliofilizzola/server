DROP TABLE IF EXISTS rooms;
CREATE TABLE rooms(
    id uuid primary key not null default gen_random_uuid(),
    theme varchar not null,
    name varchar not null
)