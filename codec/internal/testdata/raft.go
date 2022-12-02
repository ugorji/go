package testdata

import "time"

// copy of all the raft data types

type ProtocolVersion int
type LogType uint8
type SnapshotVersion int
type ServerID string
type ServerAddress string
type ServerSuffrage int
type Log struct {
	Index      uint64
	Term       uint64
	Type       LogType
	Data       []byte
	Extensions []byte
	AppendedAt time.Time
}
type RPCHeader struct {
	ProtocolVersion ProtocolVersion
	ID              []byte
	Addr            []byte
}
type AppendEntriesRequest struct {
	RPCHeader
	Term              uint64
	Leader            []byte
	PrevLogEntry      uint64
	PrevLogTerm       uint64
	Entries           []*Log
	LeaderCommitIndex uint64
}
type AppendEntriesResponse struct {
	RPCHeader
	Term           uint64
	LastLog        uint64
	Success        bool
	NoRetryBackoff bool
}
type InstallSnapshotRequest struct {
	RPCHeader
	SnapshotVersion    SnapshotVersion
	Term               uint64
	Leader             []byte
	LastLogIndex       uint64
	LastLogTerm        uint64
	Peers              []byte
	Configuration      []byte
	ConfigurationIndex uint64
	Size               int64
}
type InstallSnapshotResponse struct {
	RPCHeader
	Term    uint64
	Success bool
}
type RequestVoteRequest struct {
	RPCHeader
	Term               uint64
	Candidate          []byte
	LastLogIndex       uint64
	LastLogTerm        uint64
	LeadershipTransfer bool
}
type RequestVoteResponse struct {
	RPCHeader
	Term    uint64
	Peers   []byte
	Granted bool
}
type Server struct {
	Suffrage ServerSuffrage
	ID       ServerID
	Address  ServerAddress
}
type Configuration struct {
	Servers []Server
}

type SerializeTest struct {
	Name            string
	EncodedBytesHex string
	ExpectedData    interface{}
}
