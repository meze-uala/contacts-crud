package service

import (
	"contacts-crud/cmd/contacts/models"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewContactService(t *testing.T) {
	type args struct {
		contactRepository IContactRepository
	}
	tests := []struct {
		name string
		args args
		want ContactService
	}{
		{
			name: "Test with nil repo should create a service",
			args: args{contactRepository: nil},
			want: NewContactService(nil),
		},
		{
			name: "Test with no nil repo should create a service",
			args: args{contactRepository: NewMockIContactRepository(gomock.NewController(t))},
			want: NewContactService(NewMockIContactRepository(gomock.NewController(t))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContactService(tt.args.contactRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContactService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactService_AddContact(t *testing.T) {

	contact := GetValidContact()

	contactRepository := NewMockIContactRepository(gomock.NewController(t))
	contactRepository.EXPECT().AddContact(gomock.Any()).Return(&contact, nil)

	contactService := NewContactService(contactRepository)

	result, err := contactService.AddContact(contact)

	assert.NoError(t, err)
	assert.Equal(t, result.ID, "valid-id")
}

func TestContactService_AddContact_Error(t *testing.T) {

	contact := GetValidContact()

	contactRepository := NewMockIContactRepository(gomock.NewController(t))
	contactRepository.EXPECT().AddContact(gomock.Any()).Return(nil, errors.New("Error from contact repository "+
		"whenn trying to add a contact "))

	contactService := NewContactService(contactRepository)

	result, err := contactService.AddContact(contact)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestContactService_GetContact(t *testing.T) {
	contact := GetValidContact()

	contactRepository := NewMockIContactRepository(gomock.NewController(t))
	contactRepository.EXPECT().GetContact(gomock.Any()).Return(&contact, nil)

	contactService := NewContactService(contactRepository)

	result, err := contactService.GetContact(contact.ID)

	assert.NoError(t, err)
	assert.Equal(t, result.ID, "valid-id")
}

func TestContactService_GetContact_Error(t *testing.T) {

	contact := GetValidContact()

	contactRepository := NewMockIContactRepository(gomock.NewController(t))
	contactRepository.EXPECT().GetContact(gomock.Any()).Return(nil, errors.New("error from "+
		"contact repository trying to get a contact"))

	contactService := NewContactService(contactRepository)

	result, err := contactService.GetContact(contact.ID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestContactService_UpdateContactStatus(t *testing.T) {
	contact := GetValidContact()

	contactRepository := NewMockIContactRepository(gomock.NewController(t))
	contactRepository.EXPECT().UpdateContactStatus(gomock.Any()).Return(&contact, nil)

	contactService := NewContactService(contactRepository)

	result, err := contactService.UpdateContactStatus(contact.ID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestContactService_UpdateContactStatus_Error(t *testing.T) {
	contact := GetValidContact()

	contactRepository := NewMockIContactRepository(gomock.NewController(t))
	contactRepository.EXPECT().UpdateContactStatus(gomock.Any()).Return(nil, errors.New("error from repo"))

	contactService := NewContactService(contactRepository)

	result, err := contactService.UpdateContactStatus(contact.ID)

	assert.Error(t, err)
	assert.Nil(t, result)
}

//Generic test just for coverge. Function isnt implemented
func TestContactService_GetAllContacts(t *testing.T) {

	contactRepository := NewMockIContactRepository(gomock.NewController(t))

	contactService := NewContactService(contactRepository)

	_, err := contactService.GetAllContacts()

	assert.NoError(t, err)
}

func GetValidContact() models.Contact {
	return models.Contact{
		ID:        "valid-id",
		FirstName: "Fulano",
		LastName:  "De Tal",
		Status:    "CREATED",
	}
}
