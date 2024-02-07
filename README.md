![Version](https://img.shields.io/badge/Version-1.0.0-green)
# プロジェクトの概要説明
Go言語の技術的な検証や動作確認用として作成しました。
## 言語
* Go1.20
## DB
* MySQL8.0
## 環境構築
* Docker
* DevContainer
## 開発支援ツール
* Air(ホットリロード)
* delve(デバッガ)
## サポートされているIDE
* VsCode
* GoLand
## ディレクトリ構成の方針
* クリーンアーキテクチャ
# 環境構築手順
[DevContainerを使う手順](./docs/DevContainer.md)  
 DevContainerを使用して開発環境を構築する方法について説明しています。  
[WSL2配下配置していてGoLandを使っている方](./docs/GolandRemoteDev.md)  
WSL2上でGoLandをリモート開発環境として使用する方法について説明しています。
# テーブル作成とデータ投入手順
テーブル作成とデータ投入は[こちら](./docs/migration.md) を参照してください。
# デバッグ手順
プロジェクトのデバッグ方法については、[デバッグ手順](./docs/debug.md) を参照してください。このドキュメントでは、delveを使用した効率的なデバッグプロセスを紹介しています。