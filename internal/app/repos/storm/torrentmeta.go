package storm

import (
	"github.com/asdine/storm"
	"github.com/diericx/iceetime/internal/app"
)

type TorrentMeta struct {
	DB *storm.DB
}

func (r *TorrentMeta) Store(meta app.TorrentMeta) error {
	return r.DB.Save(&meta)
}

func (r *TorrentMeta) GetByInfoHashStr(infoHashStr string) (app.TorrentMeta, error) {
	var meta app.TorrentMeta
	err := r.DB.One("InfoHash", infoHashStr, &meta)
	return meta, err
}

func (r *TorrentMeta) RemoveByInfoHashStr(hashStr string) error {
	err := r.DB.DeleteStruct(app.TorrentMeta{InfoHash: hashStr})
	return err
}