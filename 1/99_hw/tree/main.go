package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func getFolders(entries []os.DirEntry) []os.DirEntry {
	var folders []os.DirEntry

	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry)
		}
	}
	return folders
}

func attachmentTree(
	out io.Writer, path string, printFiles bool, attachment string,
) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	if !printFiles {
		files = getFolders(files)
	}
	//tabs := strings.Repeat("\t", attachment)

	sort.Slice(files, func(i, j int) bool {
		return files[i].Name() < files[j].Name()
	})

	for i, file := range files {
		leftChar := '└'
		nextAttachment := attachment + "\t"
		if i < len(files)-1 {
			leftChar = '├'
			nextAttachment = attachment + "│\t"
		}
		writeText := fmt.Sprintf("%s%c───%s", attachment, leftChar, file.Name())
		if file.IsDir() {
			_, _ = out.Write([]byte(writeText + "\n"))
			err = attachmentTree(
				out, path+string(os.PathSeparator)+file.Name(), printFiles, nextAttachment)
			if err != nil {
				return err
			}
		} else {
			fileInfo, err := file.Info()
			if err != nil {
				return err
			}
			size := fileInfo.Size()
			if size == 0 {
				writeText = fmt.Sprintf("%s (empty)\n", writeText)
			} else {
				writeText = fmt.Sprintf("%s (%db)\n", writeText, size)
			}
			_, _ = out.Write([]byte(writeText))
		}
	}
	return nil
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	return attachmentTree(out, path, printFiles, "")
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
