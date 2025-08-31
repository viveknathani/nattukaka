#!/bin/bash

BASE_DIR=".ssh"

mkdir -p "$BASE_DIR"

CURRENT_INDEX=$(ls "$BASE_DIR"/id-*.pem 2>/dev/null | sed -E 's/.*id-([0-9]+)\.pem/\1/' | sort -n | tail -1 || echo 0)
if [ -z "$CURRENT_INDEX" ]; then
  CURRENT_INDEX=0
fi
NEXT_INDEX=$((CURRENT_INDEX + 1))

OLD_KEY="$BASE_DIR/id-$CURRENT_INDEX.pem"
NEW_KEY="$BASE_DIR/id-$NEXT_INDEX.pem"

SERVERS=(
  "root@daya-0"
  "root@daya-1"
)

echo "[*] generating new key: id-$NEXT_INDEX.pem"
ssh-keygen -t rsa -b 4096 -m PEM -f "$NEW_KEY" -N ""
chmod 400 "$NEW_KEY"

if [ "$CURRENT_INDEX" -gt 0 ]; then
  echo "[*] deploying to servers using old key (id-$CURRENT_INDEX.pem)..."
  for server in "${SERVERS[@]}"; do
    echo "   -> $server"
    cat "$NEW_KEY.pub" | ssh -i "$OLD_KEY" -o StrictHostKeyChecking=accept-new "$server" \
      "mkdir -p ~/.ssh && chmod 700 ~/.ssh && cat > ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys"
  done
else
  echo "[!] no old key found, skipping server deploy, install id-$NEXT_INDEX.pem.pub manually"
fi

echo ""
echo "✅ rotation complete!"
echo "   old key (if any): $OLD_KEY"
echo "   current key: $NEW_KEY"
echo ""
echo "👉 connect with:"
echo "   ssh -i $NEW_KEY user@server"
