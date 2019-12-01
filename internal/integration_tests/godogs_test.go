package main

import (
	"ImgCrop/downloader"
	lg "ImgCrop/internal/logger"
	"ImgCrop/internal/structs"
	"fmt"
	"github.com/DATA-DOG/godog"
	"go.uber.org/zap"
	"net/http"
)

//You can implement step definitions for undefined steps with these snippets:

type TestImgCropService struct {
	Downloader downloader.Downloader
	Config     structs.Config
	Logger     *zap.Logger
	Code       int
	CodeList   List
}
type List struct {
	Value int
	Next  *List
}

func NewTestImgCropService(logger *zap.Logger, config structs.Config) TestImgCropService {
	CodeList := List{
		Value: 500,
		Next: &List{
			Value: 404,
			Next: &List{
				Value: 400,
				Next: &List{
					Value: 503,
					Next:  nil,
				},
			},
		},
	}

	return TestImgCropService{
		Downloader: downloader.NewDownloader(logger, config),
		Config:     structs.Config{},
		Logger:     logger,
		Code:       0,
		CodeList:   CodeList,
	}
}

func (tics *TestImgCropService) addImageToCache(url string) bool {
	code, _, err := tics.Downloader.DownloadImage(url)
	if err != nil {
		tics.Logger.Fatal("can't create lru_cache.NewLRUCache")
	}
	if code == 200 {
		return true
	}
	return false
}

func (tics *TestImgCropService) clientMakeGet_imageRequestToImage_serverViaImgCropService() error {
	code, _, err := tics.Downloader.DownloadImage("http://imgcrop:8008/crop/300/400/http:/nginx/images/png/2.png")
	if err != nil {
		return err
	}
	if code == 200 {
		return nil
	}
	return fmt.Errorf(fmt.Sprintf("Server return non 200 code: %d", code))
}

func (tics *TestImgCropService) imgCropFindImageInCacheAndReturnThemToUser() error {
	resp, err := http.Get("http://imgcrop:8008/cache/http:/nginx/images/png/2.png")
	if err == nil && resp.StatusCode == 200 {
		return nil
	} else {
		return fmt.Errorf("Status code non 200 - Code: %d", resp.StatusCode)
	}
}

func (tics *TestImgCropService) imgCropServiceMakeGet_imageRequestToImage_server1() error {
	code, _, err := tics.Downloader.DownloadImage("http://imgcrop:8008/crop/300/400/http:/somestrangelocalhost/images/png/2.png")
	if err != nil && code == 500 {

		return nil
	} else {
		return fmt.Errorf("Server founded!")
	}
}

func (tics *TestImgCropService) imgCropServiceMakeGet_imageRequestToImage_server2() error {
	code, _, _ := tics.Downloader.DownloadImage("http://imgcrop:8008/crop/300/400/http://nginx/images/jpg/2.png")
	if code == 404 {
		tics.Code = code
		return nil
	} else {
		return fmt.Errorf("Status code non 404 - Code: %d", code)
	}
}

func (tics *TestImgCropService) imgCropShouldReturnCodeToUser(code int) error {
	tics.Logger.Info(fmt.Sprint(tics.CodeList.Value, code))
	if tics.CodeList.Value == code {
		if tics.CodeList.Next != nil {
			tics.CodeList = *tics.CodeList.Next
		}
		return nil
	} else {
		err := fmt.Errorf("%d != %d", tics.CodeList.Value, code)
		if tics.CodeList.Next != nil {
			tics.CodeList = *tics.CodeList.Next
		}

		return err
	}
}

func (tics *TestImgCropService) imgCropServiceMakeGet_fileRequestToFile_server() error {
	code, _, _ := tics.Downloader.DownloadImage("http://imgcrop:8008/crop/300/400/http://nginx/files/telnet.exe")
	if code == 400 {
		tics.Code = code
		return nil
	} else {
		return fmt.Errorf("Status code non 400 - Code: %d", code)
	}
}

func (tics *TestImgCropService) imgCropServiceMakeGet_imageRequestToImage_serverAndReturnErr_codeAndErrorToImgCrop() error {
	code, _, _ := tics.Downloader.DownloadImage("http://imgcrop:8008/crop/300/400/http://nginx/error/")
	if code == 503 {
		tics.Code = code
		return nil
	} else {
		return fmt.Errorf("Status code non 503 - Code: %d", code)
	}
}

func (tics *TestImgCropService) moveScenario(interface{}) {
	fmt.Printf("Code was: %v\n", tics.CodeList.Value)
	tics.CodeList = *tics.CodeList.Next
}

func FeatureContext(s *godog.Suite) {
	config := structs.Config{
		Logger: structs.LoggerConfig{
			Level:    "INFO",
			LogsPath: "../../logs/",
			FileName: "godog_test.log",
			Name:     "ImgCrop",
		},
		Cache: structs.CacheConfig{
			Path: "../../files/",
			Size: 100,
		},
	}
	logger := lg.GetLogger(config)
	tics := NewTestImgCropService(logger, config)

	//s.BeforeScenario(tics.moveScenario)

	s.Step(`^Client make get_image request to image_server via ImgCrop service$`, tics.clientMakeGet_imageRequestToImage_serverViaImgCropService)
	s.Step(`^ImgCrop find image in cache and return them to user$`, tics.imgCropFindImageInCacheAndReturnThemToUser)
	s.Step(`^ImgCrop service make get_image request to image_server1$`, tics.imgCropServiceMakeGet_imageRequestToImage_server1)
	s.Step(`^ImgCrop service make get_image request to image_server2$`, tics.imgCropServiceMakeGet_imageRequestToImage_server2)
	s.Step(`^ImgCrop should return (\d+) code to User$`, tics.imgCropShouldReturnCodeToUser)
	s.Step(`^ImgCrop service make get_file request to file_server$`, tics.imgCropServiceMakeGet_fileRequestToFile_server)
	s.Step(`^ImgCrop service make get_image request to image_server and return err_code and error to ImgCrop$`, tics.imgCropServiceMakeGet_imageRequestToImage_serverAndReturnErr_codeAndErrorToImgCrop)
}
