package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

/////////////////////////////////////////////////
// DSL definition                              //
/////////////////////////////////////////////////

// @Name: users
// @Desc: Prints detailed user information including groups, shell, and last login
// @Param: search "" search term to match against all fields, empty string returns all users
func users(search string) (string, error) {
	search = strings.TrimSpace(search)
	fmt.Println("Searching for", "`", search, "`")
	// Read /etc/passwd to get user information
	passwd, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return "", err
	}

	// Read /etc/group to get group information
	groupFile, err := os.ReadFile("/etc/group")
	if err != nil {
		return "", err
	}

	// Parse group information
	groupMap := make(map[string][]string)
	lines := strings.Split(string(groupFile), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) >= 4 {
			groupName := fields[0]
			members := strings.Split(fields[3], ",")
			for _, member := range members {
				if member != "" {
					groupMap[member] = append(groupMap[member], groupName)
				}
			}
		}
	}

	// Parse lastlog information
	lastLogins, err := parseLastlog()
	if err != nil {
		// If we can't parse lastlog, we'll just show "N/A" for last login times
		lastLogins = make(map[string]time.Time)
	}

	// Parse user information
	var users []UserInfo
	lines = strings.Split(string(passwd), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) >= 7 {
			username := fields[0]
			uid := fields[2]
			gid := fields[3]
			home := fields[5]
			shell := fields[6]

			// Get user's groups
			groups := groupMap[username]
			sort.Strings(groups)
			groupsStr := strings.Join(groups, ", ")
			if groupsStr == "" {
				groupsStr = "N/A"
			}

			// Get last login time
			lastLogin := "N/A"
			if loginTime, exists := lastLogins[uid]; exists {
				lastLogin = loginTime.Format("2006-01-02 15:04:05")
			}

			users = append(users, UserInfo{
				Username:  username,
				UID:       uid,
				GID:       gid,
				Home:      home,
				Groups:    groupsStr,
				Shell:     shell,
				LastLogin: lastLogin,
			})
		}
	}

	// Filter users based on search term
	var filteredUsers []UserInfo
	for _, user := range users {
		// If search term is empty, include all users
		if search == "" {
			filteredUsers = append(filteredUsers, user)
			continue
		}

		// Check if search term matches any field
		searchLower := strings.ToLower(search)
		if strings.Contains(strings.ToLower(user.Username), searchLower) ||
			strings.Contains(strings.ToLower(user.UID), searchLower) ||
			strings.Contains(strings.ToLower(user.GID), searchLower) ||
			strings.Contains(strings.ToLower(user.Home), searchLower) ||
			strings.Contains(strings.ToLower(user.Groups), searchLower) ||
			strings.Contains(strings.ToLower(user.Shell), searchLower) ||
			strings.Contains(strings.ToLower(user.LastLogin), searchLower) {
			filteredUsers = append(filteredUsers, user)
		}
	}

	// Sort users by username
	sort.Slice(filteredUsers, func(i, j int) bool {
		return filteredUsers[i].Username < filteredUsers[j].Username
	})

	// Prepare template data
	data := UsersTemplateData{
		Title: "Users",
		Users: filteredUsers,
	}

	// Render using template
	return renderTemplate(usersTemplate, data)
}

// @Name: groups
// @Desc: Prints detailed group information including GID, members, and type
// @Param: search "" search term to match against all fields, empty string returns all groups
func groups(search string) (string, error) {
	search = strings.TrimSpace(search)
	// Read /etc/group to get group information
	groupFile, err := os.ReadFile("/etc/group")
	if err != nil {
		return "", err
	}

	// Read /etc/passwd to get user primary groups
	passwdFile, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return "", err
	}

	// Parse user primary groups and count usage
	primaryGroupCount := make(map[string]int)
	lines := strings.Split(string(passwdFile), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) >= 4 {
			gid := fields[3]
			primaryGroupCount[gid]++
		}
	}

	// Parse group information
	var groups []GroupInfo
	lines = strings.Split(string(groupFile), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Split(line, ":")
		if len(fields) >= 4 {
			name := fields[0]
			gid := fields[2]
			members := fields[3]

			// Determine group type
			gidNum, _ := strconv.Atoi(gid)
			groupType := "User"
			if gidNum < 1000 {
				groupType = "System"
			}

			// Count members
			membersList := strings.Split(members, ",")
			if members == "" {
				membersList = []string{}
			}
			memberCount := len(membersList)

			// Get primary group count
			primaryCount := primaryGroupCount[gid]

			// Determine if group is in use
			inUse := primaryCount > 0 || memberCount > 0

			// Format counts
			memberCountStr := "-"
			if memberCount > 0 {
				memberCountStr = fmt.Sprintf("%d", memberCount)
			}

			primaryCountStr := "-"
			if primaryCount > 0 {
				primaryCountStr = fmt.Sprintf("%d", primaryCount)
			}

			groups = append(groups, GroupInfo{
				Name:         name,
				GID:          gid,
				Type:         groupType,
				MemberCount:  memberCountStr,
				PrimaryCount: primaryCountStr,
				InUse:        inUse,
			})
		}
	}

	// Filter groups based on search term
	var filteredGroups []GroupInfo
	for _, group := range groups {
		// If search term is empty, include all groups
		if search == "" {
			filteredGroups = append(filteredGroups, group)
			continue
		}

		// Check if search term matches any field
		searchLower := strings.ToLower(search)
		if strings.Contains(strings.ToLower(group.Name), searchLower) ||
			strings.Contains(strings.ToLower(group.GID), searchLower) ||
			strings.Contains(strings.ToLower(group.Type), searchLower) ||
			strings.Contains(strings.ToLower(group.MemberCount), searchLower) ||
			strings.Contains(strings.ToLower(group.PrimaryCount), searchLower) {
			filteredGroups = append(filteredGroups, group)
		}
	}

	// Sort groups by name
	sort.Slice(filteredGroups, func(i, j int) bool {
		return filteredGroups[i].Name < filteredGroups[j].Name
	})

	// Prepare template data
	data := GroupsTemplateData{
		Title:  "Groups",
		Groups: filteredGroups,
	}

	// Render using template
	return renderTemplate(groupsTemplate, data)
}

// @Name: sw-info
// @Desc: Prints detailed software information including services, processes, and installed packages
func swInfo() (string, error) {
	// Read process information
	processes, err := os.ReadDir("/proc")
	if err != nil {
		return "", err
	}

	// Read environment variables
	env := os.Environ()

	// Prepare data
	data := SWInfoData{
		Title: "Software Information",
	}

	// Get system uptime for CPU calculation
	uptimeContent, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	}
	uptimeFields := strings.Fields(string(uptimeContent))
	if len(uptimeFields) < 1 {
		return "", fmt.Errorf("invalid uptime format")
	}
	uptime, err := strconv.ParseFloat(uptimeFields[0], 64)
	if err != nil {
		return "", err
	}

	// Get total memory for memory calculation
	memInfo, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}
	var totalMemory int64
	for _, line := range strings.Split(string(memInfo), "\n") {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				totalMemory, _ = strconv.ParseInt(fields[1], 10, 64)
				break
			}
		}
	}

	// Parse processes (limited to top 10 by CPU usage)
	processList := make([]ProcessInfo, 0)

	for _, process := range processes {
		if process.IsDir() {
			pid := process.Name()
			if _, err := strconv.Atoi(pid); err == nil {
				// Read process status
				status, err := os.ReadFile(fmt.Sprintf("/proc/%s/stat", pid))
				if err != nil {
					continue
				}

				// Parse process name and stats
				fields := strings.Fields(string(status))
				if len(fields) >= 24 {
					name := strings.Trim(fields[1], "()")

					// Calculate CPU percentage
					utime, _ := strconv.ParseFloat(fields[13], 64)
					stime, _ := strconv.ParseFloat(fields[14], 64)
					starttime, _ := strconv.ParseFloat(fields[21], 64)
					hertz := 100.0 // Linux default

					// Calculate CPU usage percentage
					totalTime := utime + stime
					seconds := uptime - (starttime / hertz)
					cpuUsage := 100.0 * (totalTime / hertz) / seconds

					// Calculate memory percentage
					rss, _ := strconv.ParseInt(fields[23], 10, 64)
					memUsage := float64(rss*100) / float64(totalMemory)

					processList = append(processList, ProcessInfo{
						Name:   name,
						PID:    pid,
						CPU:    fmt.Sprintf("%.1f", cpuUsage),
						Memory: fmt.Sprintf("%.1f", memUsage),
					})
				}
			}
		}
	}

	// Create separate lists for CPU and memory sorting
	cpuList := make([]ProcessInfo, len(processList))
	memList := make([]ProcessInfo, len(processList))
	copy(cpuList, processList)
	copy(memList, processList)

	// Sort by CPU usage and take top 10
	sort.Slice(cpuList, func(i, j int) bool {
		cpuI, _ := strconv.ParseFloat(cpuList[i].CPU, 64)
		cpuJ, _ := strconv.ParseFloat(cpuList[j].CPU, 64)
		return cpuI > cpuJ
	})

	if len(cpuList) > 10 {
		data.ProcessesByCPU = cpuList[:10]
	} else {
		data.ProcessesByCPU = cpuList
	}

	// Sort by memory usage
	sort.Slice(memList, func(i, j int) bool {
		memI, _ := strconv.ParseFloat(memList[i].Memory, 64)
		memJ, _ := strconv.ParseFloat(memList[j].Memory, 64)
		return memI > memJ
	})

	if len(memList) > 10 {
		data.ProcessesByMemory = memList[:10]
	} else {
		data.ProcessesByMemory = memList
	}

	// Parse environment variables (only show important ones)
	importantVars := map[string]bool{
		"PATH": true, "HOME": true, "USER": true, "SHELL": true,
		"TERM": true, "LANG": true, "PWD": true,
	}

	for _, e := range env {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) == 2 && importantVars[pair[0]] {
			data.Environment = append(data.Environment, struct {
				Name  string
				Value string
			}{
				Name:  pair[0],
				Value: pair[1],
			})
		}
	}

	return renderTemplate(swInfoTemplate, data)
}

// @Name: hw-info
// @Desc: Prints comprehensive hardware information including CPU, memory, disk, and network details
func hwInfo() (string, error) {
	// Read CPU information
	cpuInfo, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	// Parse CPU information
	cpuDetails := make(map[string]string)
	lines := strings.Split(string(cpuInfo), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			cpuDetails[key] = value
		}
	}

	// Read memory information
	memInfo, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return "", err
	}

	// Parse memory information
	memory := MemoryInfo{}
	swap := MemoryInfo{}
	for _, line := range strings.Split(string(memInfo), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, _ := strconv.ParseInt(fields[1], 10, 64)
		switch fields[0] {
		case "MemTotal:":
			memory.Total = fmt.Sprintf("%d MB", value/1024)
		case "MemFree:":
			memory.Free = fmt.Sprintf("%d MB", value/1024)
		case "SwapTotal:":
			swap.Total = fmt.Sprintf("%d MB", value/1024)
		case "SwapFree:":
			swap.Free = fmt.Sprintf("%d MB", value/1024)
		}
	}

	// Calculate memory usage
	memTotal, _ := strconv.ParseInt(strings.Fields(memory.Total)[0], 10, 64)
	memFree, _ := strconv.ParseInt(strings.Fields(memory.Free)[0], 10, 64)
	memUsed := memTotal - memFree
	memory.Used = fmt.Sprintf("%d MB", memUsed)
	memory.Usage = fmt.Sprintf("%.1f%%", float64(memUsed)/float64(memTotal)*100)

	// Calculate swap usage
	swapTotal, _ := strconv.ParseInt(strings.Fields(swap.Total)[0], 10, 64)
	swapFree, _ := strconv.ParseInt(strings.Fields(swap.Free)[0], 10, 64)
	swapUsed := swapTotal - swapFree
	swap.Used = fmt.Sprintf("%d MB", swapUsed)
	swap.Usage = fmt.Sprintf("%.1f%%", float64(swapUsed)/float64(swapTotal)*100)

	// Get disk information
	cmd := exec.Command("df", "-h")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	var disks []DiskInfo
	lines = strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 6 {
			disks = append(disks, DiskInfo{
				Device:     fields[0],
				MountPoint: fields[5],
				Total:      fields[1],
				Used:       fields[2],
				Free:       fields[3],
				Usage:      fields[4],
			})
		}
	}

	// Get network interface information
	cmd = exec.Command("ip", "addr")
	output, err = cmd.Output()
	if err != nil {
		return "", err
	}

	var interfaces []NetworkInterface
	currentInterface := NetworkInterface{}
	lines = strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "inet ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				currentInterface.IP = strings.Split(fields[1], "/")[0]
			}
		} else if strings.HasPrefix(line, "link/ether ") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				currentInterface.MAC = fields[1]
			}
		} else if strings.HasPrefix(line, "UP") {
			currentInterface.Status = "UP"
		} else if strings.HasPrefix(line, "DOWN") {
			currentInterface.Status = "DOWN"
		} else if strings.HasPrefix(line, "2:") {
			if currentInterface.Name != "" {
				interfaces = append(interfaces, currentInterface)
			}
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				currentInterface = NetworkInterface{
					Name: strings.TrimSuffix(fields[1], ":"),
				}
			}
		}
	}
	if currentInterface.Name != "" {
		interfaces = append(interfaces, currentInterface)
	}

	// Prepare template data
	data := HWInfoData{
		Title:             "Hardware Information",
		CPUInfo:           cpuDetails,
		Memory:            memory,
		Swap:              swap,
		Disks:             disks,
		NetworkInterfaces: interfaces,
	}

	// Render using template
	return renderTemplate(hwInfoTemplate, data)
}

// @Name: os-info
// @Desc: Prints detailed operating system information including kernel, distribution, and system details
func osInfo() (string, error) {
	// Read kernel information
	uname, err := exec.Command("uname", "-a").Output()
	if err != nil {
		return "", err
	}

	// Parse distribution information
	var distroName, distroVersion string
	if _, err := os.Stat("/etc/os-release"); err == nil {
		content, err := os.ReadFile("/etc/os-release")
		if err == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "PRETTY_NAME=") {
					distroName = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				} else if strings.HasPrefix(line, "VERSION_ID=") {
					distroVersion = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
				}
			}
		}
	}

	// If we couldn't get the distribution info, try lsb_release
	if distroName == "" {
		if output, err := exec.Command("lsb_release", "-d").Output(); err == nil {
			distroName = strings.TrimSpace(strings.TrimPrefix(string(output), "Description:\t"))
		}
	}

	// Combine distro info
	distroInfo := "N/A"
	if distroName != "" {
		distroInfo = distroName
		if distroVersion != "" {
			distroInfo += " " + distroVersion
		}
	}

	// Get system uptime
	uptimeContent, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return "", err
	}
	uptimeFields := strings.Fields(string(uptimeContent))
	if len(uptimeFields) < 1 {
		return "", fmt.Errorf("invalid uptime format")
	}
	uptimeSeconds, err := strconv.ParseFloat(uptimeFields[0], 64)
	if err != nil {
		return "", err
	}
	uptime := time.Duration(uptimeSeconds * float64(time.Second))

	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}

	// Get timezone
	timezone, err := os.ReadFile("/etc/timezone")
	if err != nil {
		timezone = []byte("N/A")
	}

	// Get system load average
	loadAvg, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "", err
	}

	// Prepare template data
	data := OSInfoData{
		Title:    "Operating System Information",
		Kernel:   strings.TrimSpace(string(uname)),
		Distro:   distroInfo,
		Uptime:   uptime.String(),
		Hostname: hostname,
		Timezone: strings.TrimSpace(string(timezone)),
		LoadAvg:  strings.TrimSpace(string(loadAvg)),
	}

	// Render using template
	return renderTemplate(osInfoTemplate, data)
}
