STACK_NAME=colibri-dev

start:
	docker compose -p ${STACK_NAME} up -d --remove-orphans

stop:
	docker compose -p ${STACK_NAME} stop

clean:
	docker compose -p ${STACK_NAME} down -v
	docker rmi school-module finantial-module

build:
	docker build --no-cache -t school-module . --build-arg APP_SRC=school-module
	docker build --no-cache -t finantial-module . --build-arg APP_SRC=finantial-module

logs:
	docker compose -p ${STACK_NAME} logs -f

stats:
	docker stats