package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func sendFile(fileName string) {
	var (
		host   = "127.0.0.1"  
		//host   = "192.168.56.101"
		port   = "9090"  
		remote = host + ":" + port 

		mergeFileName = fileName 
		coroutine     = 10        
		bufsize       = 1024  
	)

	for index, sargs := range os.Args {
		switch index {
		case 1:
			fileName = sargs
			mergeFileName = sargs
		case 2:
			mergeFileName = sargs
		case 3:
			bufsize, _ = strconv.Atoi(sargs)
		case 4:
			coroutine, _ = strconv.Atoi(sargs)
		}

	}
	fileInfo, _ := os.Stat(fileName)
	coroutine = int(fileInfo.Size()/4098) + 1
	if(fileInfo.IsDir()){
		sendCreateFolderCommand(remote,fileName);
	}
	fl, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("userFile", err)
		return
	}

	stat, err := fl.Stat()
	if err != nil {
		panic(err)
	}
	var size int64
	size = stat.Size()
	fl.Close()

	littleSize := size / int64(coroutine)

	fmt.Printf("Size: %d  %d \n", size, littleSize)

	begintime := time.Now().Unix()
	c := make(chan string)
	var begin int64 = 0
	for i := 0; i < coroutine; i++ {

		if i == coroutine-1 {
			go splitFile(remote, c, i, bufsize, fileName, mergeFileName, begin, size)
			fmt.Println(begin, size, bufsize)
		} else {
			go splitFile(remote, c, i, bufsize, fileName, mergeFileName, begin, begin+littleSize)
			fmt.Println(begin, begin+littleSize)
		}

		begin += littleSize
	}

	for j := 0; j < coroutine; j++ {
		fmt.Println(<-c)
	}

	midtime := time.Now().Unix()
	sendtime := midtime - begintime
	fmt.Printf("Sending time:%d min %d sec \n", sendtime/60, sendtime%60)

	sendMergeCommand(remote, mergeFileName, coroutine) 
	endtime := time.Now().Unix()

	mergetime := endtime - midtime
	fmt.Printf("Merging time:%d min %d sec \n", mergetime/60, mergetime%60)

	tot := endtime - begintime
	fmt.Printf("Total time:%d min %d sec \n", tot/60, tot%60)
}

func splitFile(remote string, c chan string, coroutineNum int, size int, fileName, mergeFileName string, begin int64, end int64) {

	con, err := net.Dial("tcp", remote)
	defer con.Close()
	if err != nil {
		fmt.Println("Connection fails.")
		os.Exit(-1)
		return
	}
	fmt.Println(coroutineNum, "Connection established...")

	var by [1]byte
	by[0] = byte(coroutineNum)
	var bys []byte
	databuf := bytes.NewBuffer(bys) 
	databuf.Write(by[:])
	databuf.WriteString(mergeFileName)
	bb := databuf.Bytes()
	// bb := by[:]
	// fmt.Println(bb)
	in, err := con.Write(bb)
	if err != nil {
		fmt.Printf("Send Error: %d\n", in)
		os.Exit(0)
	}

	var msg = make([]byte, 1024)
	lengthh, err := con.Read(msg) 
	if err != nil {
		fmt.Printf("Read error.\n", lengthh)
		os.Exit(0)
	}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		fmt.Println(fileName, "-can not send the file.")
		os.Exit(0)
	}

	file.Seek(begin, 0)

	buf := make([]byte, size)

	var sendDtaTolNum int = 0 
	for i := begin; int64(i) < end; i += int64(size) {
		length, err := file.Read(buf) 
		if err != nil {
			fmt.Println("Cannot read the target file", i, coroutineNum, end)
		}

		if length == size {
			if int64(i)+int64(size) >= end {
				sendDataNum, err := con.Write(buf[:size-int((int64(i)+int64(size)-end))])
				if err != nil {
					fmt.Printf("Send Error: %d\n", sendDataNum)
					os.Exit(0)
				}
				sendDtaTolNum += sendDataNum
			} else {
				sendDataNum, err := con.Write(buf)
				if err != nil {
					fmt.Printf("Send Error: %d\n", sendDataNum)
					os.Exit(0)
				}
				sendDtaTolNum += sendDataNum
			}

		} else {
			sendDataNum, err := con.Write(buf[:length])
			if err != nil {
				fmt.Printf("Send Error: %d\n", sendDataNum)
				os.Exit(0)
			}
			sendDtaTolNum += sendDataNum
		}

		lengths, err := con.Read(msg)
		if err != nil {
			fmt.Printf("Send Error.\n", lengths)
			os.Exit(0)
		}

	}

	fmt.Println(coroutineNum, "Sending Data...:", sendDtaTolNum)

	c <- strconv.Itoa(coroutineNum) + " Routine exit"
}

func sendMergeCommand(remote, mergeFileName string, coroutine int) {

	con, err := net.Dial("tcp", remote)
	defer con.Close()
	if err != nil {
		fmt.Println("connection fails.")
		os.Exit(-1)
		return
	}
	fmt.Println("Sent a request to create files.\nMerging...")

	var by [1]byte
	by[0] = byte(coroutine)
	var bys []byte
	databuf := bytes.NewBuffer(bys) 
	databuf.WriteString("fileover")
	databuf.Write(by[:])
	databuf.WriteString(mergeFileName)
	cmm := databuf.Bytes()

	in, err := con.Write(cmm)
	if err != nil {
		fmt.Printf("Send Error: %d\n", in)
	}

	var msg = make([]byte, 1024)
	lengthh, err := con.Read(msg)
	if err != nil {
		fmt.Printf("Send Error.\n", lengthh)
		os.Exit(0)
	}
	str := string(msg[0:lengthh])
	fmt.Println("Send mission complished: ", str)
}
func sendCreateFolderCommand(remote, folderpath string) {

	con, err := net.Dial("tcp", remote)
	defer con.Close()
	if err != nil {
		fmt.Println("connection fails.")
		os.Exit(-1)
		return
	}
	fmt.Println("Sent a request to create folders.",folderpath,"\nMerging...")
	var bys []byte
	databuf := bytes.NewBuffer(bys)
	databuf.WriteString("crtfoler:")
	databuf.WriteString(folderpath)
	cmm := databuf.Bytes()

	in, err := con.Write(cmm)
	if err != nil {
		fmt.Printf("Send Error: %d\n", in)
	}

	var msg = make([]byte, 1024)
	lengthh, err := con.Read(msg)
	if err != nil {
		fmt.Printf("Send Error.\n", lengthh)
		os.Exit(0)
	}
	str := string(msg[0:lengthh])
	fmt.Println("Send mission complished: ", str)
}

func getServerInfo(remote string) string {

	con, err := net.Dial("tcp", remote)
	defer con.Close()
	if err != nil {
		fmt.Println("connection fails.")
		os.Exit(-1)
	}
	fmt.Println("Sent a request to get server's folder information...")
	var bys []byte
	databuf := bytes.NewBuffer(bys)
	databuf.WriteString("updtDict:")
	databuf.WriteString("some path")
	cmm := databuf.Bytes()

	in, err := con.Write(cmm)
	if err != nil {
		fmt.Printf("Send Error: %d\n", in)
	}

	var msg = make([]byte, 4096*3)
	lengthh, err := con.Read(msg)
	if err != nil {
		fmt.Printf("Send Error.\n", lengthh)
		os.Exit(0)
	}
	str := string(msg[0:lengthh])
	//fmt.Println("Send mission complished: ", str)
	return str
}
func ask4file(remote string,filelist string){
	fmt.Print("Sending file list :", filelist)
	con, err := net.Dial("tcp", remote)
	defer con.Close()
	if err != nil {
		fmt.Println("connection fails.")
		os.Exit(-1)
	}
	fmt.Println("Sent a request to get server's folder information...")
	var bys []byte
	databuf := bytes.NewBuffer(bys)
	databuf.WriteString("ask4file:")
	fmt.Print("Client asking for files: ",filelist);
	databuf.WriteString(filelist)
	cmm := databuf.Bytes()
	in, err := con.Write(cmm)
	if err != nil {
		fmt.Printf("Send Error: %d\n", in)
	}
	con.Close()
}