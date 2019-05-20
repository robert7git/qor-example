// +build enterprise

package migrations

import "github.com/dfang/qor-example/app/enterprise"

func init() {
	AutoMigrate(&enterprise.QorMicroSite{})
}
