package vpc_client

// Internal test package so we can override the unexported retry vars
// (sshMaxAttempts, sshRetryInterval, sshDialTimeout) without export_test.go.

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/ssh"
)

// writeTempEd25519Key generates an Ed25519 key (microseconds vs ~100ms for RSA)
// and writes it to a temp file in OpenSSH PEM format.
// The caller is responsible for removing the file.
func writeTempEd25519Key() (path string, err error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}
	block, err := ssh.MarshalPrivateKey(priv, "")
	if err != nil {
		return "", err
	}
	f, err := os.CreateTemp("", "ocm-test-key-*.pem")
	if err != nil {
		return "", err
	}
	if err := pem.Encode(f, block); err != nil {
		f.Close()
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return f.Name(), nil
}

// rejectingListener starts a TCP listener that accepts every connection and
// immediately closes it, simulating a host whose SSH daemon is not yet ready.
// It increments *count on each accepted connection.
func rejectingListener(count *int32) (net.Listener, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			atomic.AddInt32(count, 1)
			conn.Close()
		}
	}()
	return ln, nil
}

var _ = Describe("Exec_CMD SSH retry", func() {
	var (
		keyFile          string
		savedMaxAttempts int
		savedRetryIvl    time.Duration
		savedDialTimeout time.Duration
	)

	BeforeEach(func() {
		var err error
		keyFile, err = writeTempEd25519Key()
		Expect(err).NotTo(HaveOccurred())

		// Save originals and install fast values so tests complete in milliseconds.
		savedMaxAttempts = sshMaxAttempts
		savedRetryIvl = sshRetryInterval
		savedDialTimeout = sshDialTimeout

		sshRetryInterval = 0
		sshDialTimeout = 2 * time.Second
	})

	AfterEach(func() {
		os.Remove(keyFile)
		sshMaxAttempts = savedMaxAttempts
		sshRetryInterval = savedRetryIvl
		sshDialTimeout = savedDialTimeout
	})

	It("retries exactly sshMaxAttempts times when the server rejects every connection", func() {
		sshMaxAttempts = 3
		var dialCount int32

		ln, err := rejectingListener(&dialCount)
		Expect(err).NotTo(HaveOccurred())
		defer ln.Close()

		_, err = Exec_CMD("testuser", keyFile, ln.Addr().String(), "echo hello")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("3 attempts"))

		Eventually(func() int32 { return atomic.LoadInt32(&dialCount) }, "2s", "50ms").
			Should(BeEquivalentTo(3))
	})

	It("returns immediately when no server is listening (connection refused)", func() {
		sshMaxAttempts = 2

		ln, err := net.Listen("tcp", "127.0.0.1:0")
		Expect(err).NotTo(HaveOccurred())
		addr := ln.Addr().String()
		ln.Close()

		start := time.Now()
		_, err = Exec_CMD("testuser", keyFile, addr, "echo hello")
		Expect(err).To(HaveOccurred())
		Expect(time.Since(start)).To(BeNumerically("<", 10*time.Second))
	})

	It("wraps the original dial failure in the returned error", func() {
		sshMaxAttempts = 1

		ln, err := rejectingListener(new(int32))
		Expect(err).NotTo(HaveOccurred())
		defer ln.Close()

		_, err = Exec_CMD("testuser", keyFile, ln.Addr().String(), "echo hello")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("failed to dial"))
	})
})

var _ = Describe("isSSHExitError", func() {
	It("returns false for a plain connection error", func() {
		Expect(isSSHExitError(fmt.Errorf("connection refused"))).To(BeFalse())
	})

	It("returns true for an ssh.ExitError", func() {
		exitErr := &ssh.ExitError{Waitmsg: ssh.Waitmsg{}}
		Expect(isSSHExitError(exitErr)).To(BeTrue())
	})

	It("returns true for a wrapped ssh.ExitError", func() {
		exitErr := &ssh.ExitError{Waitmsg: ssh.Waitmsg{}}
		wrapped := fmt.Errorf("failed to run command: %w", exitErr)
		Expect(isSSHExitError(wrapped)).To(BeTrue())
	})
})
