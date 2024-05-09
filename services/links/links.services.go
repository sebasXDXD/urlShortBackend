package services

import (
	"urlShortenerBack/auth"
	"urlShortenerBack/entities"
	"urlShortenerBack/repositories"
)

type LinkService struct {
	LinkRepository repositories.LinkRepository
	AuthService    auth.AuthService
}

func NewLinkService(linkRepo repositories.LinkRepository) LinkService {
	return LinkService{LinkRepository: linkRepo}
}

func (ls LinkService) GetLinks() ([]entities.Link, error) {
	links, err := ls.LinkRepository.GetLinks()
	if err != nil {
		return nil, err
	}
	return links, nil
}

func (ls LinkService) CreateLink(newLink entities.Link) (entities.Link, error) {

	// Llamar al método CreateLink del repositorio y pasar el nuevo enlace
	createdLink, err := ls.LinkRepository.CreateLink(newLink)
	if err != nil {
		return entities.Link{}, err
	}

	return createdLink, nil
}
