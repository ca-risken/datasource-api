#!/bin/sh -e

if [ -z "$GOOGLE_SERVICE_ACCOUNT_JSON" ]; then
  exec "$@"
fi

CREDENTIAL_FILE="/tmp/credential.json"
if [ -n "$GOOGLE_CREDENTIAL_PATH" ]; then
  CREDENTIAL_FILE=$GOOGLE_CREDENTIAL_PATH
fi

cat <<_EOF > $CREDENTIAL_FILE
$GOOGLE_SERVICE_ACCOUNT_JSON
_EOF

chmod 600 ${CREDENTIAL_FILE}
exec "$@"
