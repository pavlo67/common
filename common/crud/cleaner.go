package crud

type Cleaner interface {
	Clean(*Options) error
}
