# newsapi
This repository contains a news fetching service and its api.

## functionality
Core service fetches news into a database, from which news are read by the API. The API exposes given endpoints:
```
GET /auth/google/login
GET /auth/google/logout

GET /news
GET /news/<id>
GET /news/<id>/tags

GET /tags/<id>
GET /tags/<id>/news

GET /user/tags
POST,DELETE /user/tags/<id>
GET /user/news
```

## features
- OAuth2 login
- Cookie based user session management

## Tech Stack

![Golang](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Python](https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)

## development
create .env files. in [api/](api/):
```
POSTGRES_CONN_STR=
GOOGLE_OAUTH_CLIENT_ID=
GOOGLE_OAUTH_CLIENT_SECRET=
SESSION_SECRET=
```
in [core/](core/):
```
NEWSAPI_KEY=
POSTGRES_CONN_STR=
```
run the [initial migration](db/migrations/01_initial_schema.sql). optionally, [fill the database](db/fill_db.sql).

run api and core fetcher:
```sh
air
```
```sh
uv sync
source .venv/bin/activate
python src/main.py
```
