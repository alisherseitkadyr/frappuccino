-- ENUM Types
CREATE TYPE order_status AS ENUM ('pending', 'processing', 'closed', 'cancelled');
CREATE TYPE inventory_unit AS ENUM ('kg', 'g', 'liter', 'ml', 'unit');
CREATE TYPE transaction_type AS ENUM ('initial_stock', 'purchase', 'waste', 'adjustment');

DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS menu_items;
DROP TABLE IF EXISTS inventory;
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS menu_item_ingredients;
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS price_history;
DROP TABLE IF EXISTS inventory_transactions;

--
-- Orders Table
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL DEFAULT 0,
    status order_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

--
-- Menu Items Table
CREATE TABLE menu_items (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    categories TEXT[] NOT NULL,
    price DECIMAL(10, 2) NOT NULL CHECK(price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--
-- Inventory Table
CREATE TABLE inventory ( 
    ingredient_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0 CHECK(quantity >= 0),
    unit inventory_unit,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

--
-- Order Items Table
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES menu_items(product_id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK(quantity >= 0),
    customization JSONB DEFAULT '{}'::JSONB
);

--
-- Menu Item Ingredients 
CREATE TABLE menu_item_ingredients (
    ingredient_id INT NOT NULL REFERENCES inventory(ingredient_id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES menu_items(product_id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK(quantity >= 0),
    PRIMARY KEY (product_id, ingredient_id)
);

--
-- Order Status History
CREATE TABLE order_status_history (
    id SERIAL PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    status order_status NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

--
-- Price History
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES menu_items(product_id) ON DELETE CASCADE,
    old_price DECIMAL(10,2) NOT NULL,
    new_price DECIMAL(10,2) NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

--
-- Inventory Transactions
CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id INT NOT NULL REFERENCES inventory(ingredient_id) ON DELETE CASCADE,
    quantity_change INT NOT NULL,
    transaction_type transaction_type NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);
CREATE INDEX idx_order_status_history_order_id ON order_status_history(order_id);
CREATE INDEX idx_menu_item_ingredients_product_id ON menu_item_ingredients(product_id);
CREATE INDEX idx_menu_item_ingredients_ingredient_id ON menu_item_ingredients(ingredient_id);
CREATE INDEX idx_price_history_product_id ON price_history(product_id);
CREATE INDEX idx_inventory_transactions_inventory_id ON inventory_transactions(inventory_id);


-- Insert inventory items
INSERT INTO inventory (name, quantity, unit) VALUES
('Espresso Beans', 5000, 'g'),
('Milk', 10000, 'ml'),
('Sugar', 3000, 'g'),
('Vanilla Syrup', 2000, 'ml'),
('Caramel Syrup', 1500, 'ml'),
('Chocolate Syrup', 1500, 'ml'),
('Whipped Cream', 500, 'ml'),
('Ice Cubes', 2000, 'g'),
('Paper Cups', 300, 'unit'),
('Lids', 300, 'unit'),
('Straws', 300, 'unit'),
('Green Tea Leaves', 2000, 'g'),
('Matcha Powder', 1000, 'g'),
('Lemon', 100, 'unit'),
('Honey', 800, 'ml'),
('Cinnamon', 300, 'g'),
('Oat Milk', 5000, 'ml'),
('Coconut Milk', 4000, 'ml'),
('Espresso Shot', 100, 'unit'),
('Cold Brew Concentrate', 2000, 'ml');

-- Insert menu items
INSERT INTO menu_items (product_name, description, categories, price) VALUES
('Latte', 'Espresso with steamed milk', ARRAY['coffee', 'hot'], 4.50),
('Cappuccino', 'Espresso with steamed milk and foam', ARRAY['coffee', 'hot'], 4.20),
('Americano', 'Espresso with hot water', ARRAY['coffee', 'hot'], 3.00),
('Iced Latte', 'Chilled latte served over ice', ARRAY['coffee', 'cold'], 4.80),
('Mocha', 'Espresso with chocolate and steamed milk', ARRAY['coffee', 'hot'], 5.00),
('Caramel Macchiato', 'Vanilla, milk, espresso, caramel', ARRAY['coffee', 'hot'], 5.30),
('Matcha Latte', 'Matcha with steamed milk', ARRAY['tea', 'hot'], 4.70),
('Green Tea', 'Brewed green tea leaves', ARRAY['tea', 'hot'], 3.50),
('Iced Americano', 'Espresso over ice and water', ARRAY['coffee', 'cold'], 3.20),
('Cold Brew', 'Cold brewed coffee concentrate', ARRAY['coffee', 'cold'], 4.00);

-- Insert menu_item_ingredients (mapping menu items to inventory)
INSERT INTO menu_item_ingredients (product_id, ingredient_id, quantity) VALUES
-- Latte: espresso beans + milk
(1, 1, 18), (1, 2, 240),
-- Cappuccino: espresso + milk + foam
(2, 1, 18), (2, 2, 180),
-- Americano
(3, 1, 18),
-- Iced Latte
(4, 1, 18), (4, 2, 240), (4, 8, 100),
-- Mocha
(5, 1, 18), (5, 2, 240), (5, 6, 30),
-- Caramel Macchiato
(6, 1, 18), (6, 2, 240), (6, 4, 20), (6, 5, 20),
-- Matcha Latte
(7, 13, 10), (7, 2, 240),
-- Green Tea
(8, 12, 5),
-- Iced Americano
(9, 1, 18), (9, 8, 100),
-- Cold Brew
(10, 20, 240);

-- Insert 30 orders
INSERT INTO orders (customer_name, total_price, status, created_at) VALUES
('Alice', 4.50, 'pending', '2024-12-01'),
('Bob', 3.00, 'closed', '2024-12-02'),
('Charlie', 5.00, 'processing', '2024-12-03'),
('Diana', 5.30, 'cancelled', '2024-12-04'),
('Eve', 4.70, 'closed', '2024-12-05'),
('Frank', 4.00, 'pending', '2024-12-06'),
('Grace', 4.20, 'closed', '2024-12-07'),
('Hank', 4.80, 'processing', '2024-12-08'),
('Isabel', 3.20, 'closed', '2024-12-09'),
('Jake', 4.50, 'closed', '2024-12-10'),
('Karen', 5.30, 'closed', '2024-12-11'),
('Leo', 5.00, 'closed', '2024-12-12'),
('Mia', 4.70, 'closed', '2024-12-13'),
('Nick', 4.50, 'pending', '2024-12-14'),
('Olivia', 3.20, 'closed', '2024-12-15'),
('Paul', 4.80, 'processing', '2024-12-16'),
('Quinn', 4.00, 'closed', '2024-12-17'),
('Rose', 3.00, 'closed', '2024-12-18'),
('Steve', 4.20, 'closed', '2024-12-19'),
('Tina', 4.50, 'closed', '2024-12-20'),
('Uma', 4.80, 'closed', '2024-12-21'),
('Victor', 5.00, 'cancelled', '2024-12-22'),
('Wendy', 5.30, 'pending', '2024-12-23'),
('Xavier', 3.20, 'closed', '2024-12-24'),
('Yara', 4.00, 'closed', '2024-12-25'),
('Zane', 4.50, 'closed', '2024-12-26'),
('Amy', 5.00, 'closed', '2024-12-27'),
('Ben', 4.80, 'closed', '2024-12-28'),
('Cara', 4.70, 'processing', '2024-12-29'),
('Dan', 5.30, 'closed', '2024-12-30');

-- Order items (just for some of the above)
INSERT INTO order_items (order_id, product_id, quantity) VALUES
(1, 1, 1),
(2, 3, 1),
(3, 5, 2),
(4, 6, 1),
(5, 7, 1),
(6, 10, 2),
(7, 2, 1),
(8, 4, 1),
(9, 9, 1),
(10, 1, 2),
(11, 8, 1),
(12, 3, 1),
(13, 5, 1),
(14, 6, 1),
(15, 7, 2),
(16, 10, 1),
(17, 2, 1),
(18, 4, 1),
(19, 9, 1),
(20, 1, 1),
(21, 8, 1),
(22, 3, 1),
(23, 5, 1),
(24, 6, 1),
(25, 7, 1),
(26, 10, 2),
(27, 2, 1),
(28, 4, 1),
(29, 9, 1),
(30, 1, 1);


-- Order status history
INSERT INTO order_status_history (order_id, status, changed_at) VALUES
-- Order 1 (closed)
(1, 'pending', '2024-12-01 08:00'),
(1, 'processing', '2024-12-01 08:05'),
(1, 'closed', '2024-12-01 08:10'),

-- Order 2 (closed)
(2, 'pending', '2024-12-02 08:00'),
(2, 'closed', '2024-12-02 08:10'),

-- Order 3 (processing)
(3, 'pending', '2024-12-03 08:00'),
(3, 'processing', '2024-12-03 08:05'),

-- Order 4 (cancelled)
(4, 'pending', '2024-12-04 08:00'),
(4, 'cancelled', '2024-12-04 08:03'),

-- Order 5 (closed)
(5, 'pending', '2024-12-05 08:00'),
(5, 'processing', '2024-12-05 08:05'),
(5, 'closed', '2024-12-05 08:10'),

-- Order 6 (pending)
(6, 'pending', '2024-12-06 08:00'),

-- Order 7 (closed)
(7, 'pending', '2024-12-07 08:00'),
(7, 'processing', '2024-12-07 08:05'),
(7, 'closed', '2024-12-07 08:15'),

-- Order 8 (processing)
(8, 'pending', '2024-12-08 08:00'),
(8, 'processing', '2024-12-08 08:10'),

-- Order 9 (closed)
(9, 'pending', '2024-12-09 08:00'),
(9, 'closed', '2024-12-09 08:05'),

-- Order 10 (closed)
(10, 'pending', '2024-12-10 08:00'),
(10, 'closed', '2024-12-10 08:10'),

-- Order 11 (closed)
(11, 'pending', '2024-12-11 08:00'),
(11, 'processing', '2024-12-11 08:05'),
(11, 'closed', '2024-12-11 08:15'),

-- Order 12 (closed)
(12, 'pending', '2024-12-12 08:00'),
(12, 'closed', '2024-12-12 08:10'),

-- Order 13 (closed)
(13, 'pending', '2024-12-13 08:00'),
(13, 'processing', '2024-12-13 08:05'),
(13, 'closed', '2024-12-13 08:10'),

-- Order 14 (pending)
(14, 'pending', '2024-12-14 08:00'),

-- Order 15 (closed)
(15, 'pending', '2024-12-15 08:00'),
(15, 'closed', '2024-12-15 08:08'),

-- Order 16 (processing)
(16, 'pending', '2024-12-16 08:00'),
(16, 'processing', '2024-12-16 08:10'),

-- Order 17 (closed)
(17, 'pending', '2024-12-17 08:00'),
(17, 'processing', '2024-12-17 08:05'),
(17, 'closed', '2024-12-17 08:15'),

-- Order 18 (closed)
(18, 'pending', '2024-12-18 08:00'),
(18, 'closed', '2024-12-18 08:05'),

-- Order 19 (closed)
(19, 'pending', '2024-12-19 08:00'),
(19, 'closed', '2024-12-19 08:08'),

-- Order 20 (closed)
(20, 'pending', '2024-12-20 08:00'),
(20, 'processing', '2024-12-20 08:04'),
(20, 'closed', '2024-12-20 08:10'),

-- Order 21 (closed)
(21, 'pending', '2024-12-21 08:00'),
(21, 'closed', '2024-12-21 08:06'),

-- Order 22 (cancelled)
(22, 'pending', '2024-12-22 08:00'),
(22, 'cancelled', '2024-12-22 08:01'),

-- Order 23 (pending)
(23, 'pending', '2024-12-23 08:00'),

-- Order 24 (closed)
(24, 'pending', '2024-12-24 08:00'),
(24, 'closed', '2024-12-24 08:07'),

-- Order 25 (closed)
(25, 'pending', '2024-12-25 08:00'),
(25, 'closed', '2024-12-25 08:10'),

-- Order 26 (closed)
(26, 'pending', '2024-12-26 08:00'),
(26, 'closed', '2024-12-26 08:10'),

-- Order 27 (closed)
(27, 'pending', '2024-12-27 08:00'),
(27, 'closed', '2024-12-27 08:12'),

-- Order 28 (closed)
(28, 'pending', '2024-12-28 08:00'),
(28, 'closed', '2024-12-28 08:08'),

-- Order 29 (processing)
(29, 'pending', '2024-12-29 08:00'),
(29, 'processing', '2024-12-29 08:10'),

-- Order 30 (closed)
(30, 'pending', '2024-12-30 08:00'),
(30, 'closed', '2024-12-30 08:09');


-- Price history for all items
INSERT INTO price_history (product_id, old_price, new_price, changed_at) VALUES
(1, 4.00, 4.50, '2024-10-01'),  -- Latte
(2, 3.80, 4.20, '2024-10-10'),  -- Cappuccino
(3, 2.70, 3.00, '2024-10-15'),  -- Americano
(4, 4.50, 4.80, '2024-10-20'),  -- Iced Latte
(5, 4.80, 5.00, '2024-11-05'),  -- Mocha
(6, 5.00, 5.30, '2024-11-15'),  -- Caramel Macchiato
(7, 4.40, 4.70, '2024-11-20'),  -- Matcha Latte
(8, 3.20, 3.50, '2024-11-25'),  -- Green Tea
(9, 2.90, 3.20, '2024-12-01'),  -- Iced Americano
(10, 3.70, 4.00, '2024-12-05'); -- Cold Brew

-- Inventory transactions (one per item)
INSERT INTO inventory_transactions (inventory_id, quantity_change, transaction_type, transaction_date) VALUES
(1, 1000, 'initial_stock', '2024-10-01'),      -- Espresso Beans
(2, 2000, 'initial_stock', '2024-10-01'),      -- Milk
(3, 1000, 'initial_stock', '2024-10-01'),      -- Sugar
(4, 500, 'purchase', '2024-10-02'),            -- Vanilla Syrup
(5, 300, 'purchase', '2024-10-02'),            -- Caramel Syrup
(6, 400, 'purchase', '2024-10-03'),            -- Chocolate Syrup
(7, 200, 'purchase', '2024-10-03'),            -- Whipped Cream
(8, 1000, 'initial_stock', '2024-10-01'),      -- Ice Cubes
(9, 100, 'purchase', '2024-10-04'),            -- Paper Cups
(10, 100, 'purchase', '2024-10-04'),           -- Lids
(11, 100, 'purchase', '2024-10-04'),           -- Straws
(12, 500, 'initial_stock', '2024-10-01'),      -- Green Tea Leaves
(13, 300, 'purchase', '2024-10-05'),           -- Matcha Powder
(14, 50, 'purchase', '2024-10-05'),            -- Lemon
(15, 200, 'purchase', '2024-10-06'),           -- Honey
(16, 150, 'waste', '2024-10-07'),              -- Cinnamon
(17, 1000, 'adjustment', '2024-10-08'),        -- Oat Milk
(18, 800, 'purchase', '2024-10-08'),           -- Coconut Milk
(19, 50, 'initial_stock', '2024-10-01'),       -- Espresso Shot
(20, 700, 'purchase', '2024-10-09');           -- Cold Brew Concentrate


