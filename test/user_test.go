package test

import (
	"im-System/model"
	"net"
	"testing"
)

type mockConn struct {
	net.Conn
}

func (c *mockConn) RemoteAddr() net.Addr {
	// Return a dummy address
	return &net.TCPAddr{IP: net.IPv4(192, 168, 1, 1), Port: 5678}
}

func TestNewUser(t *testing.T) {
	conn := &mockConn{}
	if conn == nil {
		t.Errorf("Error in connection")
	}
	user := model.NewUser(conn, "testuser")
	if user == nil {
		t.Errorf("NewUser returned nil")
	}

	if user.Name() != "testuser" {
		t.Errorf("Expected user name 'testuser', got '%s'", user.Name())
	}

	if user.Address() == "" {
		t.Errorf("User address should not be empty")
	}
}

func TestSetName(t *testing.T) {
	conn := &mockConn{}
	user := model.NewUser(conn, "oldname")

	user.SetName("newname")

	if user.Name() != "newname" {
		t.Errorf("Expected user name 'newname', got '%s'", user.Name())
	}
}
