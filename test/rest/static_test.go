package rest

// Static file serving is covered by TestIntegrationStaticFiles in
// integration_files_test.go (self-contained: temp dir + real requests +
// assertions). This file previously held TestStatic, which only registered
// a route and called api.StartService directly — no request, no assertion.
