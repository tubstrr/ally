PROJECT_NAME=ally
APP_NAME=ally

DOCKER_LOCATION=example/docker
COMPOSE_OPTIONS=-p ${PROJECT_NAME}
EXPORT=export PROJECT_NAME=${PROJECT_NAME} && export APP_NAME=${APP_NAME} &&

.PHONY: up
up:
	${EXPORT} docker-compose -f ${DOCKER_LOCATION}/development.yaml ${COMPOSE_OPTIONS} up

.PHONY: postgresql
postgresql:
	${EXPORT} docker-compose -f ${DOCKER_LOCATION}/development.yaml build --no-cache postgresql
	
.PHONY: ally
ally:
	${EXPORT} docker-compose -f ${DOCKER_LOCATION}/development.yaml build --no-cache ally