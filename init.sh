#!/bin/sh

JWT_SECRET=$(openssl rand -base64 32)

SECRETS_FILE="${SECRETS_FILE:-config/secrets.yaml}"

if [ ! -f "$SECRETS_FILE" ]; then
    echo "Creating $SECRETS_FILE..."
    cat > "$SECRETS_FILE" <<EOF
jwt:
  secret_key: ${JWT_SECRET}
EOF
else
    if grep -q "secret_key:" "$SECRETS_FILE"; then
        sed -i "s|secret_key:.*|secret_key: ${JWT_SECRET}|" "$SECRETS_FILE"
    else
        if ! grep -q "jwt:" "$SECRETS_FILE"; then
            echo "" >> "$SECRETS_FILE"
            echo "jwt:" >> "$SECRETS_FILE"
        fi
        echo "  secret_key: ${JWT_SECRET}" >> "$SECRETS_FILE"
    fi
fi

echo "JWT secret key has been generated and written to $SECRETS_FILE"

exec "$@"
