package models

type TransportMode struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Factor float64 `json:"factor"`
}

type CarbonCalculation struct {
	Mode     string  `json:"mode"`
	Distance float64 `json:"distance"`
	Result   float64 `json:"result"`
}

type CarbonFactors struct {
	Transports struct {
		Train  float64 `json:"train"`
		Flight float64 `json:"flight"`
		Car    struct {
			Small  float64 `json:"small"`
			Medium float64 `json:"medium"`
			Big    float64 `json:"big"`
		} `json:"car"`
	} `json:"Transports"`
	LogementElectromenagers struct {
		Electricity float64 `json:"electricity"`
		Gas         float64 `json:"gas"`
		Apartment   float64 `json:"apartment"`
		House       float64 `json:"house"`
		Appliance   float64 `json:"appliance"`
		Electronic  float64 `json:"electronic"`
	} `json:"Logement_electromenagers"`
	Alimentation struct {
		RedMeat          float64 `json:"redMeat"`
		WhiteMeat        float64 `json:"whiteMeat"`
		Pork             float64 `json:"pork"`
		BulkFoodPurchase struct {
			None    float64 `json:"none"`
			Partial float64 `json:"partial"`
			Total   float64 `json:"total"`
		} `json:"bulkFoodPurchase"`
		ShortCircuit struct {
			None     float64 `json:"none"`
			Partial  float64 `json:"partial"`
			Majority float64 `json:"majority"`
		} `json:"shortCircuit"`
	} `json:"Alimentation"`
	Vetements struct {
		Large  float64 `json:"large"`
		Small  float64 `json:"small"`
		Madein struct {
			France float64 `json:"france"`
			Autre  float64 `json:"autre"`
		} `json:"madein"`
	} `json:"Vetements"`
	Numerique struct {
		GoogleSearch float64 `json:"googleSearch"`
		ChatGPT      float64 `json:"chatGPT"`
		SocialMedia  float64 `json:"socialMedia"`
		Smartphone   struct {
			Small float64 `json:"small"`
			Large float64 `json:"large"`
			Used  float64 `json:"used"`
			Old   float64 `json:"old"`
		} `json:"smartphone"`
	} `json:"Numerique"`
	Consommation struct {
		Ecommerce struct {
			Amazon    float64 `json:"amazon"`
			LeBonCoin float64 `json:"leboncoin"`
			Artisanat float64 `json:"artisanat"`
		} `json:"ecommerce"`
		Commerce struct {
			Brocante   float64 `json:"brocante"`
			LocalShops float64 `json:"localShops"`
		} `json:"commerce"`
	} `json:"Consommation"`
}
