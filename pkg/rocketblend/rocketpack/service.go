package rocketpack

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/flowshot-io/x/pkg/logger"
	"github.com/rocketblend/rocketblend/pkg/jot"
	"github.com/rocketblend/rocketblend/pkg/jot/reference"
	"github.com/rocketblend/rocketblend/pkg/rocketblend/runtime"
)

const PackgeFile = "rocketpack.yaml"

type (
	Service interface {
		DescribeByReference(reference reference.Reference) (*RocketPack, error)
		FindByReference(ref reference.Reference) (*RocketPack, error)
		InstallByReference(reference reference.Reference, force bool) error
		InstallByReferenceWithContext(ctx context.Context, reference reference.Reference, force bool) error
		UninstallByReference(reference reference.Reference) error
	}

	ServiceOptions struct {
		Storage  jot.Storage
		Logger   logger.Logger
		Platform runtime.Platform
	}

	ServiceOption func(*ServiceOptions)

	service struct {
		logger   logger.Logger
		storage  jot.Storage
		platform runtime.Platform
	}
)

func WithLogger(logger logger.Logger) ServiceOption {
	return func(o *ServiceOptions) {
		o.Logger = logger
	}
}

func WithStorage(storage jot.Storage) ServiceOption {
	return func(o *ServiceOptions) {
		o.Storage = storage
	}
}

func WithPlatform(platform runtime.Platform) ServiceOption {
	return func(o *ServiceOptions) {
		o.Platform = platform
	}
}

func NewService(opts ...ServiceOption) (Service, error) {
	options := &ServiceOptions{
		Logger:   logger.NoOp(),
		Platform: runtime.Undefined,
	}

	for _, opt := range opts {
		opt(options)
	}

	if options.Storage == nil {
		return nil, fmt.Errorf("storage is required")
	}

	return &service{
		logger:   options.Logger,
		storage:  options.Storage,
		platform: options.Platform,
	}, nil
}

func (srv *service) DescribeByReference(reference reference.Reference) (*RocketPack, error) {
	url, err := url.JoinPath(reference.Url(), PackgeFile)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pack, err := load(bodyBytes)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func (srv *service) InstallByReference(ref reference.Reference, force bool) error {
	return srv.InstallByReferenceWithContext(context.Background(), ref, force)
}

func (srv *service) InstallByReferenceWithContext(ctx context.Context, ref reference.Reference, force bool) error {
	// Check if already installed.
	pack, _ := srv.FindByReference(ref)

	// Pack found but force is true, delete it.
	if pack != nil && force {
		err := srv.storage.DeleteAll(ref)
		if err != nil {
			return err
		}

		pack = nil
	}

	// Pack found is a build pack, also check it's addons.
	if pack != nil && pack.Build != nil {
		err := srv.installBuildAddons(ctx, pack.Build, force)
		if err != nil {
			return err
		}

		return nil
	}

	// Pack was not found installed, try to install it.
	if pack == nil {
		fmt.Println(ref.String())

		err := srv.fetchByReference(ref)
		if err != nil {
			return err
		}

		err = srv.pullByReference(ctx, ref, force)
		if err != nil {
			return err
		}
	}

	return nil
}

func (srv *service) UninstallByReference(ref reference.Reference) error {
	_, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	if err := srv.storage.DeleteAll(ref); err != nil {
		return err
	}

	return nil
}

func (srv *service) FindByReference(ref reference.Reference) (*RocketPack, error) {
	bytes, err := srv.storage.Read(ref, PackgeFile)
	if err != nil {
		return nil, err
	}

	pack, err := load(bytes)
	if err != nil {
		return nil, err
	}

	return pack, err
}

func (srv *service) fetchByReference(ref reference.Reference) error {
	// Validates reference is a valid pack.
	_, err := srv.DescribeByReference(ref)
	if err != nil {
		return err
	}

	downloadUrl, err := url.JoinPath(ref.Url(), PackgeFile)
	if err != nil {
		return err
	}

	err = srv.storage.Write(ref, PackgeFile, downloadUrl)
	if err != nil {
		return err
	}

	return nil
}

func (srv *service) pullByReference(ctx context.Context, ref reference.Reference, force bool) error {
	pack, err := srv.FindByReference(ref)
	if err != nil {
		return err
	}

	if pack.Addon != nil {
		return srv.writeAddon(ctx, ref, pack.Addon)
	}

	if pack.Build != nil {
		err := srv.writeBuild(ctx, ref, pack.Build)
		if err != nil {
			return err
		}

		err = srv.installBuildAddons(ctx, pack.Build, force)
		if err != nil {
			return err
		}

		return nil
	}

	return fmt.Errorf("no build or addon found in rocketpack %s", ref)
}

func (srv *service) writeAddon(ctx context.Context, ref reference.Reference, addon *Addon) error {
	// Don't write if no source is provided. Addon might be preinstalled or local only.
	if addon.Source == nil || addon.Source.URL == "" {
		return nil
	}

	err := srv.storage.WriteWithContext(ctx, ref, addon.Source.File, addon.Source.URL)
	if err != nil {
		return err
	}

	return nil
}

func (srv *service) writeBuild(ctx context.Context, ref reference.Reference, build *Build) error {
	source := build.GetSourceForPlatform(srv.platform)
	if source == nil {
		return fmt.Errorf("no source found for platform %s", (srv.platform))
	}

	err := srv.storage.WriteWithContext(ctx, ref, jot.GetFilenameFromURL(source.URL), source.URL)
	if err != nil {
		return err
	}

	return nil
}

func (srv *service) installBuildAddons(ctx context.Context, build *Build, force bool) error {
	for _, pack := range build.Addons {
		err := srv.InstallByReferenceWithContext(ctx, reference.Reference(pack), force)
		if err != nil {
			return err
		}
	}

	return nil
}
