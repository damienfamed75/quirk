package quirk

// error messages.
const (
	msgTooManyMutationFields = "Too many fields filled in QuirkMutation"
	msgTransactionFailure    = "Transaction failure"
	msgInvalidSchemaRead     = "Invalid schema caused reading error"
	msgTooManyResponses      = "Too many responses from query for unique nodes"
	msgMutationHadNoUID      = "UID was not found in the mutation response"
	msgBuilderWriting        = "invalid pred[%#v] or val[%#v]"
	msgNilUID                = "*string was nil in response"
)

// defaults.
const (
	// templateDefault is used as a default setup for the Quirk progress bar.
	// Note: This can be changed when setting up the client WithTemplate()
	templateDefault = `{{ cyan "Inserting Nodes:" }} {{counters .}} {{ bar . "[" "=" (cycle . ">" ) " " "]"}} [{{etime . | cyan }}:{{rtime . | cyan }}] {{percent .}}`
	// maxWorkers is used as a default for the maximum concurrent insert workers
	// that are allowed to run when using a Multi operation.
	maxWorkers = 50
	// blank identifier default.
	// This is used when the data given doesn't contain any
	// set predicate key in the Quirk client.
	blankDefault = "data"
)

const tagUnique tagOptions = "unique"

// rdf related constants for building the mutation.
const (
	rdfBase      = "%s <%s> \"%v\""
	rdfReference = "%s <%s> <%v>"
	rdfEnd       = " .\n"
)

const (
	// quirkTag is for identifying a structure tag.
	// Example:
	// type MyStruct struct {
	// 		field string `quirk:"field"`
	// }
	quirkTag = "quirk"
	// Used to see if the query building process returned an empty query.
	emptyQuery = "{}"
	// The beginning of the GraphQL+ query function signature.
	queryfunc = "%s(func: eq(%s, %q), first: 1){uid}\n"
)

const (
	// xsInit is used to indicate to Dgraph that we are explicitly
	// using a certain datatype in the RDF.
	xsInit = "^^"

	// XML Datatypes.
	xsInt   = xsInit + "<xs:int>"
	xsBool  = xsInit + "<xs:boolean>"
	xsFloat = xsInit + "<xs:float>"

	// unused at the moment.
	xsString   = xsInit + "<xs:string>"
	xsDateTime = xsInit + "<xs:date>"

	// notifier to fix byte slice.
	xsByte = xsInit + "<xs:byte>"
)
