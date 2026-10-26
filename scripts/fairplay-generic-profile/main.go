package main

import (
 "bytes"
 "compress/gzip"
 "crypto/sha256"
 "encoding/binary"
 "encoding/hex"
 "encoding/json"
 "flag"
 "fmt"
 "go/ast"
 "go/parser"
 "os"
 "strconv"
 "strings"
)

const recordBytes=16
var terminators=[]string{"compactBlockDirect","compactBlockConditional","compactBlockRegister","compactBlockCallRegister","compactBlockCallDirect","compactBlockReturn","compactBlockCompareZero","compactBlockTestBit","compactBlockStopDirect"}
var conditions=[]string{"compactConditionEQ","compactConditionNE","compactConditionHS","compactConditionLO","compactConditionMI","compactConditionPL","compactConditionVS","compactConditionVC","compactConditionHI","compactConditionLS","compactConditionGE","compactConditionLT","compactConditionGT","compactConditionLE"}

func namesMap(names []string) map[string]uint64 {m:=map[string]uint64{};for i,n:=range names {m[n]=uint64(i)};return m}
func literalFields(literal string, symbols map[string]uint64)(map[string]uint64,error){
 expression,err:=parser.ParseExpr("T{"+strings.Trim(strings.TrimSpace(literal),"{},")+"}")
 if err!=nil{return nil,err}; composite,ok:=expression.(*ast.CompositeLit);if !ok{return nil,fmt.Errorf("invalid semantic literal")}
 result:=map[string]uint64{}
 for _,entry:=range composite.Elts{
  field,ok:=entry.(*ast.KeyValueExpr);if !ok{return nil,fmt.Errorf("unnamed semantic operand")}
  key,ok:=field.Key.(*ast.Ident);if !ok{return nil,fmt.Errorf("invalid semantic operand")}
  if _,exists:=result[key.Name];exists{return nil,fmt.Errorf("duplicate operand")}
  var value uint64
  switch atom:=field.Value.(type){
  case *ast.Ident:var exists bool;value,exists=symbols[atom.Name];if !exists{return nil,fmt.Errorf("unknown semantic symbol %s",atom.Name)}
  case *ast.BasicLit:value,err=strconv.ParseUint(atom.Value,0,64);if err!=nil{return nil,err}
  default:return nil,fmt.Errorf("nonliteral semantic operand")
  };result[key.Name]=value
 };return result,nil
}
func normalRecord(literal string)([recordBytes]byte,error){
 var record [recordBytes]byte
 fields,err:=literalFields(literal,namesMap(opcodeNames));if err!=nil{return record,err}
 allowed:=map[string]bool{"op":true,"destination":true,"left":true,"right":true,"accumulator":true,"shift":true,"immediate":true}
 for key:=range fields{if !allowed[key]{return record,fmt.Errorf("unknown normal field")}}
 opcode,exists:=fields["op"];if !exists||opcode>=uint64(len(opcodeNames)){return record,fmt.Errorf("invalid opcode")}
 record[0]=1;binary.LittleEndian.PutUint16(record[1:3],uint16(opcode))
 for i,name:=range []string{"destination","left","right","accumulator","shift"}{if fields[name]>255{return record,fmt.Errorf("operand overflow")};record[3+i]=byte(fields[name])}
 if fields["immediate"]>65535{return record,fmt.Errorf("immediate overflow")};binary.LittleEndian.PutUint16(record[8:10],uint16(fields["immediate"]))
 return record,nil
}
func branchRecord(literal string,pc uint64)([recordBytes]byte,error){
 var record [recordBytes]byte
 symbols:=namesMap(terminators);for name,value:=range namesMap(conditions){symbols[name]=value}
 fields,err:=literalFields(literal,symbols);if err!=nil{return record,err}
 allowed:=map[string]bool{"terminator":true,"condition":true,"register":true,"width":true,"bit":true,"target":true,"alternate":true}
 for key:=range fields{if !allowed[key]{return record,fmt.Errorf("unknown branch field")}}
 term,exists:=fields["terminator"];if !exists||term>7{return record,fmt.Errorf("invalid native terminator")}
 if alternate,exists:=fields["alternate"];exists&&alternate!=pc+4{return record,fmt.Errorf("nonsequential alternate")}
 record[0]=2;record[1]=byte(term)
 for i,name:=range []string{"condition","register","width","bit"}{if fields[name]>255{return record,fmt.Errorf("branch operand overflow")};record[4+i]=byte(fields[name])}
 binary.LittleEndian.PutUint64(record[8:],fields["target"]);return record,nil
}
func main(){
 binaryPath:=flag.String("binary","","exact thin daemon")
 expectedHash:=flag.String("binary-sha256","","required source hash")
 output:=flag.String("output","","gzip semantic profile")
 reportPath:=flag.String("report","","coverage census")
 flag.Parse()
 if *binaryPath==""||*output==""||*reportPath==""||len(*expectedHash)!=64{panic("binary, binary-sha256, output and report required")}
 data,err:=os.ReadFile(*binaryPath);if err!=nil{panic(err)}
 digest:=sha256.Sum256(data);if hex.EncodeToString(digest[:])!=*expectedHash{panic("binary hash mismatch")}
 text,base,err:=openMachoARM64eText(*binaryPath);if err!=nil{panic(err)}
 if len(text)==0||len(text)%4!=0||base%4!=0||len(text)>32<<20{panic("invalid executable section")}
 var profile bytes.Buffer;profile.WriteString("CRIR1\x00\x00\x00")
 binary.Write(&profile,binary.LittleEndian,uint32(1));binary.Write(&profile,binary.LittleEndian,uint32(recordBytes))
 profile.Write(digest[:]);abi,err:=hex.DecodeString(kernelABI);if err!=nil{panic(err)};profile.Write(abi)
 binary.Write(&profile,binary.LittleEndian,base);binary.Write(&profile,binary.LittleEndian,uint32(len(text)/4));binary.Write(&profile,binary.LittleEndian,uint32(0))
 supported,branches:=0,0;missing:=map[string]int{};used:=map[string]int{}
 for offset:=0;offset<len(text);offset+=4{
  pc:=base+uint64(offset);assembly:=machoInstructionText(text[offset:offset+4],pc)
  terminal,terminalErr:=compactTerminator(annotatedInstruction{pc:pc,disassembly:assembly})
  var record [recordBytes]byte
  if terminalErr==nil{record,err=branchRecord(terminal,pc);if err==nil{branches++}}else{
   var literal string;var stop bool
   literal,stop,err=compactInstruction(assembly)
   if err==nil&&!stop{record,err=normalRecord(literal)}else if stop{err=fmt.Errorf("unexpected semantic stop")}
  }
  if err!=nil{record=[recordBytes]byte{};missing[strings.Fields(assembly)[0]]++}else{supported++;if record[0]==1{used[opcodeNames[binary.LittleEndian.Uint16(record[1:3])]]++}}
  profile.Write(record[:])
 }
 checksum:=sha256.Sum256(profile.Bytes());profile.Write(checksum[:])
 file,err:=os.Create(*output);if err!=nil{panic(err)};zip:=gzip.NewWriter(file)
 if _,err=zip.Write(profile.Bytes());err!=nil{panic(err)};if err=zip.Close();err!=nil{panic(err)};if err=file.Close();err!=nil{panic(err)}
 // Full census, including every unsupported slot: no skipped addresses or implied coverage.
 report:=map[string]any{"schema":"crumb-generic-instruction-profile-v1","binary_sha256":*expectedHash,"kernel_abi":kernelABI,"text_base":base,"instruction_slots":len(text)/4,"supported_slots":supported,"branch_slots":branches,"unsupported_slots":len(text)/4-supported,"unsupported_mnemonics":missing,"used_semantic_opcodes":used,"profile_payload_sha256":hex.EncodeToString(checksum[:]),"jit_required":false,"account_state_used":false,"decryption_validated":false}
 encoded,err:=json.MarshalIndent(report,"","  ");if err!=nil{panic(err)};if err=os.WriteFile(*reportPath,append(encoded,'\n'),0644);err!=nil{panic(err)}
 fmt.Printf("Enumerated %d instruction slots; supported %d; unsupported %d. No address allowlist.\n",len(text)/4,supported,len(text)/4-supported)
}
