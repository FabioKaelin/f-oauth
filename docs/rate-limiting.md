# API rate limiting

Release1.4.0 uses fcommon1.3.0 and the shared private Redis service.
The approved authenticated allowance is1,200 requests/minute per verified user,
with burst120. Counters are shared across sessions and replicas. Ordinary API
quotas start in observe mode; sensitive login, registration and password/OAuth
quotas remain enforced. See fcommon/pkg/ratelimit and the infrastructure runbook
in minikube-monitoring/resources/rate-limiting for exact policies and provisioning.

Provision the api-ratelimit Secret before deploying the source manifests.
Deploy the hardened public proxy configuration before enabling the limiter so
forwarded client IP is accepted only from the authenticated edge. Other callers
are keyed by their actual TCP peer. Internal session validation only receives a
separate quota after the calling service secret has been verified; user sessions
are still validated. No private-IP or submitted-user-ID exception is accepted.

The frontend displays429 and temporary503 responses, an explicit retry countdown,
and retained form input. It never automatically replays rejected writes or OAuth
callback codes. Password-reset recipient suppression keeps a generic200 response
for known, unknown and external-provider recipients. Redis outages fail sensitive
operations temporarily with503/Retry-After; ordinary APIs use bounded per-process
fallback protection with weaker aggregate guarantees and an alert.

Validation: shared-library tests/race checks include concurrency, expiry/refill,
forged headers, recipient suppression and store outages. App tests/builds and
browser fixtures cover320/393/768/1440 widths, retained login input,429/503,
countdown expiry and absence of automatic write replay.
