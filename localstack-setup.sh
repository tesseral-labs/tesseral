#!/bin/sh
echo "Initializing LocalStack KMS..."

# Function to find a key by description and return its ID or ARN
get_key_by_description() {
  DESCRIPTION=$1
  awslocal kms list-keys --query "Keys[*].KeyId" --output text | while read KEY_ID; do
    if [ "$(awslocal kms describe-key --key-id $KEY_ID --query 'KeyMetadata.Description' --output text)" = "$DESCRIPTION" ]; then
      echo $KEY_ID
      return 0
    fi
  done
  return 1 # Key not found
}

# Function to update or add a key-value pair in the .env file
update_env_file() {
  FILE=$1
  KEY=$2
  VALUE=$3
  TEMP_FILE=$(mktemp)

  if grep -q "^$KEY=" "$FILE"; then
    if ! grep -q "^$KEY=$VALUE$" "$FILE"; then
      # Update the existing key-value pair
      sed "s|^$KEY=.*|$KEY=$VALUE|" "$FILE" > "$TEMP_FILE"
      cat "$TEMP_FILE" > "$FILE"
      echo "Updated $KEY in $FILE"
    else
      echo "$KEY already up-to-date in $FILE"
    fi
  else
    # Add the new key-value pair
    cat "$FILE" > "$TEMP_FILE"
    echo "$KEY=$VALUE" >> "$TEMP_FILE"
    cat "$TEMP_FILE" > "$FILE"
    echo "Added $KEY to $FILE"
  fi

  rm "$TEMP_FILE"
}

# Ensure intermediate-session-signing-keys-key exists
INTERMEDIATE_KEY_DESCRIPTION="intermediate-session-signing-keys-key"
INTERMEDIATE_KEY_ID=$(get_key_by_description "$INTERMEDIATE_KEY_DESCRIPTION")
if [ -z "$INTERMEDIATE_KEY_ID" ]; then
  INTERMEDIATE_KEY_ID=$(awslocal kms create-key --description "$INTERMEDIATE_KEY_DESCRIPTION" --query KeyMetadata.KeyId --output text)
  echo "Created key '$INTERMEDIATE_KEY_DESCRIPTION' with ID: $INTERMEDIATE_KEY_ID"
else
  echo "Found existing key '$INTERMEDIATE_KEY_DESCRIPTION' with ID: $INTERMEDIATE_KEY_ID"
fi
# Write to the .env file
update_env_file /etc/localstack/init-output.env "API_INTERMEDIATE_SESSION_KMS_KEY_ID" "$INTERMEDIATE_KEY_ID"


# Ensure session-signing-keys-key exists
SESSION_KEY_DESCRIPTION="session-signing-keys-key"
SESSION_KEY_ID=$(get_key_by_description "$SESSION_KEY_DESCRIPTION")
if [ -z "$SESSION_KEY_ID" ]; then
  SESSION_KEY_ID=$(awslocal kms create-key --description "$SESSION_KEY_DESCRIPTION" --query KeyMetadata.KeyId --output text)
  echo "Created key '$SESSION_KEY_DESCRIPTION' with ID: $SESSION_KEY_ID"
else
  echo "Found existing key '$SESSION_KEY_DESCRIPTION' with ID: $SESSION_KEY_ID"
fi
# Write to the .env file
update_env_file /etc/localstack/init-output.env "API_SESSION_KMS_KEY_ID" "$SESSION_KEY_ID"

echo "LocalStack KMS initialization complete."
