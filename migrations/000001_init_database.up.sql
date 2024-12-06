CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE TABLE IF NOT EXISTS users_auth(
    id uuid DEFAULT uuid_generate_v4 (),
    user_id uuid DEFAULT uuid_generate_v4 (),
    rt TEXT,
    PRIMARY KEY (id)
);