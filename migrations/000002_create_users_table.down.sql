-- 刪除索引
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- 刪除表格
DROP TABLE IF EXISTS users; 