package quirk

import (
	"bufio"
	"io"
	"strings"
)

func (c *Client) setQuirkSchema(schema string) {
	var tag string
	if c.quirkReverse {
		tag = tagReverse
	}
	c.schemaString = schema + "\n" + c.quirkRel + ":" + " " +
		dTypeUID + " " + tag + " ."
}

func (c *Client) setSchema(schema string) error {
	c.setQuirkSchema(schema)
	r := bufio.NewReader(strings.NewReader(schema))

	for {
		s, err := r.ReadString(schemaDelimiter)
		if err == io.EOF {
			break
		} else if err != nil {
			return &Error{
				ExtErr:   err,
				Msg:      msgInvalidSchemaRead,
				File:     "schema.go",
				Function: "quirk.Client.SetSchema",
			}
		}

		predName, predProps, err := readSchemaLine(strings.TrimSpace(s))
		if err != nil {
			return err
		}

		c.schemaCache[predName] = predProps
	}

	c.logger.Debug("Finished reading schema.")

	return nil
}

func readSchemaLine(s string) (string, Properties, error) {
	r := bufio.NewReader(strings.NewReader(s))

	predName, err := r.ReadBytes(predicateDelimiter)
	if err != nil {
		return "", Properties{}, &Error{
			ExtErr:   err,
			Msg:      msgInvalidSchemaRead,
			File:     "schema.go",
			Function: "quirk.Client.readSchemaLine",
		}
	}

	var predProps Properties

	if strings.Contains(s, dirUpsert) {
		predProps.Upsert = true
	}

	return string(predName[:len(predName)-1]), predProps, nil
}
