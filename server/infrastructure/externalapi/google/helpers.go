package google

import (
	"fmt"
	"strings"
)

const application = "applications/"

func (e *external) initWebTokenParent(enterpriseID string) string {
	return fmt.Sprintf("enterprises/%s", enterpriseID)
}

func (e *external) initApplicationFullName(enterpriseID, packageName string) string {
	return fmt.Sprintf("enterprises/%s/applications/%s", enterpriseID, packageName)
}

func (e *external) parserApplicationName(name string) string {
	first := strings.Index(name, application)
	if first == -1 {
		return ""
	}

	posFirstAdjusted := first + len(application)

	return name[posFirstAdjusted:]
}
