CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(128) NOT NULL UNIQUE,
    full_name VARCHAR(100) NOT NULL,  
    username VARCHAR(20) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,  
    phone VARCHAR(20) NOT NULL,
    role ENUM('customer', 'owner', 'admin') DEFAULT 'customer',  
    is_active BOOLEAN DEFAULT TRUE,  
    last_password_change DATETIME NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    INDEX idx_email (email),
    INDEX idx_username (username),
    INDEX idx_role (role)
);

CREATE TABLE IF NOT EXISTS provinces (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL UNIQUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS cities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    province_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(10) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (province_id) REFERENCES provinces(id) ON DELETE RESTRICT,
    INDEX idx_province (province_id)
);

CREATE TABLE IF NOT EXISTS districts (
    id INT AUTO_INCREMENT PRIMARY KEY,
    city_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (city_id) REFERENCES cities(id) ON DELETE RESTRICT,
    INDEX idx_city (city_id)
);


CREATE TABLE IF NOT EXISTS kos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    owner_id INT NOT NULL,
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,  
    description TEXT,  
    
    street_address VARCHAR(255) NOT NULL, 
    district_id INT NOT NULL,
    postal_code VARCHAR(10),
    latitude DECIMAL(10, 8) NULL,  
    longitude DECIMAL(11, 8) NULL,
    
    gender_type ENUM('male', 'female', 'mixed') NOT NULL, 
    status ENUM('active', 'inactive', 'full', 'maintenance') DEFAULT 'active',
    
    is_verified BOOLEAN DEFAULT FALSE,
    average_rating DECIMAL(3, 2) DEFAULT 0.00,  
    total_reviews INT DEFAULT 0,  
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE RESTRICT,
    FOREIGN KEY (district_id) REFERENCES districts(id) ON DELETE RESTRICT,
    
    INDEX idx_owner (owner_id),
    INDEX idx_district (district_id),
    INDEX idx_gender_type (gender_type),
    INDEX idx_status (status),
    INDEX idx_slug (slug),
    INDEX idx_location (latitude, longitude),
    INDEX idx_rating (average_rating),
    
    FULLTEXT INDEX ft_name_description (name, description) 
);


CREATE TABLE IF NOT EXISTS room_types (
    id INT AUTO_INCREMENT PRIMARY KEY,
    kos_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    
    length DECIMAL(5, 2) NOT NULL,  
    width DECIMAL(5, 2) NOT NULL,   
    
    capacity INT NOT NULL DEFAULT 1,  
    total_rooms INT NOT NULL DEFAULT 0,  
    available_rooms INT NOT NULL DEFAULT 0,  
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    
    FOREIGN KEY (kos_id) REFERENCES kos(id) ON DELETE CASCADE,
    
    INDEX idx_kos (kos_id),
    INDEX idx_availability (kos_id, available_rooms)
);


CREATE TABLE IF NOT EXISTS room_prices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    room_type_id INT NOT NULL,
    period ENUM('daily', 'weekly', 'monthly', 'yearly') NOT NULL,
    price DECIMAL(12, 2) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    valid_from DATE NOT NULL,
    valid_until DATE NULL, 
    
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (room_type_id) REFERENCES room_types(id) ON DELETE CASCADE,
    
    INDEX idx_room_type (room_type_id),
    INDEX idx_period_price (period, price),
    INDEX idx_active (is_active, valid_from, valid_until)
);

CREATE TABLE IF NOT EXISTS facility_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,  
    icon VARCHAR(50) NULL,
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS facilities (
    id INT AUTO_INCREMENT PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(100) NOT NULL,  
    description TEXT NULL,
    icon VARCHAR(50) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (category_id) REFERENCES facility_categories(id) ON DELETE RESTRICT,
    
    INDEX idx_category (category_id),
    INDEX idx_active (is_active)
);

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
    INDEX idx_room_type (room_type_id)
);


CREATE TABLE IF NOT EXISTS rule_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL, 
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

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
    INDEX idx_category (category_id)
);

CREATE TABLE IF NOT EXISTS image_categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,  
    display_order INT DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

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
    INDEX idx_cover (kos_id, is_cover)
);


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
    
    INDEX idx_kos (kos_id),
    INDEX idx_user (user_id),
    INDEX idx_rating (rating),
    INDEX idx_verified (is_verified)
) ;

