# Proposal: Create Users Endpoint

## Intent

Add an authenticated API operation for creating users while enforcing normalized, unique email addresses and a minimum eight-character password. Passwords must be stored only as bcrypt hashes, and creation must not alter the existing login policy or seeded admin credential.

## Scope

### In Scope
- Add `POST /api/v1/users` behind the existing bearer middleware.
- Implement user creation through the user repository, service, and handler layers.
- Validate `required,email` and `required,min=8`; normalize email with trim/lowercase.
- Return a non-sensitive user representation with `201 Created`; map duplicate email to `409 Conflict`.
- Add unit/handler coverage, an environment-dependent PostgreSQL uniqueness test where feasible, and regenerate checked-in Swagger artifacts.

### Out of Scope
- Admin/role/permission modeling or restricting creation to the seeded admin.
- Public self-registration (`/auth/register`).
- Changing login validation from `min=5` or replacing `admin@example.com` / `admin`.
- Database migrations: `users.email` already has the required unique constraint.

## Capabilities

### New Capabilities
- `user-creation`: Authenticated creation of users with validation, normalization, bcrypt storage, and duplicate-email conflict behavior.

### Modified Capabilities
- None.

## Approach

Wire `POST /api/v1/users` to the existing protected API group. The service normalizes and validates input, hashes the plaintext password with bcrypt, and returns a safe entity. The repository inserts the record and translates PostgreSQL unique violations into a domain conflict; the database constraint remains the race-safe authority. Update application wiring and regenerate `docs/swagger.yaml`, `docs/swagger.json`, and `docs/docs.go`.

## API and Error Semantics

- Request: email and plaintext password.
- `201`: created user data excluding `password_hash`.
- `400`: malformed JSON or validation failure, including passwords shorter than eight characters.
- `401`: missing or invalid bearer token.
- `409`: normalized email already exists, including concurrent insert races.
- Other persistence failures retain the project’s existing server-error response convention.

## Affected Areas

| Area | Impact | Description |
|------|--------|-------------|
| `internal/user/*` | Modified | Repository, service, errors, handler, and tests. |
| `internal/http/router.go`, `cmd/api/main.go` | Modified | Route and dependency wiring. |
| `docs/swagger.*`, `docs/docs.go` | Modified | Generated API contract. |
| `migrations/000001_create_users.up.sql` | No change | Existing unique constraint is sufficient. |

## Risks

| Risk | Likelihood | Mitigation |
|------|------------|------------|
| Any authenticated user can create accounts without roles | High | Confirm authorization scope before specs; document as explicit boundary. |
| Application pre-check permits duplicates under concurrency | High | Rely on and test the database unique constraint. |
| Case variants diverge from login behavior | Med | Normalize before insert and lookup. |
| No configured PostgreSQL test harness | Med | Keep integration coverage skippable/environment-dependent; run unit tests, `go test ./...`, and `go vet ./...`. |

## Rollback Plan

Revert the route, wiring, user-layer implementation/tests, and generated Swagger changes. No migration rollback is required.

## Dependencies

- Existing bearer middleware, PostgreSQL `users.email` unique constraint, bcrypt dependency, and `just docs` generation command.

## Success Criteria

- [ ] Authenticated valid requests create users with normalized email and bcrypt-hashed passwords.
- [ ] Invalid input returns `400`; unauthenticated requests return `401`; duplicate emails return `409`.
- [ ] Existing seeded-admin login remains functional.
- [ ] Tests, vet, and regenerated Swagger artifacts pass repository checks.
