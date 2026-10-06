-- 演示种子数据
-- 使用方式：先启动一次后端（AutoMigrate 自动建表），再执行本脚本：
--   mysql -u<用户> -p <库名> < scripts/seed.sql
-- 所有演示账号的密码均为 password123（bcrypt 存储）

SET FOREIGN_KEY_CHECKS = 0;
DELETE FROM users;
DELETE FROM items;
DELETE FROM claims;
DELETE FROM announcements;
SET FOREIGN_KEY_CHECKS = 1;

-- 账号：李四(失物招领管理员) / 王五(系统管理员) / 赵六(学生)
INSERT INTO users (id, created_at, updated_at, username, password_hash, name, role, status) VALUES
(1, NOW(), NOW(), '302026000001', '$2a$10$fVPVRZkomtC6HzcH6/ZKxeS4GSpRBTjwlIwl00WpL0WaGJ3GiKi7q', '李四', 'lost_admin', 'active'),
(2, NOW(), NOW(), '302026000002', '$2a$10$fVPVRZkomtC6HzcH6/ZKxeS4GSpRBTjwlIwl00WpL0WaGJ3GiKi7q', '王五', 'system_admin', 'active'),
(3, NOW(), NOW(), '302026000003', '$2a$10$fVPVRZkomtC6HzcH6/ZKxeS4GSpRBTjwlIwl00WpL0WaGJ3GiKi7q', '赵六', 'student', 'active');
ALTER TABLE users AUTO_INCREMENT = 100;

-- 物品：4 条审核通过（含各物品状态）+ 1 条待审核（供管理员审核演示）
INSERT INTO items (id, created_at, updated_at, type, title, description, category, location, contact, image_urls, review_status, reject_reason, item_status, user_id) VALUES
(1, NOW(), NOW(), 'lost',  '黑色校园卡',   '在图书馆一楼自习区捡到一张黑色校园卡，卡套是蓝色的，请失主联系我。', 'id_card', '图书馆一楼',     '微信: lisi2026',    '[]', 'approved', NULL, 'open',     1),
(2, NOW(), NOW(), 'found', '蓝色长柄雨伞', '在教学楼A座门口拾到一把蓝色长柄雨伞，伞面有白色条纹。',             'daily',   '教学楼A座门口',  'QQ: 123456789',     '[]', 'approved', NULL, 'open',     1),
(3, NOW(), NOW(), 'lost',  '白色无线耳机', '在二食堂二楼餐桌旁丢失一只白色无线耳机充电盒。',                     'phone',   '二食堂二楼',     '电话: 13800000000', '[]', 'approved', NULL, 'claimed',  2),
(4, NOW(), NOW(), 'found', '高等数学课本', '在体育馆看台拾到一本高等数学课本，内有姓名贴纸。',                   'book',    '体育馆',         '请通过平台联系',     '[]', 'approved', NULL, 'resolved', 2),
(5, NOW(), NOW(), 'lost',  '一串钥匙',     '在行政楼一楼大厅捡到一串钥匙，挂有绿色小熊挂件。',                   'key',     '行政楼一楼',     '请通过平台联系',     '[]', 'pending',  NULL, 'open',     3);
ALTER TABLE items AUTO_INCREMENT = 100;

-- 公告：1 条已发布 + 1 条未发布（供公告管理演示）
INSERT INTO announcements (id, created_at, updated_at, title, content, published, author_id) VALUES
(1, NOW(), NOW(), '平台使用须知', '请勿在公开内容中填写身份证号等敏感信息，认领时请携带有效证件核对。', true, 2),
(2, NOW(), NOW(), '招领信息审核规范', '管理员将在 24 小时内完成审核，审核结果会以站内方式同步。', false, 2);
ALTER TABLE announcements AUTO_INCREMENT = 100;
