-- Таблица категорий
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Заполняем базовыми категориями (добавятся только если их еще нет)
INSERT OR IGNORE INTO categories (id, name) VALUES (1, 'Еда');
INSERT OR IGNORE INTO categories (id, name) VALUES (2, 'Транспорт');
INSERT OR IGNORE INTO categories (id, name) VALUES (3, 'Жильё');
INSERT OR IGNORE INTO categories (id, name) VALUES (4, 'Развлечения');
INSERT OR IGNORE INTO categories (id, name) VALUES (5, 'Другое');

-- Таблица трат. Добавляем category_id
CREATE TABLE IF NOT EXISTS expenses (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    description TEXT NOT NULL,
    date TEXT NOT NULL,
    category_id INTEGER,
    FOREIGN KEY (category_id) REFERENCES categories(id)
);