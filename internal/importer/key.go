package importer

import "fmt"

func BuildKey(
	portal string,
	input string,
	portalType string,
) string {

	return fmt.Sprintf(
		"%s:%s:%s",
		portal,
		input,
		portalType,
	)
}