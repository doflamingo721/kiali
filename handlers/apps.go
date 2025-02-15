package handlers

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/kiali/kiali/business"
)

// appListParams holds the path and query parameters for AppList
//
// swagger:parameters appList
type appListParams struct {
	baseHealthParams
	// The target workload
	//
	// in: path
	Namespace string `json:"namespace"`
	// Optional
	Health bool `json:"health"`
}

func (p *appListParams) extract(r *http.Request) {
	vars := mux.Vars(r)
	query := r.URL.Query()
	p.baseExtract(r, vars)
	p.Namespace = vars["namespace"]
	p.Health = query.Get("health") != ""
}

// AppList is the API handler to fetch all the apps to be displayed, related to a single namespace
func AppList(w http.ResponseWriter, r *http.Request) {
	p := appListParams{}
	p.extract(r)

	criteria := business.AppCriteria{Namespace: p.Namespace, IncludeIstioResources: true, Health: p.Health, RateInterval: p.RateInterval, QueryTime: p.QueryTime}

	// Get business layer
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Apps initialization error: "+err.Error())
		return
	}

	if criteria.Health {
		rateInterval, err := adjustRateInterval(r.Context(), business, p.Namespace, p.RateInterval, p.QueryTime)
		if err != nil {
			handleErrorResponse(w, err, "Adjust rate interval error: "+err.Error())
			return
		}
		criteria.RateInterval = rateInterval
	}

	// Fetch and build apps
	appList, err := business.App.GetAppList(r.Context(), criteria)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, appList)
}

// AppDetails is the API handler to fetch all details to be displayed, related to a single app
func AppDetails(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// Get business layer
	business, err := getBusiness(r)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Services initialization error: "+err.Error())
		return
	}
	namespace := params["namespace"]
	app := params["app"]

	// Fetch and build app
	appDetails, err := business.App.GetApp(r.Context(), namespace, app)
	if err != nil {
		handleErrorResponse(w, err)
		return
	}

	RespondWithJSON(w, http.StatusOK, appDetails)
}
