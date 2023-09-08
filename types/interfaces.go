package types

type IAsyncService interface {
	Name() string
	Run() error
}
