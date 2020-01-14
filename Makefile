export GOPATH=$(shell pwd):$(shell pwd)/vendor

OBJ = k8s-kms-plugin

all: $(OBJ)

$(OBJ):
	mkdir -p build
	cd src && go build -gcflags "-N -l" -o ../build/$@

clean:
	rm -fr $(OBJ)

