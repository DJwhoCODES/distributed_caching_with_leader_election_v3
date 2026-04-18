package protocol

const Version uint8 = 1

type Command uint8

const (
	CmdGet Command = iota + 1
	CmdSet
	CmdDelete
	CmdMGet
	CmdMSet
	CmdPing

	CmdSync
	CmdVoteRequest
	CmdVoteResponse
	CmdHeartbeat
)

type Status uint8

const (
	StatusOK Status = iota
	StatusNotFound
	StatusError
	StatusRedirect
)
