-- 刪除 user_id 上的索引
DROP INDEX IF EXISTS idx_url_mappings_user_id;

-- 刪除 user_id 欄位
ALTER TABLE url_mappings
DROP COLUMN IF EXISTS user_id; 