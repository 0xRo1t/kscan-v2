package mqtt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	timeout = 5 * time.Second

	// MQTT CONNECT 固定报头：包类型=1(CONNECT), flags=0
	packetTypeConnect = 0x10
	// MQTT CONNACK 固定报头
	packetTypeConnack = 0x20

	// CONNACK 返回码
	connackAccepted          = 0x00
	connackBadUsernameOrPass = 0x04
	connackNotAuthorized     = 0x05
)

// Check 尝试使用给定凭据连接 MQTT Broker（端口默认 1883）
// 返回 nil 表示认证成功，否则返回错误
func Check(host, username, password string, port int) error {
	if port == 0 {
		port = 1883
	}
	netloc := net.JoinHostPort(host, fmt.Sprint(port))
	conn, err := net.DialTimeout("tcp", netloc, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}

	pkt, err := buildConnect(username, password)
	if err != nil {
		return err
	}

	if _, err = conn.Write(pkt); err != nil {
		return err
	}

	return readConnack(conn)
}

// buildConnect 构造 MQTT 3.1.1 CONNECT 报文
func buildConnect(username, password string) ([]byte, error) {
	// 协议名 "MQTT"
	protoName := []byte{0x00, 0x04, 'M', 'Q', 'T', 'T'}
	// 协议级别 4 = MQTT 3.1.1
	protoLevel := byte(0x04)

	// Connect Flags
	// bit7: UsernameFlag  bit6: PasswordFlag  bit1: CleanSession
	var connectFlags byte = 0x02 // CleanSession=1
	if username != "" {
		connectFlags |= 0x80
	}
	if password != "" {
		connectFlags |= 0x40
		// MQTT 3.1.1 规定：如果设置了 PasswordFlag，UsernameFlag 必须同时为 1。
		if username == "" {
			connectFlags |= 0x80
		}
	}

	// Keep-alive: 10 秒
	keepAlive := []byte{0x00, 0x0A}

	// Client Identifier（随机短串避免与在线设备冲突）
	clientID := "kscan_" + fmt.Sprintf("%d", time.Now().UnixNano()%100000)

	// 可变报头 = 协议名 + 协议级别 + 连接标志 + Keep-alive
	varHeader := append(protoName, protoLevel, connectFlags)
	varHeader = append(varHeader, keepAlive...)

	// Payload：ClientID + [Username] + [Password]
	payload := encodeMQTTStr(clientID)
	if username != "" {
		payload = append(payload, encodeMQTTStr(username)...)
	}
	if password != "" {
		payload = append(payload, encodeMQTTStr(password)...)
	}

	// 剩余长度
	remaining := append(varHeader, payload...)
	encodedLen := encodeRemainingLength(len(remaining))

	// 完整报文
	pkt := []byte{packetTypeConnect}
	pkt = append(pkt, encodedLen...)
	pkt = append(pkt, remaining...)
	return pkt, nil
}

// readConnack 读取并解析 CONNACK 报文，判断是否认证成功
func readConnack(conn net.Conn) error {
	// CONNACK 固定长度 4 字节：0x20 0x02 <sessionPresent> <returnCode>
	buf := make([]byte, 4)
	if _, err := readFull(conn, buf); err != nil {
		return err
	}

	if buf[0] != packetTypeConnack {
		return errors.New("unexpected packet type, not CONNACK")
	}
	if buf[1] != 0x02 {
		return errors.New("malformed CONNACK")
	}

	returnCode := buf[3]
	switch returnCode {
	case connackAccepted:
		return nil
	case connackBadUsernameOrPass:
		return errors.New("bad username or password")
	case connackNotAuthorized:
		return errors.New("not authorized")
	case 0x02:
		return errors.New("identifier rejected")
	case 0x03:
		return errors.New("server unavailable")
	default:
		return fmt.Errorf("connection refused, return code: 0x%02x", returnCode)
	}
}

// encodeMQTTStr 将字符串编码为 MQTT UTF-8 编码字符串（2字节长度前缀 + 内容）
func encodeMQTTStr(s string) []byte {
	b := []byte(s)
	length := make([]byte, 2)
	binary.BigEndian.PutUint16(length, uint16(len(b)))
	return append(length, b...)
}

// encodeRemainingLength 编码 MQTT 可变长度字段（最多4字节）
func encodeRemainingLength(n int) []byte {
	var encoded []byte
	for {
		digit := n % 128
		n /= 128
		if n > 0 {
			digit |= 0x80
		}
		encoded = append(encoded, byte(digit))
		if n == 0 {
			break
		}
	}
	return encoded
}

// readFull 确保读取指定长度的数据
func readFull(conn net.Conn, buf []byte) (int, error) {
	total := 0
	for total < len(buf) {
		n, err := conn.Read(buf[total:])
		total += n
		if err != nil {
			return total, err
		}
	}
	return total, nil
}
