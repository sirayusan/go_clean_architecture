# 概要
WSL2配下にプロジェクトフォルダを配置しており、GoLandを使用している方向けの構築手順です。  
WSL2配下に配置する理由としては、この手順を見に来ている方は理解していると思いますが動作が早くなるからです。  
https://qiita.com/BlueBaybridge/items/c1adcf1dab5da2b40b4f 
# 環境構築手順
1. `git clone https://github.com/sirayusan/business.git`
2. `cd remote-dev`
3. `docker compose build --no-cache`  
このような表示がでたら完了 
![image](https://github.com/sirayusan/business/assets/73060776/36d6b237-4d8e-4b06-8498-2354a371eef0)
4. `docker compose up -d`  
このような表示がでたら完了 
![image](https://github.com/sirayusan/business/assets/73060776/15593eb2-75b2-4abe-a575-1fce15fd1091)
5. Docker for Desktopを起動し、goコンテナのログを確認し、tcp:// … のリンクをコピーする
![image](https://github.com/sirayusan/business/assets/73060776/3b1a7e6f-6208-432a-aaba-b35856388c77)
6. Gatewayを開き`Remote Development`の`connect to Running IDE`にペーストする。
![image](https://github.com/sirayusan/business/assets/73060776/86afc3ff-270c-4a86-bb3f-dd15859c9bf8)
7. Connectを押下し  
8. あとは流れでボタンを押して行き開くだけ
# デバッグ手順
プロジェクトのデバッグ方法については、[デバッグ手順](./debug.md) を参照してください。このドキュメントでは、delveを使用した効率的なデバッグプロセスを紹介しています。