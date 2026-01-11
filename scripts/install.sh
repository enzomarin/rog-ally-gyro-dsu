#!/bin/bash
set -e

REPO="enzomarin/rog-ally-gyro-dsu"
VERSION="latest"

echo "🚀 Installing ROG Ally Gyro DSU Server..."

# Create directories
mkdir -p ~/.local/bin
mkdir -p ~/.config/systemd/user

# Download latest binary
echo "📥 Downloading latest version..."
DOWNLOAD_URL="https://github.com/${REPO}/releases/${VERSION}/download/rog-ally-gyro-dsu"
curl -L "${DOWNLOAD_URL}" -o ~/.local/bin/rog-ally-gyro-dsu
chmod +x ~/.local/bin/rog-ally-gyro-dsu

# Download service file
echo "📥 Downloading service file..."
SERVICE_URL="https://github.com/${REPO}/releases/${VERSION}/download/rog-ally-gyro-dsu.service"
curl -L "${SERVICE_URL}" -o ~/.config/systemd/user/rog-ally-gyro-dsu.service

# Enable and start service
echo "🔧 Configuring service..."
systemctl --user daemon-reload
systemctl --user enable rog-ally-gyro-dsu.service
systemctl --user start rog-ally-gyro-dsu.service

echo ""
echo "✅ Installation complete!"
echo ""
echo "📊 Check status: systemctl --user status rog-ally-gyro-dsu"
echo "📝 View logs: journalctl --user -u rog-ally-gyro-dsu -f"
echo ""

read -p "Press any key to continue..."