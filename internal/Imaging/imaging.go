package Imaging

import (
	"ImgCrop/internal/structs"
	"github.com/disintegration/imaging"
	"go.uber.org/zap"
	"image"
)

type Imagimg struct {
	Logger *zap.Logger
	Config structs.Config
}

func NewImaging(Logger *zap.Logger, Config structs.Config) (Imagimg, error) {
	return Imagimg{
		Logger: Logger,
		Config: Config,
	}, nil
}

func (d *Imagimg) CropImage(image structs.Image, width int, height int) (*image.NRGBA, error) {
	src, err := imaging.Open(image.Path)
	if err != nil {
		d.Logger.Error("Error imaging.Open", zap.Error(err))
		return nil, err
	}

	dst := imaging.Resize(src, width, height, imaging.Lanczos)
	return dst, nil
}
