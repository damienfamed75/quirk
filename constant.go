package quirk

// defaults.
const (
	// templateDefault is used as a default setup for the Quirk progress bar.
	// Note: This can be changed when setting up the client WithTemplate()
	_templateDefault = `{{ cyan "Inserting Nodes:" }} {{counters .}} {{ bar . "[" "=" (cycle . ">" ) " " "]"}} [{{etime . | cyan }}:{{rtime . | cyan }}] {{percent .}}`
	// maxWorkers is used as a default for the maximum concurrent insert workers
	// that are allowed to run when using a Multi operation.
	_maxWorkers = 50
	// blank identifier default.
	// This is used when the data given doesn't contain any
	// set predicate key in the Quirk client.
	_blankDefault = "data"
	// used as the default identifier when inserting a node.
	// Any node that is inserted with Quirk should have this.
	_predicateKeyDefault = "name"
)

const _tagUnique tagOptions = "unique"

// rdf related constants for building the mutation.
const (
	_rdfBase      = "%s <%s> \"%v\""
	_rdfReference = "%s <%s> <%v>"
	_rdfEnd       = " .\n"
)

const (
	// quirkTag is for identifying a structure tag.
	// Example:
	// type MyStruct struct {
	// 		field string `quirk:"field"`
	// }
	_quirkTag = "quirk"
	// Used to see if the query building process returned an empty query.
	_emptyQuery = "{}"
	// The beginning of the GraphQL+ query function signature.
	_queryfunc = "%s(func: eq(%s, %q), first: 1){uid}\n"
)

const (
	// xsInit is used to indicate to Dgraph that we are explicitly
	// using a certain datatype in the RDF.
	_xsInit = "^^"

	// XML Datatypes.
	_xsInt   = _xsInit + "<xs:int>"
	_xsBool  = _xsInit + "<xs:boolean>"
	_xsFloat = _xsInit + "<xs:float>"

	// unused at the moment.
	_xsString   = _xsInit + "<xs:string>"
	_xsDateTime = _xsInit + "<xs:date>"

	// notifier to fix byte slice.
	_xsByte = _xsInit + "<xs:byte>"
)
