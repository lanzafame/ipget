// Code generated by protoc-gen-gogo.
// source: dht.proto
// DO NOT EDIT!

/*
Package dht_pb is a generated protocol buffer package.

It is generated from these files:
	dht.proto

It has these top-level messages:
	Message
	Record
*/
package dht_pb

import proto "github.com/noffle/ipget/Godeps/_workspace/src/github.com/gogo/protobuf/proto"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = math.Inf

type Message_MessageType int32

const (
	Message_PUT_VALUE     Message_MessageType = 0
	Message_GET_VALUE     Message_MessageType = 1
	Message_ADD_PROVIDER  Message_MessageType = 2
	Message_GET_PROVIDERS Message_MessageType = 3
	Message_FIND_NODE     Message_MessageType = 4
	Message_PING          Message_MessageType = 5
)

var Message_MessageType_name = map[int32]string{
	0: "PUT_VALUE",
	1: "GET_VALUE",
	2: "ADD_PROVIDER",
	3: "GET_PROVIDERS",
	4: "FIND_NODE",
	5: "PING",
}
var Message_MessageType_value = map[string]int32{
	"PUT_VALUE":     0,
	"GET_VALUE":     1,
	"ADD_PROVIDER":  2,
	"GET_PROVIDERS": 3,
	"FIND_NODE":     4,
	"PING":          5,
}

func (x Message_MessageType) Enum() *Message_MessageType {
	p := new(Message_MessageType)
	*p = x
	return p
}
func (x Message_MessageType) String() string {
	return proto.EnumName(Message_MessageType_name, int32(x))
}
func (x *Message_MessageType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Message_MessageType_value, data, "Message_MessageType")
	if err != nil {
		return err
	}
	*x = Message_MessageType(value)
	return nil
}

type Message_ConnectionType int32

const (
	// sender does not have a connection to peer, and no extra information (default)
	Message_NOT_CONNECTED Message_ConnectionType = 0
	// sender has a live connection to peer
	Message_CONNECTED Message_ConnectionType = 1
	// sender recently connected to peer
	Message_CAN_CONNECT Message_ConnectionType = 2
	// sender recently tried to connect to peer repeatedly but failed to connect
	// ("try" here is loose, but this should signal "made strong effort, failed")
	Message_CANNOT_CONNECT Message_ConnectionType = 3
)

var Message_ConnectionType_name = map[int32]string{
	0: "NOT_CONNECTED",
	1: "CONNECTED",
	2: "CAN_CONNECT",
	3: "CANNOT_CONNECT",
}
var Message_ConnectionType_value = map[string]int32{
	"NOT_CONNECTED":  0,
	"CONNECTED":      1,
	"CAN_CONNECT":    2,
	"CANNOT_CONNECT": 3,
}

func (x Message_ConnectionType) Enum() *Message_ConnectionType {
	p := new(Message_ConnectionType)
	*p = x
	return p
}
func (x Message_ConnectionType) String() string {
	return proto.EnumName(Message_ConnectionType_name, int32(x))
}
func (x *Message_ConnectionType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Message_ConnectionType_value, data, "Message_ConnectionType")
	if err != nil {
		return err
	}
	*x = Message_ConnectionType(value)
	return nil
}

type Message struct {
	// defines what type of message it is.
	Type *Message_MessageType `protobuf:"varint,1,opt,name=type,enum=dht.pb.Message_MessageType" json:"type,omitempty"`
	// defines what coral cluster level this query/response belongs to.
	ClusterLevelRaw *int32 `protobuf:"varint,10,opt,name=clusterLevelRaw" json:"clusterLevelRaw,omitempty"`
	// Used to specify the key associated with this message.
	// PUT_VALUE, GET_VALUE, ADD_PROVIDER, GET_PROVIDERS
	Key *string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	// Used to return a value
	// PUT_VALUE, GET_VALUE
	Record *Record `protobuf:"bytes,3,opt,name=record" json:"record,omitempty"`
	// Used to return peers closer to a key in a query
	// GET_VALUE, GET_PROVIDERS, FIND_NODE
	CloserPeers []*Message_Peer `protobuf:"bytes,8,rep,name=closerPeers" json:"closerPeers,omitempty"`
	// Used to return Providers
	// GET_VALUE, ADD_PROVIDER, GET_PROVIDERS
	ProviderPeers    []*Message_Peer `protobuf:"bytes,9,rep,name=providerPeers" json:"providerPeers,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *Message) Reset()         { *m = Message{} }
func (m *Message) String() string { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()    {}

func (m *Message) GetType() Message_MessageType {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Message_PUT_VALUE
}

func (m *Message) GetClusterLevelRaw() int32 {
	if m != nil && m.ClusterLevelRaw != nil {
		return *m.ClusterLevelRaw
	}
	return 0
}

func (m *Message) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Message) GetRecord() *Record {
	if m != nil {
		return m.Record
	}
	return nil
}

func (m *Message) GetCloserPeers() []*Message_Peer {
	if m != nil {
		return m.CloserPeers
	}
	return nil
}

func (m *Message) GetProviderPeers() []*Message_Peer {
	if m != nil {
		return m.ProviderPeers
	}
	return nil
}

type Message_Peer struct {
	// ID of a given peer.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// multiaddrs for a given peer
	Addrs [][]byte `protobuf:"bytes,2,rep,name=addrs" json:"addrs,omitempty"`
	// used to signal the sender's connection capabilities to the peer
	Connection       *Message_ConnectionType `protobuf:"varint,3,opt,name=connection,enum=dht.pb.Message_ConnectionType" json:"connection,omitempty"`
	XXX_unrecognized []byte                  `json:"-"`
}

func (m *Message_Peer) Reset()         { *m = Message_Peer{} }
func (m *Message_Peer) String() string { return proto.CompactTextString(m) }
func (*Message_Peer) ProtoMessage()    {}

func (m *Message_Peer) GetId() string {
	if m != nil && m.Id != nil {
		return *m.Id
	}
	return ""
}

func (m *Message_Peer) GetAddrs() [][]byte {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func (m *Message_Peer) GetConnection() Message_ConnectionType {
	if m != nil && m.Connection != nil {
		return *m.Connection
	}
	return Message_NOT_CONNECTED
}

// Record represents a dht record that contains a value
// for a key value pair
type Record struct {
	// The key that references this record
	Key *string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	// The actual value this record is storing
	Value []byte `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	// hash of the authors public key
	Author *string `protobuf:"bytes,3,opt,name=author" json:"author,omitempty"`
	// A PKI signature for the key+value+author
	Signature []byte `protobuf:"bytes,4,opt,name=signature" json:"signature,omitempty"`
	// Time the record was received, set by receiver
	TimeReceived     *string `protobuf:"bytes,5,opt,name=timeReceived" json:"timeReceived,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Record) Reset()         { *m = Record{} }
func (m *Record) String() string { return proto.CompactTextString(m) }
func (*Record) ProtoMessage()    {}

func (m *Record) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *Record) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Record) GetAuthor() string {
	if m != nil && m.Author != nil {
		return *m.Author
	}
	return ""
}

func (m *Record) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Record) GetTimeReceived() string {
	if m != nil && m.TimeReceived != nil {
		return *m.TimeReceived
	}
	return ""
}

func init() {
	proto.RegisterEnum("dht.pb.Message_MessageType", Message_MessageType_name, Message_MessageType_value)
	proto.RegisterEnum("dht.pb.Message_ConnectionType", Message_ConnectionType_name, Message_ConnectionType_value)
}
