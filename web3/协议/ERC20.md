## 使用
```solidity
contract MyToken is ERC20 {
    constructor() ERC20("MyToken", "MTK") {
        // 初始化逻辑，如铸造代币
        _mint(msg.sender, 1000000 * 10 ** decimals());
    }

    //......其他函数
}
```

```solidity
IERC20 token = IERC20(tokenAddress);

token.transferFrom(user, pool, amount);//池子合约从用户账户扣币。
```
## 余额变化
转账，铸造，销毁都调_update()(后两者传入from或to地址为0)
## allowance
`_spendAllowance()`减去该地址的一部分allowance，消耗allowance
## approve
授权别人（spender）可以从你的地址里花费最多多少个代币。

revert:回滚操作

有两个approve，一个有emitEvent（bool）传入，另一个没有（套在第一个外面，bool写死成true）。不想发事件时，调用第一个，传入false省gas

## solidity
抽象合约 = “带部分实现的合约模板”，不能直接部署。erc20是抽象合约，从没被部署

使用方法：（使用时必须实现抽象函数）
```solidity
contract Child is AbstractContract {
    constructor(uint _init) AbstractContract(_init) {}
}
```

`unchecked`省去节省指令，省gas

***为什么一个地址转成接口类型后就能调用函数？***

因为接口只提供 ABI 函数签名，编译器据此生成 selector 和 calldata，再由 EVM 对目标地址执行 CALL。

