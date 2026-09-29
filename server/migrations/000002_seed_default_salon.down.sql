-- Down cho 000002_seed_default_salon: xóa dữ liệu demo theo thứ tự ngược phụ thuộc
-- (bảng con trước, bảng cha sau) để không vi phạm FK.

DELETE FROM notifications WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM payments WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM order_item_materials WHERE order_item_id IN (SELECT id FROM order_items WHERE organization_id = '11111111-1111-1111-1111-111111111111');
DELETE FROM order_items WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM orders WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM inventory_transactions WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM inventory_document_items WHERE document_id IN (SELECT id FROM inventory_documents WHERE organization_id = '11111111-1111-1111-1111-111111111111');
DELETE FROM inventory_documents WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM suppliers WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM appointment_items WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM appointments WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM blacklists WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM service_materials WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM products WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM services WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM categories WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM staff WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM customers WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM users WHERE organization_id = '11111111-1111-1111-1111-111111111111';
DELETE FROM organizations WHERE id = '11111111-1111-1111-1111-111111111111';
