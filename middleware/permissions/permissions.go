// Package permissions2 provides a way to keep track of users, login states and permissions.
package permissions

import (
	"github.com/xyproto/pinterface"
	"net/http"
	"strings"
)

const (
	// Version number. Stable API within major version numbers.
	Version = 2.2
)

// The structure that keeps track of the permissions for various path prefixes
type Permissions struct {
	state              *UserState
	adminPathPrefixes  []string
	userPathPrefixes   []string
	publicPathPrefixes []string
	rootIsPublic       bool
	denied             http.HandlerFunc
}

// Initialize a Permissions struct with all the default settings.
// This will also connect to the redis host at localhost:6379.
func New() *Permissions {
	return NewPermissions(NewUserStateSimple())
}

// Initialize a Permissions struct with Redis DB index and host:port
func NewWithRedisConf(dbindex int, hostPort string) *Permissions {
	return NewPermissions(NewUserState(dbindex, true, hostPort))
}

// Initialize a Permissions struct with the given UserState and
// a few default paths for admin/user/public path prefixes.
func NewPermissions(state *UserState) *Permissions {
	// default permissions
	return &Permissions{state,
		[]string{"/admin"},         // admin path prefixes
		[]string{"/repo", "/data"}, // user path prefixes
		[]string{"/", "/login", "/register", "/favicon.ico", "/style", "/img", "/js",
			"/favicon.ico", "/robots.txt", "/sitemap_index.xml"}, // public
		true,
		PermissionDenied}
}

// Specify the http.HandlerFunc for when the permissions are denied
func (perm *Permissions) SetDenyFunction(f http.HandlerFunc) {
	perm.denied = f
}

// Get the current http.HandlerFunc for when permissions are denied
func (perm *Permissions) DenyFunction() http.HandlerFunc {
	return perm.denied
}

// Retrieve the UserState struct
func (perm *Permissions) UserState() pinterface.IUserState {
	return perm.state
}

// Set everything to public
func (perm *Permissions) Clear() {
	perm.adminPathPrefixes = []string{}
	perm.userPathPrefixes = []string{}
}

// Add an url path prefix that is a page for the logged in administrators
func (perm *Permissions) AddAdminPath(prefix string) {
	perm.adminPathPrefixes = append(perm.adminPathPrefixes, prefix)
}

// Add an url path prefix that is a page for the logged in users
func (perm *Permissions) AddUserPath(prefix string) {
	perm.userPathPrefixes = append(perm.userPathPrefixes, prefix)
}

// Add an url path prefix that is a public page
func (perm *Permissions) AddPublicPath(prefix string) {
	perm.publicPathPrefixes = append(perm.publicPathPrefixes, prefix)
}

// Set all url path prefixes that are for the logged in administrator pages
func (perm *Permissions) SetAdminPath(pathPrefixes []string) {
	perm.adminPathPrefixes = pathPrefixes
}

// Set all url path prefixes that are for the logged in user pages
func (perm *Permissions) SetUserPath(pathPrefixes []string) {
	perm.userPathPrefixes = pathPrefixes
}

// Set all url path prefixes that are for the public pages
func (perm *Permissions) SetPublicPath(pathPrefixes []string) {
	perm.publicPathPrefixes = pathPrefixes
}

// The default "permission denied" http handler.
func PermissionDenied(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "Permission denied.", http.StatusForbidden)
}

// Check if a given request should be rejected.
func (perm *Permissions) Rejected(w http.ResponseWriter, req *http.Request) bool {
	reject := false
	path := req.URL.Path // the path of the url that the user wish to visit

	// If it's not "/" and set to be public regardless of permissions
	if !(perm.rootIsPublic && path == "/") {

		// Reject if it is an admin page and user does not have admin permissions
		for _, prefix := range perm.adminPathPrefixes {
			if strings.HasPrefix(path, prefix) {
				if !perm.state.AdminRights(req) {
					reject = true
					break
				}
			}
		}

		if !reject {
			// Reject if it's a user page and the user does not have user rights
			for _, prefix := range perm.userPathPrefixes {
				if strings.HasPrefix(path, prefix) {
					if !perm.state.UserRights(req) {
						reject = true
						break
					}
				}
			}
		}

		if !reject {
			// Reject if it's not a public page
			found := false
			for _, prefix := range perm.publicPathPrefixes {
				if strings.HasPrefix(path, prefix) {
					found = true
					break
				}
			}
			if !found {
				reject = true
			}
		}

	}

	return reject
}

// Middleware handler (compatible with Negroni)
func (perm *Permissions) ServeHTTP(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	// Check if the user has the right admin/user rights
	if perm.Rejected(w, req) {
		// Get and call the Permission Denied function
		perm.DenyFunction()(w, req)
		// Reject the request by not calling the next handler below
		return
	}

	// Call the next middleware handler
	next(w, req)
}
