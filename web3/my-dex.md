### 价格数据的存储
chainlink的变异系数和数据存在合约中，本地的价格的变异系数存在合约中，数据存在链下

### 链下监听
链下监听需要一个provider对象（websocket连接的url），合约的地址，事件与函数的签名，前二者我在编写时都是虚构的，不知道是不是这个原因不能测试


构造contract对象需要合约地址，abi（合约的事件与函数的签名），provider对象

#### 监听所有该合约的事件
由于pair合约会被部署很多份，不能每部署一次就更新一次链下监听代码

```solidity
const event = contract.interface.parseLog(log);
```
`log`事件对象，`log.address`可以拿到事件发生的合约地址,`log`解析成event对象就可以拿事件的参数
```solidity
const eventSignature = "newPrice(uint256)";
const eventTopic = ethers.utils.id(eventSignature);//将事件包装成主题
provider.on({
    topics: [eventTopic]  // ← 这是过滤规则，告诉节点"只给我这些主题的事件"
}, (log) => {
    // 箭头函数逻辑：创建合约，创建事件对象，调用函数
});
```
箭头函数：contract监听到了对应事件就执行的函数，箭头前括号内的是参数

provider监听到每一个事件都会执行一遍箭头函数，若监听到别人的事件会因为abi不匹配报错
### 链下调用
合约定义函数后，链下监听到相关事件自动调用该函数（无需在合约内主动调用）

合约调用要先链接一个钱包来支付gas，用私钥签名证明你的身份获得权限（有的话）
#### 钱包对象
需要私钥和provider对象，私钥从环境变量中读取
```solidity
const privateKey = process.env.PRIVATE_KEY;
```
process.env是node.js的环境变量，要先创建env文件并在里面写入私钥，安装dotenv包，最后在js文件开头加载环境变量`require('dotenv').config(); `

const privateKey = process.env.PRIVATE_KEY;

