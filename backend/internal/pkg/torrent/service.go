package torrent

import (
	"time"

	"github.com/anacrolix/torrent"
	"github.com/asdine/storm"
)

const DefaultRatioToStopAt = 1.0
const DefaultHoursToStopAt = 336
const DefaultIsStopped = false

type TorrentService struct {
	Timeout          time.Duration
	Client           *torrent.Client
	DB               *storm.DB
	TorrentsLocation string
}

type TorrentMeta struct {
	InfoHash      string `storm:"id"`
	MinutesAlive  int
	RatioToStopAt float32
	HoursToStopAt int
	IsStopped     bool
}

// func (s *TorrentService) LoadTorrentFilesFromCache() error {
// }

// // AddFromURLUknownScheme will add the torrent if it is a magnet url, will download a file if it's a
// // file or recursicely follow a redirect
// func (c *TorrentClient) AddFromURLUknownScheme(rawURL string, auth *app.BasicAuth) (*Torrent, error) {
// 	u, err := url.Parse(rawURL)
// 	if err != nil {
// 		return "", err
// 	}
// 	if u.Scheme == "magnet" {
// 		return c.AddFromMagnet(rawURL)
// 	}

// 	// Attempt to make http/s call
// 	req, err := http.NewRequest("GET", rawURL, nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	if auth != nil {
// 		req.SetBasicAuth(auth.Username, auth.Password)
// 	}
// 	client := new(http.Client)
// 	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
// 		return errors.New("Redirect")
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		if resp.StatusCode == http.StatusFound { //status code 302
// 			url, err := resp.Location()
// 			if err != nil {
// 				return "", err
// 			}
// 			return c.AddFromURLUknownScheme(url.String(), auth)
// 		}
// 		return "", err
// 	}

// 	tempFilePath := fmt.Sprintf("%s/%s", c.torrentFilePath, RandomString(10))
// 	err = downloadFileFromResponse(resp, tempFilePath)
// 	defer os.Remove(tempFilePath)
// 	if err != nil {
// 		return "", err
// 	}

// 	return c.AddFromFile(tempFilePath)
// }

// func (s *TorrentService) cachedMetaInfo(infoHash metainfo.Hash) (ret *metainfo.MetaInfo) {
// 	file := filepath.Join(s.TorrentsLocation, fmt.Sprintf("%s.torrent", infoHash.HexString()))
// 	ret, err := metainfo.LoadFromFile(file)
// 	if err != nil {
// 		ret = nil
// 		return
// 	}
// 	if ret.HashInfoBytes() != infoHash {
// 		ret = nil
// 		return
// 	}
// 	return
// }

// func (s *TorrentService) saveTorrentFile(t *torrent.Torrent) (err error) {
// 	file := filepath.Join(s.TorrentsLocation, fmt.Sprintf("%s.torrent", t.InfoHash().HexString()))
// 	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0660)
// 	if err != nil {
// 		return
// 	}
// 	defer f.Close()
// 	return t.Metainfo().Write(f)
// }

// func (s *TorrentService) getOrCreateMetaForTorrent(t *app.Torrent) (*TorrentMeta, error) {
// 	infoHashStr := t.InfoHash.HexString()
// 	var meta TorrentMeta
// 	err := s.DB.One("InfoHash", infoHashStr, &meta)
// 	if err != nil {
// 		// Not found err, save a new one
// 		if err == storm.ErrNotFound {
// 			meta = TorrentMeta{
// 				InfoHash:      infoHashStr,
// 				RatioToStopAt: DefaultRatioToStopAt,
// 				HoursToStopAt: DefaultHoursToStopAt,
// 			}
// 			err := s.DB.Save(&meta)
// 			if err != nil {
// 				return nil, err
// 			}

// 			return &meta, nil
// 		}

// 		// Some other err, return it
// 		return nil, err
// 	}

// 	return &meta, nil
// }

// func (s *TorrentService) updateMetaForTorrent(meta TorrentMeta) error {
// 	return s.DB.Update(&meta)
// }

// func anacrolixTorrentToApp(t *torrent.Torrent) app.Torrent {
// 	return app.Torrent{
// 		InfoHash:       t.InfoHash(),
// 		Stats:          t.Stats(),
// 		Length:         t.Length(),
// 		BytesCompleted: t.BytesCompleted(),
// 		Name:           t.Name(),
// 		Seeding:        t.Seeding(),
// 	}
// }

// func downloadFileFromResponse(resp *http.Response, filePath string) error {
// 	// Get the data
// 	if resp.StatusCode != 200 {
// 		return fmt.Errorf("couldn't reach file server with code: %v", resp.StatusCode)
// 	}
// 	defer resp.Body.Close()

// 	// Create the file
// 	out, err := os.Create(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	// Write the body to file
// 	_, err = io.Copy(out, resp.Body)
// 	return err
// }