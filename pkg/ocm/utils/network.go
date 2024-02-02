package utils

import "net"

func IsEmptyCIDR(cidr net.IPNet) bool {
	return cidr.String() == "<nil>"
}
