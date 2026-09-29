
-- ==========================================

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Trigger dùng chung: tự động set updated_at
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ==========================================
-- 1. ORGANIZATIONS
-- ==========================================
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(150) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    phone VARCHAR(20),
    address VARCHAR(255),
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'closed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TRIGGER trg_organizations_updated_at
    BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 2. USERS (chỉ tài khoản nội bộ)
-- ==========================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    email VARCHAR(150),
    phone VARCHAR(20),
    password_hash TEXT NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    role VARCHAR(30) NOT NULL CHECK (role IN ('owner', 'manager', 'receptionist')),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'locked', 'disabled')),
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CHECK (email IS NOT NULL OR phone IS NOT NULL)
);
CREATE TRIGGER trg_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE UNIQUE INDEX uq_users_org_email ON users(organization_id, email) WHERE email IS NOT NULL;
CREATE UNIQUE INDEX uq_users_org_phone ON users(organization_id, phone) WHERE phone IS NOT NULL;

-- ==========================================
-- 3. CUSTOMERS
-- ==========================================
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    full_name VARCHAR(150) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    gender VARCHAR(20) CHECK (gender IN ('male', 'female', 'other')),
    birth_date DATE,
    note TEXT,
    total_spent NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (total_spent >= 0),
    total_visits INTEGER NOT NULL DEFAULT 0 CHECK (total_visits >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (organization_id, phone)
);
CREATE TRIGGER trg_customers_updated_at
    BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 4. STAFF (thợ/nhân viên — CHỈ quản lý & thống kê)
-- ==========================================
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    code VARCHAR(30) NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    phone VARCHAR(20),
    position VARCHAR(50) NOT NULL DEFAULT 'technician',
    commission_rate NUMERIC(5,2) NOT NULL DEFAULT 0 CHECK (commission_rate >= 0 AND commission_rate <= 100),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'on_leave')),
    hire_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    UNIQUE (organization_id, code)
);
CREATE TRIGGER trg_staff_updated_at
    BEFORE UPDATE ON staff
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 5. CATEGORIES
-- ==========================================
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('service', 'product', 'material')),
    name VARCHAR(150) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, type, name)
);
CREATE TRIGGER trg_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 6. SERVICES (Dịch vụ)
-- ==========================================
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    duration_minutes INTEGER NOT NULL CHECK (duration_minutes > 0),
    buffer_minutes INTEGER NOT NULL DEFAULT 0 CHECK (buffer_minutes >= 0),
    price NUMERIC(14,2) NOT NULL CHECK (price >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TRIGGER trg_services_updated_at
    BEFORE UPDATE ON services
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 7. PRODUCTS (Sản phẩm bán lẻ / tiêu hao nội bộ)
-- ==========================================
CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    unit VARCHAR(20),
    net_unit Varchar(10),
    net_amount NUMERIC(12,2) CHECK (net_amount >= 0),
    product_type VARCHAR(20) NOT NULL DEFAULT 'retail' CHECK (product_type IN ('retail', 'material', 'both')),
    show_on_web BOOLEAN NOT NULL DEFAULT true,
    cost_price NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (cost_price >= 0),
    sell_price NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (sell_price >= 0),
    stock_quantity NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
    reserved_quantity NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    min_stock NUMERIC(12,2) NOT NULL DEFAULT 0 CHECK (min_stock >= 0),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, sku)
);
CREATE TRIGGER trg_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Trigger ép ẩn Sản phẩm khỏi Web nếu nó chỉ là vật tư (material)
CREATE OR REPLACE FUNCTION check_product_web_visibility()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.product_type = 'material' AND NEW.show_on_web = true THEN
        NEW.show_on_web = false; 
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_product_web
    BEFORE INSERT OR UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION check_product_web_visibility();

-- ==========================================
-- 8. SERVICE_MATERIALS (Định mức vật tư/sản phẩm cho dịch vụ)
-- ==========================================
CREATE TABLE service_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity NUMERIC(12,2) NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, service_id, product_id)
);
CREATE TRIGGER trg_service_materials_updated_at
    BEFORE UPDATE ON service_materials
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 9. BLACKLISTS (Quản lý chặn / cọc)
-- ==========================================
CREATE TABLE blacklists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('phone', 'ip', 'device_id')),
    value VARCHAR(255) NOT NULL, 
    reason TEXT,
    require_deposit BOOLEAN NOT NULL DEFAULT true, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(organization_id, type, value)
);
CREATE TRIGGER trg_blacklists_updated_at
    BEFORE UPDATE ON blacklists
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 10. APPOINTMENTS (Lịch hẹn + Logic Cọc)
-- ==========================================
CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    staff_id UUID REFERENCES staff(id) ON DELETE SET NULL,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'checked_in', 'in_progress', 'completed', 'cancelled', 'no_show')),
    source VARCHAR(20) NOT NULL DEFAULT 'online' CHECK (source IN ('online', 'walk_in', 'phone')),
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    note TEXT,
    cancelled_reason VARCHAR(255),

    -- Các trường bổ sung phục vụ Blacklist & Đặt cọc
    device_id VARCHAR(255),
    client_ip VARCHAR(50),
    required_deposit_amount NUMERIC(14,2) DEFAULT 0 CHECK (required_deposit_amount >= 0),
    deposit_status VARCHAR(20) DEFAULT 'none' CHECK (deposit_status IN ('none', 'pending', 'paid', 'refunded', 'waived', 'failed')),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, code),
    CHECK (end_time > start_time)
);
CREATE TRIGGER trg_appointments_updated_at
    BEFORE UPDATE ON appointments
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 11. APPOINTMENT_ITEMS (Chi tiết dịch vụ/sản phẩm trong lịch hẹn)
-- ==========================================
CREATE TABLE appointment_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('service', 'product')),
    service_id UUID REFERENCES services(id) ON DELETE SET NULL,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    quantity NUMERIC(12,2) NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price NUMERIC(14,2) NOT NULL CHECK (unit_price >= 0),
    line_total NUMERIC(14,2) NOT NULL CHECK (line_total >= 0),
    CHECK (
        (item_type = 'service' AND service_id IS NOT NULL AND product_id IS NULL) OR
        (item_type = 'product' AND product_id IS NOT NULL AND service_id IS NULL)
    )
);

-- ==========================================
-- 11.5 SUPPLIERS (Nhà cung cấp hàng hóa cho kho)
-- ==========================================
CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE inventory_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    
    -- Các liên kết đối tượng (Tùy loại phiếu mà điền cái tương ứng, còn lại để NULL)
    supplier_id UUID REFERENCES suppliers(id),       -- Dùng khi nhập mua từ nhà cung cấp
    -- order_id trỏ tới orders (bảng tạo ở mục 13) nên FK gắn bằng ALTER TABLE phía dưới
    order_id UUID,                                    -- Dùng khi xuất bán hoặc khách trả hàng
    
    document_code VARCHAR(50) NOT NULL,              -- Mã chứng từ (VD: NK-202609-001, XK-202609-001)
    type VARCHAR(20) NOT NULL CHECK (type IN ('import', 'export', 'adjust', 'sale', 'return')),
    status VARCHAR(20) NOT NULL DEFAULT 'completed' CHECK (status IN ('completed', 'cancelled')),
    
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0,   -- Tổng tiền của chứng từ
    note TEXT,                       
    reason TEXT,               
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, document_code),
    UNIQUE (order_id)
);

-- Index tối ưu tìm kiếm theo tổ chức và loại phiếu
CREATE INDEX idx_inventory_docs_org_type ON inventory_documents(organization_id, type);
CREATE INDEX idx_inventory_docs_order_id ON inventory_documents(order_id);

CREATE TABLE inventory_document_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES inventory_documents(id) ON DELETE CASCADE, -- Khóa ngoại trỏ lên bảng chứng từ mới
    product_id UUID NOT NULL REFERENCES products(id),                             
    
    -- 📸 PHẦN SNAPSHOT (Lưu cứng lại tại thời điểm phát sinh giao dịch)
    product_name VARCHAR(255) NOT NULL,               -- Lưu tên sản phẩm lúc đó
    sku VARCHAR(100) NOT NULL,                        -- Lưu mã SKU lúc đó
    unit VARCHAR(20) NOT NULL,                        -- Đơn vị tính lớn (chai, hộp, cái...)
    net_unit VARCHAR(20),                             -- Đơn vị định lượng nhỏ (ml, g...)
    net_amount NUMERIC(12,2),                         -- Quy cách (VD: 1000 ml/chai)
    
    quantity NUMERIC(12,2) NOT NULL,                  -- Số lượng theo đơn vị lớn
    unit_price NUMERIC(14,2) NOT NULL,                -- Đơn giá tại thời điểm đó
    total_price NUMERIC(14,2) NOT NULL                -- Thành tiền (quantity * unit_price)
);

-- Index tối ưu tìm kiếm chi tiết theo chứng từ
CREATE INDEX idx_inv_doc_items_doc_id ON inventory_document_items(document_id);

-- ==========================================
-- 12. INVENTORY_TRANSACTIONS (Giao dịch kho)
-- ==========================================
CREATE TABLE inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL CHECK (type IN ('import', 'export', 'adjust', 'sale', 'return')),
    quantity NUMERIC(12,2) NOT NULL CHECK (quantity <> 0),
    reference_type VARCHAR(50),
    reference_id UUID,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ==========================================
-- 13. ORDERS (Đơn hàng)
-- ==========================================
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    code VARCHAR(20) NOT NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    appointment_id UUID REFERENCES appointments(id) ON DELETE SET NULL,
    subtotal_amount NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (subtotal_amount >= 0),
    discount_amount NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (discount_amount >= 0),
    total_amount NUMERIC(14,2) NOT NULL DEFAULT 0 CHECK (total_amount >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'pending_payment', 'paid', 'serving', 'completed', 'cancelled', 'refunded')),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (organization_id, code),
    CHECK (total_amount = GREATEST(subtotal_amount - discount_amount, 0))
);
-- FK từ chứng từ kho sang đơn hàng (orders tạo ở mục 13, sau inventory_documents)
ALTER TABLE inventory_documents
    ADD CONSTRAINT fk_inventory_documents_order_id FOREIGN KEY (order_id) REFERENCES orders(id);
CREATE TRIGGER trg_orders_updated_at
    BEFORE UPDATE ON orders
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- ==========================================
-- 14. ORDER_ITEMS (Chi tiết đơn hàng)
-- ==========================================
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('service', 'product')),
    service_id UUID REFERENCES services(id) ON DELETE SET NULL,
    product_id UUID REFERENCES products(id) ON DELETE SET NULL,
    staff_id UUID REFERENCES staff(id) ON DELETE SET NULL,
    quantity NUMERIC(12,2) NOT NULL DEFAULT 1 CHECK (quantity > 0),
    unit_price NUMERIC(14,2) NOT NULL CHECK (unit_price >= 0),
    line_total NUMERIC(14,2) NOT NULL CHECK (line_total >= 0),
    CHECK (
        (item_type = 'service' AND service_id IS NOT NULL AND product_id IS NULL) OR
        (item_type = 'product' AND product_id IS NOT NULL AND service_id IS NULL)
    )
);
-- ===========================================
-- 14.1 ORDER_ITEM_MATERIALS (Chi tiết vật tư tiêu hao cho từng order_item)
-- ===========================================
CREATE TABLE order_item_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_item_id UUID REFERENCES order_items(id),
    product_id UUID REFERENCES products(id),
    quantity NUMERIC(12, 2) NOT NULL,       -- Số lượng thực tế dùng (gram, ml)
    unit_cost_at_time NUMERIC(12, 2) NOT NULL -- GIÁ VỐN TẠI THỜI ĐIỂM DÙNG (Cực kỳ quan trọng để tính chính xác lợi nhuận đơn đó)
);

-- ==========================================
-- 15. PAYMENTS (Thanh toán)
-- ==========================================
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    method VARCHAR(20) NOT NULL CHECK (method IN ('cash', 'bank_transfer', 'gateway', 'card')),
    provider VARCHAR(50),
    provider_txn_id VARCHAR(100),
    amount NUMERIC(14,2) NOT NULL CHECK (amount >= 0),
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'succeeded', 'failed', 'refunded')),
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE UNIQUE INDEX uq_payments_provider_txn ON payments(provider, provider_txn_id) WHERE provider_txn_id IS NOT NULL;

-- ==========================================
-- 16. NOTIFICATIONS (Thông báo)
-- ==========================================
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    customer_id UUID REFERENCES customers(id) ON DELETE CASCADE,
    channel VARCHAR(20) NOT NULL CHECK (channel IN ('sms', 'email', 'zalo', 'in_app')),
    template_code VARCHAR(50) NOT NULL,
    payload JSONB NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'queued' CHECK (status IN ('queued', 'sent', 'failed')),
    retry_count INTEGER NOT NULL DEFAULT 0 CHECK (retry_count >= 0),
    sent_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ==========================================
-- INDEXES (multi-tenant + FK + nghiệp vụ)
-- ==========================================
CREATE INDEX idx_organizations_slug ON organizations(slug);

CREATE INDEX idx_users_org_id ON users(organization_id);
CREATE INDEX idx_users_deleted ON users(organization_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_customers_org_id ON customers(organization_id);
CREATE INDEX idx_customers_phone ON customers(organization_id, phone);
CREATE INDEX idx_customers_active ON customers(organization_id) WHERE deleted_at IS NULL;

CREATE INDEX idx_staff_org_id ON staff(organization_id);
CREATE INDEX idx_staff_status ON staff(organization_id, status) WHERE deleted_at IS NULL;

CREATE INDEX idx_categories_org_type ON categories(organization_id, type);

CREATE INDEX idx_services_org_id ON services(organization_id);
CREATE INDEX idx_services_category ON services(category_id);
CREATE INDEX idx_services_active ON services(organization_id, is_active) WHERE is_active = true;

CREATE INDEX idx_products_org_id ON products(organization_id);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_type ON products(organization_id, product_type);

CREATE INDEX idx_service_materials_service ON service_materials(service_id);

CREATE INDEX idx_blacklists_lookup ON blacklists(organization_id, type, value);

CREATE INDEX idx_appointments_org_id ON appointments(organization_id);
CREATE INDEX idx_appointments_time_range ON appointments(organization_id, start_time, end_time);
CREATE INDEX idx_appointments_customer ON appointments(customer_id);
CREATE INDEX idx_appointments_staff ON appointments(staff_id);
CREATE INDEX idx_appointments_status ON appointments(organization_id, status);
CREATE INDEX idx_appointments_device_id ON appointments(organization_id, device_id);
CREATE INDEX idx_appointments_client_ip ON appointments(organization_id, client_ip);

CREATE INDEX idx_appointment_items_appointment_id ON appointment_items(appointment_id);

CREATE INDEX idx_suppliers_org_id ON suppliers(organization_id);

CREATE INDEX idx_inventory_docs_supplier ON inventory_documents(supplier_id);
CREATE INDEX idx_inv_doc_items_product ON inventory_document_items(product_id);

CREATE INDEX idx_inventory_tx_org_id ON inventory_transactions(organization_id);
CREATE INDEX idx_inventory_tx_product ON inventory_transactions(product_id);
CREATE INDEX idx_inventory_tx_ref ON inventory_transactions(reference_type, reference_id);

CREATE INDEX idx_orders_org_id ON orders(organization_id);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_appointment ON orders(appointment_id);
CREATE INDEX idx_orders_status ON orders(organization_id, status);

CREATE INDEX idx_order_items_org_id ON order_items(organization_id);
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_staff ON order_items(staff_id);

CREATE INDEX idx_payments_org_id ON payments(organization_id);
CREATE INDEX idx_payments_order_id ON payments(order_id);
CREATE INDEX idx_payments_status ON payments(organization_id, status);

CREATE INDEX idx_notifications_org_id ON notifications(organization_id);
CREATE INDEX idx_notifications_status ON notifications(organization_id, status) WHERE status = 'queued';