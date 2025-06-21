package utils

func FolderExists(folderPath string, stat StatFunc) bool {
	folderStat, err := stat(folderPath)
	if err != nil {
		return false
	}
	return folderStat.IsDir()
}

func FileExists(filePath string, stat StatFunc) bool {
	_, err := stat(filePath)
	if err != nil {
		return false
	}
	return true
}
