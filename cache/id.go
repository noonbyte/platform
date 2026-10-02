package cache

func (c *EntityCache[T]) ID(entity *T) string {
	return c.getID(entity)
}
