package main

import "fmt"

type brick struct {
	/*
		1. len(sources) == 1
		   sources will be tagged which destination defined and push then
		2. len(source) > 1
		   sources will be combined into one manifest and push then
	*/

	// sources can be multi of an image
	Sources []*source `json:"sources"`

	// manifest generic private register like: foo.register/bar-v1.0
	// used to combined multi digest to one manifest
	Manifest string `json:"manifest"`

	// Auto will automatically lift all images under this sources
	// which length should equal 1
	Auto bool `json:"auto"`

	// Amend will amend remote manifest if true
	Amend bool `json:"amend"`

	// Insecure can access a insecure remote registry
	Insecure bool `json:"insecure"`
}

// pushManifest would failed if remote register already had such manifest
func (b *brick) pushManifest() cmdInstruction {
	cmdInst := makeCmdInstruction("docker manifest push %s", b.Manifest)
	if b.Insecure {
		cmdInst.AppendOption("--insecure")
	}
	return cmdInst
}

func (b *brick) createManifest() error {
	digests := ""

	for _, s := range b.Sources {
		if s.Skip {
			continue
		}

		digests += fmt.Sprintf(" %s", s.NewTag)
	}

	if len(digests) < 1 {
		return fmt.Errorf("digest less than 1, can not create manifest %s", b.Manifest)
	}

	cmdInst := makeCmdInstruction("docker manifest create %s%s", b.Manifest, digests)

	if b.Amend {
		cmdInst.AppendOption("--amend")
	}

	if b.Insecure {
		cmdInst.AppendOption("--insecure")
	}

	return cmdInst.doExec()
}

func (b *brick) moving() error {
	if len(b.Sources) == 1 && b.Auto == false {
		return singleMove(b.Sources[0])
	}

	return multiMove(b)
}

// multiMove can move multi images at same time
// if manifest specified, would combine multi images
// into one manifest
func multiMove(b *brick) error {
	var err error

	// do auto multi move if need
	if b.Auto && len(b.Sources) == 1 {
		if err = autoBuildSources(b); err != nil {
			return err
		}
	}

	// do truly images move 1 by 1
	for _, s := range b.Sources {
		err = singleMove(s)
		if err != nil {
			return err
		}
	}

	if len(b.Manifest) > 0 {
		err = b.createManifest()
		if err != nil {
			return err
		}

		err = b.pushManifest().doExec()
		if err != nil {
			return err
		}
	}

	return nil
}

// autoBuildSources automatically build resources corresponding to all digest under given manifest
func autoBuildSources(b *brick) error {
	if len(b.Sources) != 1 {
		return fmt.Errorf("len(sources) when auto move must be 1")
	}

	s := b.Sources[0]

	if s.Skip {
		return nil
	}

	m := manifest{images: make([]image, 0)}

	cmdInst := s.inspect()
	if b.Insecure {
		cmdInst.AppendOption("--insecure")
	}

	// query manifest of given image
	err := cmdInst.doExecInto(&m.images)
	if err != nil {
		return err
	}

	// rebuild sources for each digest
	rebuildSource(b, &m)

	return nil
}

type manifest struct {
	images []image
}

type image struct {
	Ref              string      `json:"Ref"`
	Descriptor       descriptor  `json:"Descriptor"`
	SchemaV2Manifest interface{} `json:"SchemaV2Manifest"`
}

type descriptor struct {
	MediaType string                 `json:"mediaType"`
	Digest    string                 `json:"digest"`
	Size      int                    `json:"size"`
	Platform  map[string]interface{} `json:"platform"`
}

func rebuildSource(b *brick, m *manifest) {
	sources := make([]*source, 0, len(m.images))

	for i, img := range m.images {
		arch, ok := img.Descriptor.Platform["architecture"]
		if !ok {
			continue
		}

		s := &source{}
		s.Addr = img.Ref
		s.NewTag = fmt.Sprintf("%v-%v-%v", b.Manifest, i, arch)
		sources = append(sources, s)
	}

	b.Sources = sources
}
