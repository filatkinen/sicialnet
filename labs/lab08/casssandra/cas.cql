CREATE KEYSPACE IF NOT EXISTS dialog
    WITH replication = {
        'class': 'SimpleStrategy',
        'replication_factor': '1'
        };

CREATE TABLE IF NOT EXISTS dialog.dialogs (
                                user_id text,
                                friend_id text,
                                dialog_id uuid,
                                message text,
                                PRIMARY KEY (user_id, friend_id, dialog_id)
);
