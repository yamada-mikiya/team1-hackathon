-- 開発環境専用のテストデータ

-- 既存のテストデータをクリア（開発環境のみ）
TRUNCATE TABLE article_tags, articles, tags, users RESTART IDENTITY CASCADE;

-- テストユーザーの挿入
INSERT INTO users (id, name, email, affiliation, password_hash, icon_url) VALUES
(1, '田中 太郎', 'tanaka@example.com', '開発', '$2a$10$dummyhash1111111111111111111111111111111111111111', 'https://i.pravatar.cc/150?img=1'),
(2, '佐藤 花子', 'sato@example.com', 'マーケティング', '$2a$10$dummyhash2222222222222222222222222222222222222222', 'https://i.pravatar.cc/150?img=2'),
(3, '鈴木 一郎', 'suzuki@example.com', '組織管理', '$2a$10$dummyhash3333333333333333333333333333333333333333', 'https://i.pravatar.cc/150?img=3'),
(4, '高橋 美咲', 'takahashi@example.com', '開発', '$2a$10$dummyhash4444444444444444444444444444444444444444', 'https://i.pravatar.cc/150?img=4'),
(5, '山田 健太', 'yamada@example.com', 'マーケティング', '$2a$10$dummyhash5555555555555555555555555555555555555555', 'https://i.pravatar.cc/150?img=5'),
(6, '伊藤 愛', 'ito@example.com', '組織管理', '$2a$10$dummyhash6666666666666666666666666666666666666666', 'https://i.pravatar.cc/150?img=6'),
(7, '渡辺 大輔', 'watanabe@example.com', '開発', '$2a$10$dummyhash7777777777777777777777777777777777777777', 'https://i.pravatar.cc/150?img=7'),
(8, '小林 陽子', 'kobayashi@example.com', 'マーケティング', '$2a$10$dummyhash8888888888888888888888888888888888888888', 'https://i.pravatar.cc/150?img=8'),
(9, 'みきや', 'mikiya@example.com', '開発', '$2a$10$Z.INna39d0l2JmolNoMdKeFTlOsFkLTBa4buqXr040byW.NdY9aOm', 'https://i.pravatar.cc/150?img=9');
-- パスワードは "aaaaaaaa" (bcrypt hash)

-- ユーザーIDシーケンスをリセット
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));

-- タグの挿入
INSERT INTO tags (name, is_category) VALUES
('React', false),
('Go', false),
('Docker', false),
('マーケティング', false),
('インフラ', false),
('チュートリアル', true),
('ベストプラクティス', true),
('トラブルシューティング', true);

-- テスト記事の挿入
INSERT INTO articles (author_id, article_type, title, content, slug, department, status, thumbnail_url) VALUES
(
    1,
    'markdown',
    'React Hooks完全ガイド',
    E'# React Hooks完全ガイド\n\nReact Hooksの基本から応用まで、実践的な使い方を解説します。\n\n## useState\n\nコンポーネントの状態管理に使用します。\n\n```javascript\nconst [count, setCount] = useState(0);\n```\n\n## useEffect\n\n副作用を扱うためのHookです。\n\n```javascript\nuseEffect(() => {\n  console.log(\'Component mounted\');\n}, []);\n```\n\n## カスタムフック\n\n独自のフックを作成して、ロジックを再利用できます。',
    'react-hooks-guide',
    '開発',
    'public',
    'https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=800'
),
(
    2,
    'markdown',
    'マーケティング戦略2026',
    E'# マーケティング戦略2026\n\n2026年のマーケティングトレンドと戦略について解説します。\n\n## SNSマーケティング\n\n最新のSNSマーケティング手法を紹介します。\n\n- Instagram\n- TikTok\n- X (Twitter)\n\n## データ分析\n\nマーケティングデータの分析手法について説明します。',
    'marketing-strategy-2026',
    'マーケティング',
    'public',
    'https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=800'
),
(
    3,
    'markdown',
    'Kubernetes運用ガイド',
    E'# Kubernetes運用ガイド\n\nKubernetesの運用に関するベストプラクティスを紹介します。\n\n## デプロイ戦略\n\n- Rolling Update\n- Blue-Green Deployment\n- Canary Deployment',
    'kubernetes-operations',
    '組織管理',
    'internal',
    'https://images.unsplash.com/photo-1666875753105-c63a6f3bdc86?w=800'
),
(
    1,
    'markdown',
    'Go言語入門',
    E'# Go言語入門\n\nGo言語の基本的な文法と特徴を学びます。\n\n## Go言語の特徴\n\n- シンプルな文法\n- 高速なコンパイル\n- 並行処理のサポート\n\n## サンプルコード\n\n```go\npackage main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Println(\"Hello, Go!\")\n}\n```',
    'go-introduction',
    '開発',
    'internal',
    'https://images.unsplash.com/photo-1617854818583-09e7f077a156?w=800'
),
(
    1,
    'external',
    'TypeScript公式ドキュメント',
    NULL,
    'typescript-official-docs',
    '開発',
    'public',
    'https://images.unsplash.com/photo-1587620962725-abab7fe55159?w=800'
),
(
    2,
    'markdown',
    'コンテンツマーケティングの基礎',
    E'# コンテンツマーケティングの基礎\n\n効果的なコンテンツマーケティングの手法を紹介します。\n\n## コンテンツの企画\n\nターゲットオーディエンスを明確にしましょう。\n\n## 配信チャネル\n\n- ブログ\n- メールマガジン\n- SNS',
    'content-marketing-basics',
    'マーケティング',
    'draft',
    'https://images.unsplash.com/photo-1542744094-3a31f272c490?w=800'
),
(
    4,
    'markdown',
    'Docker実践入門',
    E'# Docker実践入門\n\nDockerを使ったコンテナ化の基礎を学びます。\n\n## Dockerとは\n\nアプリケーションをコンテナとして実行するプラットフォームです。\n\n## 基本コマンド\n\n```bash\ndocker run -d -p 8080:80 nginx\ndocker ps\ndocker logs <container_id>\n```',
    'docker-introduction',
    '開発',
    'public',
    'https://images.unsplash.com/photo-1605745341112-85968b19335b?w=800'
),
(
    5,
    'markdown',
    'SNS広告運用のコツ',
    E'# SNS広告運用のコツ\n\n効果的なSNS広告運用について解説します。\n\n## ターゲティング\n\n適切なターゲット設定が成功の鍵です。\n\n## クリエイティブの最適化\n\nA/Bテストを活用しましょう。',
    'sns-advertising-tips',
    'マーケティング',
    'public',
    'https://images.unsplash.com/photo-1611162617474-5b21e879e113?w=800'
),
(
    6,
    'markdown',
    'CI/CDパイプライン構築',
    E'# CI/CDパイプライン構築\n\n継続的インテグレーション・デリバリーの実践方法を紹介します。\n\n## GitHub Actionsの活用\n\n自動テストとデプロイを設定しましょう。\n\n## ベストプラクティス\n\n- テストの自動化\n- デプロイの自動化\n- ロールバック戦略',
    'cicd-pipeline',
    '組織管理',
    'internal',
    'https://images.unsplash.com/photo-1618401471353-b98afee0b2eb?w=800'
),
(
    7,
    'markdown',
    'Next.jsパフォーマンス最適化',
    E'# Next.jsパフォーマンス最適化\n\nNext.jsアプリケーションのパフォーマンスを向上させる方法を解説します。\n\n## 画像最適化\n\nNext.js Imageコンポーネントを活用しましょう。\n\n## コード分割\n\nDynamic Importでバンドルサイズを削減できます。',
    'nextjs-performance',
    '開発',
    'public',
    'https://images.unsplash.com/photo-1627398242454-45a1465c2479?w=800'
),
(
    8,
    'markdown',
    'メールマーケティング戦略',
    E'# メールマーケティング戦略\n\nメールマーケティングの効果を最大化する方法を紹介します。\n\n## セグメンテーション\n\n顧客を適切にセグメント化しましょう。\n\n## パーソナライゼーション\n\n個別化したメッセージで開封率を向上させます。',
    'email-marketing-strategy',
    'マーケティング',
    'internal',
    'https://images.unsplash.com/photo-1596526131083-e8c633c948d2?w=800'
),
(
    3,
    'markdown',
    'システム監視とアラート設定',
    E'# システム監視とアラート設定\n\nPrometheusとGrafanaを使った監視システムの構築方法を解説します。\n\n## 監視メトリクスの選定\n\n重要な指標を特定しましょう。\n\n## アラートルールの設定\n\n適切な閾値を設定してアラートを最適化します。',
    'system-monitoring',
    '組織管理',
    'public',
    'https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=800'
),
(
    9,
    'markdown',
    'GraphQL API設計のベストプラクティス',
    E'# GraphQL API設計のベストプラクティス\n\nGraphQL APIの設計と実装について解説します。\n\n## スキーマ設計\n\n効率的なスキーマ設計の原則を紹介します。\n\n```graphql\ntype User {\n  id: ID!\n  name: String!\n  email: String!\n  posts: [Post!]!\n}\n```\n\n## クエリの最適化\n\nN+1問題の解決方法とDataLoaderの活用方法を説明します。',
    'graphql-api-best-practices',
    '開発',
    'public',
    'https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=800'
),
(
    9,
    'markdown',
    'マイクロサービスアーキテクチャ入門',
    E'# マイクロサービスアーキテクチャ入門\n\nマイクロサービスの設計と運用について学びます。\n\n## サービス分割の原則\n\nドメイン駆動設計(DDD)に基づいた分割方法を解説します。\n\n## サービス間通信\n\n- REST API\n- gRPC\n- メッセージキュー\n\n## データ管理\n\n各サービスが独自のデータベースを持つパターンを紹介します。',
    'microservices-introduction',
    '開発',
    'internal',
    'https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=800'
),
(
    9,
    'external',
    'React公式ドキュメント',
    NULL,
    'react-official-docs',
    '開発',
    'public',
    'https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=800'
),
(
    9,
    'markdown',
    'WebAssemblyで高速化する',
    E'# WebAssemblyで高速化する\n\nWebAssemblyを使ったパフォーマンス向上の手法を紹介します。\n\n## WebAssemblyとは\n\nブラウザ上でネイティブに近い速度で実行できるバイナリフォーマットです。\n\n## 使用例\n\n- 画像処理\n- 動画エンコーディング\n- ゲームエンジン\n\n## RustからWebAssemblyへ\n\n```rust\n#[wasm_bindgen]\npub fn add(a: i32, b: i32) -> i32 {\n    a + b\n}\n```',
    'webassembly-performance',
    '開発',
    'internal',
    'https://images.unsplash.com/photo-1607706189992-eae578626c86?w=800'
);

-- 外部記事のURLを設定
UPDATE articles
SET external_url = 'https://www.typescriptlang.org/docs/'
WHERE slug = 'typescript-official-docs';

UPDATE articles
SET external_url = 'https://react.dev/'
WHERE slug = 'react-official-docs';

-- 記事とタグの関連付け
INSERT INTO article_tags (article_id, tag_id) VALUES
-- React Hooks完全ガイド
(1, 1),
(1, 6),
(1, 7),
-- マーケティング戦略2026
(2, 4),
-- Kubernetes運用ガイド
(3, 5),
(3, 8),
-- Go言語入門
(4, 2),
(4, 6),
-- TypeScript公式ドキュメント
(5, 1),
-- コンテンツマーケティングの基礎
(6, 4),
-- Docker実践入門
(7, 3),
(7, 6),
-- SNS広告運用のコツ
(8, 4),
-- CI/CDパイプライン構築
(9, 5),
(9, 7),
-- Next.jsパフォーマンス最適化
(10, 1),
(10, 7),
-- メールマーケティング戦略
(11, 4),
-- システム監視とアラート設定
(12, 5),
(12, 7);
-- GraphQL API設計のベストプラクティス (みきやの記事)
INSERT INTO article_tags (article_id, tag_id) VALUES
(13, 2),
(13, 7),
-- マイクロサービスアーキテクチャ入門 (みきやの記事)
(14, 2),
(14, 7),
-- React公式ドキュメント (みきやの記事)
(15, 1),
-- WebAssemblyで高速化する (みきやの記事)
(16, 7);
