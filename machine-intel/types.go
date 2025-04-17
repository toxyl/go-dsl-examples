package main

/////////////////////////////////////////////////
// Utility types                               //
/////////////////////////////////////////////////

type TemplateData struct {
	Title   string
	Items   interface{}
	Details map[string]string
}

// UserInfo holds detailed information about a user
type UserInfo struct {
	Username  string
	UID       string
	GID       string
	Home      string
	Groups    string
	Shell     string
	LastLogin string
}

// LastlogEntry represents the structure of a lastlog entry
type LastlogEntry struct {
	Time uint32
	Line [32]byte
	Host [256]byte
}

// UsersTemplateData holds the data structure for user template rendering
type UsersTemplateData struct {
	Title string
	Users []UserInfo
}

// GroupInfo holds detailed information about a group
type GroupInfo struct {
	Name         string
	GID          string
	Type         string
	MemberCount  string
	PrimaryCount string
	InUse        bool
}

// GroupsTemplateData holds the data structure for group template rendering
type GroupsTemplateData struct {
	Title  string
	Groups []GroupInfo
}

// SWInfoData holds the data structure for swInfo template rendering
type SWInfoData struct {
	Title             string
	ProcessesByCPU    []ProcessInfo
	ProcessesByMemory []ProcessInfo
	Environment       []struct {
		Name  string
		Value string
	}
}

type ProcessInfo struct {
	Name   string
	PID    string
	CPU    string
	Memory string
}

// MemoryInfo holds memory usage information
type MemoryInfo struct {
	Total string
	Used  string
	Free  string
	Usage string
}

// DiskInfo holds disk usage information
type DiskInfo struct {
	Device     string
	MountPoint string
	Total      string
	Used       string
	Free       string
	Usage      string
}

// NetworkInterface holds network interface information
type NetworkInterface struct {
	Name   string
	IP     string
	MAC    string
	Status string
}

// HWInfoData holds the data structure for hwInfo template rendering
type HWInfoData struct {
	Title             string
	CPUInfo           map[string]string
	Memory            MemoryInfo
	Swap              MemoryInfo
	Disks             []DiskInfo
	NetworkInterfaces []NetworkInterface
}

// OSInfoData represents the data structure for operating system information
type OSInfoData struct {
	Title    string
	Kernel   string
	Distro   string
	Uptime   string
	Hostname string
	Timezone string
	LoadAvg  string
}
