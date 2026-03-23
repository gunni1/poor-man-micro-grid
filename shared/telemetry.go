package shared

type Telemetry struct {
	AssetId   string  `json:"asset_id"`
	AssetType string  `json:"asset_type"`
	PkW       float64 `json:"p_kw"`
}
