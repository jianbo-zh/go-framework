package config

import (
	"goframework/internal/env"

	"github.com/jianbo-zh/go-config"
)

func appConfig() {
	config.Add("app", config.StrMap{
		"env": config.Env("APP_ENV", env.Development),
		"key": config.Env("APP_KEY", ""),
	})

	config.Add("storage", config.StrMap{
		"dir": "",
	})

	config.Add("ipfs", config.StrMap{
		"gateway": config.Env("IPFS_GATEWAY", "https://ipfs.io/ipfs/"),
		"cluster": config.Env("IPFS_CLUSTER", ""),
	})

	config.Add("sms", config.StrMap{
		"akid":     config.Env("SMS_ACCESS_KEY_ID", ""),     // 短信 access_key_id
		"aksecret": config.Env("SMS_ACCESS_KEY_SECRET", ""), // 短信 access_key_secret
		"signname": config.Env("SMS_SIGN_NAME", ""),         // 短信签名
	})

	config.Add("blockchain", config.StrMap{
		"contract":   config.Env("SRC_CONTRACT_ADDRESS", ""),
		"wallet_key": config.Env("SRC_WALLET_KEY", ""),
		"wallet_pw":  config.Env("SRC_WALLET_PW", ""),
	})

}
