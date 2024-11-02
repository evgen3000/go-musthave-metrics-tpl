make create-db:
	docker run --name admin -p "5432:5432" -e POSTGRES_DB=admin -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -d "postgres:14-bullseye"
	until docker exec admin pg_isready -U admin
	do
		echo "Waiting for postgres container..."
		sleep 1
	done