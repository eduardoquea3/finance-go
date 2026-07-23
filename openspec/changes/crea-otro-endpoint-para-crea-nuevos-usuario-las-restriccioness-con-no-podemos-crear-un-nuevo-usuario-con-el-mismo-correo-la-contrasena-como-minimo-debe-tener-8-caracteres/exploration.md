## Exploration: crea otro endpoint para crea nuevos usuario, las restriccioness con no podemos crear un nuevo usuario con el mismo correo, la contraseña como minimo debe tener 8 caracteres

### Current State
The API currently has a PostgreSQL-backed `users` table with a database-level unique constraint on `email` and a bcrypt `password_hash`. The user repository only supports `FindByEmail`; the user service and HTTP handler are empty. Authentication is wired through `internal/auth`, with public `POST /api/v1/auth/login` and `POST /api/v1/auth/refresh`, while the `/api/v1` group is protected by bearer-token middleware. Login normalizes email with trim/lowercase and the login request currently validates passwords with `min=5`.

The existing migration seeds `admin@example.com` with bcrypt password `admin` for testing. This is an established prior decision, so changing the login minimum to 8 would make that seed credential unusable; the new minimum should apply to user creation unless the product explicitly decides to replace the seed credential.

### Affected Areas
- `internal/user/repository.go` — add a create operation (and preferably preserve repository-level error translation).
- `internal/user/postgres_repository.go` — persist the entity and translate PostgreSQL unique violations into a domain conflict error.
- `internal/user/service.go` — normalize email, enforce creation rules, hash the password with bcrypt, and avoid returning the hash.
- `internal/user/http_handler.go` — bind and validate the create request, map validation/conflict/persistence errors to HTTP responses.
- `internal/user/errors.go` — add a duplicate-email/conflict error and creation-related errors as needed.
- `internal/http/router.go` — wire the user handler; the smallest safe default is a protected `POST /api/v1/users` route because no role model exists for admin-only authorization.
- `cmd/api/main.go` — construct the user service/handler and inject them into the router.
- `docs/swagger.{yaml,json}` and `docs/docs.go` — regenerate the OpenAPI contract for the endpoint and request/response schemas.
- `migrations/000001_create_users.up.sql` — already supplies the required unique database constraint; no migration is needed for this scope.
- `migrations/000002_seed_admin.up.sql`, `internal/auth/request.go` — prior admin seed/password-minimum decision must remain compatible.
- `internal/validation/*_test.go` and new user tests — extend the existing table-driven validation/unit testing pattern.

### Approaches
1. **Protected user resource endpoint** — add `POST /api/v1/users` behind the existing bearer middleware, with user-domain service/repository layers.
   - Pros: reuses existing route protection, cleanly separates user creation from authentication, preserves the database uniqueness guarantee, and avoids adding an authorization/role model.
   - Cons: any authenticated user can create users because roles/admin claims do not exist; this must be documented as out of scope or followed by a separate authorization change.
   - Effort: Medium

2. **Public auth registration endpoint** — add `POST /api/v1/auth/register` without authentication and implement creation in the auth package.
   - Pros: conventional self-registration API and no bearer token prerequisite.
   - Cons: mixes account lifecycle with authentication, increases abuse/security concerns, and conflicts with the existing admin-seed direction if creation is intended as an administrative action.
   - Effort: Medium

### Recommendation
Use approach 1: `POST /api/v1/users`, protected by the existing bearer middleware, implemented through the currently empty user service and handler. Accept an email and plaintext password, normalize the email exactly as login does, validate `required,email` and `required,min=8`, hash with bcrypt, insert through the repository, and return a non-sensitive user representation (for example, ID, email, and timestamps) with `201 Created`. Treat the database unique constraint as authoritative so concurrent requests cannot create duplicate emails; map the resulting conflict to `409 Conflict`. Keep the seeded `admin@example.com`/`admin` credential and its login `min=5` validation unchanged unless a separate decision explicitly replaces that test fixture.

The smallest viable test scope is table-driven service tests for normalization, valid creation, short password, malformed/missing email, bcrypt non-plaintext storage, and duplicate-email mapping; handler tests for `201`, `400`, `401`, and `409`; plus a repository/integration test against PostgreSQL for the unique constraint/error translation if the project test environment supports it. Run `go test ./...` and `go vet ./...`; regenerate docs with the existing `just docs` command.

### Risks
- There is no role/permission model, so a protected endpoint would authorize every authenticated account, not only the seeded admin.
- Checking `FindByEmail` before insert alone is race-prone; the existing database `UNIQUE` constraint must remain the final duplicate-email guard.
- Email case normalization is currently done in login but not enforced by the database; creation must normalize consistently or case variants can coexist and later behave unexpectedly.
- The current login request says `min=5` because the prior seed deliberately uses password `admin`; applying `min=8` globally would break the documented/test credential.
- The repository has no HTTP integration tests and no configured PostgreSQL test harness, so full duplicate-conflict verification may require an environment-dependent integration test.
- Generated Swagger files are checked into the repository and must be regenerated together, not edited in only one representation.

### Ready for Proposal
Yes. The core behavior and affected boundaries are sufficiently clear. The proposal should explicitly record the route exposure decision (recommended: authenticated `POST /api/v1/users`) and state that admin-only authorization and password-policy changes to the existing seed/login flow are out of scope.
