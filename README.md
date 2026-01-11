
# ROG Ally Gyro DSU

DSU (Cemuhook protocol) server for ROG Ally gyroscope/accelerometer on Bazzite/SteamOS.

Enables motion controls in emulators (Yuzu, Dolphin, Citra, Cemu, etc.) using the ROG Ally's built-in IMU sensor.

## Features

- ✅ Reads gyroscope and accelerometer data from ROG Ally IIO sensors
- ✅ Exposes motion data via DSU/Cemuhook protocol (UDP port 26760)
- ✅ Works with any emulator that supports Cemuhook motion
- ✅ Runs as systemd user service (auto-start on boot)
- ✅ One-click installation with .desktop file
- ✅ Lightweight and efficient (written in Go)

## Requirements

- ROG Ally (or compatible device with BMI323 IMU sensor)
- Bazzite, ChimeraOS, or SteamOS

## Installation

### Quick Install (Recommended)

Open this page in the browser in **Desktop Mode**.

1. Download [install-rog-ally-gyro-dsu.desktop](https://github.com/enzomarin/rog-ally-gyro-dsu/releases/latest/download/install-rog-ally-gyro-dsu.desktop)
2. Save it to Desktop
3. In Dolphin, right-click it → **Properties** → **Permissions** → Check **"Is executable"**
4. Double-click **Install ROG Ally Gyro DSU** to install

The installer will automatically:
- Download the latest binary
- Install as a systemd user service
- Start the server

### Uninstall

1. Download [uninstall-rog-ally-gyro-dsu.desktop](https://github.com/enzomarin/rog-ally-gyro-dsu/releases/latest/download/uninstall-rog-ally-gyro-dsu.desktop)
2. Right-click → **Properties** → **Permissions** → Check **"Is executable"**
3. Double-click to uninstall

### Manual Installation
```bash
# Download and run install script
curl -sL https://raw.githubusercontent.com/enzomarin/rog-ally-gyro-dsu/main/scripts/install.sh | bash
```

## Usage

Server runs automatically on startup. It listens on `127.0.0.1:26760` (standard Cemuhook port).

### Service Management
```bash
# Check status
systemctl --user status rog-ally-gyro-dsu

# View logs
journalctl --user -u rog-ally-gyro-dsu -f

# Restart service
systemctl --user restart rog-ally-gyro-dsu

# Stop service
systemctl --user stop rog-ally-gyro-dsu
```

### Emulator Configuration

Configure your emulator to use DSU motion at `127.0.0.1:26760`

#### Yuzu / Ryujinx
1. **Settings** → **Controls** → **Motion**
2. Motion provider: **Cemuhook UDP**
3. Server: `127.0.0.1`
4. Port: `26760`

#### Dolphin
1. **Controllers** → **Alternate Input Sources**
2. Enable **DSU Client**
3. Add server: `127.0.0.1:26760`

#### Citra
1. **Emulation** → **Configure** → **Motion/Touch**
2. Motion Provider: **CemuhookUDP**
3. Server: `127.0.0.1:26760`

#### Cemu (via Cemuhook)
1. Install Cemuhook plugin
2. **Input Settings** → Select **DSU1** as motion source

#### RPCS3
1. **Pads** → **Pad Settings**
2. **Device** → Select your controller
3. **Motion Settings** → Enable motion
4. Server: `127.0.0.1:26760`

## Testing

After installation:

1. Switch to **Gaming Mode**
2. Open your emulator
3. Configure motion controls (see above)
4. **Physically move/rotate the ROG Ally**
5. Motion should be detected in the emulator

**Note:** You must physically move the device to see gyro data - values will be near zero when stationary.

## Troubleshooting

### Server won't start - "address already in use"

Another DSU server is running on port 26760:
```bash
# Stop any running instances
systemctl --user stop rog-ally-gyro-dsu
pkill -f rog-ally-gyro-dsu

# Or kill the port directly
fuser -k 26760/udp

# Restart
systemctl --user start rog-ally-gyro-dsu
```

### Emulator doesn't detect motion

1. **Verify server is running:**
```bash
   systemctl --user status rog-ally-gyro-dsu
```

2. **Check logs for errors:**
```bash
   journalctl --user -u rog-ally-gyro-dsu -f
```

3. **Test with movement** - physically rotate/tilt the ROG Ally (gyro shows near-zero when stationary)

4. **Verify emulator configuration:**
   - Server: `127.0.0.1`
   - Port: `26760`

### Installation fails

Re-download the installer and ensure it's executable:
```bash
chmod +x install-rog-ally-gyro-dsu.desktop
./install-rog-ally-gyro-dsu.desktop
```

### No gyroscope data

Verify IIO sensor is available:
```bash
ls /sys/bus/iio/devices/iio:device0/
cat /sys/bus/iio/devices/iio:device0/name  # Should show "bmi323"
```

If sensor is missing, your device may not be supported.

## Building from Source
```bash
# Clone repository
git clone https://github.com/enzomarin/rog-ally-gyro-dsu.git
cd rog-ally-gyro-dsu

# Build
go build -ldflags="-s -w" -o rog-ally-gyro-dsu ./cmd

# Run directly (for testing)
./rog-ally-gyro-dsu
```


## License

MIT License - see [LICENSE](LICENSE) file for details
