package redis

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// 这是原本的扫描逻辑，没有设置空密码的情况
//
//	func Check(Host, Password string, Port int) error {
//		netloc := fmt.Sprintf("%s:%d", Host, Port)
//		conn, err := net.DialTimeout("tcp", netloc, 5*time.Second)
//		if err != nil {
//			return err
//		}
//		defer conn.Close()
//		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
//		if err != nil {
//			return err
//		}
//		_, err = conn.Write([]byte(fmt.Sprintf("auth %s\r\n", Password)))
//		time.Sleep(time.Millisecond * 500)
//		if err != nil {
//			return err
//		}
//		reply, err := readResponse(conn)
//		if err != nil {
//			return err
//		}
//		if strings.Contains(reply, "+OK") == false {
//			return errors.New("login failed")
//		}
//		return nil
//	}
func Check(Host, Password string, Port int) error {
	netloc := fmt.Sprintf("%s:%d", Host, Port)
	conn, err := net.DialTimeout("tcp", netloc, 5*time.Second)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 如果密码为空，跳过AUTH命令，直接发送PING
	if Password != "" {
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		if err != nil {
			return err
		}

		// 发送AUTH命令
		_, err = conn.Write([]byte(fmt.Sprintf("AUTH %s\r\n", Password)))
		if err != nil {
			return err
		}

		// 等待响应
		reply, err := readResponse(conn)
		if err != nil {
			return err
		}
		if strings.Contains(reply, "+OK") == false {
			return errors.New("login failed")
		}
		// 如果AUTH成功，接着发送PING
		if strings.Contains(reply, "+OK") {
			// Fall through to the PING check below
		} else {
			return errors.New("auth failed")
		}
	}
	// 对于有密码和无密码两种情况，都发送PING来确认是否登录成功
	err = conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		return err
	}

	_, err = conn.Write([]byte("PING\r\n"))
	if err != nil {
		return err
	}

	reply, err := readResponse(conn)
	if err != nil {
		return err
	}

	if strings.Contains(reply, "+PONG") == false {
		// 如果没有收到PONG，可能认证失败
		return errors.New("ping failed, may need authentication")
	}

	return nil
}
func readResponse(conn net.Conn) (r string, err error) {
	buf := make([]byte, 4096)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			break
		}
		r += string(buf[0:count])
		if count < 4096 {
			break
		}
	}
	return r, err
}
