## TODO
cancel前端要自动调getid获取订单id给用户用

## POSTMAN
GET：参数在 URL，用于获取/取消数据，不修改主体

POST：参数在 Body，用于创建/提交数据

## 步骤
1. **确保程序运行**
```bash
go run main.go
```
看到无报错，终端卡住即启动成功

2. **打开 Postman**
   - 不用创建 workspace，直接点 `+` 新标签

3. **测试登录**
   - Method: `POST`
   - URL: `http://localhost:8080/api/login`
   - Body → raw → JSON:
   ```json
   {"user_id": "123"}
   ```
   - 点 Send

4. **看结果**
   返回 `{"code":0,"msg":"ok",...}` 即成功