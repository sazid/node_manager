package app

import "io/fs"

func FileExistsInDir(dirEntries []fs.DirEntry, fileName string) bool {
	for _, f := range dirEntries {
		if f.Name() == fileName {
			return true
		}
	}
	return false
}
