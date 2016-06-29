package archive

import (
	"fmt"
	"framework/util/shell"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func checkFolder(path string) {
	_, err := os.Stat(path)
	if !(err == nil || os.IsExist(err)) {
		fmt.Println("Create Folder")
		os.MkdirAll(path, 0777)
	}
}

// 将文件夹解压到filePath下面
func ArchiveBufferUnderPath(content string, filePath string) error {
	folerPath := filepath.Dir(filePath)
	tmpFileName := strconv.Itoa(rand.Int())
	tmpFilePath := filepath.Join(folerPath, tmpFileName)
	// write content to file
	f, err := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	f.Write([]byte(content))
	f.Close()

	// uncompress file
	_, _, err = shell.RunShell(folerPath, "tar", "-zxvf", tmpFileName)
	if err != nil {
		return err
	}
	_, _, err = shell.RunShell(folerPath, "rm", tmpFileName)
	if err != nil {
		return err
	}
	return err
}

// 将文件解压到filePath下
func ArchiveBufferToPath(content string, filePath string) error {
	folerPath := filepath.Dir(filePath)
	fmt.Println(folerPath)
	tmpFileName := strconv.Itoa(rand.Int())
	tmpFilePath := filepath.Join(folerPath, tmpFileName)
	// write content to file
	f, err := os.OpenFile(tmpFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	f.Write([]byte(content))
	f.Close()

	// uncompress file
	_, _, err = shell.RunShell(folerPath, "tar", "-zxvf", tmpFileName)
	if err != nil {
		return err
	}

	fmt.Println("remote: ", tmpFilePath)
	os.Remove(tmpFilePath)
	return err
}
