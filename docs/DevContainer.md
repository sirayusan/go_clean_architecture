# 概要
DevContainerを使用する際の手順です。  
Vscodeでも開けます。  

### ⚠️注意⚠️
GoLandでDevContainerを使うには、Windows側のディレクトリへ配置している必要があります。  
intelliJのDevContaierに原因があり失敗するためです。
https://blog.jetbrains.com/ja/2022/02/15/automate-intellij-remotedevenv/

# 環境構築手順
1. `git clone https://github.com/sirayusan/business.git`
2.  Golandでプロジェクトを開く。
3.  `.devcontainer/devcontainer.json`を開く。
4. 左上のキューブみたいなアイコンをクリック。
![タイトルなし](https://github.com/sirayusan/business/assets/73060776/e40f04b5-158d-4e97-8694-95f62ed9ae8a)
5. Create Dev Container and Mount Sourcesをクリック。
![タイトルなし](https://github.com/sirayusan/business/assets/73060776/9b01aad6-2abb-4690-b690-c184764c22d2)
6. コンテナ作成完了するとこのような表示になるのでcontinueをクリック。
![タイトルなし](https://github.com/sirayusan/business/assets/73060776/690b8084-340b-4c43-baf5-4fef6d11efed)
7. このような表示がでて新しくウインドウが開かれる。以降は新しいウインドウで操作する。  
![image](https://github.com/sirayusan/business/assets/73060776/989e02ae-9595-451a-93e6-d637a33fb0aa)  
![image](https://github.com/sirayusan/business/assets/73060776/739bd03a-b40d-4fc6-a209-474225fbb41c)  
8. Alt + F12でターミナルを開く  
![image](https://github.com/sirayusan/business/assets/73060776/26fc15e0-09d3-43be-afa1-120889d1aa24)  
9. `air`とターミナルに入力しEnter  
このような表示が出たら環境構築完了。  
![image](https://github.com/sirayusan/business/assets/73060776/54a74657-e32a-42ab-9c1d-64fea294b58d)  
## デバッガの実行方法
1. Golandの右上のデバッガのVみたいなマークを押下する。  
![image](https://github.com/sirayusan/business/assets/73060776/85d61cf2-b5af-4b09-b522-834e2012402b)
2. Edit Configurations...をクリック  
![image](https://github.com/sirayusan/business/assets/73060776/43477116-2fd8-481e-97a8-6e768b021750)  
3. +マークから Go Remoteを選択し、スクションの内容にする  
Name: お好みで  
Host: localhost  
Port: 2345(.envのGO_DEBUG_PORT参照。)  
![image](https://github.com/sirayusan/business/assets/73060776/41a90cde-b2e8-430d-8526-285ba889a515)  
4. app/main.goにブレークポイントを設置  
5. curl http://localhost:8080/user/index