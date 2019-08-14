package robber

import (
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/format/diff"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// DiffObject holds everything that is needed to analyze a diff.
type DiffObject struct {
	Commit   *object.Commit
	Diff     *string
	Reponame *string
	Filepath *string
}

// NewDiffObject returns a new DiffObject.
func NewDiffObject(commit *object.Commit, diff, reponame, filepath *string) *DiffObject {
	return &DiffObject{
		Commit:   commit,
		Diff:     diff,
		Reponame: reponame,
		Filepath: filepath,
	}
}

// cloneRepo creates a temp directory in the OS's temp directory
// and clones the given URL into it.
func cloneRepo(url string, depth int) (*git.Repository, error) {
	dir, err := ioutil.TempDir("", "yar")
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:   url,
		Depth: depth,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// OpenRepo opens a repository found at the given path.
// If the path points to a nonexistant repository it assumes that an URL
// was given and tries to clone it instead.
func OpenRepo(path string, depth int) (*git.Repository, error) {
	var repo *git.Repository
	if _, err := os.Stat(path); err == nil {
		repo, err = git.PlainOpen(path)
		if err != nil {
			return nil, err
		}
		return repo, nil
	}

	repo, err := cloneRepo(path, depth)
	if err != nil {
		return nil, err
	}
	return repo, nil
}

// GetCommits simply traverses a given repository, gathering all commits
// and then returns a list of them.
func GetCommits(repo *git.Repository) ([]*object.Commit, error) {
	var commits []*object.Commit
	ref, err := repo.Head()
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}

	commitIter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})
	return commits, nil
}

func getParentTree(commit *object.Commit) (*object.Tree, error) {
	// Bit of a hack to handle the edge case of 0 parents.
	var emptyTree *object.Tree
	if commit.NumParents() == 0 {
		emptyTree = &object.Tree{Entries: []object.TreeEntry{}}
		return emptyTree, nil
	}

	parent, err := commit.Parents().Next()
	if err != nil {
		return nil, err
	}
	parentTree, err := parent.Tree()
	if err != nil {
		return nil, err
	}
	return parentTree, nil
}

// GetCommitChanges gets the changes of a commit by comparing it to its'
// parent commit tree.
func GetCommitChanges(commit *object.Commit) (object.Changes, error) {
	commitTree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	parentTree, err := getParentTree(commit)
	if err != nil {
		return nil, err
	}

	changes, err := object.DiffTree(commitTree, parentTree)
	if err != nil {
		return nil, err
	}
	return changes, nil
}

func getFilepath(file diff.FilePatch) string {
	_, to := file.Files()
	if to == nil {
		return ""
	}
	return to.Path()
}

// GetDiffs gets all diffs which are either of type addage or removal
// for a change in a commit.
func GetDiffs(change *object.Change) ([]string, string, error) {
	patch, err := change.Patch()
	if err != nil {
		return nil, "", err
	}

	var diffs []string
	var filepath string
	for _, file := range patch.FilePatches() {
		if file.IsBinary() {
			continue
		}
		filepath = getFilepath(file)
		for _, chunk := range file.Chunks() {
			if chunk.Type() != 1 {
				continue
			}
			diff := strings.Trim(chunk.Content(), " \n")
			diffs = append(diffs, diff)
		}
	}
	return diffs, filepath, nil
}
