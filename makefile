
run:
	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env up --build  test-runner
# 	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env run -T --rm test-runner go test -v git.iu7.bmstu.ru/ped22u691/PPO.git/internal/tests/integration 

up:
	docker compose -f ./deployment/docker-compose.ci.yml --env-file ./configs/test_db.env up --build -d postgres migrator

clear:
	docker compose -f ./deployment/docker-compose.ci.yml down -v

