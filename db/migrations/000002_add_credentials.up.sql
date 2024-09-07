CREATE TABLE IF NOT EXISTS credentials (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    service_name VARCHAR NOT NULL,
    identity VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    user_id uuid references users(id),
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);