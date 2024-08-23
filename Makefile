swagger:
	swagger generate spec -o ./swagger.yaml --scan-models

swagger_serve:
	swagger serve -F=swagger --port 4004 --no-open  swagger.yaml
