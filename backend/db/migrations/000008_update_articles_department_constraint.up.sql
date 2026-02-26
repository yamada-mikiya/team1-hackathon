-- articlesテーブルのdepartment制約を日本語に更新
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_department_check;
ALTER TABLE articles ADD CONSTRAINT articles_department_check CHECK (department IN ('開発', 'マーケティング', '組織管理'));
