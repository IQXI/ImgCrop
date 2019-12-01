package main

import (
	cfg "ImgCrop/internal/config"
	lg "ImgCrop/internal/logger"
	m "ImgCrop/internal/services"
	"ImgCrop/internal/structs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func rmDirAll(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func initialization(config structs.Config) bool {
	//чекаем папку логов
	if config.Logger.LogsPath != "" {
		if _, err := os.Stat(config.Logger.LogsPath); os.IsNotExist(err) {
			e := os.Mkdir(config.Logger.LogsPath, 0777)
			if e != nil {
				log.Fatal(e.Error())
				return false
			}
		}
	}
	//чекаем папку кеша
	if config.Cache.Path != "" {
		if _, err := os.Stat(config.Cache.Path); os.IsNotExist(err) {
			e := os.Mkdir(config.Cache.Path, 0777)
			if e != nil {
				log.Fatal(e.Error())
				return false
			}
		} else {
			e := rmDirAll(config.Cache.Path)
			if e != nil {
				log.Fatal(e.Error())
				return false
			}
		}
	}

	return true
}

func main() {
	config, err := cfg.GetConfig()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	cacheSize := os.Getenv("CACHE_SIZE")
	if cacheSize != "" {
		intCacheSize, err := strconv.Atoi(cacheSize)
		if err != nil {
			log.Fatal(err.Error())
			return
		}

		config.Cache.Size = int64(intCacheSize)
	}

	ok := initialization(config)
	if !ok {
		return
	}

	logger := lg.GetLogger(config)
	logger.Info("Service loading!")
	newMux, err := m.NewMux(logger, config)
	if err != nil {
		logger.Fatal(err.Error())
	}

	r := mux.NewRouter()
	r.HandleFunc("/crop/{width}/{height}/{url:(?:.+)}", newMux.CropHandler)
	r.HandleFunc("/cache/{url:(?:.+)}", newMux.CacheChecker)
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8008", nil))

}
