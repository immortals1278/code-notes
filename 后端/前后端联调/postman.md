要同时启动前端和后端`npm run dev`
ai ide要先按”保留“，他给的修改才会生效，修改过程把后端停止运行不然修改不会生效

## TODO
u1是买家，u2是卖家，每次下单等量等价
u1下两单u2下两单，u1剩一单的冻结金额不扣（交易后双方余额符合预期）
u2下一单，u1下一单，交易后双方余额符合预期，过程中u1产生冻结金额不扣掉

后下的订单交易完成不改状态成filled

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