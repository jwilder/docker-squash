package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"path/filepath"
)

type Layer struct {
	Id      string
	TarPath string
	TarHash string
}

type JsonFile struct {
	data *simplejson.Json
	path string
}

type Manifestor struct {
	basePath string
	tag      string
	layers   []Layer
	export   *Export
	manifest *JsonFile
	config   *JsonFile
}

func NewManifestor(e *Export, tag string) (*Manifestor, error) {
	manifest := &JsonFile{
		path: filepath.Join(e.Path, "manifest.json"),
	}

	manifestContent, err := ioutil.ReadFile(manifest.path)
	if err != nil {
		return nil, err
	}

	manifest.data, err = simplejson.NewJson(manifestContent)
	if err != nil {
		return nil, err
	}

	configPath, err := manifest.data.GetIndex(0).Get("Config").String()
	if err != nil {
		return nil, err
	}

	config := &JsonFile{
		path: filepath.Join(e.Path, configPath),
	}

	configContent, err := ioutil.ReadFile(config.path)
	if err != nil {
		return nil, err
	}

	config.data, err = simplejson.NewJson(configContent)
	if err != nil {
		return nil, err
	}

	manifestor := &Manifestor{
		tag:      tag,
		basePath: e.Path,
		manifest: manifest,
		config:   config,
		export:   e,
	}

	return manifestor, nil
}

func (m *Manifestor) GenerateLayers() {
	current := m.export.Root()
	order := []*ExportedImage{}
	for {
		order = append(order, current)
		current = m.export.ChildOf(current.LayerConfig.Id)
		if current == nil {
			break
		}
	}

	m.layers = []Layer{}

	for i := 0; i < len(order); i++ {
		data, err := ioutil.ReadFile(order[i].LayerTarPath)

		if err != nil {
			debug("Cannot read layer file: ", order[i].LayerTarPath, err)
			return
		}

		sha := fmt.Sprintf("%x", sha256.Sum256(data))

		layer := Layer{}
		layer.Id = order[i].LayerConfig.Id
		layer.TarPath = order[i].LayerConfig.Id + "/layer.tar"
		layer.TarHash = sha

		m.layers = append(m.layers, layer)
	}
}

func (m *Manifestor) UpdateManifest() {
	newLayers := []string{}
	for i := 0; i < len(m.layers); i++ {
		newLayers = append(newLayers, m.layers[i].TarPath)
	}

	// Replace the layers section
	m.manifest.data.GetIndex(0).Set("Layers", newLayers)

	if m.tag != "" {
		m.manifest.data.GetIndex(0).Set("RepoTags", []string{m.tag})
	}
}

func (m *Manifestor) UpdateConfig() {
	newDiffs := []string{}
	for i := 0; i < len(m.layers); i++ {
		newDiffs = append(newDiffs, "sha256:"+m.layers[i].TarHash)
	}

	// Set new hashes and clear history
	m.config.data.Get("rootfs").Set("diff_ids", newDiffs)
	m.config.data.Set("history", []string{})
}

func (m *Manifestor) SaveChanges() error {
	newManifest, err := m.manifest.data.MarshalJSON()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(m.manifest.path, newManifest, 0644)
	if err != nil {
		return err
	}

	newConfig, err := m.config.data.MarshalJSON()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(m.config.path, newConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}
