package quirk

const (
	schemaDelimiter    = '.'
	predicateDelimiter = ':'
)

const (
	msgTooManyMutationFields = "Too many fields filled in QuirkMutation"
	msgInvalidSchemaRead     = "Invalid schema caused reading error"
	msgTransactionFailure    = "Transaction failure"
)

const (
	maxWorkers = 50
)

const tagUnique tagOptions = "unique"

const (
	emptyQuery = "{}"
	whenRDF = `<%s> <when> "%d"^^<xs:int> .`
	blankDefault = "data"
	rdfBase = "_:%s <%s> %q .\n"
	queryfunc = "%s(func: eq(%s, %q), first: 1){uid}\n"
)
