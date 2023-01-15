package core

func (d *Driver) Initialize() error {
	err := d.resource.SaveOut()
	if err != nil {
		return err
	}

	return nil
}
