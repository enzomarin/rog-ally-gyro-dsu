package dsu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"time"
	
	"rog-ally-gyro-dsu/internal/iio"
)

type Server struct {
	conn          *net.UDPConn
	clients       map[string]*net.UDPAddr
	packetCounter uint32
}

func NewServer(port int) (*Server, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return nil, err
	}
	
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	
	log.Printf("Server listening on 127.0.0.1:%d", port)
	
	server := &Server{
		conn:          conn,
		clients:       make(map[string]*net.UDPAddr),
		packetCounter: 0,
	}
	
	go server.listenClients()
	
	return server, nil
}

func (s *Server) listenClients() {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := s.conn.ReadFromUDP(buffer)
		if err != nil {
			continue
		}
		
		if n < 4 {
			continue
		}
		
		clientKey := addr.String()
		s.clients[clientKey] = addr
		
		log.Printf("📡 Request from client: %s (%d bytes)", addr, n)
		
		s.sendPortInfo(addr)
	}
}

func (s *Server) sendPortInfo(addr *net.UDPAddr) {
	buf := new(bytes.Buffer)
	
	// Slot info (11 bytes)
	binary.Write(buf, binary.LittleEndian, uint8(0))   // Slot
	binary.Write(buf, binary.LittleEndian, uint8(2))   // State (connected)
	binary.Write(buf, binary.LittleEndian, uint8(2))   // Model (full gyro)
	binary.Write(buf, binary.LittleEndian, uint8(1))   // Connection (USB)
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // MAC
	binary.Write(buf, binary.LittleEndian, uint8(0x05)) // Battery (full)
	
	// Padding
	binary.Write(buf, binary.LittleEndian, uint8(0))
	
	packet := s.createPacket(0x00100001, buf.Bytes())
	s.conn.WriteToUDP(packet, addr)
	
	log.Printf("Sen port info to %s", addr)
}

func (s *Server) createPacket(packetType uint32, payload []byte) []byte {
	buf := new(bytes.Buffer)
	
	// Header
	buf.WriteString("DSUS")
	binary.Write(buf, binary.LittleEndian, uint16(0x03E9)) // Version 1001 in hex
	binary.Write(buf, binary.LittleEndian, uint16(len(payload)+4)) // Size (payload + messageType)
	binary.Write(buf, binary.LittleEndian, uint32(0)) // CRC placeholder
	binary.Write(buf, binary.LittleEndian, uint32(0)) // Server ID
	binary.Write(buf, binary.LittleEndian, packetType) // Message type
	buf.Write(payload)
	
	packet := buf.Bytes()
	
	// Calculate CRC32 over entire packet with CRC field zeroed
	crc := crc32.ChecksumIEEE(packet)
	
	// Write CRC at position 8
	binary.LittleEndian.PutUint32(packet[8:12], crc)
	
	return packet
}

func (s *Server) createDataPacket(data *iio.SensorData) []byte {
	buf := new(bytes.Buffer)
	
	// Slot info (11 bytes)
	binary.Write(buf, binary.LittleEndian, uint8(0))   // Slot
	binary.Write(buf, binary.LittleEndian, uint8(2))   // State (connected)
	binary.Write(buf, binary.LittleEndian, uint8(2))   // Model (full gyro)
	binary.Write(buf, binary.LittleEndian, uint8(1))   // Connection (USB)
	buf.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}) // MAC
	binary.Write(buf, binary.LittleEndian, uint8(0x05)) // Battery (full)
	
	// Is connected (1 byte)
	binary.Write(buf, binary.LittleEndian, uint8(1))
	
	// Packet counter (4 bytes)
	binary.Write(buf, binary.LittleEndian, s.packetCounter)
	s.packetCounter++
	
	// Buttons bitmask 1 (1 byte)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	
	// Buttons bitmask 2 (1 byte)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	
	// HOME button (1 byte)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	
	// Touch button (1 byte)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	
	// Left stick (2 bytes)
	binary.Write(buf, binary.LittleEndian, uint8(128)) // X centered
	binary.Write(buf, binary.LittleEndian, uint8(128)) // Y centered
	
	// Right stick (2 bytes)
	binary.Write(buf, binary.LittleEndian, uint8(128)) // X centered
	binary.Write(buf, binary.LittleEndian, uint8(128)) // Y centered
	
	// Analog buttons (12 bytes)
	for i := 0; i < 12; i++ {
		binary.Write(buf, binary.LittleEndian, uint8(0))
	}
	
	// Touch 1 (6 bytes)
	binary.Write(buf, binary.LittleEndian, uint8(0))   // active
	binary.Write(buf, binary.LittleEndian, uint8(0))   // id
	binary.Write(buf, binary.LittleEndian, uint16(0))  // x
	binary.Write(buf, binary.LittleEndian, uint16(0))  // y
	
	// Touch 2 (6 bytes)
	binary.Write(buf, binary.LittleEndian, uint8(0))
	binary.Write(buf, binary.LittleEndian, uint8(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	binary.Write(buf, binary.LittleEndian, uint16(0))
	
	// Timestamp (8 bytes - microseconds)
	timestamp := uint64(time.Now().UnixNano() / 1000)
	binary.Write(buf, binary.LittleEndian, timestamp)
	
	// Accelerometer (12 bytes - 3 floats in G's)
	binary.Write(buf, binary.LittleEndian, float32(data.AccelX/9.81))
	binary.Write(buf, binary.LittleEndian, float32(-data.AccelY/9.81))
	binary.Write(buf, binary.LittleEndian, float32(-data.AccelZ/9.81))
	
	// Gyroscope (12 bytes - 3 floats in deg/s)
	binary.Write(buf, binary.LittleEndian, float32(data.GyroX*57.2958))
	binary.Write(buf, binary.LittleEndian, float32(-data.GyroY*57.2958))
	binary.Write(buf, binary.LittleEndian, float32(-data.GyroZ*57.2958))
	
	return s.createPacket(0x00100002, buf.Bytes())
}

func (s *Server) SendData(data *iio.SensorData) error {
	if len(s.clients) == 0 {
		return nil
	}
	
	packet := s.createDataPacket(data)
	
	for _, addr := range s.clients {
		s.conn.WriteToUDP(packet, addr)
	}
	
	return nil
}

func (s *Server) Close() {
	s.conn.Close()
}