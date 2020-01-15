package plugin

import (
	"encoding/base64"
	core "github.com/jdcloud-api/jdcloud-sdk-go/core"
	kms "github.com/jdcloud-api/jdcloud-sdk-go/services/kms/apis"
	client "github.com/jdcloud-api/jdcloud-sdk-go/services/kms/client"
)

type KmsClient struct {
	keyID  string
	client *client.KmsClient
}

func NewKmsClient(accessKey, secretKey, kmsEndpoint, kmsKeyId, schema string) *KmsClient {
	/** 设置Credentials对象 **/
	credentials := core.NewCredentials(accessKey, secretKey)

	/** 设置Config对象 **/
	config := core.NewConfig()
	config.SetEndpoint(kmsEndpoint)
	config.SetScheme(schema)

	client := client.NewKmsClient(credentials)
	client.SetConfig(config)
    client.SetLogger(core.NewDummyLogger())

	return &KmsClient{
		keyID:  kmsKeyId,
		client: client,
	}
}

func (c *KmsClient) Encrypt(plain []byte) ([]byte, error) {
	plainData := base64.StdEncoding.EncodeToString(plain)
	reqEnc := kms.NewEncryptRequest(c.keyID)
	reqEnc.SetPlaintext(plainData)

	resp, err := c.client.Encrypt(reqEnc)
	if err != nil {
		return nil, err
	}

	cipher, err := base64.StdEncoding.DecodeString(resp.Result.CiphertextBlob)
	if err != nil {
		return nil, err
	}
	return cipher, nil
}

func (c *KmsClient) Decrypt(cipher []byte) ([]byte, error) {
	cipherData := base64.StdEncoding.EncodeToString(cipher)
	reqDec := kms.NewDecryptRequest(c.keyID)
	reqDec.SetCiphertextBlob(cipherData)

	resp, err := c.client.Decrypt(reqDec)
	if err != nil {
		return nil, err
	}

	plain, err := base64.StdEncoding.DecodeString(resp.Result.Plaintext)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
