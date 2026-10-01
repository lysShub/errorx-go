package neterr

func ConnectReset(err error) bool   { return connectReset(err) }
func ConnectRefused(err error) bool { return connectRefused(err) }
func ConnectAborted(err error) bool { return connectAborted(err) }
func NetworkUnreach(err error) bool { return networkUnreach(err) }
func BuffSize(err error) bool       { return buffSize(err) }
func AddrNotAvail(err error) bool   { return addrNotAvail(err) }
func NetTimeout(err error) bool     { return netTimeout(err) }
func AddrInuse(err error) bool      { return addrInuse(err) }
func WouldBlock(err error) bool     { return wouldBlock(err) }
