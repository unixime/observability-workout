# AGENTS.md

## Project overview
Simple Golang web application to manage the tenats on a block of flat. The application uses PostgreSQL as 
databaase backend

## Modules
- project uses [Gin Web Framwork](https://gin-gonic.com/en/) to implement API definition
- project oses [Gorm](https://gorm.io/index.html) to connect to PostgreSQL database

OpenAPI 3.0.3 specification for the tenant management API. The whole project is [spec.yml](./spec.yml); there is no application code yet.

## Repository structure
- `api/spec.yml`: API definition (paths, schemas).
- `AGENTS.md`: instructions for AI agents.

## API summary
- `GET /tenants`: list all tenants.
- `POST /tenants`: register a tenant (`201`, `400`, `409` if it already exists).
- `GET /tenants/search`: find tenants by `id`, `name` or `lastname` query parameters.

## Conventions
- Define reusable models under `components/schemas` and reference them with `$ref`.
- Schema names are lowercase (e.g. `tenant`).
- Every operation has a `summary`, a camelCase `operationId` and documented responses.
- `tenant.id` is `readOnly` (server-assigned); `name` and `lastname` are required.
- Use 2-space indentation and single-quoted response codes (`'200'`).

## Validation
Validate the spec after every change:
```sh
npx @redocly/cli lint spec.yml
```

## Rules for agents
- Make minimal, targeted edits; don't reformat unrelated parts of the spec.
- Keep this file and the API summary above in sync when paths change.
- Do not commit secrets or real tenant data; use fake values in examples.
