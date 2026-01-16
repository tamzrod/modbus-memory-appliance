// internal/rawingest/resolver_func.go
// PURPOSE: Functional adapter for MemoryResolver.
// ALLOWED: type adapter only
// FORBIDDEN: logic beyond delegation

package rawingest

type MemoryResolverFunc func(id uint16) (RawWritableMemory, bool)

func (f MemoryResolverFunc) ResolveMemoryByID(id uint16) (RawWritableMemory, bool) {
	return f(id)
}
