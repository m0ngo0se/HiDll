# **HiDll**

dll劫持自动化搜索工具，暂时只支持preload的搜索，可用于寻找白加黑。

## 使用说明：

1.首先直接执行工具，进入preload模式

![image-20221121171023423](https://github.com/m0ngo0se/HiDll/blob/master/img/image-20221121171023423.png)



2.创建工程，create+工程名(可随意起)+想要挖掘的软件路径。这里以某office为例

![image-20221121171351869](https://github.com/m0ngo0se/HiDll/blob/master/img/image-20221121171351869.png)

工具会自动将符合劫持条件的exe存储到你的工程目录并输出控制台，你可直接查看，如果显示全部删除，则该软件没有方便劫持的dll。

![image-20221121171936824](https://github.com/m0ngo0se/HiDll/blob/master/img/image-20221121171936824.png)

3.选择想要进行劫持的exe进行输出。输入get+对应的id即可

![image-20221121172033939](https://github.com/m0ngo0se/HiDll/blob/master/img/image-20221121172033939.png)

后续就是利用劫持dll了，网上很多方法

## 过滤条件:

一切都在源码中，根据自身需要更改吧

1. 64位(个人喜好)
2. 签名(必须的，不然白加黑无意义)
3. 非.net
4. 非系统自带dll(大部分杀软都会根据文件名查杀系统已知dll的劫持，所以这种无意义)
5. 无其他依赖(有的可劫持的exe没有静态编译，这种利用起来很麻烦，直接通过前缀去除)
