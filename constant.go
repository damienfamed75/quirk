package quirk

const (
	msgTooManyMutationFields = "Too many fields filled in QuirkMutation"
	msgTransactionFailure    = "Transaction failure"
	msgInvalidSchemaRead     = "Invalid schema caused reading error"
	msgTooManyResponses      = "Too many responses from query for unique nodes"
	msgMutationHadNoUID      = "UID was not found in the mutation response"
	msgBuilderWriting        = "invalid pred[%#v] or val[%#v]"
	msgNilUID                = "*string was nil in response"
)

const (
	maxWorkers = 50
)

const tagUnique tagOptions = "unique"

const (
	quirkTag     = "quirk"
	emptyQuery   = "{}"
	blankDefault = "data"
	whenRDF      = `<%s> <when> "%d"^^<xs:int> .`
	rdfBase      = "_:%s <%s> %q .\n"
	queryfunc    = "%s(func: eq(%s, %q), first: 1){uid}\n"
)
