-- Seed salon duy nhất (single-salon mode) — đầy đủ 16 bảng, dữ liệu demo liên kết nhau.
-- SaaS sau này: mỗi salon mới = 1 row organizations + dữ liệu riêng theo organization_id.
-- Chạy lại an toàn (idempotent): mọi INSERT đều có ON CONFLICT DO NOTHING.
-- Mật khẩu demo cho cả 3 user: admin123 (bcrypt).
-- Lưu ý: đổi slug/phone/address cho đúng salon thật trước khi chạy production.

-- 0. Biến dùng chung: organization demo
-- '11111111-1111-1111-1111-111111111111'

-- ==========================================
-- 1. ORGANIZATIONS
-- ==========================================
INSERT INTO organizations (id, name, slug, phone, address, timezone, status)
VALUES (
  '11111111-1111-1111-1111-111111111111',
  'Salon Mặc Định',
  'default-salon',
  '0900000000',
  'TP. Hồ Chí Minh',
  'Asia/Ho_Chi_Minh',
  'active'
)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 2. USERS (tài khoản nội bộ — password: admin123)
-- ==========================================
INSERT INTO users (id, organization_id, email, phone, password_hash, full_name, role, status)
VALUES
  ('22222222-2222-2222-2222-222222222201', '11111111-1111-1111-1111-111111111111',
   'owner@salon.test', '0900000001',
   '$2a$10$Db14SJVmANoblwby4Yrzfe9RL6cN7Fc5pASfrqObu29IPpp2bM1G2',
   'Chủ Salon', 'owner', 'active'),
  ('22222222-2222-2222-2222-222222222202', '11111111-1111-1111-1111-111111111111',
   'manager@salon.test', '0900000002',
   '$2a$10$Db14SJVmANoblwby4Yrzfe9RL6cN7Fc5pASfrqObu29IPpp2bM1G2',
   'Quản Lý', 'manager', 'active'),
  ('22222222-2222-2222-2222-222222222203', '11111111-1111-1111-1111-111111111111',
   'reception@salon.test', '0900000003',
   '$2a$10$Db14SJVmANoblwby4Yrzfe9RL6cN7Fc5pASfrqObu29IPpp2bM1G2',
   'Lễ Tân', 'receptionist', 'active')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 3. CUSTOMERS
-- ==========================================
INSERT INTO customers (id, organization_id, full_name, phone, gender, birth_date, note, total_spent, total_visits)
VALUES
  ('33333333-3333-3333-3333-333333333301', '11111111-1111-1111-1111-111111111111',
   'Nguyễn Thị C', '0912345678', 'female', '1995-05-20', 'Khách VIP', 600000, 2),
  ('33333333-3333-3333-3333-333333333302', '11111111-1111-1111-1111-111111111111',
   'Lê Văn D', '0987654321', 'male', '1990-11-02', NULL, 0, 0),
  ('33333333-3333-3333-3333-333333333303', '11111111-1111-1111-1111-111111111111',
   'Trần Thị E', '0933333333', 'other', NULL, 'Dị ứng thuốc nhuộm', 0, 0)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 4. STAFF
-- ==========================================
INSERT INTO staff (id, organization_id, code, full_name, phone, position, commission_rate, status, hire_date)
VALUES
  ('44444444-4444-4444-4444-444444444401', '11111111-1111-1111-1111-111111111111',
   'NV001', 'Nguyễn Văn A', '0901234567', 'Thợ chính', 15, 'active', '2023-01-15'),
  ('44444444-4444-4444-4444-444444444402', '11111111-1111-1111-1111-111111111111',
   'NV002', 'Trần Thị B', '0907654321', 'Thợ phụ', 10, 'active', '2024-03-01'),
  ('44444444-4444-4444-4444-444444444403', '11111111-1111-1111-1111-111111111111',
   'NV003', 'Lê Văn C', '0911111111', 'Lễ tân', 5, 'inactive', '2024-06-01')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 5. CATEGORIES (id cố định để services/products tham chiếu)
-- ==========================================
INSERT INTO categories (id, organization_id, type, name, description, is_active)
VALUES
  ('55555555-5555-5555-5555-555555555501', '11111111-1111-1111-1111-111111111111',
   'service', 'Cắt tóc', 'Cắt, tạo kiểu tóc', true),
  ('55555555-5555-5555-5555-555555555502', '11111111-1111-1111-1111-111111111111',
   'service', 'Nhuộm tóc', 'Nhuộm, highlight, phủ bạc', true),
  ('55555555-5555-5555-5555-555555555503', '11111111-1111-1111-1111-111111111111',
   'service', 'Chăm sóc da', 'Facial, massage mặt', true),
  ('55555555-5555-5555-5555-555555555504', '11111111-1111-1111-1111-111111111111',
   'product', 'Dầu gội', 'Sản phẩm chăm sóc tóc bán lẻ', true)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 6. SERVICES
-- ==========================================
INSERT INTO services (id, organization_id, category_id, name, description, duration_minutes, buffer_minutes, price, is_active)
VALUES
  ('66666666-6666-6666-6666-666666666601', '11111111-1111-1111-1111-111111111111',
   '55555555-5555-5555-5555-555555555501', 'Cắt tóc nam', 'Cắt + gội + sấy', 45, 10, 150000, true),
  ('66666666-6666-6666-6666-666666666602', '11111111-1111-1111-1111-111111111111',
   '55555555-5555-5555-5555-555555555502', 'Nhuộm phủ bạc', 'Nhuộm + hấp dầu', 120, 15, 500000, true),
  ('66666666-6666-6666-6666-666666666603', '11111111-1111-1111-1111-111111111111',
   '55555555-5555-5555-5555-555555555503', 'Facial cơ bản', 'Rửa mặt + massage + đắp mặt nạ', 60, 10, 300000, true)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 7. PRODUCTS
-- Lưu ý: product_type='material' bị trigger ép show_on_web=false.
-- ==========================================
INSERT INTO products (id, organization_id, category_id, sku, name, unit, product_type, show_on_web, cost_price, sell_price, stock_quantity, min_stock, is_active)
VALUES
  ('77777777-7777-7777-7777-777777777701', '11111111-1111-1111-1111-111111111111',
   '55555555-5555-5555-5555-555555555504', 'DG-BUOI-500', 'Dầu gội bưởi 500ml', 'chai',
   'retail', true, 80000, 150000, 50, 5, true),
  ('77777777-7777-7777-7777-777777777702', '11111111-1111-1111-1111-111111111111',
   '55555555-5555-5555-5555-555555555504', 'THUOC-NHUOM-DEN', 'Thuốc nhuộm đen', 'tuýp',
   'material', false, 120000, 0, 20, 2, true),
  ('77777777-7777-7777-7777-777777777703', '11111111-1111-1111-1111-111111111111',
   '55555555-5555-5555-5555-555555555504', 'KEM-U-300', 'Kem ủ tóc 300ml', 'hũ',
   'both', true, 60000, 120000, 30, 3, true)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 8. SERVICE_MATERIALS (định mức vật tư cho dịch vụ)
-- ==========================================
INSERT INTO service_materials (id, organization_id, service_id, product_id, quantity)
VALUES
  ('88888888-8888-8888-8888-888888888801', '11111111-1111-1111-1111-111111111111',
   '66666666-6666-6666-6666-666666666602', '77777777-7777-7777-7777-777777777702', 1),
  ('88888888-8888-8888-8888-888888888802', '11111111-1111-1111-1111-111111111111',
   '66666666-6666-6666-6666-666666666601', '77777777-7777-7777-7777-777777777703', 0.5)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 9. BLACKLISTS
-- ==========================================
INSERT INTO blacklists (id, organization_id, type, value, reason, require_deposit)
VALUES
  ('99999999-9999-9999-9999-999999999901', '11111111-1111-1111-1111-111111111111',
   'phone', '0999999999', 'Spam đặt lịch rồi không đến', true),
  ('99999999-9999-9999-9999-999999999902', '11111111-1111-1111-1111-111111111111',
   'ip', '192.168.1.200', 'Bot đặt lịch hàng loạt', true)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 10. APPOINTMENTS
-- ==========================================
INSERT INTO appointments (id, organization_id, code, customer_id, staff_id, start_time, end_time, status, source, total_amount, note, device_id, client_ip, required_deposit_amount, deposit_status)
VALUES
  -- Lịch đã hoàn thành trong quá khứ
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01', '11111111-1111-1111-1111-111111111111',
   'APT-0001', '33333333-3333-3333-3333-333333333301', '44444444-4444-4444-4444-444444444401',
   '2026-09-01 10:00:00+07', '2026-09-01 12:00:00+07',
   'completed', 'walk_in', 650000, 'Khách hài lòng', NULL, NULL, 0, 'none'),
  -- Lịch xác nhận trong tương lai, yêu cầu cọc
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa02', '11111111-1111-1111-1111-111111111111',
   'APT-0002', '33333333-3333-3333-3333-333333333302', '44444444-4444-4444-4444-444444444402',
   '2026-10-01 14:00:00+07', '2026-10-01 15:00:00+07',
   'confirmed', 'online', 150000, NULL, 'device-demo-001', '203.0.113.10', 100000, 'pending')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 11. APPOINTMENT_ITEMS
-- ==========================================
INSERT INTO appointment_items (id, organization_id, appointment_id, item_type, service_id, product_id, quantity, unit_price, line_total)
VALUES
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbb0001', '11111111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01', 'service', '66666666-6666-6666-6666-666666666602', NULL, 1, 500000, 500000),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbb0002', '11111111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01', 'product', NULL, '77777777-7777-7777-7777-777777777701', 1, 150000, 150000),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbb0003', '11111111-1111-1111-1111-111111111111',
   'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa02', 'service', '66666666-6666-6666-6666-666666666601', NULL, 1, 150000, 150000)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 11.5 SUPPLIERS (Nhà cung cấp)
-- ==========================================
INSERT INTO suppliers (id, organization_id, name, phone, address, note)
VALUES
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0001', '11111111-1111-1111-1111-111111111111',
   'L''Oréal Professionnel VN', '02838238638', 'Tầng 10, Tòa nhà Saigon Centre, Quận 1, TP.HCM', 'NCC hóa chất nhuộm chính'),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0002', '11111111-1111-1111-1111-111111111111',
   'Kérastase Paris Flagship VN', '02838230001', 'Đồng Khởi, Quận 1, TP.HCM', 'NCC dầu gội cao cấp'),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0003', '11111111-1111-1111-1111-111111111111',
   'Olaplex Vietnam Official', '02839123456', 'Thảo Điền, Quận 2, TP.HCM', 'NCC phục hồi tóc')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 11.6 INVENTORY_DOCUMENTS + ITEMS (Chứng từ kho mẫu)
-- total_amount = tổng total_price các dòng
-- ==========================================
INSERT INTO inventory_documents (id, organization_id, supplier_id, order_id, document_code, type, total_amount, note, created_by)
VALUES
  ('ffffffff-ffff-ffff-ffff-ffffffff0001', '11111111-1111-1111-1111-111111111111',
   'eeeeeeee-eeee-eeee-eeee-eeeeeeee0001', NULL, 'NK-202609-001', 'import', 5200000, 'Nhập bổ sung đầu tháng 9',
   '22222222-2222-2222-2222-222222222201'),
  ('ffffffff-ffff-ffff-ffff-ffffffff0002', '11111111-1111-1111-1111-111111111111',
   'eeeeeeee-eeee-eeee-eeee-eeeeeeee0002', NULL, 'NK-202609-002', 'import', 2400000, 'Nhập kem ủ + thuốc nhuộm',
   '22222222-2222-2222-2222-222222222201')
ON CONFLICT DO NOTHING;

INSERT INTO inventory_document_items (id, document_id, product_id, product_name, sku, unit, net_unit, net_amount, quantity, unit_price, total_price)
VALUES
  ('00000000-0000-4000-8000-000000000001', 'ffffffff-ffff-ffff-ffff-ffffffff0001',
   '77777777-7777-7777-7777-777777777701', 'Dầu gội bưởi 500ml', 'DG-BUOI-500', 'chai', NULL, NULL, 50, 80000, 4000000),
  ('00000000-0000-4000-8000-000000000002', 'ffffffff-ffff-ffff-ffff-ffffffff0001',
   '77777777-7777-7777-7777-777777777702', 'Thuốc nhuộm đen', 'THUOC-NHUOM-DEN', 'tuýp', NULL, NULL, 10, 120000, 1200000),
  ('00000000-0000-4000-8000-000000000003', 'ffffffff-ffff-ffff-ffff-ffffffff0002',
   '77777777-7777-7777-7777-777777777703', 'Kem ủ tóc 300ml', 'KEM-U-300', 'hũ', NULL, NULL, 30, 60000, 1800000),
  ('00000000-0000-4000-8000-000000000004', 'ffffffff-ffff-ffff-ffff-ffffffff0002',
   '77777777-7777-7777-7777-777777777702', 'Thuốc nhuộm đen', 'THUOC-NHUOM-DEN', 'tuýp', NULL, NULL, 5, 120000, 600000)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 12. INVENTORY_TRANSACTIONS
-- ==========================================
INSERT INTO inventory_transactions (id, organization_id, product_id, type, quantity, reference_type, reference_id, created_by)
VALUES
  ('cccccccc-cccc-cccc-cccc-cccccccc0001', '11111111-1111-1111-1111-111111111111',
   '77777777-7777-7777-7777-777777777701', 'import', 50, 'purchase_order', NULL,
   '22222222-2222-2222-2222-222222222201'),
  ('cccccccc-cccc-cccc-cccc-cccccccc0002', '11111111-1111-1111-1111-111111111111',
   '77777777-7777-7777-7777-777777777701', 'sale', -1, 'order', 'dddddddd-dddd-dddd-dddd-dddddddd0001',
   '22222222-2222-2222-2222-222222222203')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 13. ORDERS (total = subtotal - discount)
-- ==========================================
INSERT INTO orders (id, organization_id, code, customer_id, appointment_id, subtotal_amount, discount_amount, total_amount, status, created_by)
VALUES
  ('dddddddd-dddd-dddd-dddd-dddddddd0001', '11111111-1111-1111-1111-111111111111',
   'ORD-0001', '33333333-3333-3333-3333-333333333301', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaa01',
   650000, 50000, 600000, 'paid', '22222222-2222-2222-2222-222222222203')
ON CONFLICT DO NOTHING;

-- ==========================================
-- 14. ORDER_ITEMS
-- ==========================================
INSERT INTO order_items (id, organization_id, order_id, item_type, service_id, product_id, staff_id, quantity, unit_price, line_total)
VALUES
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0001', '11111111-1111-1111-1111-111111111111',
   'dddddddd-dddd-dddd-dddd-dddddddd0001', 'service', '66666666-6666-6666-6666-666666666602', NULL,
   '44444444-4444-4444-4444-444444444401', 1, 500000, 500000),
  ('eeeeeeee-eeee-eeee-eeee-eeeeeeee0002', '11111111-1111-1111-1111-111111111111',
   'dddddddd-dddd-dddd-dddd-dddddddd0001', 'product', NULL, '77777777-7777-7777-7777-777777777701',
   NULL, 1, 150000, 150000)
ON CONFLICT DO NOTHING;

-- ==========================================
-- 15. PAYMENTS
-- ==========================================
INSERT INTO payments (id, organization_id, order_id, method, provider, provider_txn_id, amount, status, paid_at)
VALUES
  ('ffffffff-ffff-ffff-ffff-ffffffff0001', '11111111-1111-1111-1111-111111111111',
   'dddddddd-dddd-dddd-dddd-dddddddd0001', 'cash', NULL, NULL, 600000, 'succeeded', now())
ON CONFLICT DO NOTHING;

-- ==========================================
-- 16. NOTIFICATIONS
-- ==========================================
INSERT INTO notifications (id, organization_id, customer_id, channel, template_code, payload, status, retry_count)
VALUES
  ('12345678-1234-1234-1234-123456780001', '11111111-1111-1111-1111-111111111111',
   '33333333-3333-3333-3333-333333333302', 'sms', 'appointment_reminder',
   '{"appointment_code": "APT-0002", "message": "Nhắc lịch hẹn ngày 01/10/2026 lúc 14:00"}', 'queued', 0),
  ('12345678-1234-1234-1234-123456780002', '11111111-1111-1111-1111-111111111111',
   '33333333-3333-3333-3333-333333333301', 'in_app', 'order_paid',
   '{"order_code": "ORD-0001", "message": "Đơn hàng đã thanh toán thành công"}', 'sent', 0)
ON CONFLICT DO NOTHING;
