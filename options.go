package pctk

// AppOption is a function that can be used to configure the application.
type AppOption func(*App)

// WithScreenCaption sets the screen caption of the application.
func WithScreenCaption(caption string) AppOption {
	return func(a *App) { a.screenCaption = caption }
}

// WithScreenZoom sets the screen zoom of the application.
func WithScreenZoom(zoom int32) AppOption {
	return func(a *App) { a.screenZoom = zoom }
}

var defaultAppOptions = []AppOption{
	WithScreenCaption("Point&Click Toolkit"),
	WithScreenZoom(4),
}
