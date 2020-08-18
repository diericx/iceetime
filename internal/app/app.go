package app

import (
	"log"
	"time"
)

type Error struct {
	OrigionalError error
	Code           int
	Message        string
}

func NewError(origionalError error, code int, message string) *Error {
	log.Printf("%s\n%s", message, origionalError)
	return &Error{
		OrigionalError: origionalError,
		Code:           code,
		Message:        message,
	}
}

func (e Error) Error() string {
	return e.Message
}

type MediaMeta struct {
	Title       string `json:"title"`
	ReleaseDate string `json:"release_date"`
}

type Quality struct {
	Name    string `yaml:"name"`
	Regex   string `yaml:"regex"`
	MinSize int64  `yaml:"minSize"`
	MaxSize int64  `yaml:"maxSize"`
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Indexer struct {
	Name                 string     `yaml:"name"`
	URL                  string     `yaml:"url"`
	BasicAuth            *BasicAuth `yaml:"basicAuth"`
	SupportsImdbIDSearch bool       `yaml:"supportsImdbIDSearch"`
	APIKey               string     `yaml:"apiKey"`
	Categories           string     `yaml:"categories"`
}

// Torrent metadata for a certain torrent
type Torrent struct {
	ID          int    `storm:"id,increment" json:"id"`
	Type        string // Movie, Episode, Season, Season Pack, etc.
	ImdbID      string `storm:"unique" json:"imdbID"`
	Title       string `json:"title"` // Note: Quality is inferred from this
	Size        int64  `json:"size"`
	InfoHash    string `storm:"unique" json:"infoHash"`
	Grabs       int    `json:"grabs"`
	Link        string
	Seeders     int `json:"seeders"` // Note: subject to change
	Tracker     string
	MinRatio    float32
	MinSeedTime int
	CreatedAt   time.Time `json:"createdAt"`
}

type Tmdb struct {
	ApiKey string `yaml:"apiKey"`
}

type Config struct {
	Indexers  []Indexer `yaml:"indexers"`
	Qualities []Quality `yaml:"qualities"`
	Tmdb      Tmdb      `yaml:"tmdb"`
}

type TorrentDAO interface {
	Save(Torrent)
	GetByImdbIDAndMinQuality(imdbID string, minQuality int)
	GetByImdbIDAndInfoHash(imdbID string, infoHash string)
	GetByID(id int)
}

type TorrentClient interface {
	AddFromMagnet(magnet string) (hash string, err error)
	AddFromFileURL(fileURL string, name string) (hash string, err error)
	RemoveByHash(hash string) error
}

// IndexerQueryHandler given inputs will handle querying indexers for torrents
type IndexerQueryHandler interface {
	QueryMovie(imdbID string, title string, year string, minQuality Quality) (*Torrent, *Error)
}

type MediaMetaManager interface {
	GetByImdbID(imdbID string) (*MediaMeta, error)
}

type Transcoder interface{}
