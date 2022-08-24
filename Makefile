genpb:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./pb/*.proto

auth_service:
	go run services/auth/*.go

orders_service:
	go run services/orders/*.go

suppliers_service:
	go run services/suppliers/*.go

customers_service:
	go run services/customers/*.go

.PHONY: genpb auth_service orders_service suppliers_service customers_service
