CREATE TABLE users (
    id INTEGER NOT NULL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL
);

CREATE TABLE threads (
    id INTEGER NOT NULL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    author_id INTEGER NOT NULL,
    date_added DATETIME NOT NULL,

    FOREIGN KEY(author_id) REFERENCES users(id)
);

CREATE INDEX idx_threads_date ON threads(date_added);

CREATE TABLE messages (
    id INTEGER NOT NULL PRIMARY KEY,
    body TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    thread_id INTEGER NOT NULL,
    date_added DATETIME NOT NULL,
    
    FOREIGN KEY(author_id) REFERENCES users(id)
    FOREIGN KEY(thread_id) REFERENCES threads(id),
);

CREATE INDEX idx_messages_date ON messages(date_added);