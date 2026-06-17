# Task: Implement Email-Based Password Reset Feature

## Overview

Implement a complete self-service password reset flow. When a user forgets their password, they provide their email address, receive a reset link via Gmail, and can set a new password. The link expires after a set time and can only be used once.

---

## Current State

The basic infrastructure already exists and should be reused/extended:

- `pkg/password/resetpassword.go` — `CreateResetPassword()` generates a 16-character random secret and stores it in the DB
- `controllers/password.go` — has two relevant endpoints:
  - `POST /password/reset` — creates a token (currently **admin-only**, needs to become user-facing)
  - `POST /password/reset/:secret` — validates the secret and sets the new password
- **Missing**: token expiration, email sending, and a user-facing "forgot password" flow

---

## Email Service: Gmail SMTP

Emails are sent using a dedicated Gmail account via SMTP with an App Password. No external service or paid plan is required.

### Setup (manual steps — do before implementation)

1. Create a Gmail account dedicated to this app (e.g., `noreply.yourapp@gmail.com`)
2. Enable 2-Factor Authentication on that Google account
3. Go to **Google Account → Security → App Passwords** and generate an App Password for "Mail"
4. Add the credentials to the `.env` file (see below)

> ⚠️ Never commit the actual credentials to git. They go only in `app.env` (which is in `.gitignore`).

### Environment Variables

Add to `app.env` and document the keys (without values) in `example.env`:

```env
# Email (Gmail SMTP)
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_FROM_ADDRESS=noreply.yourapp@gmail.com
EMAIL_FROM_NAME=F-OAuth
EMAIL_SMTP_PASSWORD=xxxx xxxx xxxx xxxx   # Google App Password (16 chars with spaces)

# Password Reset
RESET_PASSWORD_TOKEN_EXPIRY=3600          # seconds — 1 hour
FRONTEND_RESET_PASSWORD_URL=http://localhost:3000/reset-password  # update for production
```

---

## Backend Changes

### 1. Config (`config/config.go`)

Add the new environment variables to the config struct so they are available across the app.

### 2. Extend the reset password DB model

Add an `expires_at` timestamp column to the `reset_password` table. The token should be invalidated after `RESET_PASSWORD_TOKEN_EXPIRY` seconds.

Also track whether the token has already been used (e.g., a `used` boolean column), so each link works only once.

### 3. Update `pkg/password/resetpassword.go`

- Increase token length from 16 to **32 characters** for better security
- Set `expires_at` to `now + RESET_PASSWORD_TOKEN_EXPIRY` on creation
- In the validation function (used by `POST /password/reset/:secret`), check:
  - Token exists
  - Token is not expired
  - Token has not been used
  - Mark token as used after successful password change

### 4. Email sending (`pkg/notification/email.go`) — new file

Create a function `SendPasswordResetEmail(toAddress string, resetLink string) error` that:

- Connects to `smtp.gmail.com:587` using `STARTTLS`
- Authenticates with `EMAIL_FROM_ADDRESS` + `EMAIL_SMTP_PASSWORD`
- Sends a plain-text (and optionally HTML) email to the user
- Email subject: `Reset your password`
- Email body contains the reset link and an expiry notice (e.g., "This link is valid for 1 hour")

Use Go's standard library `net/smtp` for the SMTP connection. No external packages required.

### 5. New user-facing endpoint: `POST /api/password/forgot`

This replaces the admin-only reset flow with a public one:

- **Input** (JSON body): `{ "email": "user@example.com" }`
- Lookup user by email in DB
- If the user does not exist: return HTTP 200 with a generic message (do not reveal if email exists)
- If user exists:
  - Call `CreateResetPassword(userId)` to generate token
  - Build reset link: `FRONTEND_RESET_PASSWORD_URL + "?token=" + secret + "&id=" + resetPasswordId`
  - Call `SendPasswordResetEmail(user.Email, resetLink)`
- Return HTTP 200 with `{ "status": "success", "message": "If that email is registered, you will receive a reset link shortly." }`
- Apply rate limiting: max **3 requests per email address per hour**

### 6. Keep existing endpoint `POST /api/password/reset/:secret`

This endpoint already handles token validation and applying the new password. Extend it to:

- Check `expires_at` (reject if expired)
- Check `used` flag (reject if already used)
- Mark token as `used = true` after success
- Return clear error messages: `"token expired"` vs `"token already used"` vs `"invalid token"`

---

## Frontend Changes

### 1. New routes

Create the following views and register them in `src/router/index.ts`:

- `/forgot-password` → `src/views/ForgotPasswordView.vue`
- `/reset-password` → `src/views/ResetPasswordView.vue`

### `ForgotPasswordView.vue`

- Simple form: email input + submit button
- On submit: `POST /api/password/forgot` with `{ email }`
- Always show success message after submit (do not reveal if email exists)
- Disable button and show loading state during request

### `ResetPasswordView.vue`

- Reads `token` and `id` query parameters from the URL
- Shows a form: new password + confirm password
- On submit: `POST /api/password/reset/:secret` with `{ password }` and the `id` from URL
- On success: redirect to login page with a success message
- On error: display user-friendly message (`"Link expired"`, `"Link already used"`, etc.)

### 2. Login page

Add a `"Forgot password?"` link below the login form that navigates to `/forgot-password`.

---

## Security Checklist

- [ ] Tokens are 32 characters (upgraded from 16)
- [ ] Tokens expire after 1 hour
- [ ] Tokens are single-use (marked as used after first successful use)
- [ ] Generic response on forgot-password endpoint (no email enumeration)
- [ ] Rate limiting on `POST /api/password/forgot` (max 3/hour per email)
- [ ] App Password stored only in `app.env`, never committed
- [ ] SMTP connection uses STARTTLS (port 587)

---

## Questions to Clarify Before / During Review

1. **Token expiry duration** — 1 hour is assumed. Change if needed. (make it 3 hours)
2. **Email sender name** — currently set to `F-OAuth`. Adjust if needed. this is okay
3. **Rate limiting** — 3 requests/hour per email is assumed. Change if needed.
4. **Confirmation email** — should the user receive a second email after successfully changing their password? (not included by default) no additional mail
5. **Frontend route names** — `/forgot-password` and `/reset-password` are assumed; adjust to match existing router conventions. thy are okay

---

## Resources

- Go `net/smtp` docs: <https://pkg.go.dev/net/smtp>
- Google App Passwords: <https://myaccount.google.com/apppasswords>
- Existing code to extend: `pkg/password/resetpassword.go`, `controllers/password.go`
