.PHONY: proto

proto:
	protoc --go_out=. --go_opt=paths=source_relative ./proto/*.proto