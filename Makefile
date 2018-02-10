APP := aws-sign-proxy

all:
	CGO_ENABLED=0 $(MAKE) clean build compress checksums

include github.com/codekoala/make/golang
include github.com/codekoala/make/upx
include github.com/codekoala/make/docker
