package storage

type Storage interface {
	Update(string, string, string) error
	UpdateV2(Metric) error
	UpdateBatch([]Metric) error
	Get(string, string) (string, error)
	GetV2(Metric) (*Metric, error)
	GetAll() ([]byte, error)
	Save(string) error
	Load(string) error
}
