package webrpc

import "fmt"

type Error interface {
	// Code is of the valid error codes
	Code() ErrorCode

	// Msg returns a human-readable, unstructured messages describing the error
	Msg() string

	// Cause is reason for the error
	Cause() error

	// Error returns a string of the form "webrpc error <Code>: <Msg>"
	Error() string
}

func Errorf(code ErrorCode, format string, args ...interface{}) Error {
	msg := fmt.Sprintf(format, args...)
	if IsValidErrorCode(code) {
		return &rpcErr{code: code, msg: msg}
	}
	return &rpcErr{code: ErrInternal, msg: "invalid error type " + string(code)}
}

func WrapError(code ErrorCode, cause error, format string, args ...interface{}) Error {
	msg := fmt.Sprintf(format, args...)
	if IsValidErrorCode(code) {
		return &rpcErr{code: code, msg: msg, cause: cause}
	}
	return &rpcErr{code: ErrInternal, msg: "invalid error type " + string(code), cause: cause}
}

func ErrorNotFound(format string, args ...interface{}) Error {
	return Errorf(ErrNotFound, format, args...)
}

func ErrorInvalidArgument(argument string, validationMsg string) Error {
	return Errorf(ErrInvalidArgument, argument+" "+validationMsg)
}

func ErrorRequiredArgument(argument string) Error {
	return ErrorInvalidArgument(argument, "is required")
}

func ErrorInternal(format string, args ...interface{}) Error {
	return Errorf(ErrInternal, format, args...)
}

type ErrorCode string

const (
	// Canceled indicates the operation was cancelled (typically by the caller).
	ErrCanceled ErrorCode = "canceled"

	// Unknown error. For example when handling errors raised by APIs that do not
	// return enough error information.
	ErrUnknown ErrorCode = "unknown"

	// InvalidArgument indicates client specified an invalid argument. It
	// indicates arguments that are problematic regardless of the state of the
	// system (i.e. a malformed file name, required argument, number out of range,
	// etc.).
	ErrInvalidArgument ErrorCode = "invalid_argument"

	// DeadlineExceeded means operation expired before completion. For operations
	// that change the state of the system, this error may be returned even if the
	// operation has completed successfully (timeout).
	ErrDeadlineExceeded ErrorCode = "deadline_exceeded"

	// NotFound means some requested entity was not found.
	ErrNotFound ErrorCode = "not_found"

	// BadRoute means that the requested URL path wasn't routable to a webrpc
	// service and method. This is returned by the generated server, and usually
	// shouldn't be returned by applications. Instead, applications should use
	// NotFound or Unimplemented.
	ErrBadRoute ErrorCode = "bad_route"

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	ErrAlreadyExists ErrorCode = "already_exists"

	// PermissionDenied indicates the caller does not have permission to execute
	// the specified operation. It must not be used if the caller cannot be
	// identified (Unauthenticated).
	ErrPermissionDenied ErrorCode = "permission_denied"

	// Unauthenticated indicates the request does not have valid authentication
	// credentials for the operation.
	ErrUnauthenticated ErrorCode = "unauthenticated"

	// ResourceExhausted indicates some resource has been exhausted, perhaps a
	// per-user quota, or perhaps the entire file system is out of space.
	ErrResourceExhausted ErrorCode = "resource_exhausted"

	// FailedPrecondition indicates operation was rejected because the system is
	// not in a state required for the operation's execution. For example, doing
	// an rmdir operation on a directory that is non-empty, or on a non-directory
	// object, or when having conflicting read-modify-write on the same resource.
	ErrFailedPrecondition ErrorCode = "failed_precondition"

	// Aborted indicates the operation was aborted, typically due to a concurrency
	// issue like sequencer check failures, transaction aborts, etc.
	ErrAborted ErrorCode = "aborted"

	// OutOfRange means operation was attempted past the valid range. For example,
	// seeking or reading past end of a paginated collection.
	//
	// Unlike InvalidArgument, this error indicates a problem that may be fixed if
	// the system state changes (i.e. adding more items to the collection).
	//
	// There is a fair bit of overlap between FailedPrecondition and OutOfRange.
	// We recommend using OutOfRange (the more specific error) when it applies so
	// that callers who are iterating through a space can easily look for an
	// OutOfRange error to detect when they are done.
	ErrOutOfRange ErrorCode = "out_of_range"

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	ErrUnimplemented ErrorCode = "unimplemented"

	// Internal errors. When some invariants expected by the underlying system
	// have been broken. In other words, something bad happened in the library or
	// backend service. Do not confuse with HTTP Internal Server Error; an
	// Internal error could also happen on the client code, i.e. when parsing a
	// server response.
	ErrInternal ErrorCode = "internal"

	// Unavailable indicates the service is currently unavailable. This is a most
	// likely a transient condition and may be corrected by retrying with a
	// backoff.
	ErrUnavailable ErrorCode = "unavailable"

	// DataLoss indicates unrecoverable data loss or corruption.
	ErrDataLoss ErrorCode = "data_loss"

	// ErrNone is the zero-value, is considered an empty error and should not be
	// used.
	ErrNone ErrorCode = ""
)

func HTTPStatusFromErrorCode(code ErrorCode) int {
	switch code {
	case ErrCanceled:
		return 408 // RequestTimeout
	case ErrUnknown:
		return 500 // Internal Server Error
	case ErrInvalidArgument:
		return 400 // BadRequest
	case ErrDeadlineExceeded:
		return 408 // RequestTimeout
	case ErrNotFound:
		return 404 // Not Found
	case ErrBadRoute:
		return 404 // Not Found
	case ErrAlreadyExists:
		return 409 // Conflict
	case ErrPermissionDenied:
		return 403 // Forbidden
	case ErrUnauthenticated:
		return 401 // Unauthorized
	case ErrResourceExhausted:
		return 403 // Forbidden
	case ErrFailedPrecondition:
		return 412 // Precondition Failed
	case ErrAborted:
		return 409 // Conflict
	case ErrOutOfRange:
		return 400 // Bad Request
	case ErrUnimplemented:
		return 501 // Not Implemented
	case ErrInternal:
		return 500 // Internal Server Error
	case ErrUnavailable:
		return 503 // Service Unavailable
	case ErrDataLoss:
		return 500 // Internal Server Error
	case ErrNone:
		return 200 // OK
	default:
		return 0 // Invalid!
	}
}

func IsValidErrorCode(code ErrorCode) bool {
	return HTTPStatusFromErrorCode(code) != 0
}

type rpcErr struct {
	code  ErrorCode
	msg   string
	cause error
}

func (e *rpcErr) Code() ErrorCode {
	return e.code
}

func (e *rpcErr) Msg() string {
	return e.msg
}

func (e *rpcErr) Cause() error {
	return e.cause
}

func (e *rpcErr) Error() string {
	return fmt.Sprintf("webrpc error %s: %s", e.code, e.msg)
}
