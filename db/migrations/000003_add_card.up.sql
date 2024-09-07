CREATE TABLE IF NOT EXISTS cards (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    number VARCHAR NOT NULL,
    expiration_date VARCHAR NOT NULL,
    holder_name VARCHAR NOT NULL,
    cvv VARCHAR NOT NULL,
    user_id uuid references users(id),
    uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);