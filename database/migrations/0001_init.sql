-- +goose Up
-- +goose StatementBegin

-- users
CREATE TABLE users (
  id integer PRIMARY KEY,
  username text,
  first_name text,
  last_name text,
  created_at datetime NOT NULL
);

-- chats
CREATE TABLE chats (
  id integer PRIMARY KEY,
  command text NOT NULL,
  step integer NOT NULL,
  data json,
  reply_markup_1 json,
  reply_markup_2 json,
  reply_markup_3 json,
  reply_markup_4 json
);

-- prepaid_products
CREATE TABLE prepaid_products (
  id integer PRIMARY KEY AUTOINCREMENT,
  name text NOT NULL,
  category text NOT NULL,
  brand text NOT NULL,
  type text NOT NULL,
  seller_name text NOT NULL,
  price integer NOT NULL,
  buyer_sku_code text NOT NULL,
  buyer_product_status boolean NOT NULL,
  seller_product_status boolean NOT NULL,
  unlimited_stock boolean NOT NULL,
  stock integer NOT NULL,
  multi boolean NOT NULL,
  start_cut_off time,
  end_cut_off time,
  description text
);

CREATE INDEX idx_prepaid_products_category ON prepaid_products(category);
CREATE INDEX idx_prepaid_products_brand ON prepaid_products(brand);
CREATE INDEX idx_prepaid_products_type ON prepaid_products(type);
CREATE UNIQUE INDEX idx_prepaid_products_buyer_sku_code ON prepaid_products(buyer_sku_code COLLATE NOCASE);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE prepaid_products;
DROP TABLE chats;
DROP TABLE users;
-- +goose StatementEnd
