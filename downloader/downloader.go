package downloader

import (
	"ImgCrop/internal/structs"
	"crypto/md5"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Downloader struct {
	Logger *zap.Logger    `mapstructure:"logger"`
	Config structs.Config `mapstructure:"config"`
}

func NewDownloader(Logger *zap.Logger, Config structs.Config) Downloader {
	return Downloader{
		Logger: Logger,
		Config: Config,
	}
}

func (d *Downloader) DownloadImage(originalLink string) (statusCode int, image structs.Image, err error) {

	link := originalLink
	// если в пути нет http, то добавляем его
	u, _ := url.Parse(originalLink)
	if u.Scheme == "" {
		err := fmt.Errorf("URL uncorrect)", originalLink)
		return 400, structs.Image{}, err
	}
	if u.Host == "" {
		link = fmt.Sprintf("%v:/%v", u.Scheme, u.Path)
	}

	response, err := http.Get(link)
	if err != nil {
		d.Logger.Error(err.Error())
		return 500, structs.Image{}, err
	}
	//делаем разбор что пошло не так
	if response.StatusCode != 200 {
		return response.StatusCode, structs.Image{}, errors.New(response.Status)
	}
	var ext string = "jpg"
	if !strings.Contains(response.Header.Get("Content-Type"), "image") {
		d.Logger.Error("Extension is BAD")
		return 400, structs.Image{}, errors.New(("Extension is BAD (Accepted only TIFF, JPG, PNG, GIF)"))
	} else {
		ext = strings.Split(response.Header.Get("Content-Type"), "/")[1]
	}

	imageName := fmt.Sprintf("%x.%s", md5.Sum([]byte(response.Request.URL.Path)), ext)

	path := d.Config.Cache.Path + imageName

	file, err := os.Create(path)
	if err != nil {
		d.Logger.Error(err.Error())
		return 500, structs.Image{}, err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		d.Logger.Error(err.Error())
		return 500, structs.Image{}, err
	}
	d.Logger.Info("Success Downloaded!")
	fileStat, err := file.Stat()
	if err != nil {
		d.Logger.Error(err.Error())
		return 500, structs.Image{}, err
	}
	headers := make(map[string]string)

	for key, value := range response.Header {
		headers[key] = strings.Join(value, ";")
	}
	defer func() {
		err = response.Body.Close()
		if err != nil {
			d.Logger.Fatal(err.Error())
		}
		err = file.Close()
		if err != nil {
			d.Logger.Fatal(err.Error())
		}
	}()

	image = structs.Image{
		Size:       fileStat.Size(),
		Path:       path,
		Headers:    headers,
		Link:       originalLink,
		FileName:   imageName,
		Exctension: ext,
	}
	return 200, image, nil
}
