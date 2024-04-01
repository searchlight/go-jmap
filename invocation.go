package jmap

import (
	"encoding/json"
	"errors"
	"fmt"
)

// An Invocation represents method calls and responses
type Invocation struct {
	// The name of the method call or response
	Name string
	// Object containing the named arguments for the method or response
	Args interface{}
	// Arbitrary string set by client, echoed back with responses
	CallID string
}

func (i *Invocation) MarshalJSON() ([]byte, error) {
	j := []interface{}{
		i.Name,
		i.Args,
		i.CallID,
	}

	marshalled, err := json.Marshal(j)
	if err != nil {
		return []byte{}, err
	}
	if i.Name != "Email/set" {
		return marshalled, nil
	}

	var unMarshalled []interface{}
	if tmpErr := json.Unmarshal(marshalled, &unMarshalled); tmpErr != nil {
		return []byte{}, tmpErr
	}

	fmt.Println(unMarshalled)

	firstElement, ok := unMarshalled[1].(map[string]interface{})
	if !ok {
		return []byte{}, errors.New("Unable to parse Email/set method")
	}

	create, ok := firstElement["create"].(map[string]interface{})
	if !ok {
		return []byte{}, errors.New("Unable to find key `create`")
	}

	draft, ok := create["draft"].(map[string]interface{})
	if !ok {
		return []byte{}, errors.New("Unable to find key `draft`")
	}

	customHeaders, ok := draft["customHeaders"].([]*struct {
		Name  string
		value string
	})

	if !ok {
		return marshalled, nil
	}

	draft["customHeaders"] = customHeaders
	create["draft"] = draft
	unMarshalled[1] = create

	for _, v := range customHeaders {
		fmt.Println(v)
	}

	return []byte{}, nil
}

func (i *Invocation) UnmarshalJSON(data []byte) error {
	raw := []json.RawMessage{}

	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	if len(raw) != 3 {
		return fmt.Errorf("Not enough values in invocation")
	}
	if err := json.Unmarshal(raw[0], &i.Name); err != nil {
		return err
	}
	newFn, ok := methods[i.Name]
	if !ok {
		return fmt.Errorf("method '%s' not registered", i.Name)
	}
	i.Args = newFn()
	if err := json.Unmarshal(raw[1], i.Args); err != nil {
		return err
	}
	if err := json.Unmarshal(raw[2], &i.CallID); err != nil {
		return err
	}
	return nil
}
