package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	plugin "k8s-kms-plugin/plugin"
)

var (
	cfgFile = flag.String("f", "/etc/kubernetes/jdcloud-kms-plugin.json", "kms plugin configuration for kubernetes")
)

type PluginConfig struct {
	AccessKey      string `json:"AccessKey"`
	SecretKey      string `json:"SecretKey"`
	KmsEndpoint    string `json:"KmsEndpoint"`
	KmsKeyId       string `json:"KmsKeyId"`
	KmsSchema      string `json:"KmsSchema"`
	GRPCSocketPath string `json:"GRPCSocketPath"`
}

func main() {
	flag.Parse()

	/** 加载配置文件 **/
	cfgData, err := ioutil.ReadFile(*cfgFile)
	if err != nil {
		panic(err)
	}

	/** 解析配置文件 **/
	var cfg PluginConfig
	if err := json.Unmarshal(cfgData, &cfg); err != nil {
		panic(err)
	}

	fmt.Println("JdCloud K8S KMS Plugin start. ")
	fmt.Println("EndPoint:", cfg.KmsEndpoint)
	fmt.Println("KmsSchema:", cfg.KmsSchema)
	fmt.Println("KmsKeyId:", cfg.KmsKeyId)
	fmt.Println("GRPCSocketPath:", cfg.GRPCSocketPath)

	kmsClient := plugin.NewKmsClient(cfg.AccessKey, cfg.SecretKey, cfg.KmsEndpoint, cfg.KmsKeyId, cfg.KmsSchema)
	kmsPlugin := plugin.NewK8sKmsPlugin(cfg.GRPCSocketPath)
	kmsPlugin.SetKmsClient(kmsClient)

	gRPCSrv, kmsErrorChain := kmsPlugin.ServeKMSRequests()
	defer gRPCSrv.GracefulStop()

	for {
		select {
		case kmsError := <-kmsErrorChain:
			panic(kmsError)
		}
	}
}
