Bluebell Go&Vue 前后端分离

使用到的:
1.JWT token
2.go-redis
3.go-mysql
4.sony flake
5.md5
6.


用户注册测试:POST
POST ==> url:localhost:8989/api/v1/signup
post data: {
  "username":"zhangsan",
  "password:"123123",
  "re_password":"123123"
}

登录测试 POST
localhost:8989/api/v1/login

{
    "username":"zhangsan",
    "password":"123123"
}
