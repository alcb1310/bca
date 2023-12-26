# Budget Control Application

## Tech Stack

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](http://go.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)](https://jwt.io)

## Socials

[![Twitter Badge](https://img.shields.io/twitter/follow/username.svg?style=social&label=Follow)](https://twitter.com/alcb1310)
[![Github Badg](https://img.shields.io/badge/GitHub-100000?style=for-the-badge&logo=github&logoColor=white)](https://github.com/alcb1310)
[![built with Codeium](https://codeium.com/badges/main)](https://codeium.com/badges/main)

## Description

The objective of this application is to manage the budget of a construction company, to achieve this it will be able to:

- Manage Suppliers
- Manage Budget Items
- Create Budgets for each project
- Create Invoices
- Each invoice will decrease the budgeted data
- Several reports to be defined

## Table of contents

- [Routes](#routes)
    - [Login](#login)
    - [Users](#users)
    - [Settings](#settings)
        - [Projects](#projects)
        - [Suppliers](#suppliers)
        - [Budget Items](#budget-items)

- [Deployment](#deployment)

## Routes

In order to achieve the project's description, the application will have both public and protected Endpoints

### Login

The application will start at the `/login` route which allows a user to login to it by providing their email and password, then the
server will validate their credentials, and on success it will go to the protected routes and if it didn't succeed, it will display
a message indicating `invalid credentials`

### Users

The logged in user will be able to:

- Change his/her password by making a `PATCH` request to the `/api/v1/users` route

The admin users will be able to:

- Create new users by making a `POST` request to the `/api/v1/users` route
- Get all users by making a `GET` request to the `/api/v1/users` route
- Get one user by making a `GET` request to the `/api/v1/users/:id` route
- Update one user by making a `PUT` request to the `/api/v1/users/:id` route
- Delete one user by making a `DELETE` request to the `/api/v1/users/:id` route

### Settings

The settings for this application will be divided into the following:

#### Projects

The projects will be able to:

- Create new projects by making a `POST` request to the `/api/v1/projects` route
- Get all projects by making a `GET` request to the `/api/v1/projects` route
- Get one project by making a `GET` request to the `/api/v1/projects/:id` route
- Update one project by making a `PUT` request to the `/api/v1/projects/:id` route

#### Suppliers

The suppliers will be able to:

- Create new suppliers by making a `POST` request to the `/api/v1/suppliers` route
- Get all suppliers by making a `GET` request to the `/api/v1/suppliers` route
- Get one supplier by making a `GET` request to the `/api/v1/suppliers/:id` route
- Update one supplier by making a `PUT` request to the `/api/v1/suppliers/:id` route

#### Budget Items

The budget items will be able to:

- Create new budget items by making a `POST` request to the `/api/v1/budget-items` route
- Get all budget items by making a `GET` request to the `/api/v1/budget-items` route
- Get one budget item by making a `GET` request to the `/api/v1/budget-items/:id` route
- Update one budget item by making a `PUT` request to the `/api/v1/budget-items/:id` route

## Deployment


In order to be able to deploy this application, the following is needed:

1. Clone the repository using the following command

```bash
git clone https://github.com/alcb1310/bca-go-w-test.git
```

2. Download all of the project's dependencies by running 

```bash
go mod tidy
```

3. At the root of the project directory, create a `.env` file with the following fields:

```.env
PORT=<Port number where the application will listen>

PGDATABASE=<Name of the postgres database>
PGHOST=<Host where the postgres database server is running>
PGPASSWORD=<Password that the postgres server uses>
PGPORT=<Port where the postgres server is listening>
PGUSER=<Username of the postgres server>

SECRET=<Secret used to generate the JWT Token>
```

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
