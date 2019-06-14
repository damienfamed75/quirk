package quirk

const (
	msgTooManyMutationFields = "Too many fields filled in QuirkMutation"
	msgInvalidSchemaRead     = "Invalid schema caused reading error"
	msgTransactionFailure    = "Transaction failure"
	msgBuilderWriting        = "invalid pred[%#v] or val[%#v]"
	msgTooManyResponses      = "Too many responses from query for unique nodes"
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
