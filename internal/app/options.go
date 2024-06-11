package app

// Option defines the function used for
// setup an application option config
type Option func(a *Application)

// WithHTTPPort function adds the given http
// port into the application base config
func WithHTTPPort(port string) Option {
	return func(a *Application) {
		a.httpPort = port
	}
}

// WithDBHost function adds the given db
// host into the application base config
func WithDBHost(host string) Option {
	return func(a *Application) {
		a.dbHost = host
	}
}

// WithDBPort function adds the given db
// port into the application base config
func WithDBPort(port string) Option {
	return func(a *Application) {
		a.dbPort = port
	}
}

// WithDBUser function adds the given db
// user into the application base config
func WithDBUser(user string) Option {
	return func(a *Application) {
		a.dbUser = user
	}
}

// WithDBPassword function adds the given db
// password into the application base config
func WithDBPassword(password string) Option {
	return func(a *Application) {
		a.dbPassword = password
	}
}

// WithDBName function adds the given db
// name into the application base config
func WithDBName(name string) Option {
	return func(a *Application) {
		a.dbName = name
	}
}

// WithShoppingCartHTTPEndpoint function adds the given HTTP endpoint
// into the application base config
func WithShoppingCartHTTPEndpoint(port string) Option {
	return func(a *Application) {
		a.ShoppingCartHTTPEndpoint = port
	}
}
