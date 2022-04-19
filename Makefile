APP_NAME=heyui
VERSION=1.0
KEYFOLDER:=key-pair
CURFOLDER:=$(shell pwd)

.PHONY: keygen
default: build

keygen:
	mkdir $(KEYFOLDER)
	openssl genrsa -out $(KEYFOLDER)/id_rsa 4096
	openssl rsa -in $(KEYFOLDER)/id_rsa -pubout -out $(KEYFOLDER)/id_rsa.pub

