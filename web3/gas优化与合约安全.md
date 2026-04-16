# gas优化
## mapping+=
凡是看到有人对 mapping 直接用 +=、-=，99% 情况下都应该改成先读到局部变量，再写回去
```solidity
function deposit() public payable {
    uint256 current = balances[msg.sender];
    balances[msg.sender] = current + msg.value;
    }
```
## 位压缩
每个 ***storage*** slot = 256 bits = 32 bytes

每次读/写 storage slot 都要消耗 gas

如果一个变量 ≤ 32 bytes，它尽量和前后变量挤在一个 slot 里

如果放不下了，就开启新 slot
```solidity
struct Packed {
    uint128 a;   // 16 bytes → 占 slot 0 的前 16 bytes
    uint128 b;   // 16 bytes → 刚好塞进 slot 0 的后 16 bytes
}
```
## 循环优化
优化前每次循环都要读取一遍长度变量
```solidity
// ❌ 非优化
for (uint256 i = 0; i < arr.length; i++) {
    ...
}
// ✅ 优化
uint256 len = arr.length;
for (uint i = 0; i < len; ++i) {
    ...
}
```
## 函数可见性选择
external 比 public 更节省 gas，适用于仅被外部调用的函数。

仅被外部调用：这个函数只能从合约外面被调用（比如用户钱包、其他合约、DApp 前端通过交易调用），而在当前合约内部（或继承它的子合约里）不能直接调用
# 合约安全
## 重入攻击
转账后还没更新状态时，重新调用转账函数，因为状态没更新所以可以一直转账

先更新状态，再转账

## 权限控制缺失
所有管理函数应使用 onlyOwner 或 AccessControl 修饰符保护。
## 预言机操纵
比如攻击者用一大笔资金让某个dex池子里的某币价格瞬间暴涨，让预言机得到假价格

***解决方案***：

使用权威预言机（chainlink，pyth）

增加时序约束（价格必须在几分钟内更新）

同时查询多个oracle，然后加权平均
## 整数溢出
使用 unchecked {} 时需确保逻辑安全。（unchecked 用来省gas）

推荐使用Solidity 0.8+ 的内建溢出检查或 SafeMath
## 前置交易/三明治攻击
***前置交易***：你要大额买入一个低流动性代币 → 价格会上涨。攻击者看到你的 tx，先用更高 gas 费抢先买入 → 把价格推高。你的 tx 执行时，你被迫用更高价买入。（***三明治***）攻击者随后卖出，赚差价。

攻击者知道价格一定会升高所以买入

tx是交易的意思