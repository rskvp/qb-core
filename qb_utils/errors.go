package qb_utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ErrorsHelper struct {
}

var Errors *ErrorsHelper

func init(){
	Errors = new(ErrorsHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	E r r o r
//----------------------------------------------------------------------------------------------------------------------

type Error struct {
	Errors      []error
	ErrorFormat ErrorFormatFunc
}

func (instance *Error) Error() string {
	fn := instance.ErrorFormat
	if fn == nil {
		fn = Errors.ListFormatFunc
	}

	return fn(instance.Errors)
}

// ErrorOrNil returns an error interface if this Error represents
// a list of errors, or returns nil if the list of errors is empty. This
// function is useful at the end of accumulation to make sure that the value
// returned represents the existence of errors.
func (instance *Error) ErrorOrNil() error {
	if instance == nil {
		return nil
	}
	if len(instance.Errors) == 0 {
		return nil
	}

	return instance
}

func (instance *Error) GoString() string {
	return fmt.Sprintf("*%#v", *instance)
}

// WrappedErrors returns the list of errors that this Error is wrapping.
// It is an implementation of the errwrap.Wrapper interface so that
// lygo_errors.Error can be used with that library.
//
// This method is not safe to be called concurrently and is no different
// than accessing the Errors field directly. It is implemented only to
// satisfy the errwrap.Wrapper interface.
func (instance *Error) WrappedErrors() []error {
	return instance.Errors
}

// Unwrap returns an error from Error (or nil if there are no errors).
// This error returned will further support Unwrap to get the next error,
// etc. The order will match the order of Errors in the lygo_errors.Error
// at the time of calling.
//
// The resulting error supports errors.As/Is/Unwrap so you can continue
// to use the stdlib errors package to introspect further.
//
// This will perform a shallow copy of the errors slice. Any errors appended
// to this error after calling Unwrap will not be available until a new
// Unwrap is called on the lygo_errors.Error.
func (instance *Error) Unwrap() error {
	// If we have no errors then we do nothing
	if instance == nil || len(instance.Errors) == 0 {
		return nil
	}

	// If we have exactly one error, we can just return that directly.
	if len(instance.Errors) == 1 {
		return instance.Errors[0]
	}

	// Shallow copy the slice
	errs := make([]error, len(instance.Errors))
	copy(errs, instance.Errors)
	return chain(errs)
}

//----------------------------------------------------------------------------------------------------------------------
//	chain
//----------------------------------------------------------------------------------------------------------------------

// chain implements the interfaces necessary for errors.Is/As/Unwrap to
// work in a deterministic way with lygo_errors. A chain tracks a list of
// errors while accounting for the current represented error. This lets
// Is/As be meaningful.
//
// Unwrap returns the next error. In the cleanest form, Unwrap would return
// the wrapped error here but we can't do that if we want to properly
// get access to all the errors. Instead, users are recommended to use
// Is/As to get the correct error type out.
//
// Precondition: []error is non-empty (len > 0)
type chain []error

// Error implements the error interface
func (instance chain) Error() string {
	return instance[0].Error()
}

// Unwrap implements errors.Unwrap by returning the next error in the
// chain or nil if there are no more errors.
func (instance chain) Unwrap() error {
	if len(instance) == 1 {
		return nil
	}

	return instance[1:]
}

// As implements errors.As by attempting to map to the current value.
func (instance chain) As(target interface{}) bool {
	return errors.As(instance[0], target)
}

// Is implements errors.Is by comparing the current value directly.
func (instance chain) Is(target error) bool {
	return errors.Is(instance[0], target)
}

//----------------------------------------------------------------------------------------------------------------------
//	constructors
//----------------------------------------------------------------------------------------------------------------------


// New returns an error with the supplied message.
// New also records the stack trace at the point it was called.
func (instance *ErrorsHelper) NewError(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// Errorf formats according to a format specifier and returns the string
// as a value that satisfies error.
// Errorf also records the stack trace at the point it was called.
func (instance *ErrorsHelper) Errorf(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}


// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func  (instance *ErrorsHelper) WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}

// WrapWithStack returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func  (instance *ErrorsHelper) WrapWithStack(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	return &withStack{
		err,
		callers(),
	}
}

// WrapWithStackf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func  (instance *ErrorsHelper) WrapWithStackf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return &withStack{
		err,
		callers(),
	}
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func  (instance *ErrorsHelper) WithMessage(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func  (instance *ErrorsHelper) WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	f o r m a t t i n g
//----------------------------------------------------------------------------------------------------------------------

// ErrorFormatFunc turn the list of errors into a string.
type ErrorFormatFunc func([]error) string

// ListFormatFunc is a basic formatter that outputs the number of errors
// that occurred along with a bullet point list of the errors.
func (instance *ErrorsHelper) ListFormatFunc(list []error) string {
	if len(list) == 1 {
		return fmt.Sprintf("1 error occurred:\n\t* %s\n\n", list[0])
	}

	points := make([]string, len(list))
	for i, err := range list {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d errors occurred:\n\t%s\n\n",
		len(list), strings.Join(points, "\n\t"))
}

//----------------------------------------------------------------------------------------------------------------------
//	Append
//----------------------------------------------------------------------------------------------------------------------

// Append is a helper function that will append more errors
// onto an Error in order to create a larger multi-error.
//
// If err is not a lygo_errors.Error, then it will be turned into
// one. If any of the errs are multierr.Error, they will be flattened
// one level into err.
// Any nil errors within errs will be ignored. If err is nil, a new
// *Error will be returned.
func (instance *ErrorsHelper) Append(err error, errs ...error) *Error {
	switch err := err.(type) {
	case *Error:
		// Typed nils can reach here, so initialize if we are nil
		if err == nil {
			err = new(Error)
		}

		// Go through each error and flatten
		for _, e := range errs {
			switch e := e.(type) {
			case *Error:
				if e != nil {
					err.Errors = append(err.Errors, e.Errors...)
				}
			default:
				if e != nil {
					err.Errors = append(err.Errors, e)
				}
			}
		}

		return err
	default:
		newErrs := make([]error, 0, len(errs)+1)
		if err != nil {
			newErrs = append(newErrs, err)
		}
		newErrs = append(newErrs, errs...)

		// recursive
		return instance.Append(&Error{}, newErrs...)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	Flatten
//----------------------------------------------------------------------------------------------------------------------


// Flatten flattens the given error, merging any *Errors together into
// a single *Error.
func (instance *ErrorsHelper) Flatten(err error) error {
	// If it isn't an *Error, just return the error as-is
	if _, ok := err.(*Error); !ok {
		return err
	}

	// Otherwise, make the result and flatten away!
	flatErr := new(Error)
	flatten(err, flatErr)
	return flatErr
}

func flatten(err error, flatErr *Error) {
	switch err := err.(type) {
	case *Error:
		for _, e := range err.Errors {
			flatten(e, flatErr)
		}
	default:
		flatErr.Errors = append(flatErr.Errors, err)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	Wrap
//----------------------------------------------------------------------------------------------------------------------

// WalkFunc is the callback called for Walk.
type WalkFunc func(error)

// Wrapper is an interface that can be implemented by custom types to
// have all the Contains, Get, etc. functions in errwrap work.
//
// When Walk reaches a Wrapper, it will call the callback for every
// wrapped error in addition to the wrapper itself. Since all the top-level
// functions in errwrap use Walk, this means that all those functions work
// with your custom type.
type Wrapper interface {
	WrappedErrors() []error
}

// wrappedError is an implementation of error that has both the
// outer and inner errors.
type wrappedError struct {
	Outer error
	Inner error
}

func (w *wrappedError) Error() string {
	return w.Outer.Error()
}

func (w *wrappedError) WrappedErrors() []error {
	return []error{w.Outer, w.Inner}
}

// Wrap defines that outer wraps inner, returning an error type that
// can be cleanly used with the other methods in this package, such as
// Contains, GetAll, etc.
//
// This function won't modify the error message at all (the outer message
// will be used).
func (instance *ErrorsHelper) Wrap(outer, inner error) error {
	return &wrappedError{
		Outer: outer,
		Inner: inner,
	}
}

// Wrapf wraps an error with a formatting message. This is similar to using
// `fmt.Errorf` to wrap an error. If you're using `fmt.Errorf` to wrap
// errors, you should replace it with this.
//
// format is the format of the error message. The string '{{err}}' will
// be replaced with the original error message.
func (instance *ErrorsHelper) Wrapf(format string, err error) error {
	outerMsg := "<nil>"
	if err != nil {
		outerMsg = err.Error()
	}

	outer := errors.New(strings.Replace(
		format, "{{err}}", outerMsg, -1))

	return instance.Wrap(outer, err)
}

// Contains checks if the given error contains an error with the
// message msg. If err is not a wrapped error, this will always return
// false unless the error itself happens to match this msg.
func (instance *ErrorsHelper) Contains(err error, msg string) bool {
	return len(instance.GetAll(err, msg)) > 0
}

// ContainsType checks if the given error contains an error with
// the same concrete type as v. If err is not a wrapped error, this will
// check the err itself.
func (instance *ErrorsHelper) ContainsType(err error, v interface{}) bool {
	return len(instance.GetAllType(err, v)) > 0
}

// Get is the same as GetAll but returns the deepest matching error.
func (instance *ErrorsHelper) Get(err error, msg string) error {
	es := instance.GetAll(err, msg)
	if len(es) > 0 {
		return es[len(es)-1]
	}

	return nil
}

// GetType is the same as GetAllType but returns the deepest matching error.
func (instance *ErrorsHelper) GetType(err error, v interface{}) error {
	es := instance.GetAllType(err, v)
	if len(es) > 0 {
		return es[len(es)-1]
	}

	return nil
}

// GetAll gets all the errors that might be wrapped in err with the
// given message. The order of the errors is such that the outermost
// matching error (the most recent wrap) is index zero, and so on.
func (instance *ErrorsHelper) GetAll(err error, msg string) []error {
	var result []error

	instance.Walk(err, func(err error) {
		if err.Error() == msg {
			result = append(result, err)
		}
	})

	return result
}

// GetAllType gets all the errors that are the same type as v.
//
// The order of the return value is the same as described in GetAll.
func (instance *ErrorsHelper) GetAllType(err error, v interface{}) []error {
	var result []error

	var search string
	if v != nil {
		search = reflect.TypeOf(v).String()
	}
	instance.Walk(err, func(err error) {
		var needle string
		if err != nil {
			needle = reflect.TypeOf(err).String()
		}

		if needle == search {
			result = append(result, err)
		}
	})

	return result
}

// Walk walks all the wrapped errors in err and calls the callback. If
// err isn't a wrapped error, this will be called once for err. If err
// is a wrapped error, the callback will be called for both the wrapper
// that implements error as well as the wrapped error itself.
func (instance *ErrorsHelper) Walk(err error, cb WalkFunc) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *wrappedError:
		cb(e.Outer)
		instance.Walk(e.Inner, cb)
	case Wrapper:
		cb(err)

		for _, err := range e.WrappedErrors() {
			instance.Walk(err, cb)
		}
	default:
		cb(err)
	}
}
//----------------------------------------------------------------------------------------------------------------------
//	Prefix
//----------------------------------------------------------------------------------------------------------------------

// Prefix is a helper function that will prefix some text
// to the given error. If the error is a lygo_errors.Error, then
// it will be prefixed to each wrapped error.
//
// This is useful to use when appending multiple lygo_errors
// together in order to give better scoping.
func (instance *ErrorsHelper) Prefix(err error, prefix string) error {
	if err == nil {
		return nil
	}

	format := fmt.Sprintf("%s {{err}}", prefix)
	switch err := err.(type) {
	case *Error:
		// Typed nils can reach here, so initialize if we are nil
		if err == nil {
			err = new(Error)
		}

		// Wrap each of the errors
		for i, e := range err.Errors {
			err.Errors[i] = instance.Wrapf(format, e)
		}

		return err
	default:
		return instance.Wrapf(format, err)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	Cause
//----------------------------------------------------------------------------------------------------------------------

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func (instance *ErrorsHelper) Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}


//----------------------------------------------------------------------------------------------------------------------
//	Sort
//----------------------------------------------------------------------------------------------------------------------

// Len implements sort.Interface function for length
func (instance Error) Len() int {
	return len(instance.Errors)
}

// Swap implements sort.Interface function for swapping elements
func (instance Error) Swap(i, j int) {
	instance.Errors[i], instance.Errors[j] = instance.Errors[j], instance.Errors[i]
}

// Less implements sort.Interface function for determining order
func (instance Error) Less(i, j int) bool {
	return instance.Errors[i].Error() < instance.Errors[j].Error()
}

