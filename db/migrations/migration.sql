-- ============================================
-- SCHEMA DATABASE KOS - OPTIMIZED VERSION
-- ============================================

-- ============================================
-- 1. TABEL USERS (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(128) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,  
    username VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,  
    phone VARCHAR(20) NOT NULL,
    role ENUM('customer', 'owner', 'admin') DEFAULT 'customer',  
    is_verified BOOLEAN DEFAULT FALSE,  
    last_password_change DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_role (role),
    INDEX idx_phone (phone),  -- ADDED: untuk search by phone
    INDEX idx_deleted (deleted_at)  -- ADDED: untuk soft delete queries
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 2. TABEL PROVINCES (No changes needed)
-- ============================================
CREATE TABLE IF NOT EXISTS provinces (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_code (code),
    INDEX idx_name (name)  -- ADDED: untuk search
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 3. TABEL CITIES (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS cities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    province_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE RESTRICT,
    
    UNIQUE KEY unique_city_code (province_id, code),  -- ADDED: untuk mencegah duplikasi
    INDEX idx_province (province_id),
    INDEX idx_name (name)  -- ADDED: untuk search
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 4. TABEL DISTRICTS (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS districts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    city_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE RESTRICT,
    
    INDEX idx_city (city_id),
    INDEX idx_name (name),  -- ADDED: untuk search
    INDEX idx_city_name (city_id, name)  -- ADDED: composite index
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 5. TABEL KOS (Diperbaiki - MAJOR CHANGES)
-- ============================================
CREATE TABLE IF NOT EXISTS kos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    owner_id INT NOT NULL,
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    description TEXT,
    gender_type ENUM('male', 'female', 'mixed') NOT NULL,  -- FIXED: Missing column
    status ENUM('active', 'inactive', 'full', 'maintenance') DEFAULT 'active',
    is_verified BOOLEAN DEFAULT FALSE,
    -- REMOVED: average_rating dan total_reviews (akan dihitung via VIEW)
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE RESTRICT,
    
    INDEX idx_owner (owner_id),
    INDEX idx_gender_type (gender_type),
    INDEX idx_status (status),
    INDEX idx_slug (slug),
    INDEX idx_verified (is_verified),  -- ADDED
    INDEX idx_deleted (deleted_at),  -- ADDED
    INDEX idx_owner_status (owner_id, status),  -- ADDED: composite
    
    FULLTEXT INDEX ft_name_description (name, description) 
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 6. TABEL ADDRESSES (OPTIMIZED - Menghilangkan redundansi)
-- ============================================
CREATE TABLE IF NOT EXISTS addresses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kos_id INT NOT NULL UNIQUE,
    street_address VARCHAR(255) NOT NULL,
    district_id INT NOT NULL,
    -- REMOVED: city_id dan province_id (redundan, bisa JOIN via district)
    postal_code VARCHAR(10) NULL,
    latitude DECIMAL(10,8) NULL,
    longitude DECIMAL(11,8) NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,

    FOREIGN KEY (kos_id) REFERENCES kos(id) ON DELETE CASCADE,
    FOREIGN KEY (district_id) REFERENCES districts(id) ON DELETE RESTRICT,

    INDEX idx_kos (kos_id),
    INDEX idx_district (district_id),
    INDEX idx_location (latitude, longitude),  -- ADDED: untuk geolocation search
    INDEX idx_postal (postal_code)  -- ADDED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 7. TABEL ROOM_TYPES (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS room_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kos_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    length DECIMAL(5, 2) NOT NULL CHECK (length > 0),  -- ADDED: validation
    width DECIMAL(5, 2) NOT NULL CHECK (width > 0),  -- ADDED: validation
    capacity INT NOT NULL DEFAULT 1 CHECK (capacity > 0),  -- ADDED: validation
    gender_type ENUM('male', 'female', 'mixed') NOT NULL,
    total_rooms INT NOT NULL DEFAULT 0 CHECK (total_rooms >= 0),  -- ADDED: validation
    available_rooms INT NOT NULL DEFAULT 0 CHECK (available_rooms >= 0),  -- ADDED: validation
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,

    FOREIGN KEY (kos_id) REFERENCES kos(id) ON DELETE CASCADE,

    INDEX idx_kos (kos_id),
    INDEX idx_gender (gender_type),  -- ADDED
    INDEX idx_availability (kos_id, available_rooms),
    INDEX idx_deleted (deleted_at),  -- ADDED
    
    -- ADDED: Constraint untuk memastikan available_rooms <= total_rooms
    CHECK (available_rooms <= total_rooms)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 8. TABEL ROOM_PRICES (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS room_prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    room_type_id INT NOT NULL,
    period ENUM('daily', 'weekly', 'monthly', 'yearly') NOT NULL,
    price DECIMAL(12, 2) NOT NULL CHECK (price >= 0),  -- ADDED: validation
    is_active BOOLEAN DEFAULT TRUE,
    valid_from DATE NOT NULL,
    valid_until DATE NULL, 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (room_type_id) REFERENCES room_types(id) ON DELETE CASCADE,

    UNIQUE KEY unique_room_period_date (room_type_id, period, valid_from),  -- ADDED
    INDEX idx_room_type (room_type_id),
    INDEX idx_period_price (period, price),
    INDEX idx_active (is_active, valid_from, valid_until),
    
    -- ADDED: Constraint untuk validasi tanggal
    CHECK (valid_until IS NULL OR valid_until >= valid_from)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 9. TABEL FACILITY_CATEGORIES (No major changes)
-- ============================================
CREATE TABLE IF NOT EXISTS facility_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,  -- ADDED: unique
    icon VARCHAR(50) NULL,
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_order (display_order)  -- ADDED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 10. TABEL FACILITIES (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS facilities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT NULL,
    icon VARCHAR(50) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (category_id) REFERENCES facility_categories(id) ON DELETE RESTRICT,

    UNIQUE KEY unique_facility_name (category_id, name),  -- ADDED
    INDEX idx_category (category_id),
    INDEX idx_active (is_active),
    INDEX idx_name (name)  -- ADDED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 11. TABEL KOS_FACILITIES (No major changes)
-- ============================================
CREATE TABLE IF NOT EXISTS kos_facilities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kos_id INT NOT NULL,
    facility_id INT NOT NULL,
    scope ENUM('kos', 'room_type') NOT NULL DEFAULT 'kos',
    room_type_id INT NULL,
    notes TEXT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (kos_id) REFERENCES kos(id) ON DELETE CASCADE,
    FOREIGN KEY (facility_id) REFERENCES facilities(id) ON DELETE CASCADE,
    FOREIGN KEY (room_type_id) REFERENCES room_types(id) ON DELETE CASCADE,

    UNIQUE KEY unique_kos_facility (kos_id, facility_id, scope, room_type_id),
    INDEX idx_kos (kos_id),
    INDEX idx_facility (facility_id),
    INDEX idx_room_type (room_type_id),
    INDEX idx_scope (scope)  -- ADDED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 12. TABEL RULE_CATEGORIES (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS rule_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,  -- ADDED: unique
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_order (display_order)  -- ADDED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 13. TABEL KOS_RULES (No major changes)
-- ============================================
CREATE TABLE IF NOT EXISTS kos_rules (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kos_id INT NOT NULL,
    category_id INT NOT NULL,
    rule TEXT NOT NULL,
    is_mandatory BOOLEAN DEFAULT TRUE,
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    FOREIGN KEY (kos_id) REFERENCES kos(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES rule_categories(id) ON DELETE RESTRICT,

    INDEX idx_kos (kos_id),
    INDEX idx_category (category_id),
    INDEX idx_mandatory (is_mandatory),  -- ADDED
    INDEX idx_kos_category (kos_id, category_id)  -- ADDED: composite
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 14. TABEL IMAGE_CATEGORIES (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS image_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,  -- ADDED: unique
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_order (display_order)  -- ADDED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 15. TABEL KOS_IMAGES (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS kos_images (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kos_id INT NOT NULL,
    category_id INT NOT NULL,
    room_type_id INT NULL,
    image_url VARCHAR(512) NOT NULL,
    title VARCHAR(100) NULL,
    description TEXT NULL,
    is_cover BOOLEAN DEFAULT FALSE,
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (kos_id) REFERENCES kos(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES image_categories(id) ON DELETE RESTRICT,
    FOREIGN KEY (room_type_id) REFERENCES room_types(id) ON DELETE CASCADE,

    INDEX idx_kos (kos_id),
    INDEX idx_category (category_id),
    INDEX idx_room_type (room_type_id),
    INDEX idx_cover (kos_id, is_cover),
    INDEX idx_kos_order (kos_id, display_order)  -- ADDED
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 16. TABEL REVIEWS (Diperbaiki)
-- ============================================
CREATE TABLE IF NOT EXISTS reviews (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kos_id INT NOT NULL,
    user_id INT NOT NULL,
    rating INT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    title VARCHAR(200) NULL,
    content TEXT NULL,
    stay_from DATE NULL,
    stay_until DATE NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    is_anonymous BOOLEAN DEFAULT FALSE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,

    FOREIGN KEY (kos_id) REFERENCES kos(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    -- ADDED: Satu user hanya bisa review satu kos sekali (kecuali dihapus)
    UNIQUE KEY unique_user_kos_review (user_id, kos_id, deleted_at),
    
    INDEX idx_kos (kos_id),
    INDEX idx_user (user_id),
    INDEX idx_rating (rating),
    INDEX idx_verified (is_verified),
    INDEX idx_deleted (deleted_at),  -- ADDED
    INDEX idx_kos_rating (kos_id, rating),  -- ADDED: composite
    
    -- ADDED: Validasi stay_until >= stay_from
    CHECK (stay_until IS NULL OR stay_from IS NULL OR stay_until >= stay_from)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================
-- 17. VIEW: KOS STATISTICS (MENGGANTIKAN DENORMALISASI)
-- ============================================
CREATE OR REPLACE VIEW vw_kos_statistics AS
SELECT 
    k.id,
    k.owner_id,
    k.name,
    k.slug,
    k.gender_type,
    k.status,
    k.is_verified,
    COALESCE(AVG(r.rating), 0) as average_rating,
    COUNT(DISTINCT r.id) as total_reviews,
    SUM(rt.total_rooms) as total_rooms,
    SUM(rt.available_rooms) as available_rooms,
    MIN(rp.price) as min_price,
    MAX(rp.price) as max_price,
    k.created_at,
    k.updated_at
FROM kos k
LEFT JOIN reviews r ON k.id = r.kos_id AND r.deleted_at IS NULL
LEFT JOIN room_types rt ON k.id = rt.kos_id AND rt.deleted_at IS NULL
LEFT JOIN room_prices rp ON rt.id = rp.room_type_id AND rp.is_active = TRUE
WHERE k.deleted_at IS NULL
GROUP BY k.id;


-- ============================================
-- 18. VIEW: ADDRESS WITH COMPLETE LOCATION
-- ============================================
CREATE OR REPLACE VIEW vw_addresses_complete AS
SELECT 
    a.id,
    a.kos_id,
    a.street_address,
    a.postal_code,
    a.latitude,
    a.longitude,
    d.id as district_id,
    d.name as district_name,
    c.id as city_id,
    c.name as city_name,
    p.id as province_id,
    p.name as province_name,
    p.code as province_code
FROM addresses a
JOIN districts d ON a.district_id = d.id
JOIN cities c ON d.city_id = c.id
JOIN provinces p ON c.province_id = p.id
WHERE a.deleted_at IS NULL;


-- ============================================
-- 19. TRIGGER: UPDATE KOS STATISTICS (OPTIONAL)
-- Jika ingin tetap denormalisasi dengan trigger
-- ============================================
DELIMITER $$

-- Trigger setelah INSERT review
CREATE TRIGGER after_review_insert
AFTER INSERT ON reviews
FOR EACH ROW
BEGIN
    UPDATE kos SET 
        average_rating = (
            SELECT COALESCE(AVG(rating), 0) 
            FROM reviews 
            WHERE kos_id = NEW.kos_id AND deleted_at IS NULL
        ),
        total_reviews = (
            SELECT COUNT(*) 
            FROM reviews 
            WHERE kos_id = NEW.kos_id AND deleted_at IS NULL
        )
    WHERE id = NEW.kos_id;
END$$

-- Trigger setelah UPDATE review
CREATE TRIGGER after_review_update
AFTER UPDATE ON reviews
FOR EACH ROW
BEGIN
    UPDATE kos SET 
        average_rating = (
            SELECT COALESCE(AVG(rating), 0) 
            FROM reviews 
            WHERE kos_id = NEW.kos_id AND deleted_at IS NULL
        ),
        total_reviews = (
            SELECT COUNT(*) 
            FROM reviews 
            WHERE kos_id = NEW.kos_id AND deleted_at IS NULL
        )
    WHERE id = NEW.kos_id;
END$$

-- Trigger setelah DELETE review
CREATE TRIGGER after_review_delete
AFTER DELETE ON reviews
FOR EACH ROW
BEGIN
    UPDATE kos SET 
        average_rating = (
            SELECT COALESCE(AVG(rating), 0) 
            FROM reviews 
            WHERE kos_id = OLD.kos_id AND deleted_at IS NULL
        ),
        total_reviews = (
            SELECT COUNT(*) 
            FROM reviews 
            WHERE kos_id = OLD.kos_id AND deleted_at IS NULL
        )
    WHERE id = OLD.kos_id;
END$$

DELIMITER ;


-- ============================================
-- 20. CONTOH QUERY OPTIMIZED
-- ============================================

-- Query untuk mendapatkan kos dengan statistik lengkap
SELECT * FROM vw_kos_statistics 
WHERE status = 'active' 
ORDER BY average_rating DESC, total_reviews DESC
LIMIT 10;

-- Query untuk search kos dengan lokasi lengkap
SELECT 
    k.*,
    v.*
FROM kos k
JOIN vw_addresses_complete v ON k.id = v.kos_id
WHERE k.status = 'active'
    AND v.city_name LIKE '%Jakarta%'
    AND k.gender_type = 'mixed'
ORDER BY k.created_at DESC;