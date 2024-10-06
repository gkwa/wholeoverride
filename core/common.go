package core

import "fmt"

func formatImage(name, url string, isRemote bool) string {
	if isRemote {
		return fmt.Sprintf("![%s](%s)", name, url)
	}
	return fmt.Sprintf("![[%s]]", url)
}
