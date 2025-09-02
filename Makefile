db-login:
	docker compose exec db psql -U enomo -d enomo_dev
db-users:
	docker compose exec db psql -U enomo -d enomo_dev -c "SELECT * FROM users"
db-groups:
	docker compose exec db psql -U enomo -d enomo_dev -c "SELECT * FROM groups"
db-group_members:
	docker compose exec db psql -U enomo -d enomo_dev -c "SELECT * FROM group_members"
db-sessions:
	docker compose exec db psql -U enomo -d enomo_dev -c "SELECT * FROM sessions"