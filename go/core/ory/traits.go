package ory

type Traits map[string]interface{}

func (t Traits) Email() string {
	return t["email"].(string)
}

func (t Traits) FirstName() string {
	return t["first_name"].(string)
}

func (t Traits) LastName() string {
	res, ok := t["last_name"]
	if ok {
		return res.(string)
	}
	return ""
}

func (t Traits) Domain() string {
	res, ok := t["hd"]
	if ok {
		return res.(string)
	}
	return ""
}

func (t Traits) Profile() string {
	res, ok := t["profile"]
	if ok {
		return res.(string)
	}
	return ""
}
