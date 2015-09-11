package models

import (
	"sort"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
)

// DockerImage represents a Image
type DockerImage struct {
	ID          string            `json:"id"`
	RepoTags    []string          `json:"repoTags,omitempty"`
	Created     int64             `json:"created,omitempty"`
	Size        int64             `json:"size,omitempty"`
	VirtualSize int64             `json:"virtualSize,omitempty"`
	ParentID    string            `json:"parentId,omitempty"`
	RepoDigests []string          `json:"repoDigests,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
}

// DockerImages represents list of DockerImage
type DockerImages []DockerImage

// SearchImages checks whether it contains key word or not
func SearchImages(images []docker.APIImages, words []string) DockerImages {
	results := DockerImages{}
	for _, i := range images {
		image := convertImage(i)
		if image.contains(words) {
			results = append(results, image)
		}
	}
	sort.Sort(results)
	return results
}

func convertImage(i docker.APIImages) DockerImage {
	image := DockerImage{
		ID:          i.ID,
		RepoTags:    i.RepoTags,
		Created:     i.Created,
		Size:        i.Size,
		VirtualSize: i.VirtualSize,
		ParentID:    i.ParentID,
		RepoDigests: i.RepoDigests,
		Labels:      i.Labels,
	}
	return image
}

func (i DockerImage) contains(words []string) bool {
	image := i.toUpperFields()
	match := true
	for _, word := range words {
		match = match && (strings.Contains(image.ID, word) ||
			inStringArray(image.RepoTags, word) ||
			strings.Contains(image.ParentID, word) ||
			inStringArray(image.RepoDigests, word) ||
			inMapString(image.Labels, word))
	}
	return match
}

func inMapString(m map[string]string, word string) bool {
	match := false
	for _, value := range m {
		match = match || strings.Contains(value, word)
	}
	return match
}

func (i DockerImage) toUpperFields() DockerImage {
	image := DockerImage{}
	image.ID = strings.ToUpper(i.ID)
	image.RepoTags = make([]string, len(i.RepoTags))
	for idx, repo := range i.RepoTags {
		image.RepoTags[idx] = strings.ToUpper(repo)
	}
	image.ParentID = strings.ToUpper(i.ParentID)
	image.RepoDigests = make([]string, len(i.RepoDigests))
	for idx, repo := range i.RepoDigests {
		image.RepoDigests[idx] = strings.ToUpper(repo)
	}
	image.Labels = map[string]string{}
	for key, value := range i.Labels {
		image.Labels[key] = strings.ToUpper(value)
	}
	return image
}

func (imgs DockerImages) Len() int {
	return len(imgs)
}

func (imgs DockerImages) Swap(i, j int) {
	imgs[i], imgs[j] = imgs[j], imgs[i]
}

func (imgs DockerImages) Less(i, j int) bool {
	a, b := imgs[i], imgs[j]
	if len(a.RepoTags) > 0 && len(b.RepoTags) > 0 {
		return a.RepoTags[0] < b.RepoTags[0]
	}
	return a.Created < b.Created
}
