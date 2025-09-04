-- name: GetCategories :many
SELECT DISTINCT category
FROM prepaid_products;

-- name: GetBrandsByCategory :many
SELECT DISTINCT brand
FROM prepaid_products
WHERE category = ?;

-- name: GetTypesByCategoryAndBrand :many
SELECT DISTINCT type
FROM prepaid_products
WHERE category = ? AND brand = ?;

-- name: InsertPrepaidProduct :exec
INSERT INTO prepaid_products (
  name, 
  category, 
  brand, 
  type, 
  seller_name, 
  price,
  buyer_sku_code, 
  buyer_product_status, 
  seller_product_status,
  unlimited_stock,
  stock,
  multi,
  start_cut_off,
  end_cut_off,
  description
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: GetPrepaidProducts :many
SELECT
  pp.name,
  pp.buyer_sku_code,
  pp.price,
  pp.seller_name,
  pp.buyer_product_status,
  pp.seller_product_status
FROM prepaid_products pp
WHERE pp.category = ?
  AND pp.brand = ?
  AND pp.type = ?
ORDER BY pp.price ASC;

-- name: GetPrepaidProductBySKUCode :one
SELECT 
  id,
  name,
  seller_name,
  price,
  buyer_sku_code,
  buyer_product_status,
  seller_product_status,
  description
FROM prepaid_products
WHERE buyer_sku_code = ? COLLATE NOCASE
LIMIT 1;

-- name: DeleteAllPrepaidProducts :exec
DELETE FROM prepaid_products;