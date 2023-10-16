// Copyright (c) 2015-2023 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package grid

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/minio/minio/internal/logger/target/testlogger"
)

func TestDisconnect(t *testing.T) {
	defer testlogger.T.SetLogTB(t)()
	hosts, listeners := getHosts(2)
	dialer := &net.Dialer{
		Timeout: 1 * time.Second,
	}
	errFatal := func(err error) {
		t.Helper()
		if err != nil {
			t.Fatal(err)
		}
	}
	wrapServer := func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			t.Logf("Got a %s request for: %v", r.Method, r.URL)
			handler.ServeHTTP(w, r)
		})
	}
	connReady := make(chan struct{})
	// We fake a local and remote server.
	localHost := hosts[0]
	remoteHost := hosts[1]
	local, err := NewManager(context.Background(), ManagerOptions{
		Dialer:       dialer.DialContext,
		Local:        localHost,
		Hosts:        hosts,
		AddAuth:      func(aud string) string { return aud },
		AuthRequest:  dummyRequestValidate,
		BlockConnect: connReady,
	})
	errFatal(err)
	defer local.debugMsg(debugShutdown)

	// 1: Echo
	errFatal(local.RegisterSingle(handlerTest, func(payload []byte) ([]byte, *RemoteErr) {
		t.Log("1: server payload: ", len(payload), "bytes.")
		return append([]byte{}, payload...), nil
	}))
	// 2: Return as error
	errFatal(local.RegisterSingle(handlerTest2, func(payload []byte) ([]byte, *RemoteErr) {
		t.Log("2: server payload: ", len(payload), "bytes.")
		err := RemoteErr(payload)
		return nil, &err
	}))

	remote, err := NewManager(context.Background(), ManagerOptions{
		Dialer:       dialer.DialContext,
		Local:        remoteHost,
		Hosts:        hosts,
		AddAuth:      func(aud string) string { return aud },
		AuthRequest:  dummyRequestValidate,
		BlockConnect: connReady,
	})
	errFatal(err)
	defer remote.debugMsg(debugShutdown)
	defer local.debugMsg(debugWaitForExit)
	defer remote.debugMsg(debugWaitForExit)

	localServer := startServer(t, listeners[0], wrapServer(local.Handler()))
	defer localServer.Close()
	remoteServer := startServer(t, listeners[1], wrapServer(remote.Handler()))
	defer remoteServer.Close()
	close(connReady)

	cleanReqs := make(chan struct{})
	gotCall := make(chan struct{})
	defer close(cleanReqs)
	// 1: Block forever
	h1 := func(payload []byte) ([]byte, *RemoteErr) {
		gotCall <- struct{}{}
		<-cleanReqs
		return nil, nil
	}
	// 2: Also block, but with streaming.
	h2 := StreamHandler{
		Handle: func(ctx context.Context, payload []byte, request <-chan []byte, resp chan<- []byte) *RemoteErr {
			gotCall <- struct{}{}
			select {
			case <-ctx.Done():
				gotCall <- struct{}{}
			case <-cleanReqs:
				panic("should not be called")
			}
			return nil
		},
		OutCapacity: 1,
		InCapacity:  1,
	}
	errFatal(remote.RegisterSingle(handlerTest, h1))
	errFatal(remote.RegisterStreamingHandler(handlerTest2, h2))
	errFatal(local.RegisterSingle(handlerTest, h1))
	errFatal(local.RegisterStreamingHandler(handlerTest2, h2))

	// local to remote
	remoteConn := local.Connection(remoteHost)
	errFatal(remoteConn.WaitForConnect(context.Background()))
	const testPayload = "Hello Grid World!"

	gotResp := make(chan struct{})
	go func() {
		start := time.Now()
		t.Log("Roundtrip: sending request")
		resp, err := remoteConn.Request(context.Background(), handlerTest, []byte(testPayload))
		t.Log("Roundtrip:", time.Since(start), resp, err)
		gotResp <- struct{}{}
	}()
	<-gotCall
	remote.debugMsg(debugKillInbound)
	local.debugMsg(debugKillInbound)
	<-gotResp

	// Must reconnect
	errFatal(remoteConn.WaitForConnect(context.Background()))

	stream, err := remoteConn.NewStream(context.Background(), handlerTest2, []byte(testPayload))
	errFatal(err)
	go func() {
		for resp := range stream.responses {
			t.Log("Resp:", resp, err)
		}
		gotResp <- struct{}{}
	}()

	<-gotCall
	remote.debugMsg(debugKillOutbound)
	local.debugMsg(debugKillOutbound)
	<-gotResp
	// Killing should cancel the context on the request.
	<-gotCall
}

func dummyRequestValidate(r *http.Request) error {
	return nil
}

func TestShouldConnect(t *testing.T) {
	var c Connection
	var cReverse Connection
	hosts := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	for x := range hosts {
		should := 0
		for y := range hosts {
			if x == y {
				continue
			}
			c.Local = hosts[x]
			c.Remote = hosts[y]
			cReverse.Local = hosts[y]
			cReverse.Remote = hosts[x]
			if c.shouldConnect() == cReverse.shouldConnect() {
				t.Errorf("shouldConnect(%q, %q) != shouldConnect(%q, %q)", hosts[x], hosts[y], hosts[y], hosts[x])
			}
			if c.shouldConnect() {
				should++
			}
		}
		if should < 10 {
			t.Errorf("host %q only connects to %d hosts", hosts[x], should)
		}
		t.Logf("host %q should connect to %d hosts", hosts[x], should)
	}
}
