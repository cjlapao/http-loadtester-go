package domain

import "strings"

// JobOPerationType Enum
type JobOPerationType int

// JobOPerationType Enum Definition
const (
	Constant JobOPerationType = iota
	Increasing
	Fuzzy
	VirtualUser
)

// String Converts JobOPerationType Enum to String
func (t JobOPerationType) String() string {
	switch t {
	case Constant:
		return "constant"
	case Increasing:
		return "increasing"
	case Fuzzy:
		return "fuzzy"
	case VirtualUser:
		return "vu"
	default:
		return "fuzzy"
	}
}

// Get Converts JobOPerationType String to Enum
func (t JobOPerationType) Get(value string) JobOPerationType {
	switch value {
	case "constant":
		return Constant
	case "increasing":
		return Increasing
	case "fuzzy":
		return Fuzzy
	default:
		return Fuzzy
	}
}

// TargetMethod Enum
type TargetMethod int

// TargetMethod Enum Definition
const (
	GET TargetMethod = iota
	POST
	PUT
	DELETE
	OPTIONS
	PATCH
)

// String Converts TargetMethod Enum to String
func (t TargetMethod) String() string {
	switch t {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "PUT"
	case DELETE:
		return "DELETE"
	case OPTIONS:
		return "OPTIONS"
	case PATCH:
		return "PATCH"
	default:
		return "GET"
	}
}

// Get Converts TargetMethod String to Enum
func (t TargetMethod) Get(value string) TargetMethod {
	value = strings.ToUpper(value)
	switch value {
	case "GET":
		return GET
	case "POST":
		return POST
	case "PUT":
		return PUT
	case "DELETE":
		return DELETE
	case "OPTIONS":
		return OPTIONS
	case "PATCH":
		return PATCH
	default:
		return GET
	}
}

// BlockType Enum
type BlockType int

// BlockType Enum Definition
const (
	SequentialBlock BlockType = iota
	ParallelBlock
)

// String Converts the BlockType to string
func (l BlockType) String() string {
	switch l {
	case SequentialBlock:
		return "sequential"
	case ParallelBlock:
		return "parallel"
	default:
		return "parallel"
	}
}

// Get Converts a String to a BlockType
func (l BlockType) Get(value string) BlockType {
	switch value {
	case "sequential":
		return SequentialBlock
	case "parallel":
		return ParallelBlock
	default:
		return ParallelBlock
	}

}
