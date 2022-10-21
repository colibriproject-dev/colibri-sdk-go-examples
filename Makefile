STACK_NAME=colibri-dev

start:
	docker-compose -p ${STACK_NAME} up -d

stop:
	docker-compose -p ${STACK_NAME} stop

clean:
	docker-compose -p ${STACK_NAME} down -v
	docker rmi school-module finantial-module

build:
	docker buildx build --progress=plain --no-cache -t school-module . --build-arg SRC_PATH=/school-module
	docker buildx build --progress=plain --no-cache -t finantial-module . --build-arg SRC_PATH=/finantial-module

logs:
	docker-compose -p ${STACK_NAME} logs -f

stats:
	docker stats