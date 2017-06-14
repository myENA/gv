package rest

import (
	// core
	"net/http"

	// gv package
	"github.com/myENA/gv"

	// RESTful interface
	"github.com/emicklei/go-restful"
)

// BuildInfo is a simple wrapper around gv.BuildInfo
type BuildInfo struct {
	*gv.BuildInfo
}

// New builds and returns a restful version object
func New(bi *gv.BuildInfo) *BuildInfo {
	return &BuildInfo{
		BuildInfo: bi,
	}
}

// return json feed of our version struct
func (b *BuildInfo) getVersion(request *restful.Request, response *restful.Response) {
	if err := response.WriteEntity(b); err != nil {
		response.WriteServiceError(http.StatusInternalServerError,
			restful.NewError(http.StatusInternalServerError, err.Error()))
	}
}

// return yaml configuration suitable for reproducing an identical build
func (b *BuildInfo) getYAML(request *restful.Request, response *restful.Response) {
	var yml []byte // returned yaml
	var err error  // error holder

	// nothing we can do if config is undefined
	if b.GlideConfig == nil {
		response.WriteErrorString(http.StatusNoContent, "missing glide config")
		return
	}

	// attempt to generate yaml
	if yml, err = b.GlideConfig.Marshal(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	// return our response
	response.ResponseWriter.Header().Set("Content-Type", "application/x-yaml")
	response.ResponseWriter.Header().Set("Content-Disposition", "inline; filename=\"glide.yaml\"")
	response.ResponseWriter.Write(y)
}

// return yaml configuration suitable for reproducing an identical build
func (b *BuildInfo) getLockfile(request *restful.Request, response *restful.Response) {
	var yml []byte // returned yaml
	var err error  // error holder

	// nothing we can do if config is undefined
	if b.GlideLockfile == nil {
		response.WriteErrorString(http.StatusNoContent, "missing glide Lockfile")
		return
	}

	// attempt to generate yaml
	if yml, err = b.GlideLockfile.Marshal(); err != nil {
		response.WriteError(http.StatusInternalServerError, err)
		return
	}

	// return our response
	response.ResponseWriter.Header().Set("Content-Type", "application/x-yaml")
	response.ResponseWriter.Header().Set("Content-Disposition", "inline; filename=\"glide.lock\"")
	response.ResponseWriter.Write(yml)
}

// Routes returns restful service routes
func (b *BuildInfo) Routes(ws *restful.WebService) {
	// json route
	ws.Route(ws.GET("/version").
		To(b.getVersion).
		Doc("Get available version data").
		Notes("Return JSON formatted version data specific "+
			"to this application and it's dependencies.").
		Operation("getVersion").
		Returns(http.StatusAccepted, "Okay", nil).
		Returns(http.StatusInternalServerError, "Error", restful.ServiceError{}).
		Writes(BuildInfo{}).
		Produces(restful.MIME_JSON))

	// yaml route
	ws.Route(ws.GET("/glide/config").
		To(b.getYAML).
		Doc("Get glide YAML").
		Notes("Return glide YAML configuration suitable "+
			"for reproducing a build with matching dependencies").
		Operation("getYAML").
		Returns(http.StatusAccepted, "Okay", nil).
		Returns(http.StatusNoContent, "Missing", nil).
		Returns(http.StatusInternalServerError, "Error", nil).
		Produces("application/x-yaml"))

	// yaml route
	ws.Route(ws.GET("/glide/lock").
		To(b.getLockfile).
		Doc("Get glide Lockfile YAML").
		Notes("Return glide YAML lockfile suitable "+
			"for reproducing a build with matching dependencies").
		Operation("getLockfile").
		Returns(http.StatusAccepted, "Okay", nil).
		Returns(http.StatusNoContent, "Missing", nil).
		Returns(http.StatusInternalServerError, "Error", nil).
		Produces("application/x-yaml"))
}
