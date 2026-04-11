package main

import "github.com/gin-gonic/gin"

// Error standardizes error responses for the API
func (_ *API) Error(err string) gin.H {
	return gin.H{"error": err}
}

// ErrorResponse sends a JSON response with the given error message and HTTP status code
func (api *API) ErrorResponse(c *gin.Context, code int, err string) {
	c.Set("error", err)
	c.JSON(code, api.Error(err))
}

/*
	Successful responses
	https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#successful_responses
*/

// OK is a request that has succeeded and resulted in a response being returned, such as a
// successful login or data retrieval
func (_ *API) OK(c *gin.Context, data any) {
	c.JSON(200, data)
}

// Created is a request that has been fulfilled and resulted in a new resource being created, such
// as a new user or food entry
func (_ *API) Created(c *gin.Context, data any) {
	c.JSON(201, data)
}

// Accepted is a request that has been accepted for processing asynchronously
func (_ *API) Accepted(c *gin.Context) {
	c.Status(202)
}

// NoContent is a request that has been successfully processed but is not returning any content,
// such as a successful logout or deletion of a resource
func (_ *API) NoContent(c *gin.Context) {
	c.Status(204)
}

// PartialContent is a request that has been successfully processed but is only returning a portion
// of the requested resource
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Range_requests
func (_ *API) PartialContent(c *gin.Context, data any) {
	c.JSON(206, data)
}

/*
	Client error responses
	https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#client_error_responses
*/

// BadRequest is a request perceived to be a client error, like invalid input or missing fields
func (api *API) BadRequest(c *gin.Context, err string) {
	api.ErrorResponse(c, 400, err)
}

// Unauthorized is a request that lacks valid authentication credentials for the target resource
func (api *API) Unauthorized(c *gin.Context) {
	api.ErrorResponse(c, 401, "Unauthorized")
}

// Forbidden is a request that is valid but the server is refusing action, like insufficient or
// incorrect permissions - such as trying to access another user's data
func (api *API) Forbidden(c *gin.Context) {
	api.ErrorResponse(c, 403, "Forbidden")
}

// Conflict is a request that could not be completed due to a conflict with the current state of the
// target resource, such as trying to create a user with an existing username
func (api *API) Conflict(c *gin.Context, err string) {
	api.ErrorResponse(c, 409, err)
}

// TooManyRequests is a request that has been rate limited by the server
func (api *API) TooManyRequests(c *gin.Context) {
	api.ErrorResponse(c, 429, "Too Many Requests")
}

/*
	Server error responses
	https://developer.mozilla.org/en-US/docs/Web/HTTP/Status#server_error_responses
*/

// InternalServerError is a request that has encountered an unexpected condition that prevented it
// from being fulfilled
func (api *API) InternalServerError(c *gin.Context) {
	api.ErrorResponse(c, 500, "Internal Server Error")
}

// NotImplemented is a request that has not been implemented by the server
func (api *API) NotImplemented(c *gin.Context) {
	api.ErrorResponse(c, 501, "Not Implemented")
}
