package main

import (
	_ "embed"
)

/////////////////////////////////////////////////
// Utility templates                           //
/////////////////////////////////////////////////

//go:embed templates/groups.tmpl
var groupsTemplate string

//go:embed templates/users.tmpl
var usersTemplate string

//go:embed templates/sw_info.tmpl
var swInfoTemplate string

//go:embed templates/hw_info.tmpl
var hwInfoTemplate string

//go:embed templates/os_info.tmpl
var osInfoTemplate string
