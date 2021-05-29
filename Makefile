.DEFAULT_GOAL := build
.PHONY: build run test test_backend test_frontend

COMPOSE = docker-compose -f docker/docker-compose.yml --project-directory .

build:
	${COMPOSE} build

up:
	${COMPOSE} up

test: test_backend test_frontend

test_backend:
	docker build -f docker/Dockerfile_test -t tp_int:test_backend .
	docker run tp_int:test_backend

test_frontend:
	docker build -f docker/Dockerfile_webtest -t tp_int_web:test_frontend .
	docker run tp_int_web:test_frontend
