# Security Policy

## Reporting Security Vulnerabilities

If you discover a security vulnerability in this project, please report it by emailing the maintainer directly rather than opening a public issue.

## Security Best Practices

### Configuration

1. **JWT Secret**: Always use a strong, randomly generated JWT secret in production. Generate one using:
   ```bash
   openssl rand -base64 64
   ```

2. **HTTPS**: Always use HTTPS in production. The application sets the `Secure` flag on cookies when not running on localhost.

3. **Environment Variables**: Never commit sensitive credentials to version control. Use environment variables for all secrets.

### OAuth2 Redirect Security

The application implements strict redirect URL validation to prevent open redirect attacks:
- Only redirects to the configured `FRONTEND_ORIGIN` are allowed
- Relative paths (starting with `/`) are allowed and prepended with the frontend origin
- All other redirect attempts default to `/profile`

### CORS Policy

CORS is configured to only allow requests from:
- The configured `FRONTEND_ORIGIN`
- Local development origins (`http://localhost:3000`, `http://localhost:5173`)

### Password Security

- Passwords are hashed using bcrypt with a cost factor of 12
- Plain text passwords are never stored or logged

### Database Security

- All database queries use parameterized statements to prevent SQL injection
- Database credentials should be stored securely and rotated regularly

## Security Improvements Made

1. **Fixed Open Redirect Vulnerability**: Added strict validation for OAuth2 redirect URLs
2. **Fixed CORS Misconfiguration**: Changed from accepting any origin to an allowlist approach
3. **Improved Cookie Security**: Enabled `Secure` flag for cookies in non-localhost environments
4. **Removed Information Disclosure**: Removed debug logging of sensitive tokens
5. **Enhanced Configuration**: Added security warnings to example configuration files

## Recommended Production Settings

```env
# Use HTTPS
FRONTEND_ORIGIN=https://your-domain.com
TOKEN_URL=your-domain.com

# Strong JWT secret (use openssl rand -base64 64)
JWT_SECRET=your-strong-random-secret-here

# Reasonable token expiration
TOKEN_EXPIRED_IN=60m
TOKEN_MAXAGE=60

# Database over TLS (recommended)
DATABASE_HOST=your-secure-db-host
```

## Regular Security Maintenance

- Keep all dependencies up to date
- Regularly review and rotate secrets
- Monitor application logs for suspicious activity
- Conduct regular security audits
- Enable and monitor CodeQL scans
