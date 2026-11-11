// Package identity defines saved accounts and effective Git identities.
package identity

// Account is a saved identity template. It does not contain login credentials.
type Account struct {
	Label string `json:"label"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Resolved records the effective Git values and their configuration scopes.
type Resolved struct {
	Name, Email           string
	NameScope, EmailScope string
}

// Matches compares the commit identity, ignoring the account's display label.
func (v Resolved) Matches(a Account) bool { return v.Name == a.Name && v.Email == a.Email }

// RepositoryScope reports whether a value overrides the global identity.
func RepositoryScope(scope string) bool { return scope == "local" || scope == "worktree" }

// Source returns the source label that the UI can translate for display.
func (v Resolved) Source() string {
	if v.NameScope == "" && v.EmailScope == "" {
		return "未设置"
	}
	if v.NameScope != v.EmailScope {
		return "混合"
	}
	switch v.NameScope {
	case "local", "worktree":
		return "专用"
	case "global":
		return "继承"
	default:
		return v.NameScope
	}
}
