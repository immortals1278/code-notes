新建一个文件夹，然后对着空文件夹forge init（vscode就有终端）
## src
里面写源码，开头定义开源协议和编译器版本

实现一个合约前先创建该合约的接口：抽象函数，事件，报错

然后让合约实现该接口。实现时安装（用终端）并导入相关库（接口不需要）wizard.openzeppelin.com

安装完相关库要在toml里配置
```solidity
constructor{
    string memory name_,
    string memory symbol_,
    address[] memory tokens_,
    uint256[] memory initTokenAmountPerShare_,
    uint256 minMintAmount_,
    address swapRouter_
} ERC20(name_, symbol_) Ownable(msg.sender) {
    // 构造函数体
}
```
Solidity 中继承合约构造函数的初始化方式,如`ERC20(name_, symbol_)`调用 ERC20 父合约的构造函数，传入 name_ 和 symbol_ 参数来初始化代币名称和符号

## 部署在本地网
### 如何在wsl里使用vpn
在wsl里：ipconfig，找到两个本地配置下的第一段里的ipV4，替换到
```
export http_proxy=http://10.253.62.130:7890
export https_proxy=http://10.253.62.130:7890
export ALL_PROXY=socks5://10.253.62.130:7891
```
然后直接执行：（是不是直接执行这个就行？？？？？）
```
export http_proxy=http://10.253.62.130:7890
export https_proxy=http://10.253.62.130:7890
export ALL_PROXY=socks5://10.253.62.130:7891
```
然后查询当前wsl的ip是不是vpn的ip：(这些命令任意都可以)
```
curl ifconfig.me
curl ip.sb
curl ipapi.co/json
curl https://api.ipify.org
curl https://ifconfig.co/json
```
然后就可以下载很多东西了（下载指令看官方文档不要看ai，ai老乱说）