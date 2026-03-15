USE users_db;

CREATE TABLE user_friends (
    user_id INT,
    friend_id INT,

    PRIMARY KEY (user_id, friend_id),

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (friend_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT no_self_friend CHECK ( user_id <> friend_id )

);