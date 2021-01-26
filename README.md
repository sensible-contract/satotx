
## 签名器 API接口

### SatoTx

我们部署了一个签名器 [https://api.satotx.com](https://api.satotx.com) ，可暂做测试，获取utxo的签名。支持的API如下：

### 一、 对“某UTXO”签名


URL中需要的参数为：

1. txid: 产生UTXO的txid
2. index: UTXO的output index

Body中需要的json参数为：

1. txHex: 产生UTXO的rawtx内容


#### Request
- Method: **POST**
- URL:  ```/utxo/{txid}/{index}```
    - 签名:  ```/utxo/4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b/0```
- Headers：`Content-Type: application/json`
- Body:
```
{
  "txHex": "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff4d04ffff001d0104455468652054696d65732030332f4a616e2f32303039204368616e63656c6c6f72206f6e206272696e6b206f66207365636f6e64206261696c6f757420666f722062616e6b73ffffffff0100f2052a01000000434104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac00000000"
}
```

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- pubKey: rabin签名的公钥，hex编码
- txId: 同输入参数
- index: 同输入参数
- byTxId: 空
- sigBE: 签名，大端字节序，hex编码
- padding: 签名的padding，hex编码
- payload: 签名的内容，hex编码

其中payload字节内容为：

    txid, index, value, script

txid在payload中为原始字节序, index是小端4字节，value是小端8字节。

- Body
```
{
  "code": 0,
  "msg": "ok",
  "data": {
    "pubKey": "25108ec89eb96b99314619eb5b124f11f00307a833cda48f5ab1865a04d4cfa567095ea4dd47cdf5c7568cd8efa77805197a67943fe965b0a558216011c374aa06a7527b20b0ce9471e399fa752e8c8b72a12527768a9fc7092f1a7057c1a1514b59df4d154df0d5994ff3b386a04d819474efbd99fb10681db58b1bd857f6d5",
    "txId": "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
    "index": 0,
    "byTxId": "",
    "sigBE": "196a308c5393bd4b2aed57b4aa8891f1628b98bdf05776b8f79519c09f7515a564566421243aafbb6993c2088e87d73f1debf0acf58362d9f9b63289e548c8827a6a4fd1eb20d5bed80ff2e28dac42cd18865806d3c5bdd31ea515c7c16de89f43f7fe388fce885c0a1da9a3156d116c86afe08f2a4a7a74fcc66b280a8913a4",
    "padding": "0100",
    "payload": "3ba3edfd7a7b12b27ac72c3e67768f617fc81bc3888a51323a9fb8aa4b1e5e4a0000000000f2052a010000004104678afdb0fe5548271967f1a67130b7105cd6a828e03909a67962e0ea1f61deb649f6bc3f4cef38c4f35504e51ec112de5c384df7ba0b8d578a4c702b6bf11d5fac"
  }
}
```


### 二、对“某UTXO被下一个Tx花费”签名

URL中需要的参数为：

- txid: 产生UTXO的txid
- index: UTXO的output index
- byTxid: 花费UTXO的txid

Body中需要的json参数为：
- txHex: 产生UTXO的rawtx内容
- byTxHex: 花费UTXO的rawtx内容

#### Request
- Method: **POST**
- URL:  ```/utxo-spend-by/{txid}/{index}/{byTxid}```
    - 签名:  ```/utxo-spend-by/0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9/0/f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16```
- Headers：`Content-Type: application/json`
- Body:
```
{
  "txHex": "01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff0704ffff001d0134ffffffff0100f2052a0100000043410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac00000000",
  "byTxHex": "0100000001c997a5e56e104102fa209c6a852dd90660a20b2d9c352423edce25857fcd3704000000004847304402204e45e16932b8af514961a1d3a1a25fdf3f4f7732e9d624c6c61548ab5fb8cd410220181522ec8eca07de4860a4acdd12909d831cc56cbbac4622082221a8768d1d0901ffffffff0200ca9a3b00000000434104ae1a62fe09c5f51b13905f07f06b99a2f7159b2225f374cd378d71302fa28414e7aab37397f554a7df5f142c21c1b7303b8a0626f1baded5c72a704f7e6cd84cac00286bee0000000043410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac00000000"
}
```

#### Response

返回值 code == 0 为正确，其他都是错误

data包括字段为：

- pubKey: rabin签名的公钥，hex编码
- txId: 同输入参数
- index: 同输入参数
- byTxId: 同输入参数
- sigBE: 签名，大端字节序，hex编码
- padding: 签名的padding，hex编码
- payload: 签名的内容，hex编码

其中payload字节内容为：

    txid, index, value, script, bytxid

txid在payload中为原始字节序, index是小端4字节，value是小端8字节。

- Body
```
{
  "code": 0,
  "msg": "ok",
  "data": {
    "pubKey": "25108ec89eb96b99314619eb5b124f11f00307a833cda48f5ab1865a04d4cfa567095ea4dd47cdf5c7568cd8efa77805197a67943fe965b0a558216011c374aa06a7527b20b0ce9471e399fa752e8c8b72a12527768a9fc7092f1a7057c1a1514b59df4d154df0d5994ff3b386a04d819474efbd99fb10681db58b1bd857f6d5",
    "txId": "0437cd7f8525ceed2324359c2d0ba26006d92d856a9c20fa0241106ee5a597c9",
    "index": 0,
    "byTxId": "f4184fc596403b9d638783cf57adfe4c75c605f6356fbc91338530e9831e9e16",
    "sigBE": "215ea6362e87203e8eec7ce28724185d4f5436245e9796b2b843368d8d9f0f98f75e5f402e87469584ac3a24b7fdb1eae937db9bfbd96692fd721479647506fee1c69141dd8c6793d27898416cb5c5b23b658780fc55cf56f75c3e196849034397b9de986301b4dc10880bc5cae54576dec45a92a95ffd5dd95c333325449d09",
    "padding": "0100",
    "payload": "c997a5e56e104102fa209c6a852dd90660a20b2d9c352423edce25857fcd37040000000000f2052a01000000410411db93e1dcdb8a016b49840f8c53bc1eb68a382e97b1482ecad7b148a6909a5cb2e0eaddfb84ccf9744464f82e160bfa9b8b64f9d4c03f999b8643f656b412a3ac169e1e83e930853391bc6f35f605c6754cfead57cf8387639d3b4096c54f18f4"
  }
}
```
