CREATE TABLE IF NOT EXISTS categories (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `name` VARCHAR(20) NOT NULL UNIQUE,
  `title` VARCHAR(100) NOT NULL,
  `emoji` VARCHAR(10) NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;

CREATE TABLE IF NOT EXISTS products (
  `id` BIGINT UNSIGNED NOT NULL PRIMARY KEY,
  `category_id` BIGINT UNSIGNED NOT NULL,
  `name` VARCHAR(100) NOT NULL,
  `rating` DECIMAL(2, 1) NOT NULL DEFAULT 0,
  `url` VARCHAR(100) NOT NULL,
  `brand` VARCHAR(100) NOT NULL,
  `brand_id` BIGINT UNSIGNED NOT NULL,
  `colors` JSON NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL,
  INDEX `index_category_id` (category_id),
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
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

CREATE TABLE IF NOT EXISTS tracking_settings (
  `chat_id` BIGINT SIGNED NOT NULL,
  `size` VARCHAR(100) NOT NULL,
  `category` VARCHAR(20) NOT NULL,
  `diff_value` TINYINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL,
  INDEX `index_chat_id` (chat_id),
  UNIQUE KEY `uk_tracking_chat_product_size` (`chat_id`, `size`, `category`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;

insert into
  categories (name, title, emoji)
values
  ('dresses', 'Платья', '👗');