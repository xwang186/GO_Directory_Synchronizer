package main

import (
//	"fmt"
    "os"
    "path/filepath"
    "strings"
    "container/list"
    "io/ioutil"
)
func gothrough(){
	filepath.Walk("./vcv", walkfunc)
}
func walkfunc(path string, info os.FileInfo, err error) error {
    linuxPath := strings.Replace(path, "\\", "/", -1)
    if(!info.IsDir()){
	    sendFile(linuxPath)
    }
    return nil
}
func getFileList() *list.List{
	l := list.New()
    files, _ := ioutil.ReadDir("./vcv")
    for _, f := range files {
       l.PushBack(f.Name())
    }
	//filepath.Walk("./vcv",walkGetName,l)
	return l
}
func walkGetName(l list.List,path string, info os.FileInfo, err error) error{
	linuxPath:=strings.Replace(path, "\\", "/", -1)
    if(!info.IsDir()){
	    l.PushBack(linuxPath)
    }
    return nil
}
func logDeep(root string){
	filepath.Walk(root, logFolder)
}
func logFolder(path string, info os.FileInfo, err error) error{
	if(!info.IsDir()){
	 	return nil
	}	 
	f, err := os.Create(path+"/fileInfo")
	files, _ := ioutil.ReadDir(path)
	for _, fi := range files {
		if(fi.Name()!="fileInfo"){
		md5:=SumMd5(path+"/"+fi.Name())
       f.WriteString(fi.Name()+"\t"+fi.ModTime().String()+"\t"+md5+"\n")
		}
    }
	f.Close()
	return err
}
func LogAllFile(pathname string)string{
	m:=GetAllFile(pathname)
	f,_ := os.Create(pathname+"/fileInfo")
	f.WriteString(m)
	f.Close()
	return m
}
func GetAllFile(pathname string) string {
    rd, _ := ioutil.ReadDir(pathname)
    result:="";
    for _, fi := range rd {
        if fi.IsDir() {
        	result+=pathname+"/"+fi.Name()+"\t"+fi.ModTime().String()+"\t"+"---"+"\n"
            result+=GetAllFile(pathname +"/"+ fi.Name())
        } else {
            md5:=SumMd5(pathname+"/"+fi.Name())
	        result+=pathname+"/"+fi.Name()+"\t"+fi.ModTime().String()+"\t"+md5+"\n"
        }
    }
    return result
}