# Smarter Contract 
## 简介
将功能委托给当前最具成本效益的区块链，从而优化智能合约的性能
## mailbox
Mailbox 是 Hyperlane 在每个链上已经部署好的智能合约，可以发送/接收跨链信息（将消息发到另一条链的某个合约上）
```solidity
uint256 handleType = (uint256)((bytes32)(payload[0:32]));
```
从 payload 字节数组中提取 前32个字节（索引0到31）并转化成uint256类型，常用于解码跨链信息的头部消息
## IInterchainAccountRouter
可以调别的链上的合约的函数，可以获得源链上对应的地址
```solidity
import {Call} from "./OwnableMulticall.sol";
```
调用别的合约的结构体
## IInterchainQueryRouter
跨链查询路由器，用于在不同的区块链之间执行查询操作并获取返回数据
## OwnableMulticall
包含函数：1.对批量地址发送交易2。对批量地址发送交易后将返回的数据和对应的回调函数的数据连起来3.对一个地址批量发送交易


`_transferOwnership()` 是一个所有权转移的内部函数，用于将合约的所有权从一个地址转移到另一个地址。使用前要先import OwnableUpgradeable 


`to.call(calls[i].data)`使用低级call函数向to地址发送交易`calls[i].data`包含函数选择器和参数的字节数据

`bytes.concat(callbacks[i], returnData)`将()内的两个字节数组连接起来
## MockInterchainAccountRouter
合约作用：1.将:序列化后的"调用"组、调用者、链ID打包并放到一个列表里2.计算某用户在某链将或已部署的合约将出现在哪个地址3.将前文列表里的"调用"数据拿出来批量调用

openzzeplin工具：Address：提供地址相关的安全操作。Create2

`type(合约名).creationCode` 返回该合约的构造函数字节码,用于在其他地方复制部署相同的合约

`_salt(_origin, _sender)`通过 `_origin` (链ID) 和 `_sender` (用户地址) 的组合，为每个用户在每条链生成唯一的确定性地址。

```solidity
Call[] memory calls = abi.decode(pendingCall.serializedCalls, (Call[]));
```
将``pendingCall.serializedCalls``(序列化后的字节数组)解码成``Call[]``






//读完剩下的一点代码
//TODO将所有“交易”换成“调用”
//梳理整个项目工作过程