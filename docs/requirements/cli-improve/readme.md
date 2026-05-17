之前部署cli主要通mock的形式实现的，这次我们完善cli，采用真实的生产环境。

发布游戏需要登录，获取到accessToken,使用如下接口：

curl --request POST \
  --url https://api.3sdk.yu3.co/common/v1/auth/login \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Length: 83' \
  --header 'Content-Type: application/json' \
  --header 'Host: api.3sdk.yu3.co' \
  --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
  --data '{
	"identifier": "15869155220@163.com",
	"type": "password",
	"data": "Dazz-0001"
}'

然后需要上传文件，如sql文件加，游戏程序包等，使用如下接口获取凭证，然后采用阿里云oss客户端直传：

curl --request GET \
  --url https://api.3sdk.yu3.co/developer/v1/file/policy-token \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM0OTg3NDIsImp0aSI6IjYifQ.uxu5SQ4c7EQu_ZmXlD4PIo6DfPQbYZnBuu4D3IOg_lk' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Host: api.3sdk.yu3.co' \
  --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0'

可以创建游戏，通过以下接口：

curl --request POST \
  --url https://api.3sdk.yu3.co/developer/v1/game/create-with-version \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM0OTg3NDIsImp0aSI6IjYifQ.uxu5SQ4c7EQu_ZmXlD4PIo6DfPQbYZnBuu4D3IOg_lk' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Length: 890' \
  --header 'Content-Type: application/json' \
  --header 'Host: api.3sdk.yu3.co' \
  --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
  --data '{
	"name": "测试游戏审核4",
	"description": "测试游戏审4",
	"logo": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261107SZUNOT.jpeg",

	"version": {
		"version": "1.0.0",
		"changeLog": "发布新版本",
		"fileUrl": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261109kQriFw.zip",
		"initSqlUrl": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261109r2HqVX.sql",
		"screenConfig":{
			"screenType":1,
			"halfSupport":1,
			"halfRatio":"0.75"
		},
		"buildConfig": {
			"backend": {
				"workDir": "lucky77pro_1.0.7_20250625/admin",
				"cmd": "./server_lucky77pro -type admin"
			},
			"frontend": {
				"workDir": "lucky77pro_1.0.7_20250625/h5/20250624143413",
				"cmd": ""
			},
			"socket": {
				"workDir": "lucky77pro_1.0.7_20250625/logic",
				"cmd": "./server_lucky77pro -type logic"
			}
		}
	}
}'




更新游戏基本信息，使用以下接口：

curl --request PUT \
  --url https://api.3sdk.yu3.co/developer/v1/game/2506301759369iqh \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM0OTg3NDIsImp0aSI6IjYifQ.uxu5SQ4c7EQu_ZmXlD4PIo6DfPQbYZnBuu4D3IOg_lk' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Length: 963' \
  --header 'Content-Type: application/json' \
  --header 'Host: api.3sdk.yu3.co' \
  --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
  --data '{
	"name": "测试游戏审核3",
	"description": "测试游戏审核",
	"logo": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261107SZUNOT.jpeg",
	"resourceItems": [
		{
			"img": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261107zqceiw.jpeg",
			"description": "1234"
		}
	],
	"version": {
		"version": "1.0.0",
		"changeLog": "发布新版本",
		"fileUrl": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261109kQriFw.zip",
		"initSqlUrl": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261109r2HqVX.sql",
		"buildConfig": {
			"backend": {
				"workDir": "lucky77pro_1.0.7_20250625/admin",
				"cmd": "./server_lucky77pro -type admin"
			},
			"frontend": {
				"workDir": "lucky77pro_1.0.7_20250625/h5/20250624143413",
				"cmd": ""
			},
			"socket": {
				"workDir": "lucky77pro_1.0.7_20250625/logic",
				"cmd": "./server_lucky77pro -type logic"
			}
		}
	}
}'

更新游戏版本，使用如下接口：

curl --request POST \
  --url https://api.3sdk.yu3.co/developer/v1/game/update-with-version \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTM0OTg3NDIsImp0aSI6IjYifQ.uxu5SQ4c7EQu_ZmXlD4PIo6DfPQbYZnBuu4D3IOg_lk' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Length: 1019' \
  --header 'Content-Type: application/json' \
  --header 'Host: api.3sdk.yu3.co' \
  --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
  --data '{
	"uri":"25062611270HPptQ",
	"name": "测试游戏审核1",
	"description": "测试游戏审核",
	"logo": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261107SZUNOT.jpeg",
	"resourceItems": [
		{
			"uri":"2506261430FCAW6D",
			"img": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261107zqceiw.jpeg",
			"description": "1234"
		}
	],
	"version": {
		"version": "1.0.2",
		"changeLog": "发布新版本",
		"fileUrl": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261109kQriFw.zip",
		"initSqlUrl": "https://paas-3os.oss-cn-shanghai.aliyuncs.com/uploads/2025/06/26/2506261109r2HqVX.sql",
		"buildConfig": {
			"backend": {
				"workDir": "lucky77pro_1.0.7_20250625/admin",
				"cmd": "./server_lucky77pro -type admin"
			},
			"frontend": {
				"workDir": "lucky77pro_1.0.7_20250625/h5/20250624143413",
				"cmd": ""
			},
			"socket": {
				"workDir": "lucky77pro_1.0.7_20250625/logic",
				"cmd": "./server_lucky77pro -type logic"
			}
		}
	}
}'


发起审核，使用下面的接口

curl --request POST \
  --url https://api.3sdk.yu3.co/developer/v1/game/apply-review \
  --header 'Accept: */*' \
  --header 'Accept-Encoding: gzip, deflate, br' \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NTU3NjI5NjksImp0aSI6IjYifQ.y9hSr-Gg7AX3UbzORWFDqM32iGDoitLK-ZbNEnbUEhs' \
  --header 'Cache-Control: no-cache' \
  --header 'Connection: keep-alive' \
  --header 'Content-Length: 32' \
  --header 'Content-Type: application/json' \
  --header 'Host: api.3sdk.yu3.co' \
  --header 'User-Agent: PostmanRuntime-ApipostRuntime/1.1.0' \
  --data '{
    "uri":"2506261515XA0Jsc"
}'

可以参考项目 /Users/liuxuanyao/work/paas-backend，查看以上接口的具体逻辑