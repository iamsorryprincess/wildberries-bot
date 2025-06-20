CREATE TABLE IF NOT EXISTS categories (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `name` VARCHAR(20) NOT NULL UNIQUE,
  `title` VARCHAR(100) NOT NULL,
  `emoji` VARCHAR(10) NOT NULL,
  `request_url` VARCHAR(200) NOT NULL,
  `product_url` VARCHAR(200) NOT NULL,
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

CREATE TABLE IF NOT EXISTS sizes (
  `id` BIGINT UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,
  `name` VARCHAR(100) NOT NULL COLLATE utf8mb4_bin,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NULL,
  UNIQUE KEY `uk_size_name` (`name`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;

CREATE TABLE IF NOT EXISTS products_sizes (
  `product_id` BIGINT UNSIGNED NOT NULL,
  `size_id` BIGINT UNSIGNED NOT NULL,
  `first_price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
  `previous_price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
  `current_price` DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
  `current_price_int` BIGINT UNSIGNED NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NULL,
  INDEX `index_product_id` (product_id),
  INDEX `index_size_id` (size_id),
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (size_id) REFERENCES sizes(id) ON DELETE CASCADE,
  UNIQUE KEY `uk_product_size` (`product_id`, `size_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;

CREATE TABLE IF NOT EXISTS tracking_settings (
  `chat_id` BIGINT SIGNED NOT NULL,
  `size_id` BIGINT UNSIGNED NOT NULL,
  `category_id` BIGINT UNSIGNED NOT NULL,
  `diff_value` TINYINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL,
  INDEX `index_chat_id` (chat_id),
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE,
  FOREIGN KEY (size_id) REFERENCES sizes(id) ON DELETE CASCADE,
  UNIQUE KEY `uk_tracking_chat_product_size` (`chat_id`, `size_id`, `category_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;

CREATE TABLE IF NOT EXISTS tracking_logs (
  `chat_id` BIGINT SIGNED NOT NULL,
  `size_id` BIGINT UNSIGNED NOT NULL,
  `product_id` BIGINT UNSIGNED NOT NULL,
  `price` BIGINT UNSIGNED NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT NOW(),
  `updated_at` DATETIME NULL,
  INDEX `index_chat_id` (chat_id),
  FOREIGN KEY (size_id) REFERENCES sizes(id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  UNIQUE KEY `uk_tracking_logs_chat_size_product` (`chat_id`, `size_id`, `product_id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = COMPRESSED KEY_BLOCK_SIZE = 8;

insert into
  categories (name, title, emoji, request_url, product_url)
values
  (
    'dresses',
    'Платья',
    '👗',
    'https://catalog.wb.ru/catalog/%s/v2/catalog?ab_testing=false&appType=1&cat=8137&curr=rub&dest=-1257786&hide_dtype=13&lang=ru&page=%d&sort=popular&spp=30',
    'https://www.wildberries.ru/catalog/%d/detail.aspx'
  ),
  (
    'bl_shirts',
    'Блузки и рубашки',
    '👚',
    'https://catalog.wb.ru/catalog/%s/v2/catalog?ab_testing=false&appType=1&cat=8126&curr=rub&dest=-5892277&hide_dtype=13&lang=ru&page=%d&sort=popular&spp=30',
    'https://www.wildberries.ru/catalog/%d/detail.aspx'
  );