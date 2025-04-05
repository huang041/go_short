-- 創建 url_mappings 表
CREATE TABLE IF NOT EXISTS url_mappings (
  id SERIAL PRIMARY KEY,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  short_url VARCHAR(255) UNIQUE,
  original_url VARCHAR(255) NOT NULL,
  algorithm VARCHAR(50) DEFAULT 'base62',
  visits INTEGER DEFAULT 0,
  expires_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- 創建 url_mappings 索引
CREATE INDEX IF NOT EXISTS idx_url_mappings_short_url ON url_mappings(short_url);
CREATE INDEX IF NOT EXISTS idx_url_mappings_original_url ON url_mappings(original_url);
CREATE INDEX IF NOT EXISTS idx_url_mappings_expires_at ON url_mappings(expires_at);
CREATE INDEX IF NOT EXISTS idx_url_mappings_deleted_at ON url_mappings(deleted_at); 