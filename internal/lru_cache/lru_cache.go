package lru_cache

import (
	"ImgCrop/internal/structs"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"os"
)

type LRUCache struct {
	cache   []structs.Cache
	Size    int64
	MaxSize int64
	Logger  *zap.Logger    `mapstructure:"logger"`
	Config  structs.Config `mapstructure:"config"`
}

func NewLRUCache(size int64, logger *zap.Logger, config structs.Config) (LRUCache, error) {

	logger.Info(fmt.Sprintf("Creating Cache, Size: %d", size))

	return LRUCache{
		cache:   make([]structs.Cache, 0),
		Size:    0,
		MaxSize: size * 1024 * 1024,
		Logger:  logger,
		Config:  config,
	}, nil
}

func (c *LRUCache) Get(link string) (cachedImage structs.Image, err error) {
	c.Logger.Info(fmt.Sprintf("Search in cache, Link: %s", link))
	// ищем по линку, если такой есть то возвращаем
	for index, cache := range c.cache {
		if cache.Image.Link == link {
			c.Logger.Info(fmt.Sprintf("Image in cache, Link %s", link))

			//чекам наличие изображения на диске

			c.Logger.Info(fmt.Sprintf("Search file in path: %s", cache.Image.Path))
			if _, err := os.Stat(cache.Image.Path); os.IsNotExist(err) {
				if err != nil {
					c.Logger.Error(err.Error())
					c.cache = append(c.cache[:index], c.cache[index+1:]...)
					return structs.Image{}, fmt.Errorf("Image not fount on disc")
				}
			}

			//увеличиваем возраст у остальных
			for i, ch := range c.cache {
				if i != index {
					ch.Age += 1
				}
			}

			return cache.Image, nil
		}
	}

	return structs.Image{}, errors.New("Image not found in Cache")
}

func (c *LRUCache) Add(image structs.Image) error {
	// вернет false если не сможет положить в кеш

	s := image.Size
	c.Logger.Info(fmt.Sprintf("Image size: %v Cache size: %v/%v", image.Size, c.Size, c.MaxSize))
	// если файл не влезает в кеш. то удаляем самый старый
	if s > c.MaxSize {
		c.Logger.Info(fmt.Sprintf("Image size: %v more > than max Cache size: %v", image.Size, c.MaxSize))
		return fmt.Errorf("Image size: %v more > than max Cache size: %v", image.Size, c.MaxSize)
	}
	if s > c.MaxSize-c.Size {

		for {
			c.Logger.Info(fmt.Sprintf("ImageSize: %d > CacheAvailiableSize: %d ", s, c.MaxSize-c.Size))
			err := c.RemoveOldest()
			if err == nil {
				if s > c.MaxSize-c.Size {
					continue
				}
				// увеличиваем возраст всех кто в кеше
				for _, cache := range c.cache {
					cache.Age += 1
				}
				c.cache = append(c.cache, structs.Cache{
					Age:   0,
					Image: image,
				})
				c.Size += image.Size
				c.Logger.Info(fmt.Sprintf("Cache size: %v/%v", c.Size, c.MaxSize))
				return nil
			} else {
				return err
			}
		}

	} else {
		// увеличиваем возраст всех кто в кеше
		for _, cache := range c.cache {
			cache.Age += 1
		}
		// если влзает, то кладем его
		c.cache = append(c.cache, structs.Cache{
			Age:   0,
			Image: image,
		})

		c.Size += image.Size
	}
	c.Logger.Info(fmt.Sprintf("Cache size: %v/%v", c.Size, c.MaxSize))
	return nil
}

func (c *LRUCache) RemoveOldest() error {

	// ищем макс возраст
	var maxAge int64 = 0
	var oldesIndex int
	for index, cache := range c.cache {
		if cache.Age > maxAge {
			maxAge = cache.Age
			oldesIndex = index
		}
	}
	c.Logger.Info(fmt.Sprintf("Try remove Image from Cache: Name: %s URL: %s", c.cache[oldesIndex].Image.FileName, c.cache[oldesIndex].Image.Link))

	go func(image structs.Image) {
		err := os.Remove(image.Path)
		if err != nil {
			c.Logger.Error(fmt.Sprintf("I can't remove file %s in path %s", image.FileName, image.Path))
		}
	}(c.cache[oldesIndex].Image)

	if oldesIndex == len(c.cache) {
		c.Size -= c.cache[oldesIndex].Image.Size
		c.cache = c.cache[:oldesIndex]

	} else {
		c.Size -= c.cache[oldesIndex].Image.Size
		c.cache = append(c.cache[:oldesIndex], c.cache[oldesIndex+1:]...)

	}
	//todo как же вернуть false ????
	c.Logger.Info("Successfully remove Image from Cache")
	return nil
}
