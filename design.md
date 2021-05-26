# Design Doc
​
## Features/Scope
- Login page
  - On successful login, go to dashboard
  - On failed login, alert user and do not change pages
- Dashboard
  - For Level 1, just a dummy page
  - On click logout, log user out and go to login page
​
## Tech
​
### Frontend
Choices here are mostly down to personal preference/familiarity.
​
**Language: TypeScript**
​
**Framework: React**
​
### Backend
**Language: Go**
​
**Framework: [Gin](https://github.com/gin-gonic/gin)**
​
- While the Go standard library's HTTP capabilities are plenty to support the app, using a framework like Gin can reduce boilerplate.
​
### Database
**Database: MySQL**
​
The app really only needs some basic CRUD operations; MySQL is simple and gets the job done.
​
Why not:
- PostgreSQL - wouldn't hurt, but it's more complicated than we need.
- SQLite - although a file-based DB would allow us to run one fewer process and simplify setup, SQLite doesn't have user management and isn't very concurrency friendly.
​
Schema:
- Accounts (account_id PK, email, pw_hash, salt)
- Sessions (account_id (nullable), session_token PK, csrf_token, expire_idle, expire_abs)
​
### Other
docker-compose for launching multiple containers and [secrets management](https://docs.docker.com/engine/swarm/secrets/#use-secrets-in-compose).
​
## API
POST `/api/login`
​
- Logs in user. Expects Basic Auth.
  - No post body; username, password, and csrf token will all be sent in headers
- Return 200 and set session token on successful login.
  - If user is logged in, do nothing.
- Return 401 on failed login.
  - Bad user/pass combo
  - Bad csrf token
​
POST `/api/logout`
​
- Logs out user. Invalidates session token.
  - If user is not logged in, do nothing.
  - A user is considered logged in if the Sessions table contains a tuple with their account ID that hasn’t expired (see Sessions for info about expiration).
- Return 200.
​
GET `/*`
​
- Fetch static files.
  - If file is index.html, inject CSRF token.
- Return 200.
​
## Security
​
### Sessions
Following [OWASP recommendations](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html):
- 16 byte random token
- Set cookie on authenticate
  - Make sure to set these attributes: Secure, HttpOnly, SameSite
- Session creation
  - When a user first visits the site, a pre-auth session (null account_id, new session token and csrf token, and appropriate expire times) is created, and the created csrf token is sent back to the user.
  - When a user attempts to log in, the server looks for a Sessions entry that matches the sent csrf token and pre-auth session ID. If a match is found, a new session replaces the pre-auth session, with account_id determined from the username and password.
- Expiring sessions (the backend will delete an active session if any of the following occur):
  - A logged in user makes a logout request
  - A specified amount of time passes after a logged in user makes an authenticated request (Idle expire)
  - A specified amount of time passes after a logged in user's initial login (Absolute expire)
    - Note: a background task will remove expired sessions periodically; however, the server will still need to verify that exire_idle and expire_abs aren't in the past.
​
### Passwords
Passwords will be salted and hashed with argon2id, [as recommended by OWASP](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html#argon2id).
- Check OWASP for ideal parameter settings
- random per-user salt
​
### Attack Mitigation
**XSS**: React should handle HTML escaping when rendering. We shouldn't ever need to use any [dangerous rendering methods](https://stackoverflow.com/a/51852579).
​
**CSRF**: Inject CSRF token into index.html when it's fetched; when performing future requests, client must put token in custom header.
​
**SSRF**: Server shouldn't ever need to make any HTTP requests, much less requests to user-specified URLs.
​
**SQL Injection**: Use prepared statements/parameterized queries.
