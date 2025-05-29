package utils

func FolderExists(folderPath string, stat StatFunc) bool {
	folderStat, err := stat(folderPath)
	if err != nil {
		return false
	}
	return folderStat.IsDir()
}
