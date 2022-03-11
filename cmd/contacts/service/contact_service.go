package service

//go:generate mockgen -source contact_service.go -destination mock_contact_service.go -package service

import (
	"contacts-crud/cmd/contacts/models"
	"github.com/google/uuid"
)

type IContactRepository interface {
	AddContact(contact models.Contact) (*models.Contact, error)
	GetContact(id string) (*models.Contact, error)
	UpdateContactStatus(id string) (*models.Contact, error)
	GetAllContacts() ([]*models.Contact, error)
}

type ContactService struct {
	contactRepository IContactRepository
}

func NewContactService(contactRepository IContactRepository) ContactService {
	return ContactService{contactRepository: contactRepository}
}

func (cs *ContactService) AddContact(contact models.Contact) (*models.Contact, error) {

	uniqueID, _ := uuid.NewUUID()
	contact.ID = uniqueID.String()
	contact.Status = "CREATED"

	return cs.contactRepository.AddContact(contact)

}

func (cs *ContactService) GetContact(id string) (*models.Contact, error) {
	result, err := cs.contactRepository.GetContact(id)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *ContactService) UpdateContactStatus(id string) (*models.Contact, error) {
	result, err := cs.contactRepository.UpdateContactStatus(id)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cs *ContactService) GetAllContacts() ([]*models.Contact, error) {
	return nil, nil
}
