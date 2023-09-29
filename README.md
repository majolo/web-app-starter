# Web-app-starter

An opinionated web application starter kit.

## Backend

* API server built using Golang
* APIs defined using Buf and protobuf, exposed as REST using grpc-gateway
* Database access using GORM

## Frontend

* React app built using Typescript, Next.js, and TailwindCSS
* API client generated using OpenAPI Generator, implemented as hooks using ReactQuery
* Authentication using Supabase

## Notes

Designed to be very minimal, doesn't include necessary things such as:

* Dockerfiles + deployment scripts
* Testing
* CI/CD
