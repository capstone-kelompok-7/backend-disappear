package dto

import "github.com/capstone-kelompok-7/backend-disappear/module/entities"

type CarouselFormatter struct {
	ID    uint64 `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

func FormatCarousel(carousel *entities.CarouselModels) *CarouselFormatter {
	carouselFormatter := &CarouselFormatter{}
	carouselFormatter.ID = carousel.ID
	carouselFormatter.Name = carousel.Name
	carouselFormatter.Photo = carousel.Photo

	return carouselFormatter
}

func FormatterCarousel(carousels []*entities.CarouselModels) []*CarouselFormatter {
	var carouselFormatters []*CarouselFormatter

	for _, carousel := range carousels {
		formattedCarousel := FormatCarousel(carousel)
		carouselFormatters = append(carouselFormatters, formattedCarousel)
	}

	return carouselFormatters
}
