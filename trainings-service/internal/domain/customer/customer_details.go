package customer

type Details struct {
	uuid string
	name string
}

func (c *Details) UUID() string {
	return c.uuid
}

func (c *Details) Name() string {
	return c.name
}

func NewCustomerDetails(UUID, name string) (Details, error) {
	if UUID == "" {
		return Details{}, ErrEmptyCustomerUUID
	}
	if name == "" {
		return Details{}, ErrEmptyCustomerName
	}
	c := Details{
		uuid: UUID,
		name: name,
	}
	return c, nil
}

func UnmarshalCustomerDetails(UUID, name string) (Details, error) {
	c, err := NewCustomerDetails(UUID, name)
	if err != nil {
		return Details{}, err
	}
	return c, nil
}
