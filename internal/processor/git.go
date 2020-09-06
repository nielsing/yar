package processor

import (
	"strings"

	"github.com/nielsing/yar/internal/robber"
	"github.com/nielsing/yar/internal/utils"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// DiffObject holds everything that is needed to analyze a diff.
type DiffObject struct {
	Commit   *object.Commit
	Diff     *string
	Filepath *string
}

// NewDiffObject returns a new DiffObject.
func NewDiffObject(commit *object.Commit, diff, filepath *string) *DiffObject {
	return &DiffObject{
		Commit:   commit,
		Diff:     diff,
		Filepath: filepath,
	}
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

// GetDiffObjects goes through all commits in a given repository and searches all diffs within each
// commit (both Add and Remove diffs) and finally returns a list of DiffObjects containing all diffs
// for the repository.
func GetDiffObjects(r *robber.Robber, repo *git.Repository, reponame string) ([]*DiffObject, error) {
	defer func() {
		if err := recover(); err != nil {
			r.Logger.LogFail("The cache folder %s is corrupted please run yar clear %s and try again\n",
				reponame[9:], reponame[9:])
		}
	}()

	var diffObjects []*DiffObject
	ref, err := repo.Head()
	commitIter, err := repo.Log(&git.LogOptions{From: ref.Hash(), Order: git.LogOrderCommitterTime})
	if err != nil {
		return nil, err
	}

	count := 0
	commit, err := commitIter.Next()
	for commit != nil {
		if count == r.Args.Git.Depth {
			break
		}
		changes, err := GetCommitChanges(commit)
		if err != nil {
			r.Logger.LogVerbose("Unable to get commit changes for commit %s in repo %s. Skipping...",
				commit.Hash.String(), reponame)
			continue
		}
		diffObjects = append(diffObjects, getDiff(r, changes, commit, reponame)...)
		commit, _ = commitIter.Next()
	}
	return diffObjects, nil
}

// getDiff is a helper function for the GetDiffObjects function here above.
func getDiff(r *robber.Robber, changes object.Changes, commit *object.Commit, reponame string) []*DiffObject {
	// This is done to handle the following inevitable error https://github.com/sergi/go-diff/issues/89
	// If you run into this error a bunch of times then please take a look at the issue and see if you can
	// contribute a fix :).
	defer func() {
		if err := recover(); err != nil {
			r.Logger.LogWarn("Encountered a file that is too large to handle in %s! Skipping...\n", reponame)
		}
	}()

	var diffObjects []*DiffObject
	for _, change := range changes {
		patch, err := change.Patch()
		if err != nil {
			r.Logger.LogVerbose("Unable to get patches in a change in commit %s in repo %s. Skipping...",
				commit.Hash.String(), reponame)
			continue
		}

		for _, file := range patch.FilePatches() {
			if file.IsBinary() {
				continue
			}
			filename := getFilepath(file)
			if utils.BlacklistedFile(r, filename) {
				continue
			}
			for _, chunk := range file.Chunks() {
				if chunk.Type() != diff.Equal {
					continue
				}
				content := strings.Trim(chunk.Content(), " \n")
				diffObjects = append(diffObjects, NewDiffObject(commit, &content, &filename))
			}
		}
	}
	return diffObjects
}

// GetDiffs helper
func getFilepath(file diff.FilePatch) string {
	from, to := file.Files()
	if from != nil {
		return from.Path()
	}
	return to.Path()
}
