CREATE SCHEMA IF NOT EXISTS `development` DEFAULT CHARACTER SET utf8mb4 ;
CREATE SCHEMA IF NOT EXISTS `test` DEFAULT CHARACTER SET utf8mb4 ;

-- 'user'@'%' が未作成の場合、ユーザーを作成
CREATE USER IF NOT EXISTS 'user'@'%' IDENTIFIED BY 'password';
-- 'user'@'%' に対して `test` データベースの全権限を付与
GRANT ALL PRIVILEGES ON `test`.* TO 'user'@'%';
FLUSH PRIVILEGES;

-- セッションの文字セットをutf8mb4に設定
SET NAMES 'utf8mb4';