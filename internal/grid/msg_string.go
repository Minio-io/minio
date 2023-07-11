// Code generated by "stringer -type=Op -output=msg_string.go -trimprefix=Op msg.go"; DO NOT EDIT.

package grid

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OpConnect-1]
	_ = x[OpConnectResponse-2]
	_ = x[OpRequest-3]
	_ = x[OpResponse-4]
	_ = x[OpConnectMux-5]
	_ = x[OpMuxConnectError-6]
	_ = x[OpDisconnectMux-7]
	_ = x[OpMuxClientMsg-8]
	_ = x[OpMuxServerMsg-9]
	_ = x[OpMuxServerErr-10]
	_ = x[OpUnblockMux-11]
	_ = x[OpAckMux-12]
	_ = x[OpDisconnect-13]
}

const _Op_name = "ConnectConnectResponseRequestResponseConnectMuxMuxConnectErrorDisconnectMuxMuxClientMsgMuxServerMsgMuxServerErrUnblockMuxAckMuxDisconnect"

var _Op_index = [...]uint8{0, 7, 22, 29, 37, 47, 62, 75, 87, 99, 111, 121, 127, 137}

func (i Op) String() string {
	i -= 1
	if i >= Op(len(_Op_index)-1) {
		return "Op(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Op_name[_Op_index[i]:_Op_index[i+1]]
}
