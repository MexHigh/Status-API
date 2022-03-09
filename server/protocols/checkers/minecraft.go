package checkers

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"status-api/protocols"
	"status-api/structs"
)

// Minecraft -
type Minecraft struct{}

// Check -
func (Minecraft) Check(name string, c *structs.ServiceConfig) (structs.CheckResult, error) {

	var hostPort string
	if hp, ok := c.ProtocolConfig["server_address"].(string); ok {
		hostPort = hp
	} else {
		hostPort = c.FriendlyURL
		if !strings.Contains(hostPort, ":") {
			hostPort = hostPort + ":25565"
		}
	}

	var res = structs.CheckResult{
		URL: c.FriendlyURL,
	}

	pong, err := mineping(hostPort)
	if err != nil {
		res.Status = structs.Down
		if e := err.Error(); strings.Contains(e, "i/o timeout") {
			res.Reason = "I/O timeout"
		} else if strings.Contains(e, "connection refused") {
			res.Reason = "connection refused"
		} else if strings.Contains(e, "no route to host") {
			res.Reason = "no route to host"
		} else {
			res.Reason = e
		}
		return res, nil
	}

	res.Status = structs.Up
	res.Misc = map[string]string{
		"version":        pong.Version.Name,
		"players_online": fmt.Sprintf("%d/%d", pong.Players.Online, pong.Players.Max),
	}

	return res, nil

}

// Register checker
func init() {
	protocols.Register("minecraft", Minecraft{})
}

//
// This following lines of go were written by Andrew Tian and Syfaro.
// Taken from: https://github.com/Syfaro/minepong
//

const (
	protocolVersion = 573
)

var (
	connectionTimeout = 2 * time.Second
)

type pong struct {
	Version struct {
		Name     string
		Protocol int
	} `json:"version"`
	Players struct {
		Max    int `json:"max"`
		Online int `json:"online"`
		Sample []map[string]string
	} `json:"players"`
	Description  interface{} `json:"description"`
	FavIcon      string      `json:"favicon"`
	ResolvedHost string      `json:"resolved_host"`
}

func resolveSRV(addr string) (host string, err error) {
	h, _, err := net.SplitHostPort(addr)
	if err != nil {
		h = addr
	}

	_, addrs, err := net.LookupSRV("minecraft", "tcp", h)
	if err != nil || len(addrs) == 0 {
		return host, errors.New("unable to find SRV record")
	}

	return net.JoinHostPort(addrs[0].Target, strconv.Itoa(int(addrs[0].Port))), nil
}

func mineping(host string) (*pong, error) {
	srvHost, err := resolveSRV(host)

	if err == nil {
		host = srvHost
	}

	conn, err := net.DialTimeout("tcp", host, connectionTimeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(connectionTimeout))
	conn.SetWriteDeadline(time.Now().Add(connectionTimeout))

	if err := sendHandshake(conn, host); err != nil {
		return nil, err
	}

	if err := sendStatusRequest(conn); err != nil {
		return nil, err
	}

	p, err := readPong(conn)
	if err != nil {
		return nil, err
	}

	p.ResolvedHost = host

	return p, nil
}

func makePacket(pl *bytes.Buffer) *bytes.Buffer {
	var buf bytes.Buffer
	// get payload length
	buf.Write(encodeVarint(uint64(len(pl.Bytes()))))

	// write payload
	buf.Write(pl.Bytes())

	return &buf
}

func sendHandshake(conn net.Conn, host string) error {
	pl := &bytes.Buffer{}

	// packet id
	pl.WriteByte(0x00)

	// protocol version
	pl.Write(encodeVarint(uint64(protocolVersion)))

	// server address
	host, port, err := net.SplitHostPort(host)
	if err != nil {
		return errors.New("cannot split host and port")
	}

	pl.Write(encodeVarint(uint64(len(host))))
	pl.WriteString(host)

	// server port
	iPort, err := strconv.Atoi(port)
	if err != nil {
		return errors.New("cannot convert port to int")
	}
	binary.Write(pl, binary.BigEndian, int16(iPort))

	// next state (status)
	pl.WriteByte(0x01)

	if _, err := makePacket(pl).WriteTo(conn); err != nil {
		return errors.New("cannot write handshake")
	}

	return nil
}

func sendStatusRequest(conn net.Conn) error {
	pl := &bytes.Buffer{}

	// send request zero
	pl.WriteByte(0x00)

	if _, err := makePacket(pl).WriteTo(conn); err != nil {
		return errors.New("cannot write send status request")
	}

	return nil
}

// https://code.google.com/p/goprotobuf/source/browse/proto/encode.go#83
func encodeVarint(x uint64) []byte {
	var buf [10]byte
	var n int
	for n = 0; x > 127; n++ {
		buf[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	buf[n] = uint8(x)
	n++
	return buf[0:n]
}

func readPong(rd io.Reader) (*pong, error) {
	r := bufio.NewReader(rd)
	nl, err := binary.ReadUvarint(r)
	if err != nil {
		return nil, errors.New("could not read length")
	}

	pl := make([]byte, nl)
	_, err = io.ReadFull(r, pl)
	if err != nil {
		return nil, errors.New("could not read length given by length header")
	}

	// packet id
	_, n := binary.Uvarint(pl)
	if n <= 0 {
		return nil, errors.New("could not read packet id")
	}

	// string varint
	_, n2 := binary.Uvarint(pl[n:])
	if n2 <= 0 {
		return nil, errors.New("could not read string varint")
	}

	var p pong
	if err := json.Unmarshal(pl[n+n2:], &p); err != nil {
		return nil, errors.New("could not read pong json")
	}

	return &p, nil
}
