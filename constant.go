package quirk

const (
	msgTooManyMutationFields = "Too many fields filled in QuirkMutation"
	msgInvalidSchemaRead     = "Invalid schema caused reading error"
	msgTransactionFailure    = "Transaction failure"
	msgQueryingUnique        = "Dgraph querying error"
	msgBuilderWriting        = "invalid pred[%#v] or val[%#v]"
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
