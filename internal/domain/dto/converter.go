package dto

import (
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
)

// ToUserResponse converts a User entity to UserResponse DTO
func ToUserResponse(user *entity.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// ToArticleResponse converts an Article entity to ArticleResponse DTO
func ToArticleResponse(article *entity.Article) *ArticleResponse {
	if article == nil {
		return nil
	}
	return &ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Slug:        article.Slug,
		Content:     article.Content,
		ImageURL:    article.ImageURL,
		Category:    article.Category,
		AuthorID:    article.AuthorID,
		PublishedAt: article.PublishedAt,
		CreatedAt:   article.CreatedAt,
		UpdatedAt:   article.UpdatedAt,
	}
}

// ToArticleResponseList converts a slice of Article entities to ArticleResponse DTOs
func ToArticleResponseList(articles []entity.Article) []*ArticleResponse {
	responses := make([]*ArticleResponse, len(articles))
	for i, article := range articles {
		responses[i] = ToArticleResponse(&article)
	}
	return responses
}

// ToProgramResponse converts a Program entity to ProgramResponse DTO
func ToProgramResponse(program *entity.Program) *ProgramResponse {
	if program == nil {
		return nil
	}
	return &ProgramResponse{
		ID:           program.ID,
		Title:        program.Title,
		Description:  program.Description,
		ImageURL:     program.ImageURL,
		TargetAmount: program.TargetAmount,
		CreatedAt:    program.CreatedAt,
		UpdatedAt:    program.UpdatedAt,
	}
}

// ToProgramResponseList converts a slice of Program entities to ProgramResponse DTOs
func ToProgramResponseList(programs []entity.Program) []*ProgramResponse {
	responses := make([]*ProgramResponse, len(programs))
	for i, program := range programs {
		responses[i] = ToProgramResponse(&program)
	}
	return responses
}

// ToEventResponse converts an Event entity to EventResponse DTO
func ToEventResponse(event *entity.Event, currentParticipants int) *EventResponse {
	if event == nil {
		return nil
	}
	return &EventResponse{
		ID:                  event.ID,
		Title:               event.Title,
		Slug:                event.Slug,
		Description:         event.Description,
		ImageURL:            event.ImageURL,
		EventDate:           event.EventDate,
		Location:            event.Location,
		Quota:               event.Quota,
		CurrentParticipants: currentParticipants,
		CreatedAt:           event.CreatedAt,
		UpdatedAt:           event.UpdatedAt,
	}
}

// ToEventResponseList converts a slice of Event entities to EventResponse DTOs
func ToEventResponseList(events []entity.Event, participantCounts map[string]int) []*EventResponse {
	responses := make([]*EventResponse, len(events))
	for i, event := range events {
		count := 0
		if participantCounts != nil {
			count = participantCounts[event.ID]
		}
		responses[i] = ToEventResponse(&event, count)
	}
	return responses
}

// ToDonationResponse converts a Donation entity to DonationResponse DTO
func ToDonationResponse(donation *entity.Donation) *DonationResponse {
	if donation == nil {
		return nil
	}
	return &DonationResponse{
		ID:        donation.ID,
		OrderID:   donation.OrderID,
		UserID:    donation.UserID,
		Amount:    donation.Amount,
		Status:    donation.Status,
		PaidBy:    donation.PaidBy,
		CreatedAt: donation.CreatedAt,
		UpdatedAt: donation.UpdatedAt,
	}
}

// ToDonationResponseList converts a slice of Donation entities to DonationResponse DTOs
func ToDonationResponseList(donations []entity.Donation) []*DonationResponse {
	responses := make([]*DonationResponse, len(donations))
	for i, donation := range donations {
		responses[i] = ToDonationResponse(&donation)
	}
	return responses
}

// ToVolunteerResponse converts a Volunteer entity to VolunteerResponse DTO
func ToVolunteerResponse(volunteer *entity.Volunteer) *VolunteerResponse {
	if volunteer == nil {
		return nil
	}
	return &VolunteerResponse{
		ID:              volunteer.ID,
		FullName:        volunteer.FullName,
		Email:           volunteer.Email,
		PhoneNumber:     volunteer.PhoneNumber,
		BirthDate:       volunteer.BirthDate,
		Gender:          volunteer.Gender,
		Domicile:        volunteer.Domicile,
		Status:          volunteer.Status,
		Interest:        volunteer.Interest,
		CertificateName: volunteer.CertificateName,
		CreatedAt:       volunteer.CreatedAt,
		UpdatedAt:       volunteer.UpdatedAt,
	}
}

// ToVolunteerResponseList converts a slice of Volunteer entities to VolunteerResponse DTOs
func ToVolunteerResponseList(volunteers []entity.Volunteer) []*VolunteerResponse {
	responses := make([]*VolunteerResponse, len(volunteers))
	for i, volunteer := range volunteers {
		responses[i] = ToVolunteerResponse(&volunteer)
	}
	return responses
}

// ToRefreshTokenResponse converts a RefreshToken entity to RefreshTokenResponse DTO
func ToRefreshTokenResponse(token *entity.RefreshToken) *RefreshTokenResponse {
	if token == nil {
		return nil
	}
	return &RefreshTokenResponse{
		ID:        token.ID,
		UserID:    token.UserID,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		CreatedAt: token.CreatedAt,
	}
}

// ToRefreshTokenResponseList converts a slice of RefreshToken entities to RefreshTokenResponse DTOs
func ToRefreshTokenResponseList(tokens []entity.RefreshToken) []*RefreshTokenResponse {
	responses := make([]*RefreshTokenResponse, len(tokens))
	for i, token := range tokens {
		responses[i] = ToRefreshTokenResponse(&token)
	}
	return responses
}
