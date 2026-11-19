# “GitHub for Everything” Playbook

We should aim to use GitHub as the central platform for nearly all aspects of our workflows/processes—from code and CI/CD to planning, documentation, and internal collaboration.

By doing this, we reduce our reliance on a fragmented set of third-party tools and avoid the overhead of maintaining, integrating, and context-switching between multiple systems.

This approach brings several advantages:

* Single source of truth: Everything lives in one place—code, tasks, discussions, and documentation—making it easier to onboard new team members and track progress.  
* Tighter feedback loops: Developers can stay in flow, reviewing code, updating docs, and closing issues without leaving GitHub.  
* Automation-first mindset: GitHub Actions allows us to automate testing, deployments, and even project workflows.  
* Lower operational complexity: Fewer SaaS contracts, fewer logins, fewer integration headaches.  
* Security and auditability: Centralizing work in GitHub gives us a clear, versioned trail of decisions, changes, and access.

|  | Use Case | GitHub Features | Alternatives Avoided | Notes |
| :---- | :---- | :---- | :---- | :---- |
| **1** | **Source Code Management** | Repositories, Branching, Pull Requests, GitHub Actions | Bitbucket, GitLab, SVN | Leverage protected branches, code owners, PR reviews. |
| **2** | **CI/CD Pipelines** | GitHub Actions | Jenkins, CircleCI, GitLab CI | Deploy to AWS, test across stacks, use reusable workflows and matrix jobs. |
| **3** | **Project & Sprint Management** | GitHub Projects (v2?), Issues, Milestones, Labels | Jira, Trello, Linear | Track epics, bugs, features using labels \+ Projects. Automate status changes. |
| **4** | **Team Communication Around Work** | PR Comments, Issue Threads, Discussions | Slack (for dev work), Discord, Teams | Keep async conversations close to the work itself. |
| **5** | **Product Specs / Design Docs** | Markdown files, Wiki, Discussions | Notion, Google Docs, Confluence | Store docs in \`docs/\` or Wiki. Use Discussions for structured debates. We should “migrate” our Google Docs over to GH. Note \- it is not hard to what we do on GD on GH, even things like meeting notes. Many open source projects record meeting notes on GH.  |
| **6** | **Technical Roadmaps** | Milestones, Labels, Projects | Productboard, spreadsheets | Milestones \= release targets. Projects \= roadmap swimlanes. |
| **7** | **Release Management** | Releases, Tags, Automated Changelogs | Manual tracking, release dashboards | Use PR labels or semantic titles to generate changelogs. |
| **8** | **Security Audits & Vulnerability Scanning** | Dependabot, Code/Secret Scanning | Snyk, SonarQube | Enable across all repos from day one. Enforce updates via PRs. |
| **9** | **Internal Tooling and Scripts** | GitHub Actions, CLI Scripts in Repos | Jenkins jobs, internal web dashboards | Automate everything—backups, syncs, linting, reporting. |
| **10** | **Onboarding Docs** | README.md, CONTRIBUTING.md, Wiki | Notion, HR portals | \`onboarding/\` folder with dev setup, team culture, workflow guide. |
| **11** | **Hiring Technical Talent** | GitHub itself, Issues for trial tasks | Workable, HackerRank, Lever | Use Issues \+ PRs to simulate real work during trials. |
| **12** | **Community Building (if open source)** | Issues, Discussions, Pull Requests | Discourse, Reddit, standalone forums | Run feedback, support, and feature requests through GitHub itself. |
| **13** | **Blog / Public Content** | GitHub Pages \+ Jekyll/Hugo \+ Actions | Medium, Ghost, Substack | Static blog hosted directly from a repo, auto-published on push. |
| **14** | **Customer Support Knowledge Base (Early Stage)** | GitHub Pages \+ Markdown Docs | Zendesk, Intercom, Freshdesk | A public \`docs/\` repo published as support FAQ. |

# Caveats: When GitHub won’t be Enough

While GitHub can cover most of our early-stage needs, it has limitations. As we scale, some gaps will emerge where dedicated tools are more practical or necessary.

* **Internal communication**: We need some form of real-time chat. We will likely use Google Chat with our @opensource.lk email accounts.
* **Non-technical teams**: Marketing, sales, and ops may find GitHub too developer-centric. We will need alternatives when we reach that stage.  
* **Customer support**: Ditto. .  
* **Analytics & metrics**: GitHub lacks product analytics and business intelligence-type tools. We will need to rely on some third-party product, ideally coupled to the cloud provider we use (e.g. AWS).
* What else?
