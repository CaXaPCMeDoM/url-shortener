package storage

type Storage interface {
	SaveUrl(urlAbsolute string, alias string) error
	GetUrl(alias string) (string, error)
	GetAlias(urlAbsolute string) (string, error)
	CheckAliasURLExists(alias string) (bool, error)
}
