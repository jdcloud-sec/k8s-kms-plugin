package plugin

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
)

const (
	netProtocol    = "tcp"
	apiVersion     = "v1beta1"
	runtimeName    = "JDCloud-KMS"
	runtimeVersion = "0.0.1"
)

type K8sKmsPlugin struct {
	pathToUnixSocket string

	kmsClient *KmsClient
	net.Listener
	*grpc.Server
}

func NewK8sKmsPlugin(pathToUnixSocket string) *K8sKmsPlugin {
	return &K8sKmsPlugin{
		pathToUnixSocket: pathToUnixSocket,
	}
}

func (kms *K8sKmsPlugin) SetKmsClient(kmsClient *KmsClient) {
	kms.kmsClient = kmsClient
}

func (kms *K8sKmsPlugin) Version(ctx context.Context, req *VersionRequest) (*VersionResponse, error) {
	return &VersionResponse{Version: apiVersion, RuntimeName: runtimeName, RuntimeVersion: runtimeVersion}, nil
}

func (kms *K8sKmsPlugin) Decrypt(ctx context.Context, req *DecryptRequest) (*DecryptResponse, error) {
	plain, err := kms.kmsClient.Decrypt(req.Cipher)
	if err != nil {
		return nil, err
	}
	return &DecryptResponse{Plain: plain}, nil
}

func (kms *K8sKmsPlugin) Encrypt(ctx context.Context, req *EncryptRequest) (*EncryptResponse, error) {
	cipher, err := kms.kmsClient.Encrypt(req.Plain)
	if err != nil {
		return nil, err
	}
	return &EncryptResponse{Cipher: []byte(cipher)}, nil
}

func (kms *K8sKmsPlugin) ServeKMSRequests() (*grpc.Server, chan error) {
	errChain := make(chan error, 1)

	listener, err := net.Listen(netProtocol, kms.pathToUnixSocket)
	if err != nil {
		errChain <- fmt.Errorf("failed to start listener, error: %v", err)
		close(errChain)
		return nil, errChain
	}
	kms.Listener = listener
	kms.Server = grpc.NewServer()

	RegisterKeyManagementServiceServer(kms.Server, kms)

	go func() {
		defer close(errChain)
		errChain <- kms.Serve(kms.Listener)
	}()

	return kms.Server, errChain
}
