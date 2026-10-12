# F-OAuth Backend

## Rolling sessions

Authentication uses random opaque 256-bit tokens with a 30-day inactivity timeout. Only SHA-256 token hashes are stored in the `sessions` table, which is created idempotently at startup. Successful validation through `/api/users/me` extends database expiry to 30 days from the request and returns a renewed shared `token` cookie (Domain=fabkli.ch, Path=/, Secure, HttpOnly, SameSite=Lax). No absolute lifetime is imposed: users who keep returning within 30 days remain signed in.

Tipp, media, gallery and the gateway use `github.com/fabiokaelin/fcommon v1.2.6` to validate sessions and relay the renewal cookie back to the browser. `/api/sessions/check` validates without extending expiry and is used for passive WebSocket revocation checks. It does not send a renewal cookie. Expired sessions cannot be revived; cleanup removes expired records hourly.

Logout deletes the current session. Password changes and completed password resets update the password and revoke all user sessions in one transaction. Reset-token consumption is atomic, preventing reuse. A racing local login authenticated with an old password is rejected when creating the session. Initial password-reset requests do not revoke sessions until a valid reset is completed.

Legacy JWT authentication and TOKEN_EXPIRED_IN/TOKEN_MAXAGE/JWT_SECRET configuration are retired. Existing users must log in once after migration. `TOKEN_URL=fabkli.ch` still controls the shared cookie domain. GitHub and Google login must start through `/api/sessions/oauth/start/{provider}?from=...`; callbacks validate host-only browser-bound state cookies and accept only HTTPS FabKli return URLs. The frontend uses these start endpoints. Token values, provider codes, OAuth state and password-reset secrets are redacted from backend logs. Credentialed CORS accepts HTTPS FabKli origins only.

Tests: `go test ./...` and `go vet ./...`. To exercise real transaction/expiry behaviour, run `OAUTH_SESSION_TEST_DSN='root:password@tcp(127.0.0.1:13376)/oauth_session_test?parseTime=true' go test ./pkg/db`. The database must be an isolated `oauth_session_test` schema; its test tables are recreated. Never point this test DSN at production. The live verification scripts are in `/home/admin/code/f-gateway/deploy/`, with database fixtures cleaned up automatically.

## Links

* [File Upload](https://codevoweb.com/how-to-upload-single-and-multiple-files-in-golang/)
* [Good Tutorials](https://codevoweb.com/golang/)

## Gmail password reset configuration

The backend reads local environment variables from [app.env](app.env). For development, the Gmail values go there. For production, the same variable names go into the Kubernetes `Secret`.

### What to fill in

Add or update these values:

```env
# Email (Gmail SMTP)
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_FROM_ADDRESS=your.account@gmail.com
EMAIL_FROM_NAME=F-OAuth
EMAIL_SMTP_PASSWORD=abcd efgh ijkl mnop

# Password reset
RESET_PASSWORD_TOKEN_EXPIRY=10800
FRONTEND_RESET_PASSWORD_URL=https://your-frontend-domain/reset-password
```

### Where to find each value

#### `EMAIL_SMTP_HOST`

* For Gmail this is always `smtp.gmail.com`
* No lookup needed
* In production this defines which SMTP server is used for password reset emails

#### `EMAIL_SMTP_PORT`

* For Gmail this should be `587`
* Port `587` is the standard SMTP submission port with TLS negotiation
* In production this defines how the backend connects to Gmail; a wrong value usually means mail sending fails

#### `EMAIL_FROM_ADDRESS`

* This is the Gmail address that sends the reset emails
* Best used with a dedicated Gmail account for the app, for example `your.project@gmail.com`
* In production this changes the sender address shown in the mailbox

#### `EMAIL_FROM_NAME`

* This is the display name shown to the user, for example: `F-OAuth <your.account@gmail.com>`
* This value is chosen manually; Google does not generate it
* In production this only changes the visible sender label, not the authenticated Gmail account

#### `EMAIL_SMTP_PASSWORD`

* This is **not** your normal Gmail password
* It must be a Google **App Password**
* To get it:
  1. Sign in to the Gmail account used for the app
  2. Open your Google Account settings
  3. Go to `Security`
  4. Enable `2-Step Verification` first if it is not already enabled
  5. Open `App passwords`
  6. Create a new app password for mail usage
  7. Copy the generated 16-character password into `EMAIL_SMTP_PASSWORD`
* Direct link after login: <https://myaccount.google.com/apppasswords>
* In production this is the credential used to authenticate against Gmail; if it is invalid, email sending fails

#### `RESET_PASSWORD_TOKEN_EXPIRY`

* Value is in **seconds**
* `10800` means 3 hours
* Examples:
  * `3600` = 1 hour
  * `10800` = 3 hours
  * `86400` = 24 hours
* In production this changes how long a reset link stays valid after it is generated
* Shorter values increase security; longer values improve convenience

#### `FRONTEND_RESET_PASSWORD_URL`

* This must be the public frontend URL that opens the reset-password page
* For local development, something like `http://localhost:3000/reset-password`
* For production, something like `https://your-domain.com/reset-password`
* In production this is the link included in the email; if it is wrong, the reset email still sends but points to the wrong page

### Recommended Gmail setup

Use a dedicated Gmail account for the application.

* Better separation from your personal account
* Easier to revoke or rotate later
* Cleaner audit trail for password reset mails

### Local vs live environment

#### Local development

* Put the values into [app.env](app.env)
* Start the backend normally
* Test with your real Gmail app password

#### Production

* Put the same variable names into the Kubernetes `Secret`
* The deployment then exposes them as environment variables to the backend container

### Security notes

* Never store your normal Gmail password in `EMAIL_SMTP_PASSWORD`
* Only use a Google App Password
* If the secret leaks, revoke it in Google Account settings and generate a new one
* If you change the Gmail account, update both `EMAIL_FROM_ADDRESS` and `EMAIL_SMTP_PASSWORD`

## Notes
