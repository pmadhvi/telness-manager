package service

import "errors"

var (
	ErrInvalidIBANEmpty               = errors.New("Invalid IBAN, iban cannot be empty.")
	ErrInvalidIBANMinimumLength       = errors.New("Invalid IBAN, minimum length for IBAN should be 5.")
	ErrWrongInvalidIBANCountryCode    = errors.New("Invalid IBAN, wrong country code.")
	ErrInvalidIBANCountryCodeNotFound = errors.New("Invalid IBAN, country code not found.")
	ErrInvalidIBANCountryFormat       = errors.New("Invalid IBAN, IBAN format for country does not match.")
	ErrInvalidIBANChecksum            = errors.New("Invalid IBAN, wrong checksum for IBAN.")
)
