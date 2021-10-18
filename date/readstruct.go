package date

type ModuleInfo struct{
	Path string
	Name string
	Desk string
	Cmds []CmdsInfo
}

type CmdsInfo struct{
	Cmd string
	Desk string
}

type CmdRequest struct{
	Module string
	Cmd string
}