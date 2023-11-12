package lib

// Structure for using on WEB
type RouteDetails struct {
	RouteID          	string   `form:"route_id" 			 	json:"RouteID"`
	RouterName          string   `form:"router_name"		 	json:"RouterName"`
	RouteIP		        string   `form:"route_ip"			 	json:"RouteIP"`
	RouteMask	        string   `form:"route_mask" 		 	json:"RouteMask"`
	Description         string   `form:"description" 		 	json:"Description"`
	CSRFToken           string   `form:"csrftoken" 			 	json:"CSRFToken"`
}

type NameSorterRoutesDetails []*RouteDetails

func (a NameSorterRoutesDetails) Len() int           { return len(a) }
func (a NameSorterRoutesDetails) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorterRoutesDetails) Less(i, j int) bool { return a[i].RouterName < a[j].RouterName }
