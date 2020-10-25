package utils

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"fmt"

	"github.com/alioygur/godash"
	"github.com/chaseSpace/kit/fs"
	"github.com/spf13/viper"
	"golang.org/x/tools/imports"
)

// ToLowerFirstCamelCase returns the given string in camelcase formatted string
// but with the first letter being lowercase.
func ToLowerFirstCamelCase(s string) string {
	if s == "" {
		return s
	}
	if len(s) == 1 {
		return strings.ToLower(string(s[0]))
	}
	return strings.ToLower(string(s[0])) + godash.ToCamelCase(s)[1:]
}

// ToUpperFirst returns the given string with the first letter being uppercase.
func ToUpperFirst(s string) string {
	if s == "" {
		return s
	}
	if len(s) == 1 {
		return strings.ToLower(string(s[0]))
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}

// ToLowerSnakeCase the given string in snake-case format.
func ToLowerSnakeCase(s string) string {
	return strings.ToLower(godash.ToSnakeCase(s))
}

// ToCamelCase the given string in camelcase format.
func ToCamelCase(s string) string {
	return godash.ToCamelCase(s)
}

// GoImportsSource is used to format and optimize imports the
// given source.
func GoImportsSource(path string, s string) (string, error) {
	is, err := imports.Process(path, []byte(s), nil)
	return string(is), err
}

// GetServiceImportPath returns the import path of the service interface.
func GetServiceImportPath(name string) (string, error) {
	return getImportPath(name, "gk_service_path_format")
}

// GetCmdServiceImportPath returns the import path of the cmd service (used by cmd/main.go).
func GetCmdServiceImportPath(name string) (string, error) {
	return getImportPath(name, "gk_cmd_service_path_format")
}

// GetEndpointImportPath returns the import path of the service endpoints.
func GetEndpointImportPath(name string) (string, error) {
	return getImportPath(name, "gk_endpoint_path_format")
}

// GetGRPCTransportImportPath returns the import path of the service grpc transport.
func GetGRPCTransportImportPath(name string) (string, error) {
	return getImportPath(name, "gk_grpc_path_format")
}

// GetPbImportPath returns the import path of the generated service grpc pb.
func GetPbImportPath(name, pathByFlag string) (string, error) {
	if pathByFlag != "" {
		return pathByFlag, nil
	}
	return getImportPath(name, "gk_grpc_pb_path_format")
}

// GetHTTPTransportImportPath returns the import path of the service http transport.
func GetHTTPTransportImportPath(name string) (string, error) {
	return getImportPath(name, "gk_http_path_format")
}

// GetDockerFileProjectPath returns the path of the project.
func GetDockerFileProjectPath() (string, error) {
	gosrc := GetGOPATH() + "/src/"
	gosrc = strings.Replace(gosrc, "\\", "/", -1)
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if viper.GetString("gk_folder") != "" {
		pwd += "/" + viper.GetString("gk_folder")
	}

	pwd = strings.Replace(pwd, "\\", "/", -1)
	projectPath := strings.Replace(pwd, gosrc, "", 1)

	return projectPath, nil
}

// GetGOPATH returns the gopath.
func GetGOPATH() string {
	if viper.GetString("GOPATH") != "" {
		return viper.GetString("GOPATH")
	}
	return defaultGOPATH()
}

func defaultGOPATH() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, "go")
		if filepath.Clean(def) == filepath.Clean(runtime.GOROOT()) {
			// Don't set the default GOPATH to GOROOT,
			// as that will trigger warnings from the go tool.
			return ""
		}
		return def
	}
	return ""
}

func getImportPath(name string, key string) (string, error) {
	modName, err := getModNameFromModFile(name)
	if err != nil {
		return "", err
	}

	gosrc := GetGOPATH() + "/src/"
	gosrc = strings.Replace(gosrc, "\\", "/", -1)
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if viper.GetString("gk_folder") != "" {
		pwd += "/" + viper.GetString("gk_folder")
	}

	pwd = strings.Replace(pwd, "\\", "/", -1)
	projectPath := strings.Replace(pwd, gosrc, "", 1)

	svcPath := fmt.Sprintf(viper.GetString(key), ToLowerSnakeCase(name))

	path := strings.Replace(svcPath, "\\", "/", -1)
	if modName != "" {
		modName = strings.Replace(modName, "\\", "/", -1)
		modNameArr := strings.Split(modName, "/")
		if len(modNameArr) <= 1 {
			projectPath = ""
		} else {
			projectPath = strings.Join(modNameArr[0:len(modNameArr)-1], "/")
		}
	}
	var importPath string
	// Change: here should not use os.Getwd() as projectPath
	// Desc:It can't pass go test, on windows, projectPath will be "c:/User/xxx/...", this will cause err certainly.
	projectPath = ""
	if projectPath == "" {
		importPath = path
	} else {
		importPath = projectPath + "/" + path
	}
	return importPath, nil
}

func getModNameFromModFile(name string) (string, error) {
	modFile := "go.mod"
	filePath := ToLowerSnakeCase(name) + "/" + modFile
	exists, _ := fs.Get().Exists(filePath)
	var modFileInParentLevel bool
	if exists == false {
		//if the service level has no go.mod file, it will check the parent level
		exists, err := fs.Get().Exists(modFile)
		if exists == false {
			return "", err
		}
		filePath = modFile
		modFileInParentLevel = true
	}

	content, err := fs.Get().ReadFile(filePath)
	if err != nil {
		return "", err
	}

	modDataArr := strings.Split(content, "\n")
	if len(modDataArr) != 0 {
		modNameArr := strings.Split(modDataArr[0], " ")
		if len(modNameArr) < 2 { // go.mod file: module XXXX/XXXX/{projectName}
			return "", nil
		}
		if modFileInParentLevel == true {
			return modNameArr[1] + "/" + name, nil
		}
		return modNameArr[1], nil
	}
	return "", nil
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

const (
	InterpreterBash = "bash"
	InterpreterSh   = "sh"
	InterpreterCmd  = "cmd"
)

func GetScriptInterpreter(_filepath string) (string, error) {
	absP, err := filepath.Abs(_filepath)
	if err != nil {
		return "", fmt.Errorf("utils.go: get absolutely path err:%v -->[%s]", err, _filepath)
	}

	if filepath.Ext(absP) == ".bat" {
		return InterpreterCmd, nil
	}

	fileInfo, _ := os.Stat(absP)
	if fileInfo.IsDir() {
		return "", fmt.Errorf("utils.go: _filepath is a dir")
	}

	f, _ := os.Open(absP)
	rd := bufio.NewReader(f)
	if line, err := rd.ReadString('\n'); err != nil {
		return "", fmt.Errorf("utils.go: read first line err:%v -->[%s]", err, _filepath)
	} else {
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, "\r\n")
		line = strings.TrimRight(line, " ")

		// First form: #!/usr/bin/env sh
		s2 := strings.Split(line, " ")

		// Second form: #!/bin/sh [or] #!/bin/bash
		s1 := strings.Split(line, "/")

		for _, ss := range [][]string{s1, s2} {
			if s := validUnixInterpreter(ss[len(ss)-1]); s != "" {
				return s, nil
			}
		}
		return "", fmt.Errorf("utils.go: invalid first line for [%s], you have to be set a valid shell interpreter, like #!/bin/bash", _filepath)
	}
}

func validUnixInterpreter(s string) string {
	switch s {
	case "bash":
		return InterpreterBash
	case "sh":
		return InterpreterSh
	}
	return ""
}

func GetCurrShellInterpreter() []string {
	shellSlice := []string{InterpreterBash, InterpreterSh, InterpreterCmd}

	var (
		err error
		ret []string
	)
	for _, s := range shellSlice {
		cmd := exec.Command(s)
		if err = cmd.Run(); err == nil {
			ret = append(ret, s)
		} else {
			logrus.Debugf("utils.go: we can't inspect what interpreter is on current env, err:%s", err)
		}
	}
	return ret
}

func CanScriptExecRightly(_filepath string) error {
	interpreterFromFile, err := GetScriptInterpreter(_filepath)
	if err != nil {
		return err
	}

	its := GetCurrShellInterpreter()
	for _, it := range its {
		if interpreterFromFile == it {
			return nil
		}
	}
	err = fmt.Errorf("utils.go: interpreter of existed script:%s is not be supported to current %+v", interpreterFromFile, its)
	return err
}

func ConvertToUnixPath(_path string) string {
	if _path == "" {
		return ""
	}
	if !strings.Contains(_path, "\\") {
		return _path
	}

	ss := strings.Split(_path, "\\")
	if !strings.HasPrefix(ss[0], "/") {
		ss[0] = "/" + ss[0]
		ss[0] = strings.TrimRight(ss[0], ":")
	}
	return strings.Join(ss, "/")
}
