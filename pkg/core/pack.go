package core

import (
	"github.com/rocketblend/rocketblend/pkg/core/rocketpack"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
)

func (d *Driver) DescribePackByReference(reference reference.Reference) (*rocketpack.RocketPack, error) {
	pack, err := d.pack.DescribeByReference(reference)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (d *Driver) FetchPackByReference(ref reference.Reference) error {
	err := d.pack.FetchByReference(ref)
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) PullPackByReference(ref reference.Reference) error {
	err := d.pack.PullByReference(ref)
	if err != nil {
		return err
	}

	return nil
}

func (d *Driver) FindPackByReference(ref reference.Reference) (*rocketpack.RocketPack, error) {
	pack, err := d.pack.FindByReference(ref)
	if err != nil {
		return nil, err
	}

	return pack, nil
}
