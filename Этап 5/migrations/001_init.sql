-- 1. Категории
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

INSERT OR IGNORE INTO categories (id, name) VALUES (1, 'Еда');
INSERT OR IGNORE INTO categories (id, name) VALUES (2, 'Транспорт');
INSERT OR IGNORE INTO categories (id, name) VALUES (3, 'Жильё');
INSERT OR IGNORE INTO categories (id, name) VALUES (4, 'Развлечения');
INSERT OR IGNORE INTO categories (id, name) VALUES (5, 'Другое');

-- 2. Пользователи
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL
);

-- 3. Сессии (храним токены авторизации)
CREATE TABLE IF NOT EXISTS sessions (
    session_id TEXT PRIMARY KEY,
    user_id INTEGER NOT NULL,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- 4. Траты (добавляем user_id)
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    description TEXT NOT NULL,
    date TEXT NOT NULL,
    category_id INTEGER,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);