.PHONY: swagger

swagger:
	docker exec -it api sh -c "swag init -d cmd --parseDependency --parseInternal --parseDepth 4 -o swagger/docs"
