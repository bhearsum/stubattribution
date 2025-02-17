# Stub Service

Accepts an attribution code and bouncer parameters and returns a, potentially,
modified stub installer containing an attribution code.

## Environment Variables

### BASE_URL (Required)

Example: `BASE_URL=https://stubservice.services.mozilla.com/`

### HMAC_KEY

If set, the `attribution_code` parameter will be verified by validating that the
`attribution_sig` parameter matches the hex encoded sha256 hmac of
`attribution_code` using `HMAC_KEY`.

### HMAC_TIMEOUT (Default 10 minutes)

Will validate that the timestamp included in `attribution_code` is within
(Now-timeout) to Now. This variable should be in [duration
format](https://golang.org/pkg/time/#ParseDuration).

### SENTRY_DSN

If set, tracebacks will be sent to [Sentry](https://getsentry.com/).

### BOUNCER_URL

Bouncer root URL. The default value is: `https://download.mozilla.org/`.

### RETURN_MODE

Can be `direct` or `redirect`.

#### direct mode

Returns bytes directly to client

#### redirect mode

Writes bytes to a storage backend and returns a redirect response to the storage
location.

### STORAGE_BACKEND

The only valid value is: `gcs`.

### GCS_BUCKET (redirect mode)

The bucket where builds will be written.

### GCS_PREFIX (redirect mode)

A path prefix within the `GCS_BUCKET` where builds will be written.

### CDN_PREFIX (redirect mode)

A prefix which will be added to the storage key.
