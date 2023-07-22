package util

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

func GetFormatTime(time time.Time) string {
	return time.Format("0601021504")
}

func GenerateCode() string {
	// date := GetFormatTime(time.Now())
	r := rand.Intn(10000)
	code := fmt.Sprintf("%04d", r)
	return code
}

func IsExistsInStrSlice(row string, rows []string) bool {
	for _, r := range rows {
		if r == row {
			return true
		}
	}
	return false
}

func IsPathExists(p string) bool {
	_, err := os.Stat(p)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func EnsurePath(path string) error {
	workDir, _ := os.Getwd()
	dest := filepath.Join(workDir, path)
	exists := IsPathExists(dest)
	if !exists {
		err := os.MkdirAll(dest, os.ModePerm)
		return err
	}
	return nil
}

func GetRandomUDPAddr() (addr string, err error) {
	var a *net.UDPAddr
	if a, err = net.ResolveUDPAddr("udp", "localhost:0"); err == nil {
		var l *net.UDPConn
		if l, err = net.ListenUDP("udp", a); err == nil {
			defer l.Close()
			return l.LocalAddr().String(), nil
		}
	}
	return
}
