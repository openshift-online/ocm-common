package vpc_client

// Export internal vars for testing only.
// This file is compiled exclusively as part of the test binary.

var (
	SSHMaxAttempts   = &sshMaxAttempts
	SSHRetryInterval = &sshRetryInterval
	SSHDialTimeout   = &sshDialTimeout
)
