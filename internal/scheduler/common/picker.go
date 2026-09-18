package common

type Picker interface {
	PickAndReserve() *Backend
	Release(b *Backend)
	SetHealth(b *Backend, status bool)
	Snapshot()
}
