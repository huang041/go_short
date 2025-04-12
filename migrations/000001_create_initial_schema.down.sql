-- 刪除索引
DROP INDEX IF EXISTS idx_url_mappings_short_url;
DROP INDEX IF EXISTS idx_url_mappings_original_url;
DROP INDEX IF EXISTS idx_url_mappings_expires_at;
DROP INDEX IF EXISTS idx_url_mappings_deleted_at;

-- 刪除表格
DROP TABLE IF EXISTS url_mappings; 