package pdubase

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/golang/protobuf/proto"
	base "github.com/teamtalk-remix/test-client-go/proto/IM_BaseDefine"
)

const PduHeaderSize = 16

type PduHeader struct {
	length    uint32
	Version   uint16
	Flag      uint16
	ServiceId base.ServiceID
	CommandId base.CommandID
	SeqNum    uint32
	Reversed  uint32
}

type CImPdu struct {
	PduHeader
	buf bytes.Buffer
}

func (p CImPdu) Length() uint32 {
	return p.PduHeader.length
}

func (p CImPdu) Version() uint16 {
	return p.PduHeader.Version
}

func (p CImPdu) Flag() uint16 {
	return p.PduHeader.Flag
}

func (p CImPdu) ServiceId() base.ServiceID {
	return p.PduHeader.ServiceId
}

func (p CImPdu) CommandId() base.CommandID {
	return p.PduHeader.CommandId
}

func (p CImPdu) SeqNum() uint32 {
	return p.PduHeader.SeqNum
}

func (p CImPdu) Reversed() uint32 {
	return p.PduHeader.Reversed
}

func (p *CImPdu) WriteBuffer(b []byte) {
	p.buf.Write(b)
}

func (p *CImPdu) Buffer() *bytes.Buffer {
	return &p.buf
}

func (p *CImPdu) GetHeader(b []byte) {
	//b := make([]byte, 16)
	//length
	p.SetLength(binary.BigEndian.Uint32(b[0:]))
	//version
	binary.BigEndian.Uint16(b[4:])
	//flag
	binary.BigEndian.Uint16(b[6:])
	//service_id
	sid := int32(binary.BigEndian.Uint16(b[8:]))
	p.SetServiceId(base.ServiceID(sid))
	//command_id
	cid := int32(binary.BigEndian.Uint16(b[10:]))
	p.SetCommandId(base.CommandID(cid))
	//seq_num
	binary.BigEndian.Uint16(b[12:])
	//reversed
	binary.BigEndian.Uint16(b[14:])
}

func (p *CImPdu) GetPBMsg(buf []byte, pb proto.Message) {
	err := proto.Unmarshal(buf, pb)
	if err != nil {
		fmt.Println(err)
	}
}

func (p *CImPdu) SetPB(pb proto.Message) {
	var seq uint16
	seq = 0
	//pdu.SetHeader
	b := make([]byte, 16)
	pbData, err := proto.Marshal(pb)
	if err != nil {
		log.Fatal("marshaling error: ", err)
	}
	pbLen := len(pbData) + 16
	binary.BigEndian.PutUint32(b[0:], uint32(pbLen))
	//version == 1
	binary.BigEndian.PutUint16(b[4:], 0x0001)
	//flag == 0
	binary.BigEndian.PutUint16(b[6:], 0x0000)
	//service_id == 1
	binary.BigEndian.PutUint16(b[8:], uint16(p.ServiceId()))
	//command_id == 259
	binary.BigEndian.PutUint16(b[10:], uint16(p.CommandId()))
	//seq_num
	binary.BigEndian.PutUint16(b[12:], seq)
	//reversed == 0
	binary.BigEndian.PutUint16(b[14:], 0x0000)
	p.buf.Write(b)
	p.buf.Write(pbData)
}

func (p *CImPdu) SetLength(length uint32) {
	p.PduHeader.length = length
}

func (p *CImPdu) SetVersion(version uint16) {
	p.PduHeader.Version = version
}

func (p *CImPdu) SetFlag(flag uint16) {
	p.PduHeader.Flag = flag
}

func (p *CImPdu) SetServiceId(serviceId base.ServiceID) {
	p.PduHeader.ServiceId = serviceId
}

func (p *CImPdu) SetCommandId(commandId base.CommandID) {
	p.PduHeader.CommandId = commandId
}

func (p *CImPdu) SetSeqNum(seqNum uint32) {
	p.PduHeader.SeqNum = seqNum
}

func (p *CImPdu) SetReversed(reversed uint32) {
	p.PduHeader.Reversed = reversed
}
