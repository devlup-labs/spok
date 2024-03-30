#!/bin/bash

# Read values from GitHub Actions secrets
clientID_value=$GOOGLE_CLIENT_ID
clientSecret_value=$GOOGLE_CLIENT_SECRET

# Replace values in constants.go using sed
sed -i.bak -e "s/clientID = \".*\"/clientID = \"$clientID_value\"/" -e "s/clientSecret = \".*\"/clientSecret = \"$clientSecret_value\"/" internal/pkg/constants/constants.go

# Clean up backup file created by sed
rm -f constants.go.bak

echo "Values replaced successfully."
