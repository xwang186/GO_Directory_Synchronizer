package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	//"runtime"
	"strconv"
	//"log"
)
func receiveFile(con net.Conn) {
	var (
		res          string
		tempFileName string                    
		data         = make([]byte, 1024*1024)
		by           []byte
		databuf      = bytes.NewBuffer(by) 
		fileNum      int                  
	)
	defer con.Close()

	fmt.Println("New Connection established ", con.RemoteAddr())
	j := 0 
	for {
		length, err := con.Read(data)
		if err != nil {
			da := databuf.Bytes()
			fmt.Printf("client %v has disconnnected %2d %d \n", con.RemoteAddr(), fileNum, len(da))
			return
		}
		if 0 == j {
			res = string(data[0:8])
			if "fileover" == res { 
				xienum := int(data[8])
				mergeFileName := string(data[9:length])
				go mainMergeFile(xienum, mergeFileName) 
				res = "Received " + mergeFileName
				con.Write([]byte(res))
				fmt.Println(mergeFileName, "Received")
				return
			}	else if "crtfoler" == res{
				folderpath := string(data[9:length])
				os.MkdirAll(folderpath, os.ModePerm)
			}   else if "updtDict" == res{
				cont:=LogAllFile("./vcv")
				//cont:=sendBackDicInfo()
				con.Write([]byte(cont))
				fmt.Println("Sent Server's folder information to client")
			}else { 
				fileNum = int(data[0])
				tempFileName = string(data[1:length]) + strconv.Itoa(fileNum)
				fmt.Println("create temp file:", tempFileName)
				fout, err := os.Create(tempFileName)
				if err != nil {
					fmt.Println("Cannot create temp file", tempFileName)
					return
				}
				fout.Close()
			}
		} else {
			// databuf.Write(data[0:length])
			writeTempFileEnd(tempFileName, data[0:length])
		}

		res = strconv.Itoa(fileNum) + "received"
		con.Write([]byte(res))
		j++
	}

}
func writeTempFileEnd(fileName string, data []byte) {
	tempFile, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		// panic(err)
		fmt.Println("Open file failed", err)
		return
	}
	defer tempFile.Close()
	tempFile.Write(data)
}

func mainMergeFile(connumber int, filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error create file", err)
		return
	}
	defer file.Close()

	for i := 0; i < connumber; i++ {
		mergeFile(filename+strconv.Itoa(i), file)
	}

	for i := 0; i < connumber; i++ {
		os.Remove(filename + strconv.Itoa(i))
	}

}

func mergeFile(rfilename string, wfile *os.File) {

	// fmt.Println(rfilename, wfilename)
	rfile, err := os.OpenFile(rfilename, os.O_RDWR, 0666)
	defer rfile.Close()
	if err != nil {
		fmt.Println("Error occura when merging files", rfilename)
		return
	}

	stat, err := rfile.Stat()
	if err != nil {
		panic(err)
	}

	num := stat.Size()

	buf := make([]byte, 1024*1024)
	for i := 0; int64(i) < num; {
		length, err := rfile.Read(buf)
		if err != nil {
			fmt.Println("File Cannot Read.")
		}
		i += length

		wfile.Write(buf[:length])
	}

}