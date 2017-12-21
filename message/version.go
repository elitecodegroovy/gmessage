package message

import (
	"fmt"
	"runtime"
)

func VersionInfo(app string) string {
	return fmt.Sprintf(" %s v%s (built w/%s)", app, VERSION, runtime.Version())
}
