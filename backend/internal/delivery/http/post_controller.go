package http

import (
	"backend/internal/constant"
	"backend/internal/delivery/http/exception"
	"backend/internal/model"
	"backend/internal/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type PostController struct {
	PostUseCase *usecase.PostUseCase
}

func NewPostController(postUseCase *usecase.PostUseCase) *PostController {
	return &PostController{
		PostUseCase: postUseCase,
	}
}

func (ct *PostController) GetAll(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	userID := c.QueryParam("authorID")
	titleQuery := c.QueryParam("title")

	request := model.PostListRequest{
		Page:       page,
		PageSize:   pageSize,
		UserID:     userID,
		TitleQuery: titleQuery,
	}
	posts, err := ct.PostUseCase.List(c.Request().Context(), &request)
	if err != nil {
		return err
	}

	response := model.DataResponse[[]model.PostResponse]{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   posts,
	}
	return c.JSON(response.Code, response)
}

func (ct *PostController) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	request := model.PostGetByIDRequest{
		ID: uint64(id),
	}
	post, err := ct.PostUseCase.GetByID(c.Request().Context(), &request)
	if err != nil {
		return err
	}

	response := model.DataResponse[*model.PostResponse]{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   post,
	}
	return c.JSON(response.Code, response)
}

func (ct *PostController) Create(c echo.Context) error {
	// Get current user ID from context
	authDataInterface := c.Get(constant.USER_AUTH_DATA_CONTEXT_NAME)

	// Check if data is on CurrentUser tyoe
	currentUser, ok := authDataInterface.(*model.CurrentUser)
	if !ok {
		return echo.ErrUnauthorized
	}

	request := new(model.PostCreateRequest)
	if err := c.Bind(request); err != nil {
		return exception.NewBadRequestError(err.Error())
	}

	request.AuthorID = currentUser.ID

	newPost, err := ct.PostUseCase.Create(c.Request().Context(), request)
	if err != nil {
		return err
	}

	response := model.DataResponse[*model.PostResponse]{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   newPost,
	}
	return c.JSON(response.Code, response)
}

func (ct *PostController) Update(c echo.Context) error {
	// Get id from parameter
	id, _ := strconv.Atoi(c.Param("id"))

	// Get current user ID from context
	authDataInterface := c.Get(constant.USER_AUTH_DATA_CONTEXT_NAME)

	// Check if data is on CurrentUser type
	currentUser, ok := authDataInterface.(*model.CurrentUser)
	if !ok {
		return echo.ErrUnauthorized
	}

	request := new(model.PostUpdateRequest)
	if err := c.Bind(request); err != nil {
		return exception.NewBadRequestError(err.Error())
	}

	request.ID = uint64(id)
	request.AuthorID = currentUser.ID

	post, err := ct.PostUseCase.Update(c.Request().Context(), request)
	if err != nil {
		return err
	}

	response := model.DataResponse[*model.PostResponse]{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   post,
	}
	return c.JSON(response.Code, response)
}

func (ct *PostController) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	// Get current user ID from context
	authDataInterface := c.Get(constant.USER_AUTH_DATA_CONTEXT_NAME)

	// Check if data is on CurrentUser type
	currentUser, ok := authDataInterface.(*model.CurrentUser)
	if !ok {
		return echo.ErrUnauthorized
	}

	// Delete post with service
	request := model.PostDeleteRequest{
		ID:     uint64(id),
		UserID: currentUser.ID,
	}
	if err := ct.PostUseCase.Delete(c.Request().Context(), &request); err != nil {
		return err
	}

	response := model.DataResponse[any]{
		Code:   http.StatusOK,
		Status: "OK",
	}
	return c.JSON(response.Code, response)
}
