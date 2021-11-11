package utils

type CustomValidation struct {
	RegexStr    string
	Tag         string
	Message     string
	Description string
}

var CUSTOMVALIDATIONS = []CustomValidation{
	{
		RegexStr:    "^[a-zA-Z0-9-_() ]+$",
		Tag:         "alphanumdashbraces",
		Message:     "{0} should only have alphanumeric, (, ), -, _, spaces",
		Description: "string should be alphanum and includes (), -, _)",
	},
	{
		RegexStr:    "^[a-zA-Z0-9-_ ]+$",
		Tag:         "alphadash",
		Message:     "{0} should only have alphanumeric, -, _",
		Description: "string should be alphanum and includes -, _",
	},
	{
		RegexStr:    "^[a-zA-Z0-9-_ &]+$",
		Tag:         "alphadashampersand",
		Message:     "{0} should only have alphanumeric, -, _, and &",
		Description: "string should be alphanum and includes -, _, and &",
	},
	{
		RegexStr:    "^\\S+$",
		Tag:         "nowhitespace",
		Message:     "{0} should not have spaces",
		Description: "string should not have white space",
	},
	{
		RegexStr:    "^[a-zA-Z0-9_]+$",
		Tag:         "alphanumericunderscores",
		Message:     "{0} should only have alphanumeric, and underscores",
		Description: "string should be alphanum and underscores",
	},
	{
		RegexStr:    "^[a-zA-Z0-9.-]+$",
		Tag:         "alphanumdashdot",
		Message:     "{0} should only have alphanumeric, -, .",
		Description: "string should be alphanum and includes -, .",
	},
	{
		RegexStr:    "^[a-zA-Z0-9].*[a-zA-Z0-9]$",
		Tag:         "alphanumbound",
		Message:     "{0} should start and end with alphanumeric",
		Description: "string should start and end with alphanumeric",
	},
}
