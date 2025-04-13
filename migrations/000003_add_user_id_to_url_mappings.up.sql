-- 為 url_mappings 表添加 user_id 欄位
ALTER TABLE url_mappings
ADD COLUMN user_id INTEGER; -- 使用 INTEGER 對應 GORM 的 uint

-- 為 user_id 創建索引，方便未來根據使用者查詢其創建的連結
CREATE INDEX IF NOT EXISTS idx_url_mappings_user_id ON url_mappings(user_id); 