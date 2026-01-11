#!/bin/bash

echo "🗑️  Uninstalling ROG Ally Gyro DSU Server..."
echo ""

# Stop and disable service
systemctl --user stop rog-ally-gyro-dsu.service 2>/dev/null || true
systemctl --user disable rog-ally-gyro-dsu.service 2>/dev/null || true

# Remove files
rm -f ~/.local/bin/rog-ally-gyro-dsu
rm -f ~/.config/systemd/user/rog-ally-gyro-dsu.service

# Reload systemd
systemctl --user daemon-reload

echo ""
echo "✅ Uninstallation complete!"
echo ""
echo "Press Enter key to close this window..."
read -n 1 -s -r