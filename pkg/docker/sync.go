package docker

import (
	"github.com/vinkdong/gox/slice"
	"regexp"
	"strings"
	"fmt"
)

type NamedSync struct {
	ApiVersion string `yaml:"apiVersion"`
	Sync       Sync   `yaml:"sync"`
}

type Sync struct {
	From    Docker      `yaml:"from"`
	To      Docker      `yaml:"to"`
	Names   []string    `yaml:"names"`
	Rules   []NameValue `yaml:"rules"`
	Replace []Replace   `yaml:"replace"`
}

type Replace struct {
	Old string `yaml:"old"`
	New string `yaml:"new"`
}

type NameValue struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

func (s *Sync) Do() error {
	s.From.Login()
	s.To.Login()
	for _, name := range s.Names {
		if err := s.syncTags(name); err != nil {
			return err
		}
	}
	return nil
}

func (s *Sync) replaceName(name string) string {
	for _, r := range s.Replace {
		name = strings.Replace(name, r.Old, r.New, 0)
	}
	return name
}

func (s *Sync) syncTags(name string) error {
	fromImage, err := s.From.listTags(name)
	if err != nil {
		return err
	}

	tName := s.replaceName(name)
	toImage, err := s.To.listTags(tName)
	if err != nil {
		return err
	}
	diffTags := slice.Difference(fromImage.Tags, toImage.Tags)
	for _, tag := range diffTags {
		if s.matchRules(tag) {
			if err := s.sync(name, tName, tag); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Sync) matchRules(tagName string) bool {
	for _, rule := range s.Rules {
		reg := regexp.MustCompile(rule.Value)
		if !reg.MatchString(tagName) {
			return false
		}
	}
	return true
}

func (s *Sync) sync(name, tName, tag string) error {
	if err := s.From.pullImage(name, tag); err != nil {
		return err
	}

	source := fmt.Sprintf("%s/%s:%s", s.From.Registry, name, tag)
	target := fmt.Sprintf("%s/%s:%s", s.To.Registry, tName, tag)
	if err := s.From.tagImage(source, target); err != nil {
		return err
	}
	if err := s.From.pushImage(name, tag); err != nil {
		return err
	}
	return nil
}
