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