package service

import (
	"contacts-crud/cmd/contacts/models"
	"github.com/google/uuid"
)

type IContactRepository interface {
	AddContact(contact models.Contact) (*models.Contact, error)
	GetContact(id string) (*models.Contact, error)
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
	return nil, nil
}
func (cs *ContactService) GetAllContacts() ([]*models.Contact, error) {
	return nil, nil
}
