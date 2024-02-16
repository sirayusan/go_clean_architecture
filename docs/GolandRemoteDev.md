# 概要
WSL2配下にプロジェクトフォルダを配置しており、GoLandを使用している方向けの構築手順です。  
WSL2配下に配置する理由としては、この手順を見に来ている方は理解していると思いますが動作が早くなるからです。  
https://qiita.com/BlueBaybridge/items/c1adcf1dab5da2b40b4f 
# 環境構築手順
1. ソースをクローン
```
git clone https://github.com/sirayusan/business.git
```

2. ディレクトリ移動
```
cd remote-dev
```

3. コンテナの構築  
```
docker-compose up -d
```
このような表示がでたら完了  
![image](https://github.com/sirayusan/business/assets/73060776/15593eb2-75b2-4abe-a575-1fce15fd1091)

4. コンテナに入る。
```
winpty docker container exec -it remote-dev-go-1 bash
```
5. リモートコンテナを起動し、接続URLを発行する。
```
task run-remote-dev
```
![image](https://github.com/sirayusan/go_clean_architecture/assets/73060776/736e1fd5-d090-4518-aa87-060db4c0e103)
6. tcp:// … のリンクをコピーする。
![image](https://github.com/sirayusan/go_clean_architecture/assets/73060776/ab4afcc4-b2d4-41e0-aafb-deee680bd9b7)

7. Gatewayを開き`Remote Development`の`connect to Running IDE`にペーストする。
![image](https://github.com/sirayusan/business/assets/73060776/86afc3ff-270c-4a86-bb3f-dd15859c9bf8)

8. Connectを押下し  

9. あとは流れでボタンを押して行き開くだけ
# テーブル作成とデータ投入手順
テーブル作成とデータ投入は[こちら](./migration.md) を参照してください。
# デバッグ手順
プロジェクトのデバッグ方法については、[デバッグ手順](./debug.md) を参照してください。このドキュメントでは、delveを使用した効率的なデバッグプロセスを紹介しています。
