package udm

type SubscriptionData struct {
	UEID                string   `json:"ue_id"`
	AllowedPLMNs        []PLMN   `json:"allowed_plmns"`
	AllowedSNSSAIs      []SNSSAI `json:"allowed_snssais"`
	MobilityRestrictions string   `json:"mobility_restrictions"`
}

type AuthVector struct {
	RAND  string `json:"rand"`
	AUTN  string `json:"autn"`
	XRES  string `json:"xres"`
	KASME string `json:"kasme"`
}

type PLMN struct {
	MCC string `json:"mcc"`
	MNC string `json:"mnc"`
}

type SNSSAI struct {
	SST string `json:"sst"`
	SD  string `json:"sd"`
}

