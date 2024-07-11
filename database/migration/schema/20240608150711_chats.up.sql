CREATE TABLE IF NOT EXISTS chats (
    id UUID DEFAULT gen_random_uuid(),
    chat_name VARCHAR (255) DEFAULT NULL, -- Only used for group chats, can be null for individual chats
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT NULL,

    PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS chats_chat_name_idx ON chats (chat_name);