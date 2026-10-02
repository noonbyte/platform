package cache

import (
	"fmt"

	"github.com/noonbyte/platform/specs"
)

func (c *EntityCache[T]) Key(id string) string {
	return fmt.Sprintf("%s:%s", c.prefix, id)
}

func (c *EntityCache[T]) ListKey(specs specs.Specs) string {
	discriminator := specs.String()
	if discriminator == nil {
		return fmt.Sprintf("%s:list", c.prefix)
	}

	return fmt.Sprintf("%s:list:%s", c.prefix, *discriminator)
}
