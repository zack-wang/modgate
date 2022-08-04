package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"errors"
	"flag"
	"github.com/goburrow/modbus"
	"fmt"
	"time"
)

var (
	bind=flag.String(`b`,`:8080`,`bind address`)
)

func main() {
	flag.Parse()
	server := gin.Default()
	server.GET("/:protocol/:connection/:fc/:start/:len/*format", modbusRead)
	server.GET("/",RootPage)
	server.POST("/:protocol/:connection/:fc/:start/:format", modbusWrite)
	server.Run(*bind)
}


func RootPage(c *gin.Context){
	c.Redirect(http.StatusFound, "https://github.com/zack-wang/modgate")
}


func modbusWrite(c *gin.Context){
	proto:=c.Param("protocol")
	if proto!=`modbustcp`{
		c.JSON(http.StatusOK, gin.H{
			"code":-100,
			"message":proto+` not supported`,
		})
		return
	}
	host:=c.Param("connection")
	//client:=modbus.TCPClient(host)
	fc:=c.Param("fc")
	start:=c.Param("start")
	s:=0
	fmt.Sscanf(start,"%d",&s)
	format:=c.Param("format") // two characters BE or LE
	//format=format[3:]
	fmt.Println(`Endian=`,format,"function code=",fc)

	dat:=[]interface{}{}
	if err := c.ShouldBindJSON(&dat); err != nil {
		c.JSON(200, gin.H{
		  "code": -500,
		  "message":  err.Error(),
		})
		return
	}

	client:=modbus.TCPClient(host)

	ln:=0

	switch fc{
	case `15`:
		//w:=dataWidth(format)
			//ln=w*(len(dat))
			_, err := client.WriteMultipleCoils(uint16(s), uint16(1), []byte{4, 3})
			if err!=nil{
				fmt.Println(err.Error())
			}
			/*
				c.JSON(200, gin.H{
					"code": -400,
					"message":  `Invalid data format`,
				})
				return
			*/
	case `16`:
		b,n:=modbusbytes(dat,format)
		_, err := client.WriteMultipleRegisters(uint16(s), uint16(n), b)
		if err!=nil{
			fmt.Println(err.Error())
		}
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":2,
			"message":`invalidate func code`,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"message":``,
		"data":[]byte{0},
		"format":format,
		"length":ln,
		"host":host,
	})
}

func modbusRead(c *gin.Context){
	proto:=c.Param("protocol")
	if proto!=`modbustcp`{
		c.JSON(http.StatusOK, gin.H{
			"code":-100,
			"message":proto+` not supported`,
		})
		return
	}
	host:=c.Param("connection")
	fc:=c.Param("fc")
	start:=c.Param("start")
	length:=c.Param("len")
	s,ln:=0,0
	fmt.Sscanf(start,"%d",&s)
	fmt.Sscanf(length,"%d",&ln)

	format:=c.Param("format")
	format=format[1:]

	client:=modbus.TCPClient(host)

	r:=[]byte{}
	var err error

	switch fc{
	case `1`:
		r,err=client.ReadCoils(uint16(s),uint16(ln))
	case `2`:
		r,err=client.ReadDiscreteInputs(uint16(s),uint16(ln))
	case `3`:
		r,err=client.ReadHoldingRegisters(uint16(s),uint16(ln))
	case `4`:
		r,err=client.ReadInputRegisters(uint16(s),uint16(ln))
	default:
		c.JSON(http.StatusOK, gin.H{
			"code":-300,
			"message":`invalidate func code`,
		})
		return
	}

	if err!=nil{
		fmt.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code":-200,
			"message":err.Error(),
		})
		return
	}

	rslt:=formatting(r,format)
	fmt.Println(time.Now().Format(`2006-01-02 15:04:05`),`result=`,rslt)
	c.JSON(http.StatusOK, gin.H{
		"code":0,
		"message":``,
		"data":rslt,
		"format":format,
	})
}

func dataWidth(s string) int{
	ret:=-1
	switch s{
	case `I64BE`,`I64LE`,`F64BE`,`F64LE`,`S64BE`,`S64LE`:
		ret= 4
	case `I32BE`,`I32LE`,`F32BE`,`F32LE`:
		ret= 2
	case `I16BE`,`I16LE`:
		ret= 1
	}
	return ret
}