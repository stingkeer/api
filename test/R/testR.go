package r

import "gitee.com/fast_api/api"

func Test(f func()) {
	f()
	api.StartService(nil)
}
