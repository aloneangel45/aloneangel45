package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	username      = "aloneangel45"
	graphqlAPIURL = "https://api.github.com/graphql"
)

type themePalette struct {
	name          string
	bg            string
	surface       string
	surfaceAlt    string
	border        string
	hairline      string
	gridStroke    string
	gridOpacity   string
	inkPrimary    string
	inkSecondary  string
	inkMuted      string
	accentBlue    string
	accentViolet  string
	liveGreen     string
	trackBg       string
}

var darkPalette = themePalette{
	name:         "dark",
	bg:           "#0d1117",
	surface:      "#161b22",
	surfaceAlt:   "#0d1117",
	border:       "#30363d",
	hairline:     "#21262d",
	gridStroke:   "#30363d",
	gridOpacity:  "0.18",
	inkPrimary:   "#f0f6fc",
	inkSecondary: "#8b949e",
	inkMuted:     "#6e7681",
	accentBlue:   "#58a6ff",
	accentViolet: "#a371f7",
	liveGreen:    "#3fb950",
	trackBg:      "#21262d",
}

var lightPalette = themePalette{
	name:         "light",
	bg:           "#ffffff",
	surface:      "#f6f8fa",
	surfaceAlt:   "#ffffff",
	border:       "#d0d7de",
	hairline:     "#d8dee4",
	gridStroke:   "#d0d7de",
	gridOpacity:  "0.32",
	inkPrimary:   "#1f2328",
	inkSecondary: "#59636e",
	inkMuted:     "#6e7781",
	accentBlue:   "#0969da",
	accentViolet: "#8250df",
	liveGreen:    "#1a7f37",
	trackBg:      "#eaeef2",
}

type userQueryResp struct {
	Data struct {
		User struct {
			ContributionsCollection struct {
				TotalCommitContributions     int `json:"totalCommitContributions"`
				TotalPullRequestContributions int `json:"totalPullRequestContributions"`
				TotalIssueContributions      int `json:"totalIssueContributions"`
				TotalRepositoryContributions int `json:"totalRepositoryContributions"`
				ContributionCalendar         struct {
					TotalContributions int `json:"totalContributions"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
			Repositories struct {
				Nodes []struct {
					Name     string `json:"name"`
					Releases struct {
						TotalCount int `json:"totalCount"`
					} `json:"releases"`
					LatestRelease *struct {
						TagName string `json:"tagName"`
					} `json:"latestRelease"`
					Languages struct {
						Edges []struct {
							Size int `json:"size"`
							Node struct {
								Name string `json:"name"`
							} `json:"node"`
						} `json:"edges"`
					} `json:"languages"`
				} `json:"nodes"`
			} `json:"repositories"`
		} `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type prSearchResp struct {
	Data struct {
		Search struct {
			IssueCount int `json:"issueCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
			Nodes []struct {
				Repository *struct {
					NameWithOwner   string `json:"nameWithOwner"`
					StargazerCount  int    `json:"stargazerCount"`
					PrimaryLanguage *struct {
						Name string `json:"name"`
					} `json:"primaryLanguage"`
				} `json:"repository"`
			} `json:"nodes"`
		} `json:"search"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

type repoReleaseInfo struct {
	LatestTag    string
	ReleaseCount int
}

type upstreamRepoStat struct {
	NameWithOwner string
	ShortName     string
	Stars         int
	Language      string
	PRCount       int
}

type langStat struct {
	Name  string
	Short string
	Color string
	Count int
}

type profileTelemetry struct {
	TotalContributions int
	Commits            int
	PRs                int
	FlagshipCount      int
	TotalReleases      int
	FlagshipSourceMB   float64
	GoPct              int
	RepoReleases       map[string]repoReleaseInfo

	MergedPRCount     int
	DistinctRepoCount int
	TopByStars        []upstreamRepoStat
	TopByPRs          []upstreamRepoStat
	LangStats         []langStat
	SyncedDateUTC     string
}

func doGraphQL(client *http.Client, token string, query string, vars map[string]any, out any) error {
	payload, err := json.Marshal(map[string]any{
		"query":     query,
		"variables": vars,
	})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, graphqlAPIURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "aloneangel45-profile-sync")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("graphql HTTP %d: %s", resp.StatusCode, string(body))
	}
	return json.Unmarshal(body, out)
}

func fetchTelemetry(token string) (*profileTelemetry, error) {
	client := &http.Client{Timeout: 25 * time.Second}

	userQuery := `
query($login: String!) {
  user(login: $login) {
    contributionsCollection {
      totalCommitContributions
      totalPullRequestContributions
      totalIssueContributions
      totalRepositoryContributions
      contributionCalendar { totalContributions }
    }
    repositories(first: 50, ownerAffiliations: OWNER, isFork: false, privacy: PUBLIC) {
      nodes {
        name
        releases { totalCount }
        latestRelease { tagName }
        languages(first: 10, orderBy: {field: SIZE, direction: DESC}) {
          edges { size node { name } }
        }
      }
    }
  }
}`
	var uResp userQueryResp
	if err := doGraphQL(client, token, userQuery, map[string]any{"login": username}, &uResp); err != nil {
		return nil, fmt.Errorf("fetch user telemetry: %w", err)
	}
	if len(uResp.Errors) > 0 {
		return nil, fmt.Errorf("user graphql error: %s", uResp.Errors[0].Message)
	}

	cc := uResp.Data.User.ContributionsCollection
	totalContribs := cc.ContributionCalendar.TotalContributions
	if totalContribs == 0 {
		totalContribs = cc.TotalCommitContributions + cc.TotalPullRequestContributions + cc.TotalIssueContributions + cc.TotalRepositoryContributions
	}

	repoReleases := make(map[string]repoReleaseInfo)
	flagshipCount := 0
	totalReleases := 0
	langBytes := make(map[string]int)

	for _, r := range uResp.Data.User.Repositories.Nodes {
		if strings.EqualFold(r.Name, username) {
			continue
		}
		flagshipCount++
		totalReleases += r.Releases.TotalCount
		tag := ""
		if r.LatestRelease != nil {
			tag = r.LatestRelease.TagName
		}
		repoReleases[r.Name] = repoReleaseInfo{
			LatestTag:    tag,
			ReleaseCount: r.Releases.TotalCount,
		}
		for _, edge := range r.Languages.Edges {
			if edge.Node.Name == "Batchfile" {
				continue
			}
			langBytes[edge.Node.Name] += edge.Size
		}
	}

	totalBytes := 0
	for _, b := range langBytes {
		totalBytes += b
	}
	goPct := 0
	if totalBytes > 0 {
		goPct = int((float64(langBytes["Go"])/float64(totalBytes))*100.0 + 0.5)
	}

	prQuery := `
query($q: String!, $cursor: String) {
  search(query: $q, type: ISSUE, first: 100, after: $cursor) {
    issueCount
    pageInfo {
      hasNextPage
      endCursor
    }
    nodes {
      ... on PullRequest {
        repository {
          nameWithOwner
          stargazerCount
          primaryLanguage { name }
        }
      }
    }
  }
}`
	searchFilter := fmt.Sprintf("author:%s type:pr is:merged -user:%s", username, username)
	repoMap := make(map[string]*upstreamRepoStat)
	mergedTotal := 0
	var cursor any = nil

	for {
		var pResp prSearchResp
		if err := doGraphQL(client, token, prQuery, map[string]any{"q": searchFilter, "cursor": cursor}, &pResp); err != nil {
			return nil, fmt.Errorf("fetch merged PRs: %w", err)
		}
		if len(pResp.Errors) > 0 {
			return nil, fmt.Errorf("pr search graphql error: %s", pResp.Errors[0].Message)
		}
		if pResp.Data.Search.IssueCount > mergedTotal {
			mergedTotal = pResp.Data.Search.IssueCount
		}
		for _, n := range pResp.Data.Search.Nodes {
			if n.Repository == nil {
				continue
			}
			full := n.Repository.NameWithOwner
			lang := "Other"
			if n.Repository.PrimaryLanguage != nil && n.Repository.PrimaryLanguage.Name != "" {
				lang = n.Repository.PrimaryLanguage.Name
			}
			entry, ok := repoMap[full]
			if !ok {
				short := full
				if idx := strings.IndexByte(full, '/'); idx >= 0 {
					short = full[idx+1:]
				}
				entry = &upstreamRepoStat{
					NameWithOwner: full,
					ShortName:     short,
					Stars:         n.Repository.StargazerCount,
					Language:      lang,
				}
				repoMap[full] = entry
			}
			entry.PRCount++
		}
		if !pResp.Data.Search.PageInfo.HasNextPage || pResp.Data.Search.PageInfo.EndCursor == "" {
			break
		}
		cursor = pResp.Data.Search.PageInfo.EndCursor
	}

	allUpstream := make([]upstreamRepoStat, 0, len(repoMap))
	langPRCounts := make(map[string]int)
	for _, st := range repoMap {
		allUpstream = append(allUpstream, *st)
		langPRCounts[st.Language] += st.PRCount
	}

	byStars := append([]upstreamRepoStat(nil), allUpstream...)
	sort.Slice(byStars, func(i, j int) bool {
		if byStars[i].Stars != byStars[j].Stars {
			return byStars[i].Stars > byStars[j].Stars
		}
		if byStars[i].PRCount != byStars[j].PRCount {
			return byStars[i].PRCount > byStars[j].PRCount
		}
		return byStars[i].NameWithOwner < byStars[j].NameWithOwner
	})
	if len(byStars) > 8 {
		byStars = byStars[:8]
	}

	byPRs := append([]upstreamRepoStat(nil), allUpstream...)
	sort.Slice(byPRs, func(i, j int) bool {
		if byPRs[i].PRCount != byPRs[j].PRCount {
			return byPRs[i].PRCount > byPRs[j].PRCount
		}
		if byPRs[i].Stars != byPRs[j].Stars {
			return byPRs[i].Stars > byPRs[j].Stars
		}
		return byPRs[i].NameWithOwner < byPRs[j].NameWithOwner
	})
	if len(byPRs) > 4 {
		byPRs = byPRs[:4]
	}

	langOrder := []struct {
		name  string
		short string
		color string
	}{
		{"Go", "Go", "#00ADD8"},
		{"TypeScript", "TS", "#3178c6"},
		{"Rust", "Rust", "#dea584"},
		{"Python", "Python", "#3572A5"},
	}
	var langStats []langStat
	for _, spec := range langOrder {
		c := langPRCounts[spec.name]
		if c > 0 {
			langStats = append(langStats, langStat{
				Name:  spec.name,
				Short: spec.short,
				Color: spec.color,
				Count: c,
			})
		}
	}
	otherCount := 0
	for k, c := range langPRCounts {
		isKnown := false
		for _, spec := range langOrder {
			if k == spec.name {
				isKnown = true
				break
			}
		}
		if !isKnown {
			otherCount += c
		}
	}
	if otherCount > 0 && len(langStats) < 4 {
		langStats = append(langStats, langStat{
			Name:  "Other",
			Short: "Other",
			Color: "#8b949e",
			Count: otherCount,
		})
	}
	sort.Slice(langStats, func(i, j int) bool {
		if langStats[i].Count != langStats[j].Count {
			return langStats[i].Count > langStats[j].Count
		}
		return langStats[i].Name < langStats[j].Name
	})

	return &profileTelemetry{
		TotalContributions: totalContribs,
		Commits:            cc.TotalCommitContributions,
		PRs:                cc.TotalPullRequestContributions,
		FlagshipCount:      flagshipCount,
		TotalReleases:      totalReleases,
		FlagshipSourceMB:   float64(totalBytes) / (1024.0 * 1024.0),
		GoPct:              goPct,
		RepoReleases:       repoReleases,
		MergedPRCount:      mergedTotal,
		DistinctRepoCount:  len(repoMap),
		TopByStars:         byStars,
		TopByPRs:           byPRs,
		LangStats:          langStats,
		SyncedDateUTC:      time.Now().UTC().Format("2006-01-02"),
	}, nil
}

func langColor(lang string) string {
	switch lang {
	case "Go":
		return "#00ADD8"
	case "Rust":
		return "#dea584"
	case "TypeScript":
		return "#3178c6"
	case "Python":
		return "#3572A5"
	default:
		return "#8b949e"
	}
}

func formatStars(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000.0)
	}
	return fmt.Sprintf("%d", n)
}

func truncateASCII(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}

func renderHeaderSVG(t *profileTelemetry, p themePalette) string {
	var b strings.Builder
	sans := "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif"
	mono := "ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace"

	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 900 176" width="900" height="176" role="img" aria-label="Arda (@aloneangel45) — Systems, Binary Forensics and Agent Runtime">`+"\n")
	fmt.Fprintf(&b, `  <defs>`+"\n")
	fmt.Fprintf(&b, `    <pattern id="hdrGrid-%s" width="24" height="24" patternUnits="userSpaceOnUse">`+"\n", p.name)
	fmt.Fprintf(&b, `      <path d="M 24 0 L 0 0 0 24" fill="none" stroke="%s" stroke-width="0.6" stroke-opacity="%s"/>`+"\n", p.gridStroke, p.gridOpacity)
	fmt.Fprintf(&b, `    </pattern>`+"\n")
	fmt.Fprintf(&b, `    <style>`+"\n")
	fmt.Fprintf(&b, `      @keyframes livePulse { 0%%, 100%% { opacity: 1; } 50%% { opacity: 0.35; } }`+"\n")
	fmt.Fprintf(&b, `      .live-dot { animation: livePulse 3.6s ease-in-out infinite; }`+"\n")
	fmt.Fprintf(&b, `    </style>`+"\n")
	fmt.Fprintf(&b, `  </defs>`+"\n")

	// Outer frame & subtle architectural grid
	fmt.Fprintf(&b, `  <rect x="0.5" y="0.5" width="899" height="175" rx="6" fill="%s" stroke="%s" stroke-width="1"/>`+"\n", p.bg, p.border)
	fmt.Fprintf(&b, `  <rect x="1" y="1" width="898" height="174" rx="5.5" fill="url(#hdrGrid-%s)"/>`+"\n", p.name)
	fmt.Fprintf(&b, `  <rect x="0.5" y="20" width="3" height="136" fill="%s"/>`+"\n", p.accentBlue)

	// Left identity & domain pane
	fmt.Fprintf(&b, `  <g font-family="%s">`+"\n", sans)
	fmt.Fprintf(&b, `    <text x="28" y="38" fill="%s" font-family="%s" font-size="10.5" font-weight="600" letter-spacing="1.4">ARDA  &#xB7;  @ALONEANGEL45  &#xB7;  TURKEY</text>`+"\n", p.inkMuted, mono)
	fmt.Fprintf(&b, `    <text x="28" y="74" fill="%s" font-size="22" font-weight="700" letter-spacing="-0.4">Systems, Binary Forensics &amp; Agent Runtime</text>`+"\n", p.inkPrimary)
	fmt.Fprintf(&b, `    <text x="28" y="100" fill="%s" font-size="12">Pure-Go (CGO_ENABLED=0) sandboxes, stripped ELF/PE/Mach-O analyzers, and MCP proxies.</text>`+"\n", p.inkSecondary)

	// Spec tags row
	tags := []struct {
		x     int
		w     int
		label string
	}{
		{28, 74, "Go / Rust"},
		{108, 124, "ELF / PE / Mach-O"},
		{238, 188, "Landlock / Seatbelt / Win32"},
		{432, 108, "MCP / JSON-RPC"},
	}
	for _, tg := range tags {
		cx := float64(tg.x) + float64(tg.w)/2.0
		fmt.Fprintf(&b, `    <rect x="%d" y="128" width="%d" height="24" rx="3" fill="%s" stroke="%s" stroke-width="1"/>`+"\n", tg.x, tg.w, p.surface, p.border)
		fmt.Fprintf(&b, `    <text x="%.1f" y="143.5" text-anchor="middle" fill="%s" font-family="%s" font-size="10">%s</text>`+"\n", cx, p.inkSecondary, mono, html.EscapeString(tg.label))
	}
	fmt.Fprintf(&b, `  </g>`+"\n")

	// Right 2x2 live spec matrix
	fmt.Fprintf(&b, `  <g>`+"\n")
	fmt.Fprintf(&b, `    <rect x="560" y="20" width="316" height="136" rx="6" fill="%s" stroke="%s" stroke-width="1"/>`+"\n", p.surface, p.border)
	fmt.Fprintf(&b, `    <line x1="718" y1="20" x2="718" y2="156" stroke="%s" stroke-width="1"/>`+"\n", p.hairline)
	fmt.Fprintf(&b, `    <line x1="560" y1="88" x2="876" y2="88" stroke="%s" stroke-width="1"/>`+"\n", p.hairline)
	fmt.Fprintf(&b, `    <circle class="live-dot" cx="864" cy="31" r="2.5" fill="%s"/>`+"\n", p.liveGreen)

	// Q1: Upstream PRs
	fmt.Fprintf(&b, `    <text x="576" y="40" fill="%s" font-family="%s" font-size="9.5" font-weight="600" letter-spacing="0.8">UPSTREAM PRS</text>`+"\n", p.inkMuted, mono)
	fmt.Fprintf(&b, `    <text x="576" y="64" fill="%s" font-family="%s" font-size="20" font-weight="700">%d</text>`+"\n", p.accentViolet, mono, t.MergedPRCount)
	fmt.Fprintf(&b, `    <text x="576" y="79" fill="%s" font-family="%s" font-size="10">merged &#xB7; %d repos</text>`+"\n", p.inkSecondary, mono, t.DistinctRepoCount)

	// Q2: Releases
	fmt.Fprintf(&b, `    <text x="734" y="40" fill="%s" font-family="%s" font-size="9.5" font-weight="600" letter-spacing="0.8">RELEASES</text>`+"\n", p.inkMuted, mono)
	fmt.Fprintf(&b, `    <text x="734" y="64" fill="%s" font-family="%s" font-size="20" font-weight="700">%d</text>`+"\n", p.accentBlue, mono, t.TotalReleases)
	fmt.Fprintf(&b, `    <text x="734" y="79" fill="%s" font-family="%s" font-size="10">%d Go binaries</text>`+"\n", p.inkSecondary, mono, t.FlagshipCount)

	// Q3: 1Y Activity
	fmt.Fprintf(&b, `    <text x="576" y="108" fill="%s" font-family="%s" font-size="9.5" font-weight="600" letter-spacing="0.8">1Y ACTIVITY</text>`+"\n", p.inkMuted, mono)
	fmt.Fprintf(&b, `    <text x="576" y="132" fill="%s" font-family="%s" font-size="20" font-weight="700">%d</text>`+"\n", p.inkPrimary, mono, t.TotalContributions)
	fmt.Fprintf(&b, `    <text x="576" y="147" fill="%s" font-family="%s" font-size="10">%d commits &#xB7; %d PRs</text>`+"\n", p.inkSecondary, mono, t.Commits, t.PRs)

	// Q4: Flagship Code
	fmt.Fprintf(&b, `    <text x="734" y="108" fill="%s" font-family="%s" font-size="9.5" font-weight="600" letter-spacing="0.8">FLAGSHIP CODE</text>`+"\n", p.inkMuted, mono)
	fmt.Fprintf(&b, `    <text x="734" y="132" fill="%s" font-family="%s" font-size="20" font-weight="700">%.2f MB</text>`+"\n", p.inkPrimary, mono, t.FlagshipSourceMB)
	fmt.Fprintf(&b, `    <text x="734" y="147" fill="%s" font-family="%s" font-size="10">%d%% Go &#xB7; zero-dep</text>`+"\n", p.inkSecondary, mono, t.GoPct)

	fmt.Fprintf(&b, `  </g>`+"\n")
	fmt.Fprintf(&b, `</svg>`+"\n")
	return b.String()
}

func renderUpstreamSVG(t *profileTelemetry, p themePalette) string {
	var b strings.Builder
	sans := "-apple-system, BlinkMacSystemFont, 'Segoe UI', Helvetica, Arial, sans-serif"
	mono := "ui-monospace, SFMono-Regular, 'SF Mono', Menlo, Consolas, monospace"

	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 900 292" width="900" height="292" role="img" aria-label="Upstream Open Source Contributions">`+"\n")
	fmt.Fprintf(&b, `  <defs>`+"\n")
	fmt.Fprintf(&b, `    <clipPath id="langBarClip-%s">`+"\n", p.name)
	fmt.Fprintf(&b, `      <rect x="584" y="234" width="292" height="7" rx="3.5"/>`+"\n")
	fmt.Fprintf(&b, `    </clipPath>`+"\n")
	fmt.Fprintf(&b, `  </defs>`+"\n")

	// Outer panel
	fmt.Fprintf(&b, `  <rect x="0.5" y="0.5" width="899" height="291" rx="6" fill="%s" stroke="%s" stroke-width="1"/>`+"\n", p.bg, p.border)
	fmt.Fprintf(&b, `  <path d="M 0.5 6.5 Q 0.5 0.5 6.5 0.5 L 893.5 0.5 Q 899.5 0.5 899.5 6.5 L 899.5 48 L 0.5 48 Z" fill="%s"/>`+"\n", p.surface)
	fmt.Fprintf(&b, `  <line x1="0.5" y1="48" x2="899.5" y2="48" stroke="%s" stroke-width="1"/>`+"\n", p.hairline)

	// Git-merge icon (GitHub Octicon 16px scaled 0.875 at x=24, y=17)
	fmt.Fprintf(&b, `  <g transform="translate(24, 17) scale(0.875)" fill="%s">`+"\n", p.accentViolet)
	fmt.Fprintf(&b, `    <path d="M5.45 5.154A4.25 4.25 0 0 0 9.25 7.5h1.378a2.251 2.251 0 1 1 0 1.5H9.25A5.734 5.734 0 0 1 5 7.123v3.505a2.25 2.25 0 1 1-1.5 0V5.372a2.25 2.25 0 1 1 1.95-.218ZM4.25 13.5a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5Zm8.5-4.5a.75.75 0 1 0 0-1.5.75.75 0 0 0 0 1.5ZM5 3.25a.75.75 0 1 0 0 .005V3.25Z"/>`+"\n")
	fmt.Fprintf(&b, `  </g>`+"\n")

	// Header text (single <text> with <tspan> avoids font-width drift)
	fmt.Fprintf(&b, `  <text x="45" y="29" font-family="%s"><tspan fill="%s" font-size="13.5" font-weight="600">Upstream Open Source</tspan><tspan fill="%s" font-family="%s" font-size="11.5">  &#xB7;  %d merged PRs across %d repositories</tspan></text>`+"\n", sans, p.inkPrimary, p.inkSecondary, mono, t.MergedPRCount, t.DistinctRepoCount)
	fmt.Fprintf(&b, `  <text x="876" y="29" text-anchor="end" fill="%s" font-family="%s" font-size="10.5">synced %s UTC</text>`+"\n", p.inkMuted, mono, t.SyncedDateUTC)

	// Left column: Top 8 repositories by star count
	fmt.Fprintf(&b, `  <text x="24" y="72" fill="%s" font-family="%s" font-size="9.5" font-weight="600" letter-spacing="0.9">HIGHEST-STARRED TARGETS</text>`+"\n", p.inkMuted, mono)

	for i, repo := range t.TopByStars {
		y := 96 + i*24
		dotColor := langColor(repo.Language)
		name := truncateASCII(repo.NameWithOwner, 34)
		stars := formatStars(repo.Stars)
		prLabel := fmt.Sprintf("%d PR", repo.PRCount)
		if repo.PRCount > 1 {
			prLabel = fmt.Sprintf("%d PRs", repo.PRCount)
		}

		fmt.Fprintf(&b, `  <circle cx="29" cy="%d" r="3.5" fill="%s"/>`+"\n", y-4, dotColor)
		fmt.Fprintf(&b, `  <text x="42" y="%d" fill="%s" font-family="%s" font-size="11.5">%s</text>`+"\n", y, p.inkPrimary, mono, html.EscapeString(name))
		// Crisp vector star icon + count
		fmt.Fprintf(&b, `  <path transform="translate(422, %d) scale(0.625)" fill="%s" d="M8 .25a.75.75 0 0 1 .673.418l1.882 3.815 4.21.612a.75.75 0 0 1 .416 1.279l-3.046 2.97.719 4.192a.751.751 0 0 1-1.088.791L8 12.347l-3.766 1.98a.75.75 0 0 1-1.088-.79l.72-4.194L.818 6.374a.75.75 0 0 1 .416-1.28l4.21-.611L7.327.668A.75.75 0 0 1 8 .25Z"/>`+"\n", y-9, p.inkMuted)
		fmt.Fprintf(&b, `  <text x="480" y="%d" text-anchor="end" fill="%s" font-family="%s" font-size="11">%s</text>`+"\n", y, p.inkSecondary, mono, stars)
		fmt.Fprintf(&b, `  <text x="542" y="%d" text-anchor="end" fill="%s" font-family="%s" font-size="11" font-weight="600">%s</text>`+"\n", y, p.accentViolet, mono, prLabel)
	}

	// Vertical hairline divider
	fmt.Fprintf(&b, `  <line x1="564" y1="60" x2="564" y2="274" stroke="%s" stroke-width="1"/>`+"\n", p.hairline)

	// Right column upper: Top 4 repos by PR volume
	fmt.Fprintf(&b, `  <text x="584" y="72" fill="%s" font-family="%s" font-size="9.5" font-weight="600" letter-spacing="0.9">HIGHEST PR VOLUME</text>`+"\n", p.inkMuted, mono)

	maxPRs := 1
	if len(t.TopByPRs) > 0 && t.TopByPRs[0].PRCount > 0 {
		maxPRs = t.TopByPRs[0].PRCount
	}
	for i, repo := range t.TopByPRs {
		y := 96 + i*26
		short := truncateASCII(repo.ShortName, 26)
		prLabel := fmt.Sprintf("%d PR", repo.PRCount)
		if repo.PRCount > 1 {
			prLabel = fmt.Sprintf("%d PRs", repo.PRCount)
		}
		barW := (float64(repo.PRCount) / float64(maxPRs)) * 292.0
		if barW < 10.0 {
			barW = 10.0
		}
		fmt.Fprintf(&b, `  <text x="584" y="%d" fill="%s" font-family="%s" font-size="11">%s</text>`+"\n", y, p.inkPrimary, mono, html.EscapeString(short))
		fmt.Fprintf(&b, `  <text x="876" y="%d" text-anchor="end" fill="%s" font-family="%s" font-size="11" font-weight="600">%s</text>`+"\n", y, p.inkSecondary, mono, prLabel)
		fmt.Fprintf(&b, `  <rect x="584" y="%d" width="292" height="4" rx="2" fill="%s"/>`+"\n", y+6, p.trackBg)
		fmt.Fprintf(&b, `  <rect x="584" y="%d" width="%.1f" height="4" rx="2" fill="%s"/>`+"\n", y+6, barW, p.accentViolet)
	}

	// Horizontal hairline divider
	fmt.Fprintf(&b, `  <line x1="584" y1="204" x2="876" y2="204" stroke="%s" stroke-width="1"/>`+"\n", p.hairline)

	// Right column lower: By Repository Language
	fmt.Fprintf(&b, `  <text x="584" y="224" fill="%s" font-family="%s" font-size="9.5" font-weight="600" letter-spacing="0.9">BY REPOSITORY LANGUAGE</text>`+"\n", p.inkMuted, mono)
	fmt.Fprintf(&b, `  <rect x="584" y="234" width="292" height="7" rx="3.5" fill="%s"/>`+"\n", p.trackBg)

	totalLangPRs := 0
	for _, ls := range t.LangStats {
		totalLangPRs += ls.Count
	}
	if totalLangPRs == 0 {
		totalLangPRs = 1
	}

	fmt.Fprintf(&b, `  <g clip-path="url(#langBarClip-%s)">`+"\n", p.name)
	curX := 584.0
	for i, ls := range t.LangStats {
		segW := (float64(ls.Count) / float64(totalLangPRs)) * 292.0
		if segW < 4.0 {
			segW = 4.0
		}
		if i == len(t.LangStats)-1 {
			segW = (584.0 + 292.0) - curX
			if segW < 4.0 {
				segW = 4.0
			}
		}
		fmt.Fprintf(&b, `    <rect x="%.1f" y="234" width="%.1f" height="7" fill="%s"/>`+"\n", curX, segW, ls.Color)
		curX += segW
	}
	fmt.Fprintf(&b, `  </g>`+"\n")

	for i, ls := range t.LangStats {
		if i >= 4 {
			break
		}
		lx := 584 + i*74
		fmt.Fprintf(&b, `  <circle cx="%d" cy="259" r="3.5" fill="%s"/>`+"\n", lx+4, ls.Color)
		fmt.Fprintf(&b, `  <text x="%d" y="262.5" fill="%s" font-family="%s" font-size="10.5">%s %d</text>`+"\n", lx+11, p.inkSecondary, mono, html.EscapeString(ls.Short), ls.Count)
	}

	fmt.Fprintf(&b, `</svg>`+"\n")
	return b.String()
}

func updateReadmeMarkers(readmePath string, t *profileTelemetry) error {
	raw, err := os.ReadFile(readmePath)
	if err != nil {
		return err
	}
	content := string(raw)

	replacements := map[string]string{
		"merged_prs_count":     fmt.Sprintf("%d", t.MergedPRCount),
		"distinct_repos_count": fmt.Sprintf("%d", t.DistinctRepoCount),
		"total_releases":       fmt.Sprintf("%d", t.TotalReleases),
	}
	for repoName, info := range t.RepoReleases {
		if info.LatestTag != "" {
			replacements[repoName+"_tag"] = info.LatestTag
		}
		replacements[repoName+"_rels"] = fmt.Sprintf("%d", info.ReleaseCount)
	}

	for key, val := range replacements {
		pattern := fmt.Sprintf(`<!-- sync:%s -->.*?<!-- /sync:%s -->`, regexp.QuoteMeta(key), regexp.QuoteMeta(key))
		re := regexp.MustCompile(pattern)
		replacement := fmt.Sprintf(`<!-- sync:%s -->%s<!-- /sync:%s -->`, key, val, key)
		content = re.ReplaceAllString(content, replacement)
	}

	return os.WriteFile(readmePath, []byte(content), 0644)
}

func main() {
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GH_TOKEN"))
	}
	if token == "" {
		fmt.Fprintln(os.Stderr, "error: GITHUB_TOKEN or GH_TOKEN is required")
		os.Exit(1)
	}

	telemetry, err := fetchTelemetry(token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	rootDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	assetsDir := filepath.Join(rootDir, "assets")
	if err := os.MkdirAll(assetsDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: mkdir assets: %v\n", err)
		os.Exit(1)
	}

	outputs := map[string]string{
		filepath.Join(assetsDir, "header-dark.svg"):    renderHeaderSVG(telemetry, darkPalette),
		filepath.Join(assetsDir, "header-light.svg"):   renderHeaderSVG(telemetry, lightPalette),
		filepath.Join(assetsDir, "upstream-dark.svg"):  renderUpstreamSVG(telemetry, darkPalette),
		filepath.Join(assetsDir, "upstream-light.svg"): renderUpstreamSVG(telemetry, lightPalette),
	}
	for outPath, svg := range outputs {
		if err := os.WriteFile(outPath, []byte(svg), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error: write %s: %v\n", outPath, err)
			os.Exit(1)
		}
	}

	readmePath := filepath.Join(rootDir, "README.md")
	if _, err := os.Stat(readmePath); err == nil {
		if err := updateReadmeMarkers(readmePath, telemetry); err != nil {
			fmt.Fprintf(os.Stderr, "error: update README.md: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("Synced: %d merged PRs (%d repos), %d releases across %d flagships, %.2f MB (%d%% Go), %d 1Y contributions\n",
		telemetry.MergedPRCount, telemetry.DistinctRepoCount, telemetry.TotalReleases,
		telemetry.FlagshipCount, telemetry.FlagshipSourceMB, telemetry.GoPct, telemetry.TotalContributions)
}
