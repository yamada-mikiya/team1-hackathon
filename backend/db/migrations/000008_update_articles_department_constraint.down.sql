-- articlesテーブルのdepartment制約を英語に戻す
ALTER TABLE articles DROP CONSTRAINT IF EXISTS articles_department_check;
ALTER TABLE articles ADD CONSTRAINT articles_department_check CHECK (department IN ('Dev', 'MKT', 'Ops'));
