# majiang

本地运行
1. 修改 `./etc/conf.yaml`配置数据库
2. 建表 `./majiang.sql`
3.  `go build ./main.go`
4. `./main.exe`
5. > 注册   
   > POST localhost:8080/user/register  
   > username:    
   > password:
6. > 登录 localhost:8080/user/login    
   > username:    
   > password:    
7. 创建四个号连接ws http://wstool.js.org/    
   ws://127.0.0.1:8080/join/1?accessToken=[accessToken]    
8. 命令    
   {"type":1,"content":"in"}     
   {"type":1,"content":"out,1"}    
   {"type":1,"content":"ready"}    
   {"type":1,"content":"start"}    
   {"type":1,"content":"ignore"}    
   {"type":1,"content":"peng,1"}    
   {"type":1,"content":"say,1,2,3"}    
9. 效果图
   <img src="https://github.com/XiaoTe33/majiang/tree/master/xiaoguo.jpg">