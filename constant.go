package quirk

const (
	// DataTypes are a required piece of the schema.
	dTypeDefault  = "default"
	dTypeInt      = "int"
	dTypeFloat    = "float"
	dTypeString   = "string"
	dTypeBool     = "bool"
	dTypeID       = "id"
	dTypeDateTime = "datetime"
	dTypeGeo      = "geo"
	dTypeUID      = "uid"
	dTypePassword = "password"

	// Tags appear just after the Datatype.
	tagIndex   = "@index"
	tagLang    = "@lang"
	tagReverse = "@reverse"

	// Directives that appear after tags.
	dirUpsert = "@upsert"
	dirCount  = "@count"

	// string specific tokenizers.
	tokenExact    = "exact"
	tokenHash     = "hash"
	tokenTerm     = "term"
	tokenFullText = "fulltext"
	tokenTrigram  = "trigram"

	// datetime specific tokenizers.
	tokenYear  = "year"
	tokenMonth = "month"
	tokenDay   = "day"
	tokenHour  = "hour"
)

const (
	schemaDelimiter    = '.'
	predicateDelimiter = ':'
)

const (
	msgTooManyMutationFields = "Too many fields filled in QuirkMutation"
	msgInvalidSchemaRead     = "Invalid schema caused reading error"
)

const (
	maxWorkers = 50
)

const (
	hash rdfMode = iota
	incrementor
	auto
)
