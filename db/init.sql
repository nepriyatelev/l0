CREATE TABLE IF NOT EXISTS orders (
      order_uid VARCHAR(255) PRIMARY KEY,
      track_number VARCHAR(255),
      entry VARCHAR(255),
      locale VARCHAR(255),
      internal_signature VARCHAR(255),
      customer_id VARCHAR(255),
      delivery_service VARCHAR(255),
      shard_key VARCHAR(255),
      sm_id INT,
      date_created TIMESTAMP WITH TIME ZONE,
      oof_shard VARCHAR(255));

CREATE INDEX IF NOT EXISTS order_uid_idx ON orders (order_uid);

CREATE TABLE IF NOT EXISTS delivery (
    order_uid VARCHAR(255) PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(255),
    zip VARCHAR(255),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255));
--     FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE);

CREATE INDEX IF NOT EXISTS delivery_order_uid_idx ON delivery (order_uid);

CREATE TABLE IF NOT EXISTS payments (
--     order_uid VARCHAR(255),
    transaction VARCHAR(255) PRIMARY KEY,
    request_id VARCHAR(255),
    currency VARCHAR(255),
    provider VARCHAR(255),
    amount DECIMAL(10, 2),
    payment_dt INT,
    bank VARCHAR(255),
    delivery_cost DECIMAL(10, 2),
    goods_total DECIMAL(10, 2),
    custom_fee DECIMAL(10, 2));
--     FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE);

CREATE INDEX IF NOT EXISTS payments_transaction_idx ON payments (transaction);

CREATE TABLE IF NOT EXISTS items (
--      order_uid VARCHAR(255),
     chrt_id INT,
     track_number VARCHAR(255),
     price DECIMAL(10, 2),
     rid VARCHAR(255),
     name VARCHAR(255),
     sale DECIMAL(10, 2),
     size VARCHAR(255),
     total_price DECIMAL(10, 2),
     nm_id INT,
     brand VARCHAR(255),
     status INT);
--      FOREIGN KEY (order_uid) REFERENCES orders (order_uid) ON DELETE CASCADE);

CREATE INDEX IF NOT EXISTS track_number_idx ON items (track_number);