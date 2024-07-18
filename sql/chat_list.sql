SELECT c.id,
    CASE
        WHEN c.chat_name IS NOT NULL THEN c.chat_name
        ELSE (
            SELECT u.username
            FROM chat_participants cp2
                JOIN users u ON cp2.user_id = u.id
            WHERE cp2.chat_id = c.id
                AND u.id <> '78901234-5678-9012-3456-789012345678'
            LIMIT 1
        )
    END AS chat_name,
    m.content as msg_higlight,
    m.created_at as last_msg_date
FROM chats c
    JOIN chat_participants cp ON c.id = cp.chat_id
    LEFT JOIN LATERAL (
        SELECT m.content,
            m.created_at
        FROM messages m
        WHERE m.chat_id = c.id
        ORDER BY m.created_at DESC
        LIMIT 1
    ) m ON TRUE
WHERE cp.user_id = '78901234-5678-9012-3456-789012345678'
    AND m.created_at > '2024-05-15T12:34:56Z' -- the date is last date on last message on chat list
ORDER BY m.created_at DESC