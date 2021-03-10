package git

import (
	"context"
	"fmt"
	"github.com/dotWicho/utilities"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"path"
	"strings"
)

// Operations interface
type Operations interface {
	Commit(repoName, commitSHA string) *github.Commit
	Compare(repoName, base, head string) *github.CommitsComparison
	Merge(repoName, base, head, message string) *github.RepositoryCommit
	Repositories(repoType, repoSort string) []*github.Repository
	Repository(repoName string) *github.Repository
	Branches(repoName string) []*github.Branch
	Branch(repoName, branchName string) *github.Branch
	Tags(repoName string) []*github.RepositoryTag
	TagByName(repoName, tagName string) *github.RepositoryTag
	ReferenceByBranch(repoName, branchName string) *github.Reference
	ReferenceByHeads(repoName, branchName string) *github.Reference
	ReferenceByTag(repoName, tagName string) *github.Reference
	CreateRefs(repoName, branchName, SHARef string) *github.Reference
	Tree(repoName, sourceFiles string, reference *github.Reference) *github.Tree
	Users() []*github.User
	User(userName string) *github.User
	CreatePullRequest(repoName, srcBranch, dstBranch, subject, description string) *github.PullRequest
	AssignReviewers(id int, repoName string, reviewers []string) *github.PullRequest
	Download(repoName, refName, filePath string) (body io.ReadCloser, err error)
	optsPullRequest(subject, srcBranch, dstBranch, description string) *github.NewPullRequest
}

// Client encapsulate in a more simply implementation the Google's go-github
type Client struct {
	Organization string
	AllPages     bool
	token        string
	github       *github.Client
	ctx          context.Context
	tkSource     oauth2.TokenSource
	tClient      *http.Client
}

// New creates a github Client with a provided token
func New(token string) *Client {

	client := &Client{token: token, ctx: context.Background()}
	client.tkSource = oauth2.StaticTokenSource(&oauth2.Token{AccessToken: client.token})
	client.tClient = oauth2.NewClient(client.ctx, client.tkSource)
	client.github = github.NewClient(client.tClient)
	client.AllPages = false

	return client
}

// Commit returns an Object Commit based on repoName and commitSHA
func (c *Client) Commit(repoName, commitSHA string) *github.Commit {

	if commit, _, err := c.github.Git.GetCommit(c.ctx, c.Organization, repoName, commitSHA); err == nil {
		return commit
	}
	return nil
}

// Compare returns an Object Commit based on repoName and commitSHA
func (c *Client) Compare(repoName, base, head string) *github.CommitsComparison {

	if commit, _, err := c.github.Repositories.CompareCommits(c.ctx, c.Organization, repoName, base, head); err == nil {
		return commit
	}
	return nil
}

// Merge returns an Object Commit based on merge to repoName:head into repoName:base
func (c *Client) Merge(repoName, base, head, message string) *github.RepositoryCommit {

	request := &github.RepositoryMergeRequest{Base: &base, Head: &head, CommitMessage: &message}
	if commit, _, err := c.github.Repositories.Merge(c.ctx, c.Organization, repoName, request); err == nil {
		return commit
	}
	return nil
}

// Repositories list all Organization repositories
func (c *Client) Repositories(repoType, repoSort string) []*github.Repository {

	//
	opts := &github.RepositoryListByOrgOptions{Type: repoType, Sort: repoSort, ListOptions: github.ListOptions{PerPage: 128, Page: 0}}

	var repos []*github.Repository
	for {
		repo, response, err := c.github.Repositories.ListByOrg(c.ctx, c.Organization, opts)
		if err != nil {
			return nil
		}

		repos = append(repos, repo...)

		if response.NextPage == 0 || !c.AllPages {
			break
		}
		opts.Page = response.NextPage
	}
	return repos
}

// Repository return a repo selected by name
func (c *Client) Repository(repoName string) *github.Repository {

	if repo, _, err := c.github.Repositories.Get(c.ctx, c.Organization, repoName); err == nil {
		return repo
	}
	return nil
}

// Branches returns all branches for a repoName
func (c *Client) Branches(repoName string) []*github.Branch {

	//
	opts := &github.BranchListOptions{Protected: nil, ListOptions: github.ListOptions{PerPage: 4, Page: 0}}

	var branches []*github.Branch
	for {
		branch, response, err := c.github.Repositories.ListBranches(c.ctx, c.Organization, repoName, opts)
		if err != nil {
			return nil
		}

		branches = append(branches, branch...)

		if response.NextPage == 0 || !c.AllPages {
			break
		}
		opts.Page = response.NextPage
	}
	return branches
}

// Branch returns an Object branch based on repoName and branchName
func (c *Client) Branch(repoName, branchName string) *github.Branch {

	var branch *github.Branch
	var err error

	if branch, _, err = c.github.Repositories.GetBranch(c.ctx, c.Organization, repoName, branchName); err != nil {
		return nil
	}
	return branch
}

// Tags returns all tags for a repoName
func (c *Client) Tags(repoName string) []*github.RepositoryTag {
	//
	opts := &github.ListOptions{PerPage: 12, Page: 0}

	var tags []*github.RepositoryTag
	for {
		tag, response, err := c.github.Repositories.ListTags(c.ctx, c.Organization, repoName, opts)
		if err != nil {
			return nil
		}

		tags = append(tags, tag...)

		if response.NextPage == 0 || !c.AllPages {
			break
		}
		opts.Page = response.NextPage
	}
	return tags
}

// TagByName returns an Object Tag based in repoName and tagName
func (c *Client) TagByName(repoName, tagName string) *github.RepositoryTag {
	//
	opts := &github.ListOptions{PerPage: 128, Page: 0}

	var theTag *github.RepositoryTag
	for {
		tag, response, err := c.github.Repositories.ListTags(c.ctx, c.Organization, repoName, opts)
		if err != nil {
			return nil
		}

		for _, testing := range tag {
			if *testing.Name == tagName {
				theTag = testing
				break
			}
		}

		if response.NextPage == 0 || theTag != nil {
			break
		}
		opts.Page = response.NextPage
	}
	return theTag
}

// ReferenceByBranch returns an Object Reference based in repoName and branchName
func (c *Client) ReferenceByBranch(repoName, branchName string) *github.Reference {

	if ref, _, err := c.github.Git.GetRef(c.ctx, c.Organization, repoName, "refs/branch/"+branchName); err == nil {
		return ref
	}
	return nil
}

// ReferenceByHeads returns an Object Reference based in repoName and branchName from heads
func (c *Client) ReferenceByHeads(repoName, branchName string) *github.Reference {

	if ref, _, err := c.github.Git.GetRef(c.ctx, c.Organization, repoName, "refs/heads/"+branchName); err == nil {
		return ref
	}
	return nil
}

// ReferenceByTag returns an Object Reference based in repoName and tagName
func (c *Client) ReferenceByTag(repoName, tagName string) *github.Reference {

	if ref, _, err := c.github.Git.GetRef(c.ctx, c.Organization, repoName, "refs/tags/"+tagName); err == nil {
		return ref
	}
	return nil
}

// CreateRefs permits create an Object Reference based on repoName, branchName and SHAReference
func (c *Client) CreateRefs(repoName, branchName, SHARef string) *github.Reference {

	newRef := &github.Reference{Ref: github.String("refs/heads/" + branchName), Object: &github.GitObject{SHA: &SHARef}}
	if ref, _, err := c.github.Git.CreateRef(c.ctx, c.Organization, repoName, newRef); err == nil {
		return ref
	}
	return nil
}

// Tree permits create an Object Tree given a fileName list
func (c *Client) Tree(repoName, sourceFiles string, reference *github.Reference) *github.Tree {

	// Create a tree with what to commit.
	var entries []*github.TreeEntry

	// Load each file into the tree.
	for _, fileArg := range strings.Split(sourceFiles, ",") {
		content := utilities.ReadFile(fileArg)
		if content == nil {
			return nil
		}
		entries = append(entries, &github.TreeEntry{Path: github.String(fileArg), Type: github.String("blob"), Content: github.String(string(content)), Mode: github.String("100644")})
	}

	if tree, _, err := c.github.Git.CreateTree(c.ctx, c.Organization, repoName, *reference.Object.SHA, entries); err == nil {
		return tree
	}
	return nil
}

// Users returns all Users in an Organization
func (c *Client) Users() []*github.User {
	//
	opts := &github.UserListOptions{Since: 0, ListOptions: github.ListOptions{PerPage: 100, Page: 0}}

	var users []*github.User
	for {
		user, response, err := c.github.Users.ListAll(c.ctx, opts)
		if err != nil {
			return nil
		}

		users = append(users, user...)

		if response.NextPage == 0 || !c.AllPages {
			break
		}
		opts.Page = response.NextPage
	}
	return users
}

// User returns an Object User by its userName
func (c *Client) User(userName string) *github.User {

	if len(userName) == 0 {
		userName = ""
	}

	if user, _, err := c.github.Users.Get(c.ctx, userName); err == nil {
		return user
	}
	return nil
}

// CreatePullRequest permits create an PullRequest into repoName using source and destiny branches
func (c *Client) CreatePullRequest(repoName, srcBranch, dstBranch, subject, description string) *github.PullRequest {

	if len(repoName) == 0 || len(srcBranch) == 0 || len(subject) == 0 {
		return nil
	}

	if strings.Contains(srcBranch, ":") && len(c.Organization) == 0 {
		dstBranch = fmt.Sprintf("%s:%s", c.Organization, dstBranch)
	}

	newPR := c.optsPullRequest(subject, srcBranch, dstBranch, description)
	if pr, _, err := c.github.PullRequests.Create(c.ctx, c.Organization, repoName, newPR); err == nil {
		return pr
	} else {
		panic(err)
	}
	return nil
}

// AssignReviewers permits assign Reviewers to an one PullRequest
func (c *Client) AssignReviewers(id int, repoName string, reviewers []string) *github.PullRequest {

	if len(reviewers) == 0 {
		return nil
	}

	rr := github.ReviewersRequest{Reviewers: reviewers, TeamReviewers: nil}

	if pr, _, err := c.github.PullRequests.RequestReviewers(c.ctx, c.Organization, repoName, id, rr); err == nil {
		return pr
	}
	return nil
}

// Download returns body response of GET DownloadURL corresponding to filePath
func (c *Client) Download(repoName, refName, filePath string) (body io.ReadCloser, err error) {

	if len(repoName) == 0 {
		return nil, fmt.Errorf("repo cannot be null nor empty")
	}
	if len(filePath) == 0 {
		return nil, fmt.Errorf("filePath cannot be null nor empty")
	}

	dirPath := path.Dir(filePath)
	fileName := path.Base(filePath)
	opts := &github.RepositoryContentGetOptions{Ref: refName}
	dirContents := make([]*github.RepositoryContent, 0)

	_, dirContents, _, err = c.github.Repositories.GetContents(c.ctx, c.Organization, repoName, dirPath, opts)
	if err != nil {
		return nil, err
	}

	for _, contents := range dirContents {
		if contents.Type != nil && *contents.Type == "file" {
			if contents.Name != nil && *contents.Name == fileName {

				var response *http.Response
				if response, err = http.Get(*contents.DownloadURL); err != nil {
					return nil, err
				}
				return response.Body, nil
			}
		}
	}
	return nil, fmt.Errorf("filename %s not found", fileName)
}

// optsPullRequest populate a NewPullRequests with its info
func (c *Client) optsPullRequest(subject, srcBranch, dstBranch, description string) *github.NewPullRequest {

	return &github.NewPullRequest{
		Title:               &subject,
		Head:                &srcBranch,
		Base:                &dstBranch,
		Body:                &description,
		MaintainerCanModify: github.Bool(true),
	}
}
