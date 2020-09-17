package app

import (
	"github.com/anacrolix/torrent/metainfo"
	"github.com/diericx/iceetime/internal/pkg/torrent"
)

// TODO: input from config file
const DefaultResolution = "iw:ih"
const DefaultMaxBitrate = "50M"

func GetDefaultTorrentMeta() TorrentMeta {
	return TorrentMeta{
		RatioToStop: 1,
		HoursToStop: 336,
		IsStopped:   false,
	}

}

// These functions act as const arrays because go doesn't allow const arrays... I know pretty fucked up

// GetSupportedVideoFileFormats returns an array of strings that are the supported video formats
func GetSupportedVideoFileFormats() []string {
	return []string{".mkv", ".mp4"}
}

// GetBlacklistedFileNameContents returns an array of strings that are blacklisted from torrent names
func GetBlacklistedFileNameContents() []string {
	return []string{"sample"}
}

// TODO: These languages are only blacklisted because it's hard to support
func GetBlacklistedTorrentNameContents() []string {
	return []string{"fre", "french", "ita", "italian"}
}

// BasicAuth info for basic auth http requests
type BasicAuth struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

type TorrentMeta struct {
	InfoHash     string `storm:"id"`
	RatioToStop  float32
	MinutesAlive int
	HoursToStop  int
	IsStopped    bool
}

// TorrentFile represents a file in a torrent
type TorrentFile struct {
	Path string
	Size int64
}

// Release is a potential torrent for a specific piece of media. We use this info to decide whether or not we want to
// actually grab the torrent.
type Release struct {
	ImdbID      string
	Title       string
	Size        int64
	Link        string
	LinkAuth    *BasicAuth
	InfoHash    string
	Grabs       int
	Seeders     int
	MinRatio    float32
	MinSeedTime int
}

type TorrentClientConfig struct {
	MinSeeders                        int    `toml:"min_seeders"`
	TorrentInfoTimeout                int    `toml:"torrent_info_timeout"`
	TorrentFilePath                   string `toml:"torrent_file_path"`
	TorrentDataPath                   string `toml:"torrent_data_path"`
	TorrentHalfOpenConnsPerTorrent    int    `toml:"torrent_half_open_conns_per_torrent"`
	TorrentEstablishedConnsPerTorrent int    `toml:"torrent_established_conns_per_torrent"`
}

// Indexer is info we need to hit an indexer for a list of torrents
type Indexer struct {
	Name                 string     `toml:"name"`
	URL                  string     `toml:"url"`
	BasicAuth            *BasicAuth `toml:"basic_auth"`
	SupportsImdbIDSearch bool       `toml:"supports_imdb_id_search"`
	APIKey               string     `toml:"api_key"`
	Categories           string     `toml:"categories"`
}

// Quality contains specifications for a specific quality of torrent and how to infer that quality from a name
type Quality struct {
	Name       string  `toml:"name"`
	Regex      string  `toml:"regex"`
	MinSize    float64 `toml:"min_size"`
	MaxSize    float64 `toml:"max_size"`
	Resolution string  `toml:"resolution"`
}

type TranscoderConfig struct {
	Video struct {
		Format          string `yaml:"format"`
		CompressionAlgo string `yaml:"compression_algo"`
	} `yaml:"video"`
	Audio struct {
		CompressionAlgo string `yaml:"compression_algo"`
	} `yaml:"audio"`
}

// MovieTorrentLink handles linking a Movie to a specific file in a torrent
type MovieTorrentLink struct {
	ID              int `storm:"id,increment"`
	ImdbID          string
	TorrentInfoHash string
	FileIndex       int
}

type TorrentMetaRepo interface {
	Store(TorrentMeta) error
	GetByInfoHashStr(string) (TorrentMeta, error)
	RemoveByInfoHashStr(hashStr string) error
}

type ReleaseRepo interface {
	GetForMovie(imdbID string, title string, year string, minQuality int) ([]Release, error)
}

type MovieTorrentLinkRepo interface {
	Store(MovieTorrentLink) (*MovieTorrentLink, error)
	GetByImdbID(imdbID string) ([]MovieTorrentLink, error)
}

type TorrentClient interface {
	Close()
	AddMagnet(string) (torrent.Torrent, error)
	AddFile(string) (torrent.Torrent, error)
	Torrents() []torrent.Torrent
	Torrent(metainfo.Hash) (torrent.Torrent, bool)
}