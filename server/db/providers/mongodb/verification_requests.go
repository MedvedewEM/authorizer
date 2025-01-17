package mongodb

import (
	"log"
	"time"

	"github.com/authorizerdev/authorizer/server/db/models"
	"github.com/authorizerdev/authorizer/server/graph/model"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AddVerification to save verification request in database
func (p *provider) AddVerificationRequest(verificationRequest models.VerificationRequest) (models.VerificationRequest, error) {
	if verificationRequest.ID == "" {
		verificationRequest.ID = uuid.New().String()

		verificationRequest.CreatedAt = time.Now().Unix()
		verificationRequest.UpdatedAt = time.Now().Unix()
		verificationRequest.Key = verificationRequest.ID
		verificationRequestCollection := p.db.Collection(models.Collections.VerificationRequest, options.Collection())
		_, err := verificationRequestCollection.InsertOne(nil, verificationRequest)
		if err != nil {
			log.Println("error saving verification record:", err)
			return verificationRequest, err
		}
	}

	return verificationRequest, nil
}

// GetVerificationRequestByToken to get verification request from database using token
func (p *provider) GetVerificationRequestByToken(token string) (models.VerificationRequest, error) {
	var verificationRequest models.VerificationRequest

	verificationRequestCollection := p.db.Collection(models.Collections.VerificationRequest, options.Collection())
	err := verificationRequestCollection.FindOne(nil, bson.M{"token": token}).Decode(&verificationRequest)
	if err != nil {
		return verificationRequest, err
	}

	return verificationRequest, nil
}

// GetVerificationRequestByEmail to get verification request by email from database
func (p *provider) GetVerificationRequestByEmail(email string, identifier string) (models.VerificationRequest, error) {
	var verificationRequest models.VerificationRequest

	verificationRequestCollection := p.db.Collection(models.Collections.VerificationRequest, options.Collection())
	err := verificationRequestCollection.FindOne(nil, bson.M{"email": email, "identifier": identifier}).Decode(&verificationRequest)
	if err != nil {
		return verificationRequest, err
	}

	return verificationRequest, nil
}

// ListVerificationRequests to get list of verification requests from database
func (p *provider) ListVerificationRequests(pagination model.Pagination) (*model.VerificationRequests, error) {
	var verificationRequests []*model.VerificationRequest

	opts := options.Find()
	opts.SetLimit(pagination.Limit)
	opts.SetSkip(pagination.Offset)
	opts.SetSort(bson.M{"created_at": -1})

	verificationRequestCollection := p.db.Collection(models.Collections.VerificationRequest, options.Collection())

	verificationRequestCollectionCount, err := verificationRequestCollection.CountDocuments(nil, bson.M{})
	paginationClone := pagination
	paginationClone.Total = verificationRequestCollectionCount

	cursor, err := verificationRequestCollection.Find(nil, bson.M{}, opts)
	if err != nil {
		log.Println("error getting verification requests:", err)
		return nil, err
	}
	defer cursor.Close(nil)

	for cursor.Next(nil) {
		var verificationRequest models.VerificationRequest
		err := cursor.Decode(&verificationRequest)
		if err != nil {
			return nil, err
		}
		verificationRequests = append(verificationRequests, verificationRequest.AsAPIVerificationRequest())
	}

	return &model.VerificationRequests{
		VerificationRequests: verificationRequests,
		Pagination:           &paginationClone,
	}, nil
}

// DeleteVerificationRequest to delete verification request from database
func (p *provider) DeleteVerificationRequest(verificationRequest models.VerificationRequest) error {
	verificationRequestCollection := p.db.Collection(models.Collections.VerificationRequest, options.Collection())
	_, err := verificationRequestCollection.DeleteOne(nil, bson.M{"_id": verificationRequest.ID}, options.Delete())
	if err != nil {
		log.Println("error deleting verification request::", err)
		return err
	}

	return nil
}
