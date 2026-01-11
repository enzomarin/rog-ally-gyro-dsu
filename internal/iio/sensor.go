package iio


import (
    "os"
    "strconv"
    "strings"
)

type SensorData struct {
    GyroX  float64
    GyroY  float64
    GyroZ  float64
    AccelX float64
    AccelY float64
    AccelZ float64
}

func readFile(path string) (string, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return "", err

    }
    return strings.TrimSpace(string(data)), nil
}

func readValue(device, attribute string) (float64, error) {
    text, err := readFile(device + "/" + attribute)
    if err != nil {
        return 0, err

    }
    return strconv.ParseFloat(text, 64)
}

func ReadSensor(device string) (*SensorData, error) {
    data := &SensorData{}
    
    gyroX, _ := readValue(device, "in_anglvel_x_raw")
    gyroY, _ := readValue(device, "in_anglvel_y_raw")
    gyroZ, _ := readValue(device, "in_anglvel_z_raw")
    gyroScale, _ := readValue(device, "in_anglvel_scale")
    
    data.GyroX = gyroX * gyroScale
    data.GyroY = gyroY * gyroScale
    data.GyroZ = gyroZ * gyroScale
    
    accelX, _ := readValue(device, "in_accel_x_raw")
    accelY, _ := readValue(device, "in_accel_y_raw")
    accelZ, _ := readValue(device, "in_accel_z_raw")
    accelScale, _ := readValue(device, "in_accel_scale")
    

    data.AccelX = accelX * accelScale
    data.AccelY = accelY * accelScale
    data.AccelZ = accelZ * accelScale
    
    return data, nil

}