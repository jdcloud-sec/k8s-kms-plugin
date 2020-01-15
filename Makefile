export GOPROXY=https://goproxy.io,direct
export GO111MODULE=on

OBJ = k8s-kms-plugin

all: $(OBJ)

$(OBJ):
	mkdir -p ./build
	go build -gcflags "-N -l" -tags netgo -o ./build/$(OBJ)

clean:
	rm -fr ./build/$(OBJ)

