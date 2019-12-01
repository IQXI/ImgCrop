package services

import (
	"ImgCrop/downloader"
	"ImgCrop/internal/Imaging"
	"ImgCrop/internal/lru_cache"
	"ImgCrop/internal/structs"
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"
)

type Mux struct {
	Logger     *zap.Logger
	Config     structs.Config
	Imaging    Imaging.Imagimg
	LRUCache   lru_cache.LRUCache
	Downloader downloader.Downloader
}

func NewMux(Logger *zap.Logger, Config structs.Config) (Mux, error) {
	ig, err := Imaging.NewImaging(Logger, Config)
	if err != nil {
		Logger.Error("Error in NewImaging func", zap.Error(err))
		return Mux{}, err
	}

	cache, err := lru_cache.NewLRUCache(Config.Cache.Size, Logger, Config)
	if err != nil {
		Logger.Error("Error in NewLRUCache func", zap.Error(err))
		return Mux{}, err
	}

	return Mux{
		Logger:     Logger,
		Config:     Config,
		Imaging:    ig,
		LRUCache:   cache,
		Downloader: downloader.NewDownloader(Logger, Config),
	}, nil
}

func (m *Mux) CacheChecker(w http.ResponseWriter, r *http.Request) {
	url := mux.Vars(r)["url"]
	_, err := m.LRUCache.Get(url)
	if err != nil {
		http.Error(w, "Image not found in Cache", http.StatusNotFound)
	} else {
		http.Error(w, "Image found in Cache", http.StatusOK)
	}
}

func (m *Mux) CropHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var addToCacheResult bool = true

	width, err := strconv.Atoi(vars["width"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	height, err := strconv.Atoi(vars["height"])
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	url := vars["url"]
	_ = url
	m.Logger.Info(fmt.Sprintf("width: %v %T\n", vars["width"], vars["width"]))
	m.Logger.Info(fmt.Sprintf("height: %v %T\n", vars["height"], vars["height"]))
	m.Logger.Info(fmt.Sprintf("url: %v %T\n", vars["url"], vars["url"]))

	image, err := m.LRUCache.Get(url)
	var statusCode int
	if err != nil {
		// если мы тут, то не наши в кеше такого изобрадения
		m.Logger.Info(fmt.Sprintf("Image not in cache, Link %s", url))

		// тут должны добавить его к кеш
		statusCode, image, err = m.Downloader.DownloadImage(url)
		if err != nil {
			_ = statusCode
			http.Error(w, err.Error(), statusCode)
			return
		}
		err = m.LRUCache.Add(image)
		if err != nil {
			addToCacheResult = false
		} else {
			addToCacheResult = true
		}
		if addToCacheResult {
			m.Logger.Info(fmt.Sprintf("Successfully add to cache, Link %s", url))
		} else {
			m.Logger.Info(fmt.Sprintf("Faild add to cache, Link %s", url))
		}
	}

	cropped_image, err := m.Imaging.CropImage(image, width, height)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buf := new(bytes.Buffer)
	switch ext := image.Exctension; ext {
	case "tiff":
		err = tiff.Encode(buf, cropped_image, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	case "jpeg":
		err = jpeg.Encode(buf, cropped_image, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	case "png":
		err = png.Encode(buf, cropped_image)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	case "gif":
		err = gif.Encode(buf, cropped_image, nil)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	case "bmp":
		err = bmp.Encode(buf, cropped_image)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("Unknown extension %s", ext), 500)
		return
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if addToCacheResult == false {
		go func() {
			err := os.Remove(image.Path)
			if err != nil {
				m.Logger.Error(fmt.Sprintf("I can't remove file %s in path %s", image.FileName, image.Path))
			}
		}()
	}

	return
}
