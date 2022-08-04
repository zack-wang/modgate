package main

import (
	"math"
	//"encoding/json"
	"encoding/binary"
	"log"
	"bytes"
)

// from BYTE to NUMBER
// 
func F64BE(bytes []byte) float64 {
	bits := binary.BigEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func F64LE(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func F32BE(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func F32LE(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	float := math.Float32frombits(bits)
	return float
}

func I64BE(bytes []byte) uint64 {
	bits := binary.BigEndian.Uint64(bytes)
	return bits
}

func I64LE(bytes []byte) uint64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return bits
}

func I32BE(bytes []byte) uint32 {
	bits := binary.BigEndian.Uint32(bytes)
	return bits
}

func I32LE(bytes []byte) uint32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return bits
}

func I16BE(bytes []byte) uint16 {
	bits := binary.BigEndian.Uint16(bytes)
	return bits
}

func I16LE(bytes []byte) uint16 {
	bits := binary.LittleEndian.Uint16(bytes)
	return bits
}

func S64BE(bytes []byte) string {
	ret:=``
	for i:=0;i<len(bytes);i=i+2{
		ret=ret+string(bytes[i+1])+string(bytes[i])
	}
	return ret
}

func S64LE(bytes []byte) string {
	ret:=``
	for i:=0;i<len(bytes);i=i+2{
		ret=ret+string(bytes[i])+string(bytes[i+1])
	}
	return ret
}


func formatting(bytes []byte,f string) []any {
	jump:=2
	ret:=[]any{}
	l:=len(bytes)
	log.Println(`bits=`,f[1:3])
	if f[1:3]==`64`{
		jump=8
	}
	if f[1:3]==`32`{
		jump=4
	}
	for i:=0;i<l-jump+1;i=i+jump{
		num:=any(nil)
		switch f{
		case `I64LE`:
			num=I64LE(bytes[i:i+jump])
		case `I64BE`:
			num=I64BE(bytes[i:i+jump])
		case `I32LE`:
			num=I32LE(bytes[i:i+jump])
		case `I32BE`:
			num=I32BE(bytes[i:i+jump])
		case `I16LE`:
			num=I16LE(bytes[i:i+jump])
		case `I16BE`:
			num=I16BE(bytes[i:i+jump])
		case `F64LE`:
			num=F64LE(bytes[i:i+jump])
		case `F64BE`:
			num=F64BE(bytes[i:i+jump])
		case `F32LE`:
			num=F32LE(bytes[i:i+jump])
		case `F32BE`:
			num=F32BE(bytes[i:i+jump])
		case `S64LE`:
			num=S64LE(bytes[i:i+jump])
		case `S64BE`:
			num=S64BE(bytes[i:i+jump])
		}
		
		ret=append(ret,num)
	}
	return ret
}

func modbusbytes(values []interface{},format string)([]byte,int){
	totalLength:=0
	totalBytes:=[]byte{}
	newVal:=[]interface{}{}
	switch format[:3]{
	case `I16`:
		for _,v:=range values{
			newVal=append(newVal,int16(v.(float64)))
		}
	case `I32`:
		for _,v:=range values{
			newVal=append(newVal,int32(v.(float64)))
		}
	case `I64`:
		for _,v:=range values{
			newVal=append(newVal,int64(v.(float64)))
		}
	case `F32`:
		for _,v:=range values{
			newVal=append(newVal,float32(v.(float64)))
		}
	case `F64`:
		for _,v:=range values{
			newVal=append(newVal,float64(v.(float64)))
		}
	case `S64`:
		for _,v:=range values{
			newVal=append(newVal,v.(string))
		}
	}
	endian:=0
	if format[3:]==`LE`{
		endian=1
	}

	for _,val:=range newVal{
		switch val.(type){
		case int16,int32,int64,float32,float64:
			b,l:=wNUM(val,endian)
			if l>0{
				totalBytes=append(totalBytes,b...)
				totalLength+=l
			}
		case string:
			b,l:=wSTR(val,endian)
			if l>0{
				totalBytes=append(totalBytes,b...)
				totalLength+=l
			}
		}
	}
	return totalBytes,totalLength
}

// write NUMBER to BYTES
// 
func wNUM( num interface{}, endian int) ([]byte,int){
	buf := new(bytes.Buffer)
	if endian==0{
		// ==0 , BigEndian
		binary.Write(buf, binary.BigEndian, num)
	}else{
		binary.Write(buf, binary.LittleEndian, num)
	}
	b:=buf.Bytes()
	l:=len(b)/2
	return b,l
}

func wSTR(str interface{},endian int) ([]byte,int) {
	buf:=new(bytes.Buffer)
	s:=str.(string)
	ln:=len(s)
	if ln%2==1{
		s=s+` `
		ln++
	}
	if endian==0{
		rs:=Reverse(s)
		buf.WriteString(rs)
	}else{
		buf.WriteString(s)
	}
	b:=buf.Bytes()
	ln=ln/2
	return b,ln
}

func Reverse(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}