// Generic interpreter for a versioned, whole-executable-section semantic profile.
// The profile is data. No native code generation or executable mapping occurs.
package rawsource

import (
 "bytes"
 "compress/gzip"
 "crypto/sha256"
 "encoding/binary"
 "encoding/hex"
 "fmt"
 "io"
)

const genericRecordBytes=16
const genericMaxProfileBytes=128<<20

type genericInstructionProfile struct { base uint64; count uint32; records []byte }
var genericImage *genericInstructionProfile

type genericExecutionError struct { code string }
func (e *genericExecutionError) Error() string { return e.code }

func loadGenericInstructionProfile(compressed []byte, expectedBinary [32]byte, expectedABI string)(*genericInstructionProfile,error){
 stream,err:=gzip.NewReader(bytes.NewReader(compressed));if err!=nil{return nil,fmt.Errorf("generic profile compression invalid")}
 data,err:=io.ReadAll(io.LimitReader(stream,genericMaxProfileBytes+1));closeErr:=stream.Close()
 if err!=nil||closeErr!=nil||len(data)>genericMaxProfileBytes||len(data)<128{return nil,fmt.Errorf("generic profile size invalid")}
 payload,checksum:=data[:len(data)-32],data[len(data)-32:];sum:=sha256.Sum256(payload)
 if !bytes.Equal(checksum,sum[:]){return nil,fmt.Errorf("generic profile checksum mismatch")}
 if string(payload[:8])!="CRIR1\x00\x00\x00"||binary.LittleEndian.Uint32(payload[8:])!=1||binary.LittleEndian.Uint32(payload[12:])!=genericRecordBytes{return nil,fmt.Errorf("generic profile schema mismatch")}
 if !bytes.Equal(payload[16:48],expectedBinary[:])||hex.EncodeToString(payload[48:80])!=expectedABI{return nil,fmt.Errorf("generic profile source/ABI mismatch")}
 base:=binary.LittleEndian.Uint64(payload[80:]);count:=binary.LittleEndian.Uint32(payload[88:])
 if base%4!=0||count==0||uint64(count)*genericRecordBytes+96!=uint64(len(payload))||base+uint64(count)*4<base||binary.LittleEndian.Uint32(payload[92:])!=0{return nil,fmt.Errorf("generic profile geometry invalid")}
 image:=&genericInstructionProfile{base:base,count:count,records:payload[96:]}
 for i:=uint32(0);i<count;i++{
  r:=image.records[int(i)*genericRecordBytes:int(i+1)*genericRecordBytes]
  switch r[0]{
  case 0:if !bytes.Equal(r,make([]byte,genericRecordBytes)){return nil,fmt.Errorf("invalid unsupported slot")}
  case 1:if compactOpcode(binary.LittleEndian.Uint16(r[1:]))>=compactOpcodeCount{return nil,fmt.Errorf("generic opcode invalid")}
  case 2:
   term:=compactBlockTerminator(r[1]);if term>compactBlockTestBit||r[4]>byte(compactConditionLE)||r[5]>30{return nil,fmt.Errorf("generic branch operands invalid")}
   if term==compactBlockCompareZero&&r[6]!=32&&r[6]!=64{return nil,fmt.Errorf("generic branch width invalid")}
   if term==compactBlockTestBit&&r[7]>=64{return nil,fmt.Errorf("generic branch bit invalid")}
  default:return nil,fmt.Errorf("generic record kind invalid")
  }
 }
 return image,nil
}
func (p *genericInstructionProfile) contains(pc uint64)bool{return p!=nil&&pc>=p.base&&(pc-p.base)%4==0&&(pc-p.base)/4<uint64(p.count)}
func (p *genericInstructionProfile) record(pc uint64)[]byte{index:=(pc-p.base)/4;return p.records[index*genericRecordBytes:(index+1)*genericRecordBytes]}

// A chunk ends on a real control transfer, stopPC or bounded sequential budget.
// It can start at ANY aligned instruction slot, including an indirect landing
// in the middle of a region that static basic-block discovery did not find.
func runGenericInstructionBlock(s *state,runtime *rawSourceRuntime,start uint64)(uint64,bool,error){
 pc:=start
 if genericImage==nil{return pc,false,&genericExecutionError{"generic_profile_missing"}}
 runtime.genericBlocks++
 for count:=0;count<4096;count++{
  if runtime.stopPC!=0&&pc==runtime.stopPC{return pc,true,nil}
  if !genericImage.contains(pc){runtime.recordUnknownPC(pc);return pc,false,&genericExecutionError{"generic_pc_outside_profile"}}
  record:=genericImage.record(pc)
  if record[0]==0{runtime.recordUnknownPC(pc);return pc,false,&genericExecutionError{"generic_unsupported_instruction"}}
  runtime.genericInstructions++
  switch record[0]{
  case 1:
   instruction:=compactInstruction{op:compactOpcode(binary.LittleEndian.Uint16(record[1:3])),destination:record[3],left:record[4],right:record[5],accumulator:record[6],shift:record[7],immediate:binary.LittleEndian.Uint16(record[8:10])}
   if err:=runCompactKernel(s,compactKernel{start:pc,instructions:[]compactInstruction{instruction}});err!=nil{return pc,false,err}
   pc+=4
  case 2:
   block:=compactBlock{start:pc,terminator:compactBlockTerminator(record[1]),condition:compactCondition(record[4]),register:record[5],width:record[6],bit:record[7],target:binary.LittleEndian.Uint64(record[8:]),alternate:pc+4}
   return runCompactBlock(s,runtime,block)
  }
 }
 return pc,false,nil
}
