package main

import (
    "fmt"
    "sort"
    "strings"
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
    "os"
    "net"
   // "container/list"
)

type MyFile struct {
    Index   int
    Name    string
    ModiDate   string
    checked bool
}

type MyFileModel struct {
    walk.TableModelBase
    walk.SorterBase
    sortColumn int
    sortOrder  walk.SortOrder
    items      []*MyFile
}

func (m *MyFileModel) RowCount() int {
    return len(m.items)
}

func (m *MyFileModel) Value(row, col int) interface{} {
    item := m.items[row]
    switch col {
    case 0:
        return item.Index
    case 1:
        return item.Name
    case 2:
        return item.ModiDate
    }
    panic("unexpected col")
}

func (m *MyFileModel) Checked(row int) bool {
    return m.items[row].checked
}

func (m *MyFileModel) SetChecked(row int, checked bool) error {
    m.items[row].checked = checked
    return nil
}

func (m *MyFileModel) Sort(col int, order walk.SortOrder) error {
    m.sortColumn, m.sortOrder = col, order

    sort.Stable(m)

    return m.SorterBase.Sort(col, order)
}

func (m *MyFileModel) Len() int {
    return len(m.items)
}

func (m *MyFileModel) Less(i, j int) bool {
    a, b := m.items[i], m.items[j]

    c := func(ls bool) bool {
        if m.sortOrder == walk.SortAscending {
            return ls
        }

        return !ls
    }

    switch m.sortColumn {
    case 0:
        return c(a.Index < b.Index)
    case 1:
        return c(a.Name < b.Name)
    case 2:
        return c(a.ModiDate < b.ModiDate)
    }

    panic("unreachable")
}

func (m *MyFileModel) Swap(i, j int) {
    m.items[i], m.items[j] = m.items[j], m.items[i]
}

func NewMyFileModel(f string) *MyFileModel {
    m := new(MyFileModel)
   
    lines := strings.Split(f, "\n")
    m.items = make([]*MyFile, len(lines)-1)
    for  i:=0; i<len(lines)-1; i++ {
    	a := strings.Split(lines[i], "\t")
		m.items[i] = &MyFile{
	        Index: i,     
	        Name: a[0],
	        ModiDate: a[1],
		    }

	}
    return m
}

type MyFileMainWindow struct {
    *walk.MainWindow
    model *MyFileModel
    tv    *walk.TableView
}
func startAReceiver(){
	var (
		// host   = "192.168.1.5"	
		port = "9091"
		// remote = host + ":" + port

		remote = ":" + port
	)

	fmt.Println("Server initializing...")

	lis, err := net.Listen("tcp", remote)
	defer lis.Close()

	if err != nil {
		fmt.Println("Error occurs listening to the port! ", remote)
		os.Exit(-1)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			fmt.Println("An Error occurs on client: ", err.Error())
			// os.Exit(0)
			continue
		}
	go receiveFile(conn)
	}
}
func main() {
	//showFileInfo("test.txt")
	//gothrough()
	infos:=getServerInfo("127.0.0.1:9090")
	//m:="secret1.txt	2019-04-18 09:48:01.267 -0400 EDT	dd945ab221b14e3be0d31fd4026f27eb\nsecret2.txt	2019-04-18 09:48:01.282 -0400 EDT	dd945ab221b14e3be0d31fd4026f27eb\nsecret3.txt	2019-04-18 09:48:01.304 -0400 EDT	dd945ab221b14e3be0d31fd4026f27eb\nsecret4.txt	2019-04-18 09:48:01.324 -0400 EDT	dd945ab221b14e3be0d31fd4026f27eb\nsecret5.txt	2019-04-18 09:48:01.346 -0400 EDT	dd945ab221b14e3be0d31fd4026f27eb"
    m :=LogAllFile("./vcv")
    mw := &MyFileMainWindow{model: NewMyFileModel(m)}
    mserver := &MyFileMainWindow{model: NewMyFileModel(infos)}
    go  startAReceiver()
    MainWindow{
        AssignTo: &mw.MainWindow,
        Title:    "Folder Synchronizer",
        Size:     Size{800, 600},
        Layout:   VBox{},
        Children: []Widget{
            Composite{
                Layout: HBox{MarginsZero: true},
                Children: []Widget{
                    HSpacer{},
                    PushButton{
                        Text: "Refresh",
                        OnClicked: func() {
							mserver.model.items = NewMyFileModel(getServerInfo("127.0.0.1:9090")).items
                            mserver.model.PublishRowsReset()
                            mw.model.items = NewMyFileModel(LogAllFile("./vcv")).items
                            mw.model.PublishRowsReset()
                        },
                    },
                    PushButton{
                        Text: "UploadAll",
                        OnClicked: func() {
							for _, x := range mw.model.items {
								sendFile(x.Name)
							}
							mserver.model.items = NewMyFileModel(getServerInfo("127.0.0.1:9090")).items
                            mserver.model.PublishRowsReset()
                        },
                    },
                    PushButton{
                        Text: "UploadNew",
                        OnClicked: func() {
                        	infos:=getServerInfo("127.0.0.1:9090")
							for _, x := range mw.model.items {
								if(!strings.Contains(infos, x.Name)){
									sendFile(x.Name)
								}
							}
							mserver = &MyFileMainWindow{model: NewMyFileModel(getServerInfo("127.0.0.1:9090"))}
                            mserver.model.PublishRowsReset()
                        },
                    },
                    PushButton{
                        Text: "DownloadAll",
                        OnClicked: func() {
                        	filelist:=""
                        	for _, x := range mserver.model.items {
	                        	filelist+=x.Name+"\n"
                        	}                    	
							ask4file("127.0.0.1:9090",filelist)
							mw.model.items = NewMyFileModel(LogAllFile("./vcv")).items
                            mw.model.PublishRowsReset()
                        },
                    },
                    PushButton{
                        Text: "DownloadNew",
                        OnClicked: func() {
							filelist:=""
							for _, x := range mserver.model.items {
								if(!strings.Contains(m, x.Name)){
									filelist+=x.Name+"\n"
								}
							}
							ask4file("127.0.0.1:9090",filelist)
							mserver = &MyFileMainWindow{model: NewMyFileModel(getServerInfo("127.0.0.1:9090"))}
                            mserver.model.PublishRowsReset()
                        },
                    },
                    PushButton{
                        Text: "UpdateChecked",
                        OnClicked: func() {
                            for _, x := range mw.model.items {
                                if x.checked {
                                    sendFile(x.Name)
                                }
                            }
                            mserver.model.items = NewMyFileModel(getServerInfo("127.0.0.1:9090")).items
                            mserver.model.PublishRowsReset()
                            mw.tv.SetSelectedIndexes([]int{})
                        },
                    },
                    PushButton{
                        Text: "DownloadChecked",
                        OnClicked: func() {
                        	filelist:=""
                            for _, x := range mserver.model.items {
                                if x.checked {
                                  filelist+=x.Name+"\n" 
                                }
                            }
                            if(len(filelist)>0){
                             ask4file("127.0.0.1:9090",filelist)
                            }
                            mw.model.items = NewMyFileModel(LogAllFile("./vcv")).items
                            mw.model.PublishRowsReset()
                            mserver.tv.SetSelectedIndexes([]int{})
                        },
                    },
                },
            },
            Composite{
                Layout: VBox{},
                ContextMenuItems: []MenuItem{
                    Action{
                        Text:        "I&nfo",
                        OnTriggered: mw.tv_ItemActivated,
                    },
                    Action{
                        Text: "E&xit",
                        OnTriggered: func() {
                            mw.Close()
                        },
                    },
                },
                Children: []Widget{
                    TableView{
                        AssignTo:         &mw.tv,
                        CheckBoxes:       true,
                        ColumnsOrderable: true,
                        MultiSelection:   true,
                        Columns: []TableViewColumn{
                            {Title: "Index"},
                            {Title: "Name"},
                            {Title: "ModifyDate"},
                        },
                        Model: mw.model,
                        OnCurrentIndexChanged: func() {
                            i := mw.tv.CurrentIndex()
                            if 0 <= i {
                                fmt.Printf("OnCurrentIndexChanged: %v\n", mw.model.items[i].Name)
                            }
                        },
                        OnItemActivated: mw.tv_ItemActivated,
                    },
                    TableView{
                        AssignTo:         &mserver.tv,
                        CheckBoxes:       true,
                        ColumnsOrderable: true,
                        MultiSelection:   true,
                        Columns: []TableViewColumn{
                            {Title: "Index"},
                            {Title: "Name"},
                            {Title: "ModifyDate"},
                        },
                        Model: mserver.model,
                        OnCurrentIndexChanged: func() {
                            i := mserver.tv.CurrentIndex()
                            if 0 <= i {
                                fmt.Printf("OnCurrentIndexChanged: %v\n", mserver.model.items[i].Name)
                            }
                        },
                        OnItemActivated: mserver.tv_ItemActivated,
                    },
                },
            },
        },
    }.Run()

}

func (mw *MyFileMainWindow) tv_ItemActivated() {
    msg := ``
    for _, i := range mw.tv.SelectedIndexes() {
        msg = msg + "\n" + mw.model.items[i].Name
    }
    walk.MsgBox(mw, "title", msg, walk.MsgBoxIconInformation)
}