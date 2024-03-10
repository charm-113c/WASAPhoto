package api

// Close should close everything opened in the lifecycle of the `_router`; for example, background goroutines.
func (rt *_router) Close() error {
	// for the time being this jut means disabling user actions, i.e. nullifying their auth tokens
	// err := rt.db.NullifyToken()
	return nil
}
