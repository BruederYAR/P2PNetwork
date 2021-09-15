package date

type ModuleInfo struct{
	Name string
	Desk string
	Cmds []CmdsInfo
}

type CmdsInfo struct{
	Cmd string
	Desk string
}