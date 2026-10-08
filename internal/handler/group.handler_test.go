package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
	"github.com/hrnnsx/libra/internal/service"
)

type fakeGroupService struct {
	createGroupFunc func(
		ctx context.Context,
		userID int64,
		name string,
		description *string,
	) (*model.Group, error)

	getGroupsFunc func(
		ctx context.Context,
		userID int64,
	) ([]model.Group, error)

	getGroupFunc func(
		ctx context.Context,
		userID int64,
		id int64,
	) (*model.Group, error)

	updateGroupFunc func(
		ctx context.Context,
		userID int64,
		id int64,
		name *string,
		description *string,
	) (*model.Group, error)

	deleteGroupFunc func(
		ctx context.Context,
		userID int64,
		id int64,
	) error
}

func (f *fakeGroupService) CreateGroup(
	ctx context.Context,
	userID int64,
	name string,
	description *string,
) (*model.Group, error) {
	if f.createGroupFunc != nil {
		return f.createGroupFunc(
			ctx,
			userID,
			name,
			description,
		)
	}

	return nil, nil
}

func (f *fakeGroupService) GetGroups(
	ctx context.Context,
	userID int64,
) ([]model.Group, error) {
	if f.getGroupsFunc != nil {
		return f.getGroupsFunc(ctx, userID)
	}

	return nil, nil
}

func (f *fakeGroupService) GetGroup(
	ctx context.Context,
	userID int64,
	id int64,
) (*model.Group, error) {
	if f.getGroupFunc != nil {
		return f.getGroupFunc(ctx, userID, id)
	}

	return nil, nil
}

func (f *fakeGroupService) UpdateGroup(
	ctx context.Context,
	userID int64,
	id int64,
	name *string,
	description *string,
) (*model.Group, error) {
	if f.updateGroupFunc != nil {
		return f.updateGroupFunc(
			ctx,
			userID,
			id,
			name,
			description,
		)
	}

	return nil, nil
}

func (f *fakeGroupService) DeleteGroup(
	ctx context.Context,
	userID int64,
	id int64,
) error {
	if f.deleteGroupFunc != nil {
		return f.deleteGroupFunc(ctx, userID, id)
	}

	return nil
}

func setupGroupHandler(
	groupService service.GroupService,
) *GroupHandler {
	gin.SetMode(gin.TestMode)

	return NewGroupHandler(groupService)
}

func setGroupUserID(
	c *gin.Context,
	userID int64,
) {
	c.Set("user_id", userID)
}

func TestGroupHandler_CreateGroup(t *testing.T) {
	description := "Favorite anime"

	var receivedUserID int64
	var receivedName string
	var receivedDescription *string

	expectedGroup := &model.Group{
		ID:          1,
		UserID:      10,
		Name:        "Favorites",
		Description: &description,
	}

	service := &fakeGroupService{
		createGroupFunc: func(
			ctx context.Context,
			userID int64,
			name string,
			description *string,
		) (*model.Group, error) {
			receivedUserID = userID
			receivedName = name
			receivedDescription = description

			return expectedGroup, nil
		},
	}

	handler := setupGroupHandler(service)

	body := `{
		"name": "Favorites",
		"description": "Favorite anime"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/groups",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	setGroupUserID(c, 10)

	handler.CreateGroup(c)

	if rec.Code != http.StatusCreated {
		t.Fatalf(
			"expected status 201, got %d",
			rec.Code,
		)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}

	if receivedName != "Favorites" {
		t.Errorf(
			"expected name Favorites, got %s",
			receivedName,
		)
	}

	if receivedDescription == nil {
		t.Fatal("expected description, got nil")
	}

	if *receivedDescription != description {
		t.Errorf(
			"expected description %q, got %q",
			description,
			*receivedDescription,
		)
	}

	var response map[string]interface{}

	if err := json.Unmarshal(
		rec.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	if response["group"] == nil {
		t.Fatal("expected group in response")
	}
}

func TestGroupHandler_CreateGroup_InvalidRequest(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	body := `{
		"description": "missing name"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/groups",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	setGroupUserID(c, 10)

	handler.CreateGroup(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_CreateGroup_ServiceError(t *testing.T) {
	expectedErr := errors.New("database error")

	service := &fakeGroupService{
		createGroupFunc: func(
			ctx context.Context,
			userID int64,
			name string,
			description *string,
		) (*model.Group, error) {
			return nil, expectedErr
		},
	}

	handler := setupGroupHandler(service)

	body := `{"name":"Favorites"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/groups",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	setGroupUserID(c, 10)

	handler.CreateGroup(c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_CreateGroup_MissingUserContext(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	body := `{"name":"Favorites"}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/groups",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	handler.CreateGroup(c)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status 401, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_GetGroups(t *testing.T) {
	groups := []model.Group{
		{
			ID:     1,
			UserID: 10,
			Name:   "Favorites",
		},
		{
			ID:     2,
			UserID: 10,
			Name:   "Watching",
		},
	}

	var receivedUserID int64

	service := &fakeGroupService{
		getGroupsFunc: func(
			ctx context.Context,
			userID int64,
		) ([]model.Group, error) {
			receivedUserID = userID
			return groups, nil
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/groups",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	setGroupUserID(c, 10)

	handler.GetGroups(c)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}

	var response map[string]interface{}

	if err := json.Unmarshal(
		rec.Body.Bytes(),
		&response,
	); err != nil {
		t.Fatalf(
			"failed to decode response: %v",
			err,
		)
	}

	data, ok := response["data"].([]interface{})

	if !ok {
		t.Fatal("expected data array in response")
	}

	if len(data) != 2 {
		t.Errorf(
			"expected 2 groups, got %d",
			len(data),
		)
	}
}

func TestGroupHandler_GetGroups_ServiceError(t *testing.T) {
	service := &fakeGroupService{
		getGroupsFunc: func(
			ctx context.Context,
			userID int64,
		) ([]model.Group, error) {
			return nil, errors.New("database error")
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/groups",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	setGroupUserID(c, 10)

	handler.GetGroups(c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_GetGroups_MissingUserContext(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/groups",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req

	handler.GetGroups(c)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf(
			"expected status 401, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_GetGroup(t *testing.T) {
	expectedGroup := &model.Group{
		ID:     1,
		UserID: 10,
		Name:   "Favorites",
	}

	var receivedID int64
	var receivedUserID int64

	service := &fakeGroupService{
		getGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
		) (*model.Group, error) {
			receivedUserID = userID
			receivedID = id

			return expectedGroup, nil
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/groups/1",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.GetGroup(c)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if receivedID != 1 {
		t.Errorf(
			"expected ID 1, got %d",
			receivedID,
		)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}
}

func TestGroupHandler_GetGroup_InvalidID(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/groups/abc",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "abc",
		},
	}

	handler.GetGroup(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_GetGroup_NotFound(t *testing.T) {
	service := &fakeGroupService{
		getGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
		) (*model.Group, error) {
			return nil, repository.ErrGroupNotFound
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/groups/999",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "999",
		},
	}

	setGroupUserID(c, 10)

	handler.GetGroup(c)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_GetGroup_ServiceError(t *testing.T) {
	service := &fakeGroupService{
		getGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
		) (*model.Group, error) {
			return nil, errors.New("database error")
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/groups/1",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.GetGroup(c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_UpdateGroup(t *testing.T) {
	name := "Updated"
	description := "Updated description"

	var receivedUserID int64
	var receivedID int64
	var receivedName *string
	var receivedDescription *string

	expectedGroup := &model.Group{
		ID:          1,
		UserID:      10,
		Name:        name,
		Description: &description,
	}

	service := &fakeGroupService{
		updateGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
			name *string,
			description *string,
		) (*model.Group, error) {
			receivedUserID = userID
			receivedID = id
			receivedName = name
			receivedDescription = description

			return expectedGroup, nil
		},
	}

	handler := setupGroupHandler(service)

	body := `{
		"name": "Updated",
		"description": "Updated description"
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/groups/1",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.UpdateGroup(c)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected status 200, got %d",
			rec.Code,
		)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}

	if receivedID != 1 {
		t.Errorf(
			"expected ID 1, got %d",
			receivedID,
		)
	}

	if receivedName == nil {
		t.Fatal("expected name, got nil")
	}

	if *receivedName != name {
		t.Errorf(
			"expected name %q, got %q",
			name,
			*receivedName,
		)
	}

	if receivedDescription == nil {
		t.Fatal("expected description, got nil")
	}

	if *receivedDescription != description {
		t.Errorf(
			"expected description %q, got %q",
			description,
			*receivedDescription,
		)
	}
}

func TestGroupHandler_UpdateGroup_InvalidID(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	body := `{"name":"Updated"}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/groups/abc",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "abc",
		},
	}

	handler.UpdateGroup(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_UpdateGroup_InvalidRequest(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	body := `{
		"name":
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/groups/1",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.UpdateGroup(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_UpdateGroup_InvalidName(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	body := `{
		"name": ""
	}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/groups/1",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.UpdateGroup(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_UpdateGroup_NotFound(t *testing.T) {
	service := &fakeGroupService{
		updateGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
			name *string,
			description *string,
		) (*model.Group, error) {
			return nil, repository.ErrGroupNotFound
		},
	}

	handler := setupGroupHandler(service)

	body := `{"name":"Updated"}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/groups/1",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.UpdateGroup(c)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_UpdateGroup_ServiceError(t *testing.T) {
	service := &fakeGroupService{
		updateGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
			name *string,
			description *string,
		) (*model.Group, error) {
			return nil, errors.New("database error")
		},
	}

	handler := setupGroupHandler(service)

	body := `{"name":"Updated"}`

	req := httptest.NewRequest(
		http.MethodPatch,
		"/groups/1",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.UpdateGroup(c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_DeleteGroup(t *testing.T) {
	var receivedUserID int64
	var receivedID int64

	service := &fakeGroupService{
		deleteGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
		) error {
			receivedUserID = userID
			receivedID = id
			return nil
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/groups/1",
		nil,
	)

	rec := httptest.NewRecorder()

	router := gin.New()

	router.DELETE(
		"/groups/:id",
		func(c *gin.Context) {
			setGroupUserID(c, 10)
			handler.DeleteGroup(c)
		},
	)

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf(
			"expected status 204, got %d",
			rec.Code,
		)
	}

	if rec.Body.Len() != 0 {
		t.Errorf(
			"expected empty response body, got %q",
			rec.Body.String(),
		)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}

	if receivedID != 1 {
		t.Errorf(
			"expected ID 1, got %d",
			receivedID,
		)
	}
}

func TestGroupHandler_DeleteGroup_InvalidID(t *testing.T) {
	service := &fakeGroupService{}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/groups/abc",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "abc",
		},
	}

	handler.DeleteGroup(c)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status 400, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_DeleteGroup_NotFound(t *testing.T) {
	service := &fakeGroupService{
		deleteGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
		) error {
			return repository.ErrGroupNotFound
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/groups/999",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "999",
		},
	}

	setGroupUserID(c, 10)

	handler.DeleteGroup(c)

	if rec.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status 404, got %d",
			rec.Code,
		)
	}
}

func TestGroupHandler_DeleteGroup_ServiceError(t *testing.T) {
	service := &fakeGroupService{
		deleteGroupFunc: func(
			ctx context.Context,
			userID int64,
			id int64,
		) error {
			return errors.New("database error")
		},
	}

	handler := setupGroupHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/groups/1",
		nil,
	)

	rec := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(rec)
	c.Request = req
	c.Params = gin.Params{
		{
			Key:   "id",
			Value: "1",
		},
	}

	setGroupUserID(c, 10)

	handler.DeleteGroup(c)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf(
			"expected status 500, got %d",
			rec.Code,
		)
	}
}
