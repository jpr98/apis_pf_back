package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Project represents a project in the system
type Project struct {
	ID            primitive.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Owner         primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
	Title         string               `json:"title,omitempty" bson:"title,omitempty"`
	Subtitle      string               `json:"subtitle,omitempty" bson:"subtitle,omitempty"`
	Description   string               `json:"description,omitempty" bson:"desc,omitempty"`
	CreatedAt     time.Time            `json:"created_at,omitempty" bson:"created_at,omitempty"`
	Tags          []string             `json:"tags,omitempty" bson:"tags,omitempty"`
	Category      string               `json:"category,omitempty" bson:"category,omitempty"`
	Location      string               `json:"location,omitempty" bson:"location,omitempty"`
	Votes         []primitive.ObjectID `json:"votes,omitempty" bson:"votes,omitempty"`
	VotesCount    int                  `json:"votes_count,omitempty" bson:"votes_count,omitempty"`
	ImageURL      string               `json:"image_url,omitempty" bson:"image,omitempty"`
	VideoURL      string               `json:"video_url,omitempty" bson:"video,omitempty"`
	Views         int                  `json:"views,omitempty" bson:"views,omitempty"`
	Comments      []Comment            `json:"comments,omitempty" bson:"comments,omitempty"`
	Contributions []Contribution       `json:"contributions,omitempty" bson:"contributions,omitempty"`
	Duration      int                  `json:"duration,omitempty" bson:"duration,omitempty"`
}

// ProjectStore contains all the CRUD operations of Project
type ProjectStore struct {
	database   *mongo.Database
	collection *mongo.Collection
}

// NewProjectStore creates a project store with a mongo database
func NewProjectStore(database *mongo.Database) *ProjectStore {
	return &ProjectStore{database, database.Collection("projects")}
}

// Create receives a project object and tries to insert it to the project store
func (ps *ProjectStore) Create(p Project, ownerID string) (Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return Project{}, err
	}

	for index, tag := range p.Tags {
		p.Tags[index] = strings.ToLower(tag)
	}

	p.Owner = oid
	p.Views = 0
	p.VotesCount = 0
	p.CreatedAt = time.Now()
	result, err := ps.collection.InsertOne(ctx, p)
	if err != nil {
		return Project{}, err
	}

	generatedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return Project{}, errors.New("Invalid generated id on project")
	}
	p.ID = generatedID

	return p, nil
}

// EditProject helps to model de edit project data
type EditProject struct {
	Title       string   `json:"title,omitempty"`
	Subtitle    string   `json:"subtitle,omitempty"`
	Location    string   `json:"location,omitempty"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
	VideoURL    string   `json:"video_url,omitempty"`
	Duration    int      `json:"duration,omitempty"`
	Description string   `json:"description,omitempty"`
}

// Update edits a project's info
func (ps *ProjectStore) Update(project Project, editProject EditProject) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	project.Title = editProject.Title
	project.Subtitle = editProject.Subtitle
	project.Location = editProject.Location
	project.Category = editProject.Category
	project.Tags = editProject.Tags
	project.ImageURL = editProject.ImageURL
	project.VideoURL = editProject.VideoURL
	project.Duration = editProject.Duration
	project.Description = editProject.Description

	result, err := ps.collection.ReplaceOne(ctx, bson.M{"_id": project.ID}, project)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("No projects with given id")
	}

	return nil
}

// GetByID finds a project with a given id
func (ps *ProjectStore) GetByID(id string) (Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var project Project
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Project{}, err
	}

	err = ps.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&project)
	if err != nil {
		return Project{}, err
	}

	ps.getCommentsAuthors(&project)
	ps.getContributionsUsers(&project)

	return project, nil
}

// GetByTitle returns all projects with titles containing the given query string
func (ps *ProjectStore) GetByTitle(title string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{"title", primitive.Regex{Pattern: ".*" + title + ".*", Options: ""}}}
	cursor, err := ps.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	projects, err := ps.extractProjectsFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetByTags returns all projects for a given set of tags
func (ps *ProjectStore) GetByTags(tags []string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ps.collection.Find(ctx, bson.D{{"tags", bson.D{{"$in", tags}}}})
	if err != nil {
		return nil, err
	}

	projects, err := ps.extractProjectsFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetByCategory returns all projects for a given category
func (ps *ProjectStore) GetByCategory(category string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ps.collection.Find(ctx, bson.M{"category": category})
	if err != nil {
		return nil, err
	}

	projects, err := ps.extractProjectsFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetFullSearch looks for projects by title, category and returns them in a specific order
func (ps *ProjectStore) GetFullSearch(title, category, order string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var query bson.M
	if category != "todos" {
		query = bson.M{"$and": bson.A{
			bson.M{"title": primitive.Regex{Pattern: ".*" + title + ".*", Options: ""}},
			bson.M{"category": category},
		}}
	} else {
		query = bson.M{"$and": bson.A{
			bson.M{"title": primitive.Regex{Pattern: ".*" + title + ".*", Options: ""}},
		}}
	}

	options := &options.FindOptions{}
	switch order {
	case "popularity":
		options.SetSort(bson.M{"votes_count": -1})
	case "date":
		options.SetSort(bson.M{"created_at": 1})
	default:
		options = nil
	}

	cursor, err := ps.collection.Find(ctx, query, options)
	if err != nil {
		return nil, err
	}

	projects, err := ps.extractProjectsFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetByOwnerID returns projects with a given owner ID
func (ps *ProjectStore) GetByOwnerID(ownerID string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(ownerID)
	if err != nil {
		return nil, err
	}

	cursor, err := ps.collection.Find(ctx, bson.M{"owner": oid})
	if err != nil {
		return nil, err
	}

	projects, err := ps.extractProjectsFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetVotedProjects returns the projects that a user has voted for
func (ps *ProjectStore) GetVotedProjects(userID string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	query := bson.D{{"votes", bson.D{{"$in", []primitive.ObjectID{uid}}}}}
	cursor, err := ps.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	projects, err := ps.extractProjectsFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return projects, nil
}

// GetContributedProjects returns the projects that a user has contributed to
func (ps *ProjectStore) GetContributedProjects(userID string) ([]Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	query := bson.M{"contributions.user._id": uid}
	cursor, err := ps.collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	projects, err := ps.extractProjectsFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	cursor.Close(ctx)
	return projects, nil
}

// Vote appends or removes a user to the list of votes of a project
func (ps *ProjectStore) Vote(projectID string, userID string, upvote bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pid, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		return err
	}

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	var updateAction string
	var incAmount int
	if upvote {
		updateAction = "$addToSet"
		incAmount = 1
	} else {
		updateAction = "$pull"
		incAmount = -1
	}
	update := bson.M{updateAction: bson.M{"votes": uid}}
	result, err := ps.collection.UpdateOne(ctx, bson.M{"_id": pid}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("No project found with given id")
	}

	if result.ModifiedCount > 0 {
		_, err := ps.collection.UpdateOne(ctx, bson.M{"_id": pid}, bson.M{"$inc": bson.M{"votes_count": incAmount}})
		if err != nil {
			return err
		}
	}

	return nil
}

// Delete removes a project with a given id
func (ps *ProjectStore) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := ps.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("No projects with id found")
	}

	return nil
}

// View increments the views of a project by one
func (ps *ProjectStore) View(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{"$inc": bson.M{"views": 1}}
	result, err := ps.collection.UpdateOne(ctx, bson.M{"_id": oid}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("No projects with id found")
	}

	return nil
}

// AddComment appends a comment to a project
func (ps *ProjectStore) AddComment(id, authorID, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	uid, err := primitive.ObjectIDFromHex(authorID)
	if err != nil {
		return err
	}

	author := CommentAuthor{ID: uid}
	comment := Comment{primitive.NewObjectID(), author, time.Now(), text}

	update := bson.M{"$push": bson.M{"comments": comment}}
	result, err := ps.collection.UpdateOne(ctx, bson.M{"_id": pid}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("No project found with given id")
	}

	return nil
}

// AddContribution appends a contribution to a project
func (ps *ProjectStore) AddContribution(id, userID string, amount float32) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	user := ContributionUser{ID: uid}
	contribution := Contribution{primitive.NewObjectID(), user, amount, time.Now()}

	update := bson.M{"$push": bson.M{"contributions": contribution}}
	result, err := ps.collection.UpdateOne(ctx, bson.M{"_id": pid}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("No project found with given id")
	}

	return nil
}

func (ps *ProjectStore) extractProjectsFromCursor(ctx context.Context, cursor *mongo.Cursor) ([]Project, error) {
	projects := make([]Project, 0)
	for cursor.Next(ctx) {
		var project Project
		err := cursor.Decode(&project)
		if err != nil {
			return nil, err
		}
		ps.getCommentsAuthors(&project)
		ps.getContributionsUsers(&project)
		projects = append(projects, project)
	}
	return projects, nil
}

func (ps *ProjectStore) getCommentsAuthors(project *Project) {
	for index, comment := range project.Comments {
		userStore := NewUserStore(ps.database)
		user, err := userStore.GetByID(comment.Author.ID.Hex())
		if err != nil {
			project.Comments[index].Author = CommentAuthor{user.ID, "Eliminado", ""}
		}
		project.Comments[index].Author = CommentAuthor{user.ID, user.Name, user.Avatar}
	}
}

func (ps *ProjectStore) getContributionsUsers(project *Project) {
	for index, comment := range project.Contributions {
		userStore := NewUserStore(ps.database)
		user, err := userStore.GetByID(comment.User.ID.Hex())
		if err != nil {
			project.Contributions[index].User = ContributionUser{user.ID, "Eliminado", ""}
		}
		project.Contributions[index].User = ContributionUser{user.ID, user.Name, user.Avatar}
	}
}
