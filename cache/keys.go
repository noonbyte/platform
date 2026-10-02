package cache

import (
	"fmt"

	"github.com/noonbyte/platform/infrastructure/rdb"
)

func (c *EntityCache[T]) Key(id string) string {
	return fmt.Sprintf("%s:%s", c.prefix, id)
}

func (c *EntityCache[T]) ListKey(specs rdb.Specs) string {
	discriminator := rdb.GenerateKeyWithSpecs(specs)
	if discriminator == nil {
		return fmt.Sprintf("%s:list", c.prefix)
	}

	return fmt.Sprintf("%s:list:%s", c.prefix, *discriminator)
}
