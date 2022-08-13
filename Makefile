default: build-all

.PHONY: build-all
build-all: linux-intel linux-arm mac-intel mac-arm windows-intel zip

.PHONY: mac-intel
mac-intel:
	env GOOS=darwin GOARCH=amd64 go build -o build/checkout-mac-intel
	chmod u+x build/checkout-mac-intel

.PHONY: mac-arm
mac-arm:
	env GOOS=darwin GOARCH=arm64 go build -o build/checkout-mac-applesilicon
	chmod u+x build/checkout-mac-applesilicon

.PHONY: linux-intel
linux-intel:
	env GOOS=linux GOARCH=amd64 go build -o build/checkout-linux-intel
	chmod u+x build/checkout-linux-intel

.PHONY: linux-arm
linux-arm:
	env GOOS=linux GOARCH=arm64 go build -o build/checkout-linux-arm
	chmod u+x build/checkout-linux-arm

.PHONY: windows-intel
windows-intel:
	env GOOS=windows GOARCH=amd64 go build -o build/checkout.exe

.PHONY: zip
zip:
	tar -czvf build/all-checkout-binaries.tar.gz build/checkout*
