package main

import (
	"encoding/hex"
	"net"
	"time"
)

var req, _ = hex.DecodeString("474554202f20485454502f312e310d0a486f73743a203139322e3136382e322e310d0a557365722d4167656e743a20476f2d687474702d636c69656e742f312e310d0a4163636570742d456e636f64696e673a20677a69700d0a0d0a")

func scan(ip string) (error, string) {
	d := net.Dialer{Timeout: time.Millisecond * connectTimeout}
	conn, err := d.Dial("tcp", ip+":80")
	if err != nil {
		return err, ""
	}

	defer conn.Close()

	_, err = conn.Write(req)
	if err != nil {
		return err, ""
	}

	var resp []byte

	for {
		buf := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * readTimeout))
		read, err := conn.Read(buf)
		if err != nil {
			break
		} else {
			if uint64(len(buf[:read])+len(resp)) <= responseSize {
				resp = append(resp, buf[:read]...)
			} else {
				return nil, "Too big to record"
			}
		}
	}

	if string(resp) != "" {
		return nil, string(resp)
	} else {
		return &EmptyResponse{message: "Received an empty response from the server"}, ""
	}
}

func scanChunk(chunk []string, database *Db) {
	for i := uint64(0); i < chunkSize; i++ {
		if exitState {
			break
		}
		if !pauseState {
			addr := chunk[i]
			if addr != "" {
				err, resp := scan(addr)

				if err == nil {
					go database.Add(addr, resp)
				}

				incCommonVar(&scannedAddress, &countThreadsMu)
			}
		} else {
			for {
				pauseMu.Lock()
				ps := pauseState
				pauseMu.Unlock()
				if !ps {
					break
				}
			}
		}
	}
	subCommonVar(&countThreads, &countThreadsMu)
}
