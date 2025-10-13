package github_test

import (
	"fmt"
	"iter"
	"maps"
	"net/http"
	"strings"

	"github.com/AntonKosov/git-backups/internal/github"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Reader tests", func() {
	var (
		allRepos iter.Seq2[github.Repo, error]
	)

	BeforeEach(func() {
		responder := httpmock.NewStringResponder(http.StatusOK, generateResponseJSON(1, 3))
		httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(1), responder)
		responder = httpmock.NewStringResponder(http.StatusOK, generateResponseJSON(1, 0))
		httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(2), responder)
	})

	JustBeforeEach(func() {
		allRepos = github.Reader{}.AllRepos(ctx, "GH_XXX", "owner")
	})

	It("correctly reads one page of repositories", func() {
		repos := maps.Collect(allRepos)
		Expect(repos).To(Equal(map[github.Repo]error{
			{
				Name:     "Repo1Name",
				Owner:    "User",
				CloneURL: "https://Repo1Url.com",
				SSHURL:   "git:github.com/repo-owner1/hello-world.git",
			}: nil,
			{
				Name:     "Repo2Name",
				Owner:    "User",
				CloneURL: "https://Repo2Url.com",
				SSHURL:   "git:github.com/repo-owner2/hello-world.git",
			}: nil,
			{
				Name:     "Repo3Name",
				Owner:    "User",
				CloneURL: "https://Repo3Url.com",
				SSHURL:   "git:github.com/repo-owner3/hello-world.git",
			}: nil,
		}))
	})

	When("an unexpected code is returned", func() {
		BeforeEach(func() {
			responder := httpmock.NewStringResponder(http.StatusBadRequest, "")
			httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(1), responder)
		})

		It("returns an error", func() {
			repos := maps.Collect(allRepos)
			Expect(repos).To(HaveLen(1))
			Expect(repos[github.Repo{}].Error()).To(ContainSubstring("unexpected status code: 400 (400 Bad Request)"))
		})
	})

	When("an invalid response is recieved", func() {
		BeforeEach(func() {
			responder := httpmock.NewStringResponder(http.StatusOK, "Invalid json file")
			httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(1), responder)
		})

		It("returns an error", func() {
			repos := maps.Collect(allRepos)
			Expect(repos).To(HaveLen(1))
			Expect(repos[github.Repo{}]).To(MatchError("invalid character 'I' looking for beginning of value"))
		})
	})

	When("there are multiple pages", func() {
		BeforeEach(func() {
			responder := httpmock.NewStringResponder(http.StatusOK, generateResponseJSON(1, 100))
			httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(1), responder)
			responder = httpmock.NewStringResponder(http.StatusOK, generateResponseJSON(101, 200))
			httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(2), responder)
			responder = httpmock.NewStringResponder(http.StatusOK, generateResponseJSON(201, 300))
			httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(3), responder)
			responder = httpmock.NewStringResponder(http.StatusOK, generateResponseJSON(1, 0))
			httpmock.RegisterResponder(http.MethodGet, getRepositoriesURL(4), responder)
		})

		It("correctly reads all pages of repositories", func() {
			i := 1
			for repo, err := range allRepos {
				Expect(err).NotTo(HaveOccurred())
				Expect(repo).To(Equal(github.Repo{
					Name:     fmt.Sprintf("Repo%vName", i),
					Owner:    "User",
					CloneURL: fmt.Sprintf("https://Repo%vUrl.com", i),
					SSHURL:   fmt.Sprintf("git:github.com/repo-owner%v/hello-world.git", i),
				}))

				i++
			}
		})
	})
})

func getRepositoriesURL(page int) string {
	return fmt.Sprintf("https://api.github.com/user/repos?affiliation=owner&per_page=100&page=%v", page)
}

func generateResponseJSON(first, last int) string {
	sb := strings.Builder{}
	sb.WriteString(`[`)
	for i := first; i <= last; i++ {
		if i > first {
			sb.WriteString(`,`)
		}

		sb.WriteString(fmt.Sprintf(`{
			"name": "Repo%[1]vName",
			"owner": {"login": "User"},
			"clone_url": "https://Repo%[1]vUrl.com",
			"ssh_url": "git:github.com/repo-owner%[1]v/hello-world.git"
		}`, i))
	}

	sb.WriteString(`]`)

	return sb.String()
}
