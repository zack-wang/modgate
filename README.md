# Modbus Gateway (modgate): read/write modbus device by url, 使用網址來讀寫modbus設備 
## url format for reading,
```
http://0.0.0.0:8080/<protocol>/<ip and port>/<function code>/<address start>/<length>/<format>
```

1. protocol: only support modbustcp
2. ip and port: modbus tcp device host and port, for example - 10.0.0.1:502
3. function code: modbus function code, for example - 1, 2, 3, 4
4. address start: start of modbus address
5. length: number of readouts
6. format: convert bytes to (integer) I64BE,I64LE,I32BE,I32LE,I16BE,I16LE, (float) F64BE,F64LE,F32BE,F32LE, (string) S64BE,S64LE

## 讀取modbus的網址格式
1. protocol: 目前只支持 modbustcp
2. ip and port: modbus設備的地址與埠號, 例如 - 10.0.0.1:502
3. function code: modbus的功能代碼, 例如 - 1, 2, 3, 4
4. address start: modbus地址起始位置
5. length: 讀出的資料個數
6. format: 將位元組資料轉換成可以讀取的數值, (整數類)I64BE,I64LE,I32BE,I32LE,I16BE,I16LE, (浮點數) F64BE,F64LE,F32BE,F32LE, (字串) S64BE,S64LE

## output format JSON
```
{"code":0,"data":[666,123],"format":"I16BE","message":""}
```

1. code: 0 succeed, others failed
2. data: array of result
3. format: same as input
4. message: if code != 0, error message is here 

## 讀取到的資料(是JSON格式)說明
1. code: 0成功, 非0失敗
2. message: 失敗時的錯誤訊息
3. format: 轉換的格式(與輸入相同)
4. data: 資料陣列

## url format for writing (http POST)
```
http://0.0.0.0:8080/<protocol>/<ip and port>/<function code>/<address start>/<format>

[1,2,3,4]
```
1. protocol: only support modbustcp
2. ip and port: modbus tcp device host and port, for example - 10.0.0.1:502
3. function code: modbus function code, for example - 15 or 16
4. address start: start of modbus address
5. format: force convert data to these formats, (integer) I64BE,I64LE,I32BE,I32LE,I16BE,I16LE, (float) F64BE,F64LE,F32BE,F32LE, (string) S64BE,S64LE
6. Refer to No.5. It may be panic, if you force convert data from string to number. I need some help to make this better. 

## 寫入modbus的網址格式, 使用http POST方式
1. protocol: 目前只支持 modbustcp 協定
2. ip and port: modbus設備的地址與埠號, 例如 - 10.0.0.1:502
3. function code: modbus的功能代碼, 例如 - 15或16
4. address start: modbus地址起始位置
5. format: 將資料轉換成特定格式寫進去, 例如(整數類)I64BE,I64LE,I32BE,I32LE,I16BE,I16LE, (浮點數) F64BE,F64LE,F32BE,F32LE, (字串) S64BE,S64LE, 目前沒有檢查假如資料為字串, 強制轉成數字類可能會出錯