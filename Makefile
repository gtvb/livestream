GO_PATH=$(HOME)/go
SWAGGER_CMD=$(GO_PATH)/bin/swagger

swagger_check:
	which $(SWAGGER_CMD) || go get github.com/go-swagger/go-swagger/cmd/swagger@latest

swagger: swagger_check
	$(SWAGGER_CMD) generate spec -o ./swagger.yaml --scan-models

swagger_serve: swagger_check
	$(SWAGGER_CMD) serve -F=swagger swagger.yaml