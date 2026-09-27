package permissions

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type Permission string

type Permissions []Permission

func (p *Permissions) Scan(value any) error {
	if value == nil {
		*p = Permissions{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan Permissions: unsupported type %T", value)
	}

	return json.Unmarshal(bytes, p)
}

func (p Permissions) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func HasPermission(userPermissions Permissions, required Permission) bool {
	if slices.Contains(userPermissions, AllPermissions) {
		return true
	}

	requiredStr := string(required)

	for _, p := range userPermissions {
		pStr := string(p)

		if pStr == requiredStr {
			return true
		}

		_, haveSuffix := strings.CutSuffix(pStr, ".*")
		if haveSuffix {
			prefix := strings.TrimSuffix(pStr, ".*")
			if strings.HasPrefix(requiredStr, prefix) {
				return true
			}
		}
	}

	return false
}
