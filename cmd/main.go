package main

import (
    "log"
    "time"
    
    "rog-ally-gyro-dsu/internal/dsu"
    "rog-ally-gyro-dsu/internal/iio"
)

func main() {
    log.Println("Starting DSU server...")
    
    server, err := dsu.NewServer(26760)
    if err != nil {
        log.Fatal(err)
    }
    defer server.Close()
    
    device := "/sys/bus/iio/devices/iio:device0"
    
    for {
        data, err := iio.ReadSensor(device)
        if err != nil {
            log.Printf("Error reading sensor: %v", err)
            time.Sleep(1 * time.Second)
            continue
        }
        
        server.SendData(data)
        time.Sleep(10 * time.Millisecond)
    }
}