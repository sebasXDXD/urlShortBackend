package services

import (
	"errors"
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
	createdLink, err := ls.LinkRepository.CreateLink(newLink)
	if err != nil {
		return entities.Link{}, err
	}
	return createdLink, nil
}

func (ls LinkService) GetLinkByString(name string) (entities.Link, error) {
	link, err := ls.LinkRepository.GetLinkByString(name)
	if err != nil {
		return entities.Link{}, err
	}

	// Verificar si redirectTo es vacío
	if link.RedirectTo == "" {
		return entities.Link{}, errors.New("RedirectTo is empty")
	}

	return link, nil
}

func (ls LinkService) GetLinkByID(id int) (entities.Link, error) {
	link, err := ls.LinkRepository.GetLinkByID(id)
	if err != nil {
		return entities.Link{}, err
	}
	return link, nil
}

func (ls LinkService) UpdateLink(linkID int, name string, redirectTo string) error {
	return ls.LinkRepository.UpdateLink(linkID, name, redirectTo)
}
