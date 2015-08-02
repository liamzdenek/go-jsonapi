package authenticator

import . ".."
import "fmt"

type RBAC struct {
	PermissionLookup, UserPermissionLookup Resource
	GetUserId                              RBACGetUserId
}

type RBACGetUserId interface {
	GetUserId(*Request) *string
}

func NewRBAC(perm_lookup, user_perm_lookup Resource, get_user_id RBACGetUserId) *RBAC {
	return &RBAC{
		PermissionLookup:     perm_lookup,
		UserPermissionLookup: user_perm_lookup,
		GetUserId:            get_user_id,
	}
}

func (r *RBAC) Require(permission string) Authenticator {
	return &RBACAuthenticator{
		Permission: permission,
		RBAC:       r,
	}
}

// CREATE TABLE rbac_permissions (id INT NOT NULL AUTO_INCREMENT, name VARCHAR(64), `default` BOOL, PRIMARY KEY(id), INDEX(name));
type RBACPermissionLookup struct {
	Id      int    `jsonapi:"id" meddler:"id,pk"`
	Name    string `meddler:"name"`
	Default bool   `meddler:"default"`
}

//CREATE TABLE rbac_user_permissions (id INT NOT NULL AUTO_INCREMENT, user_id INT NOT NULL, permission_id INT NOT NULL, has_permission BOOL, PRIMARY KEY(id), INDEX(user_id, permission_id));
type RBACUserPermissionLookup struct {
	Id            int    `jsonapi:"id" meddler:"id,pk"`
	UserId        string `meddler:"user_id"`
	PermissionId  string `meddler:"permission_id"`
	HasPermission bool   `meddler:"has_permission"`
}

type RBACAuthenticator struct {
	RBAC       *RBAC
	Permission string
}

func (a *RBACAuthenticator) Authenticate(r *Request, permission, id string) {
	user_id := a.RBAC.GetUserId.GetUserId(r)
	if user_id == nil {
		panic("User is not logged in for a route that reuires a permission")
	}

	records, err := a.RBAC.PermissionLookup.FindManyByField(r, RequestParams{}, "Name", a.Permission)
	Check(err)
	if len(records) != 1 {
		panic(fmt.Sprintf("%d results for permission name=%s were found -- expected exactly one.", len(records), a.Permission))
	}

	default_permission := records[0].Attributes.(*RBACPermissionLookup).Default

	user_permissions, err := a.RBAC.UserPermissionLookup.FindManyByField(r, RequestParams{}, "UserId", *user_id)
	Check(err)

	for _, record := range user_permissions {
		user_permission := record.Attributes.(*RBACUserPermissionLookup)
		r.API.Logger.Debugf("USER GOT PERMISSION: %#v\n", user_permission)
		if user_permission.PermissionId == records[0].Id {
			if !user_permission.HasPermission {
				panic(InsufficientPermissions())
			}
			return
		}
	}

	if !default_permission {
		panic(InsufficientPermissions())
	}
}
