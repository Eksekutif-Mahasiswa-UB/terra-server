ALTER TABLE articles DROP FOREIGN KEY articles_ibfk_1;
ALTER TABLE articles MODIFY COLUMN author_id VARCHAR(36) NOT NULL;
ALTER TABLE articles 
ADD CONSTRAINT fk_articles_author 
FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE articles 
ADD COLUMN category VARCHAR(100) NOT NULL DEFAULT 'General' AFTER image_url;