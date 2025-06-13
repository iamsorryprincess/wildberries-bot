CREATE TABLE IF NOT EXISTS products (
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL,
  `rating` DECIMAL(2, 1) NOT NULL DEFAULT 0,
  `url` VARCHAR(100) NOT NULL,
  `brand` VARCHAR(100) NOT NULL,
  `brand_id` BIGINT UNSIGNED NOT NULL,
  `colors` JSON NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;

CREATE TABLE IF NOT EXISTS products_sizes (
  `product_id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `first_price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
  `previous_price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
  `current_price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NULL,
  INDEX `index_product_id` (product_id),
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  UNIQUE KEY `uk_product_size` (`product_id`, `name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;