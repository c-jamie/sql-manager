package sql

const (
	UNUSED = 0

	SCRIPT       = 5000
	NAME         = 5001
	DESCRIPTION  = 5002
	LAST_UPDATED = 5003
	UPDATED_BY   = 5004
	TEST         = 5005
	DEV          = 5006
	PD_BEGIN     = 5007
	PD_END       = 5008
	PD_REF       = 5009
	PD_MIG       = 5015
	COMMENT      = 5010
	PD_FILE      = 5011
	PROD      	 = 5012
	LOCAL      	 = 5013
	INT 		 = 5014

	AND             = 6001
	OR              = 6002
	NE              = 6003
	SHIFT_LEFT      = 6004
	NULL_SAFE_EQUAL = 6005
	LE              = 6006
	GE              = 6007
	SHIFT_RIGHT     = 6008

	VALUE_ARG   = 7001
	STRING      = 7002
	HEX         = 7003
	FLOAT       = 7004
	BIT_LITERAL = 7005
	LIST_ARG    = 7006
	INTEGRAL    = 7007
	HEXNUM      = 7008

	LEX_ERROR = 8000

	ID = 9000
)
