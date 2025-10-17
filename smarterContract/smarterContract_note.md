# Smarter Contract 
## 简介
将功能委托给当前最具成本效益的区块链，从而优化智能合约的性能
## 跨链调用
消息传递 + 目标链本地执行。因为不同链上的相同地址指向不同合约，所以不能直接用合约地址（合约实例）跨链调用
## mailbox
Mailbox 是 Hyperlane 在每个链上已经部署好的智能合约，可以发送/接收跨链信息（将消息发到另一条链的某个合约上）
```solidity
uint256 handleType = (uint256)((bytes32)(payload[0:32]));
```
从 payload 字节数组中提取 前32个字节（索引0到31）并转化成uint256类型，常用于解码跨链信息的头部消息.同理`payload[32,64]`(索引32到63)
## IInterchainAccountRouter
可以调别的链上的合约的函数，可以获得源链上对应的地址
```solidity
import {Call} from "./OwnableMulticall.sol";
```
调用别的合约的结构体
## IInterchainQueryRouter
跨链查询路由器，用于在不同的区块链之间执行查询操作并获取返回数据
## OwnableMulticall
在目标链上创建该合约，作为用户的跨链账户

包含函数：1.对批量地址发送调用2.对批量地址发送调用后将返回的数据和对应的回调函数的数据连起来3.对一个地址批量发送调用


`_transferOwnership()` 是一个所有权转移的内部函数，用于将合约的所有权从一个地址转移到另一个地址。使用前要先import OwnableUpgradeable 


`to.call(calls[i].data)`使用低级call函数向to地址发送调用`calls[i].data`包含函数选择器和参数的字节数据

`bytes.concat(callbacks[i], returnData)`将()内的两个字节数组连接起来
## MockInterchainAccountRouter
用户输入要调用的方法和目标链，然后帮用户在目标链调用

合约作用：1.将:序列化后的"要调用"组、调用者、目标链ID打包并放到一个列表里2.计算某用户在某链将或已部署的合约（这里是用户的跨链账户）将出现在哪个地址，若没有则帮用户部署3.将前文列表里的"调用"数据拿出来批量调用

但实际上用户的跨链地址部署并没有部署在别的链上（而是本链），该合约仅用与测试

openzzeplin工具：Address：提供地址相关的安全操作。Create2

`type(合约名).creationCode` 返回该合约的构造函数字节码,用于在其他地方复制部署相同的合约

`_salt(_origin, _sender)`通过 `_origin` (链ID) 和 `_sender` (用户地址) 的组合，为每个用户在每条链生成唯一的确定性地址。

```solidity
Call[] memory calls = abi.decode(pendingCall.serializedCalls, (Call[]));
```
将``pendingCall.serializedCalls``(序列化后的字节数组)解码成``Call[]``

import接口后还需要合约地址才能调用接口的函数
## Owner
调IInterchainAccountRouter函数：1.设置远程链上合约的手续费2.更改远程链上合约的所有权
```solidity
Call({
            to: target,
            data: abi.encodeWithSelector(0x69fe0e2d, newFee)
        });
```
target: 目标链上的合约地址,0x69fe0e2d是目标函数签名，newFee是参数
## OwnerReader
调IInterchainQueryRouter的函数查询别的链的信息

直接在合约上面定义一个接口就不需要创建一个文件再import（函数很少的情况）
```solidity
Ownable.owner.selector
```
获取ownable合约中owner函数的选择器编码，用于函数调用




//三个跨链通信合约鬼用没有




