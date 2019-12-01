package lru_cache

import (
	lg "ImgCrop/internal/logger"
	"ImgCrop/internal/structs"
	"go.uber.org/zap"
	"os"
	"reflect"
	"testing"
)

func TestLRUCache_Add(t *testing.T) {
	type fields struct {
		cache   []structs.Cache
		Size    int64
		MaxSize int64
		Logger  *zap.Logger
		Config  structs.Config
	}
	type args struct {
		image structs.Image
	}

	config := structs.Config{
		Logger: structs.LoggerConfig{
			Level:    "INFO",
			LogsPath: "",
			FileName: "..\\..\\TestLRU.log",
			Name:     "addTestLRU",
		},
		Cache: structs.CacheConfig{
			Path: "..\\..\\files",
			Size: 100,
		},
	}

	logger := lg.GetLogger(config)

	defaultFields := &fields{
		cache:   make([]structs.Cache, 0),
		Size:    0,
		MaxSize: 26000,
		Logger:  logger,
		Config:  config,
	}

	fullFields := &fields{
		cache:   make([]structs.Cache, 0),
		Size:    25000,
		MaxSize: 26000,
		Logger:  logger,
		Config:  config,
	}
	fullFields.cache = append(fullFields.cache, structs.Cache{
		Age: 5,
		Image: structs.Image{
			Size:       25000,
			Path:       "",
			Headers:    nil,
			Link:       "",
			FileName:   "",
			Exctension: "",
		},
	})

	simpleImage := structs.Image{
		Size:       25000,
		Path:       "..\\..\\files\\1.jpg",
		Headers:    nil,
		Link:       "http://lc/1.jpg",
		FileName:   "1.jpg",
		Exctension: "jpg",
	}

	bigImage := structs.Image{
		Size:       27000,
		Path:       "..\\..\\files\\2.jpg",
		Headers:    nil,
		Link:       "http://lc/2.jpg",
		FileName:   "2.jpg",
		Exctension: "jpg",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Add image", fields: *defaultFields, args: args{image: simpleImage}, wantErr: false},
		{name: "Add image which size more than cache size", fields: *defaultFields, args: args{image: bigImage}, wantErr: true},
		{name: "Add image in full cache", fields: *fullFields, args: args{image: simpleImage}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LRUCache{
				cache:   tt.fields.cache,
				Size:    tt.fields.Size,
				MaxSize: tt.fields.MaxSize,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
			}
			if err := c.Add(tt.args.image); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLRUCache_Get(t *testing.T) {
	config := structs.Config{
		Logger: structs.LoggerConfig{
			Level:    "INFO",
			LogsPath: "",
			FileName: "..\\..\\TestLRU.log",
			Name:     "addTestLRU",
		},
		Cache: structs.CacheConfig{
			Path: "..\\..\\files",
			Size: 100,
		},
	}

	logger := lg.GetLogger(config)

	type fields struct {
		cache   []structs.Cache
		Size    int64
		MaxSize int64
		Logger  *zap.Logger
		Config  structs.Config
	}
	type args struct {
		link string
	}

	simpleImage := structs.Image{
		Size:       25000,
		Path:       "..\\..\\files\\1.jpg",
		Headers:    nil,
		Link:       "http://lc/1.jpg",
		FileName:   "1.jpg",
		Exctension: "jpg",
	}

	fullFields := &fields{
		cache:   make([]structs.Cache, 0),
		Size:    25000,
		MaxSize: 26000,
		Logger:  logger,
		Config:  config,
	}
	fullFields.cache = append(fullFields.cache, structs.Cache{
		Age: 5,
		Image: structs.Image{
			Size:       25000,
			Path:       "..\\..\\files\\1.jpg",
			Headers:    nil,
			Link:       "http://lc/1.jpg",
			FileName:   "1.jpg",
			Exctension: "jpg",
		},
	})
	f, err := os.Create(simpleImage.Path)
	if err != nil {
		logger.Fatal(err.Error())
	}
	_ = f.Close()

	fullFields.cache = append(fullFields.cache, structs.Cache{
		Age: 4,
		Image: structs.Image{
			Size:       25000,
			Path:       "..\\..\\files\\3.jpg",
			Headers:    nil,
			Link:       "http://lc/3.jpg",
			FileName:   "3.jpg",
			Exctension: "jpg",
		},
	})

	tests := []struct {
		name            string
		fields          fields
		args            args
		wantCachedImage structs.Image
		wantErr         bool
	}{
		{name: "Image in cache", fields: *fullFields, args: args{link: "http://lc/1.jpg"}, wantCachedImage: simpleImage, wantErr: false},
		{name: "Image not in cache", fields: *fullFields, args: args{link: "http://lc/2.jpg"}, wantCachedImage: structs.Image{}, wantErr: true},
		{name: "Image in cache but file not found", fields: *fullFields, args: args{link: "http://lc/3.jpg"}, wantCachedImage: structs.Image{}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LRUCache{
				cache:   tt.fields.cache,
				Size:    tt.fields.Size,
				MaxSize: tt.fields.MaxSize,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
			}
			gotCachedImage, err := c.Get(tt.args.link)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCachedImage, tt.wantCachedImage) {
				t.Errorf("Get() gotCachedImage = %v, want %v", gotCachedImage, tt.wantCachedImage)
			}
		})
	}
}

func TestLRUCache_RemoveOldest(t *testing.T) {

	config := structs.Config{
		Logger: structs.LoggerConfig{
			Level:    "INFO",
			LogsPath: "",
			FileName: "..\\..\\TestLRU.log",
			Name:     "addTestLRU",
		},
		Cache: structs.CacheConfig{
			Path: "..\\..\\files",
			Size: 100,
		},
	}

	logger := lg.GetLogger(config)

	type fields struct {
		cache   []structs.Cache
		Size    int64
		MaxSize int64
		Logger  *zap.Logger
		Config  structs.Config
	}

	fullFields := &fields{
		cache:   make([]structs.Cache, 0),
		Size:    25000,
		MaxSize: 26000,
		Logger:  logger,
		Config:  config,
	}
	fullFields.cache = append(fullFields.cache, structs.Cache{
		Age: 5,
		Image: structs.Image{
			Size:       25000,
			Path:       "..\\..\\files\\1.jpg",
			Headers:    nil,
			Link:       "http://lc/1.jpg",
			FileName:   "1.jpg",
			Exctension: "jpg",
		},
	})
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "Remove oldest from cache", fields: *fullFields, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &LRUCache{
				cache:   tt.fields.cache,
				Size:    tt.fields.Size,
				MaxSize: tt.fields.MaxSize,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
			}
			if err := c.RemoveOldest(); (err != nil) != tt.wantErr {
				t.Errorf("RemoveOldest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
