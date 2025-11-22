ALTER TABLE articles DROP COLUMN category;
ALTER TABLE articles DROP FOREIGN KEY fk_articles_author;
ALTER TABLE articles MODIFY COLUMN author_id VARCHAR(36) NULL;
ALTER TABLE articles 
ADD CONSTRAINT articles_ibfk_1
FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL;