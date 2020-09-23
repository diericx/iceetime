package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/diericx/iceetime/internal/app"
	"github.com/diericx/iceetime/internal/app/http"
	"github.com/diericx/iceetime/internal/app/repos/jackett"
	"github.com/diericx/iceetime/internal/app/repos/storm"
	"github.com/diericx/iceetime/internal/app/services"

	"github.com/diericx/iceetime/internal/pkg/torrent"
)

type tomlConfig struct {
	Indexers      []app.Indexer           `toml:"indexers"`
	Qualities     []app.Quality           `toml:"qualities"`
	Transcoder    app.TranscoderConfig    `toml:"transcoder"`
	TmdbAPIKey    string                  `toml:"tmdb_api_key"`
	TorrentClient app.TorrentClientConfig `toml:"torrent_client"`
}

func main() {
	var conf tomlConfig
	if _, err := toml.DecodeFile(os.Getenv("CONFIG_FILE"), &conf); err != nil {
		panic(err)
	}
	// TODO: real check function on configuration
	if conf.TorrentClient.MetaRefreshRate < 1 {
		panic(errors.New("Torrent client refresh rate cannot be less than 1"))
	}
	log.Printf("%+v", conf)

	stormDBFilePath := filepath.Join(conf.TorrentClient.TorrentFilePath, ".iceetime.storm.db")
	stormDB, err := storm.OpenDB(stormDBFilePath)
	if err != nil {
		log.Fatalf("Couldn't open torrent file at %s. The file will be created if it doesn't exist, make sure the directory exists and user has proper permissions.", stormDBFilePath)
	}
	defer stormDB.Close()

	client, err := torrent.NewClient(conf.TorrentClient.TorrentFilePath, conf.TorrentClient.TorrentDataPath, 15, 30, 30)
	if err != nil {
		log.Fatalf("Couldn't start torrent client: %s", err.Error())
	}
	defer client.Close()

	//
	// Initialize repos
	//
	torrentMetaRepo := storm.TorrentMeta{
		DB: stormDB,
	}

	movieTorrentLinkRepo := storm.MovieTorrentLink{
		DB: stormDB,
	}

	releaseRepo := jackett.ReleaseRepo{
		Qualities: conf.Qualities,
		Indexers:  conf.Indexers,
	}

	//
	// Initialize services
	//
	torrentService := services.Torrent{
		Client:           client,
		TorrentMetaRepo:  &torrentMetaRepo,
		GetInfoTimeout:   time.Second * 15,
		MinSeeders:       conf.TorrentClient.MinSeeders,
		TorrentFilesPath: conf.TorrentClient.TorrentFilePath,
	}

	// Add and start (if running) torrents on disk
	torrentService.AddTorrentsOnDisk()
	err = torrentService.StartTorrentsAccordingToMetadata()
	if err != nil {
		log.Fatalf("Error starting existing torrents on disk: %s", err)
	}
	// Start maintinence thread which keeps meta up to date
	go func() {
		for {
			err := torrentService.UpdateMetaForAllTorrents()
			if err != nil {
				log.Fatalf("Error updating metadata for torrents: %s", err)
			}
			time.Sleep(time.Duration(conf.TorrentClient.MetaRefreshRate) * time.Second)
		}
	}()

	releaseService := services.Release{
		ReleaseRepo: releaseRepo,
		Qualities:   conf.Qualities,
	}

	torrentLinkService := services.TorrentLink{
		MovieTorrentLinkRepo: movieTorrentLinkRepo,
	}

	// TODO: Input from config file
	transcoder := services.Transcoder{
		Config: conf.Transcoder,
	}

	httpHandler := http.HTTPHandler{
		TorrentService:     torrentService,
		ReleaseService:     releaseService,
		TorrentLinkService: torrentLinkService,
		Transcoder:         transcoder,
		Qualities:          conf.Qualities,
		TorrentFilesPath:   conf.TorrentClient.TorrentFilePath,
	}

	httpHandler.Serve("secret-todo")
}
